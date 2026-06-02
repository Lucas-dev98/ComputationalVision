FROM node:20-alpine AS builder

WORKDIR /app

# Copiar tudo
COPY . .

# Instalar dependências com força
RUN npm install --force

# Build
RUN CI=false npm run build

# Final stage - usar serve para servir a aplicação
FROM node:20-alpine

WORKDIR /app

# Instalar serve
RUN npm install -g serve@14

# Copiar build do builder
COPY --from=builder /app/build ./build

# Expor porta
EXPOSE 3000

# Servir a aplicação
CMD ["serve", "-s", "build", "-l", "3000"]

# Comando de inicialização
CMD ["serve", "-s", "build", "-l", "3000"]
