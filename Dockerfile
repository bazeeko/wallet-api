# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /my-app

# COPY go.mod ./
# COPY go.sum ./

COPY ./ ./

RUN go mod download

RUN go build -o main ./main.go

EXPOSE 8080

CMD ["./main"]