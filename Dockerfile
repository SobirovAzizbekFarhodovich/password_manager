# Go Backend Dockerfile
FROM golang:1.22.3-alpine
WORKDIR /app
COPY . /app
COPY go.mod ./
COPY go.sum ./
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]
