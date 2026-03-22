srv:
	@bash build.sh build

gen:
	@bash build.sh gen

lint:
	golangci-lint run ./...

check:
	true

go-test:
	go test ./...
