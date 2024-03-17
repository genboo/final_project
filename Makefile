build:
	CGO_ENABLED=0 GOOS=linux go build -o /main main.go

run:
	./main

test:
	go test ./...