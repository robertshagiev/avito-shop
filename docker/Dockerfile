FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/merch-shop ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/merch-shop /app/merch-shop

CMD ["/app/merch-shop"]