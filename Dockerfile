FROM golang:1.22-alpine3.20

RUN apk add --no-cache gcc musl-dev pkgconfig openssl-dev

WORKDIR /draft

COPY . .

RUN go mod tidy

CMD go run cmd/kv/main.go
