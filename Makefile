boot-consul:
	@ docker run -d -p "8500:8500" -p "8600:8600" --name=discovery hashicorp/consul agent -server -ui -node=server-1  -bootstrap-expect=1 -client="0.0.0.0"

build-login-service:
	@ docker build .\Login-Service\ -t login-service

build-order-service:
	@ docker build .\Order-Service\ -t order-service

build-product-service:
	@ docker build .\Product-Service\ -t product-service

build-images:
	@echo "Building Login-Service Image..."
	@build-login-service
	@echo "Building Order-Service Image..."
	@build-order-service
	@echo "Building Product-Service Image..."
	@build-product-service
	@echo "Build Done !"

run-app:
	@docker compose up -d

.PHONY: run-app build-images