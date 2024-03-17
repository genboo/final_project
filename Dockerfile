FROM golang:1.20.4 as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN make build

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder main /bin/main
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
WORKDIR /bin
RUN make run
