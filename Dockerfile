FROM golang:latest

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux

COPY ./ ./

RUN apt install make

CMD ["make", "run"]
