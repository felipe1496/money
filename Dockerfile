# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Instala dependências necessárias para compilação
RUN apk add --no-cache git

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos de dependências
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o código fonte
COPY . .

# Compila a aplicação
# CGO_ENABLED=0 cria um binário estático
# -ldflags="-w -s" reduz o tamanho do binário
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api

# Torna o binário executável
RUN chmod +x main

# Stage 2: Runtime
FROM alpine:latest

# Instala certificados CA para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copia o binário do stage de build
COPY --from=builder /app/main .

# Se você tiver arquivos de configuração, templates, etc:
# COPY --from=builder /app/config ./config
# COPY --from=builder /app/templates ./templates

# Expõe a porta (ajuste conforme necessário)
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"]