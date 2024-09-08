generate:
	protoc --go_out=internal/pb --go_opt=paths=source_relative \
    --go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative \
    api/cdn-balancer.proto

test:
	go test ./...

build:
	sudo docker compose build

deploy:
	sudo docker compose up
