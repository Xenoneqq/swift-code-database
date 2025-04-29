FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

RUN chmod +x ./setup_main.sh
RUN chmod +x ./setup_test.sh
RUN chmod +x ./run_tests.sh

RUN apk add --no-cache curl

RUN go build -o main .

EXPOSE 8080



