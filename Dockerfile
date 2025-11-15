FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o reviewer ./cmd/app

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/reviewer .

EXPOSE 8080

CMD ["./reviewer"]
