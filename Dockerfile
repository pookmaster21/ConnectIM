FROM golang:latest

WORKDIR /app
COPY go.mod ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux

COPY ./ ./

RUN go build -o ./bin/ConnectIM-Server

CMD ["./bin/ConnectIM-Server"]
