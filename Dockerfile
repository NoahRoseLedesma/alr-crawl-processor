FROM golang:1.17-alpine
WORKDIR /app
COPY main.go go.mod go.sum /app/
RUN go install
ENTRYPOINT go run main.go
