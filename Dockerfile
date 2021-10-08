FROM golang:alpine

WORKDIR /build

COPY go.mod ./
COPY *.go ./
COPY entrypoint.sh /

RUN go mod download

RUN go build -o /app

ENTRYPOINT ["/entrypoint.sh"]
