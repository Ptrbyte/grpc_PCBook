gen:
	protoc --go_out=. proto/*.proto
	protoc --go-grpc_out=. proto/laptop_service.proto

clean:
	rm pb/*.go

run:
	go run cmd/server/main.go -port 8080
	go run cmd/client/main.go -address 0.0.0.0:8080

test:
	go test -cover -race ./...