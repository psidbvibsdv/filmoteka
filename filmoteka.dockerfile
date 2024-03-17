FROM golang:1.21 AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o filmotekaApp ./cmd/api

FROM alpine:3.14

RUN mkdir /app


COPY --from=builder /app/filmotekaApp /app/filmotekaApp

CMD ["/app/filmotekaApp"]