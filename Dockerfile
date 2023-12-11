FROM golang:latest

WORKDIR /app
COPY go.mod ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux

COPY *.go ./

RUN go build -o ./ConnectIM-Server


CMD ["./ConnectIM-Server"]
