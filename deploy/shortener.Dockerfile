FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/shortener cmd/shortener/main.go cmd/shortener/app.go

FROM alpine:3.22 AS runtime

RUN apk --no-cache add curl

WORKDIR /app
COPY --from=builder /app/bin/shortener ./
COPY --from=builder /app/config/config.yaml ./config/

EXPOSE 8082

ENTRYPOINT ["./shortener"]
