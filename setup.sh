#!/bin/bash

# Script de setup automático do projeto
# Uso: ./setup.sh

set -e

echo "🚀 Configurando Sistema de Entrada de Estoque..."
echo ""

# Verificar pré-requisitos
echo "📋 Verificando pré-requisitos..."

if ! command -v git &> /dev/null; then
    echo "❌ Git não instalado"
    exit 1
fi

if ! command -v docker &> /dev/null; then
    echo "❌ Docker não instalado"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose não instalado"
    exit 1
fi

echo "✓ Pré-requisitos OK"
echo ""

# Configurar .env
if [ ! -f .env ]; then
    echo "📝 Criando arquivo .env..."
    cp .env.example .env
    echo "✓ .env criado. Edite conforme necessário."
    echo ""
fi

# Iniciando serviços
echo "🐳 Iniciando serviços Docker..."
docker-compose -f infra/docker/docker-compose.yml up -d

echo ""
echo "⏳ Aguardando serviços ficarem prontos..."
sleep 10

# Verificar status
echo ""
echo "📊 Status dos serviços:"
docker-compose -f infra/docker/docker-compose.yml ps

echo ""
echo "✅ Setup concluído!"
echo ""
echo "📍 Acesse:"
echo "  Frontend:      http://localhost:3000"
echo "  Inventory API: http://localhost:8080/health"
echo "  OCR Service:   http://localhost:5001/health"
echo ""
echo "💡 Próximos passos:"
echo "  1. Acesse http://localhost:3000"
echo "  2. Capture uma foto com a webcam"
echo "  3. Aprove o item para registrar no estoque"
echo ""
echo "📚 Documentação: ver /docs ou QUICKSTART.md"
echo ""
