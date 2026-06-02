FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instalar dependências
RUN apk add --no-cache git gcc musl-dev

# Copiar go.mod e go.sum
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código
COPY . .

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/inventory-service .

# Final stage
FROM alpine:latest

WORKDIR /app

# Instalar CA certificates
RUN apk --no-cache add ca-certificates

# Copiar binary
COPY --from=builder /app/inventory-service .

# Expor porta
EXPOSE 8080

# Comando de inicialização
CMD ["./inventory-service"]
