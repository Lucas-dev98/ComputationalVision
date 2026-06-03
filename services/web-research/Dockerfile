FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o web-research-service .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/web-research-service ./web-research-service

EXPOSE 8083
ENV PORT=8083

CMD ["./web-research-service"]
