FROM node:20-alpine AS builder

WORKDIR /app

# Copiar package.json
COPY package*.json ./

# Instalar dependências
RUN npm ci

# Copiar código
COPY . .

# Build
RUN npm run build

# Final stage
FROM node:20-alpine

WORKDIR /app

# Instalar serve para servir a aplicação
RUN npm install -g serve

# Copiar build do builder
COPY --from=builder /app/build ./build

# Expor porta
EXPOSE 3000

# Comando de inicialização
CMD ["serve", "-s", "build", "-l", "3000"]
