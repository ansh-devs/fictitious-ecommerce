package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ansh-devs/commercelens/order-service/db"
	"github.com/ansh-devs/commercelens/order-service/endpoints"
	"github.com/ansh-devs/commercelens/order-service/repo"
	"github.com/ansh-devs/commercelens/order-service/service"
	"github.com/ansh-devs/commercelens/order-service/transport"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	_ "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func main() {
	var httpAddr = os.Getenv("HTTPPORT")
	tracer := opentracing.GlobalTracer()
	cfg := &config.Configuration{
		ServiceName: "Order Service",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			// LocalAgentHostPort - Explicitly giving jaeger host to connect as defined in docker-compose file...
			LocalAgentHostPort: "tracer:6831",
			LogSpans:           true,
		},
	}

	errs := make(chan error)
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		errs <- fmt.Errorf("%s", err)
	}
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			errs <- fmt.Errorf("%s", err)
		}
	}(closer)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "order_service",
			"time", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}
	_ = level.Info(logger).Log("msg", "service started")

	flag.Parse()
	var srv *service.OrderService
	{

		var dbSource = fmt.Sprintf("postgres://%s:%s@%s/%s",
			os.Getenv("DBUSER"),
			os.Getenv("DBPASSWORD"),
			os.Getenv("DBHOST"),
			os.Getenv("DBNAME"),
		)

		dbConn := db.MustConnectToPostgress(dbSource)
		repository := repo.NewPostgresRepository(dbConn, logger, tracer)
		srv = service.NewOrderService(repository, logger, tracer)
		srv.InitNATS()
	}

	go srv.PlaceOrder(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoint := endpoints.NewEndpoints(srv)
	srv.RegisterService(&httpAddr)
	go srv.UpdateHealthStatus()

	go func() {
		fmt.Println("listening on port", httpAddr)
		handler := transport.NewHttpServer(endpoint)
		errs <- http.ListenAndServe(httpAddr, handler)
	}()

	for sig := range errs {
		_ = level.Error(logger).Log("status", sig, "GRACEFULLY SHUTTING DOWN !")
		_ = srv.ConsulClient.Agent().ServiceDeregister(srv.SrvID)
		os.Exit(0)
	}

}
