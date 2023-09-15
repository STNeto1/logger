tcp:
	go run cmd/tcp/main.go

http:
	go run cmd/http/main.go

tcp-client:
	go run cmd/client/main.go --mode=tcp --msgs=1000

http-client:
	go run cmd/client/main.go --mode=http --msgs=1000
