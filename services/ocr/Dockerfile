FROM python:3.11-slim

WORKDIR /app

# Instalar dependências do sistema para OpenCV e build
RUN apt-get update && apt-get install -y \
    build-essential \
    make \
    libsm6 \
    libxext6 \
    libxrender-dev \
    && rm -rf /var/lib/apt/lists/*

# Copiar requirements
COPY requirements.txt .

# Instalar dependências Python
RUN pip install --no-cache-dir -r requirements.txt

# Copiar código
COPY . .

# Expor porta
EXPOSE 5001

# Comando de inicialização
CMD ["python", "main.py"]
