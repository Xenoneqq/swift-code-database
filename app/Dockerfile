FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

RUN chmod +x ./database_check.sh

RUN go build -o main .

EXPOSE 8080

CMD ["sh", "./database_check.sh", "-t", "60", "${DB_HOST}:${DB_PORT}"]



