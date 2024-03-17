build:
	CGO_ENABLED=0 GOOS=linux go build -o /main main.go

run: build
	./main

test:
	go test ./...