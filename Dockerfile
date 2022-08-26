# syntax=docker/dockerfile:1
FROM golang:1.19.0
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
COPY ./constants constants
COPY ./controllers controllers
COPY ./db db
COPY ./models models
COPY ./routers routers
COPY ./utils utils
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh
RUN go mod tidy
RUN go build -o ./funnow-url-shortener
EXPOSE 8080
CMD ["/wait-for-it.sh", "redis:6379", "--", "./funnow-url-shortener"]