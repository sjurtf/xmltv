FROM golang:1.17.4 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /app
RUN go build -o binary .

FROM alpine:latest

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /app/binary .

USER app
ENTRYPOINT ["./binary"]