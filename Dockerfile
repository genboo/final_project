FROM golang:1.20.4 as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
COPY certs/ /certs
RUN CGO_ENABLED=0 GOOS=linux go build -o /main main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder main /bin/main
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /certs /bin/certs
WORKDIR /bin
ENTRYPOINT ["/bin/main"]
