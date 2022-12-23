srv: srv.go internal/*/*.go internal/*/*/*.go
	go build -o srv -ldflags '-w -linkmode external -extldflags "-static"' ./

go-test:
	go test ./...

lint:
	golangci-lint run ./...
