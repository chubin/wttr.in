srv: srv.go internal/*/*.go internal/*/*/*.go
	go build -o srv ./

go-test:
	go test ./...
