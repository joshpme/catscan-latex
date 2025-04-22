build:
	go build -o app cmd/main.go

test:
	go test -v ./...

start:
	func start