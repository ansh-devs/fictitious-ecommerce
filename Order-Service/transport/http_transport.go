package transport

import (
	"github.com/ansh-devs/ecomm-poc/order-service/endpoints"
	transport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)
// NewHttpServer makes the http routes for the endpoints.
func NewHttpServer(endpoints *endpoints.HttpEndpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(JsonTypeReWrittermiddleware)

	/*	r.Methods("POST").Path("/v1/orders/place-order").Handler(transport.NewServer(
		endpoints.PlaceOrder,
		JsonPlaceOrderResponseDecoder,
		JsonResponseEncoder,
	))*/

	r.Methods("GET").Path("/orders/v1/get-order/{id}").Handler(transport.NewServer(
		endpoints.GetOrder,
		JsonGetOrderResponseDecoder,
		JsonResponseEncoder,
	))

	r.Methods("GET").Path("/orders/v1/cancel-order").Handler(transport.NewServer(
		endpoints.CancelOrder,
		JsonCancelOrderResponseDecoder,
		JsonResponseEncoder,
	))

	r.Methods("POST").Path("/orders/v1/get-user-all-orders").Handler(transport.NewServer(
		endpoints.GetAllUserOrders,
		JsonGetAllUserOrdersResponseDecoder,
		JsonResponseEncoder,
	))
	return r
}
