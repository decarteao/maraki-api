# Usando uma imagem leve de Go
FROM golang:1.20-alpine

# Configurando o diretório de trabalho
WORKDIR /app

# Copiando o Go module e o Go.sum para o diretório de trabalho
COPY go.mod go.sum ./

# Baixando as dependências
RUN go mod tidy

# Copiando o código para o diretório de trabalho
COPY . .

# Compilando o binário
RUN go build -o main .

# Expondo a porta do serviço
EXPOSE 8080

# Comando para rodar o serviço
CMD ["./main"]
