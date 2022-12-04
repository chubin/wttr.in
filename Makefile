srv: srv.go internal/*/*.go
	go build -o srv ./
test:
	go test ./
