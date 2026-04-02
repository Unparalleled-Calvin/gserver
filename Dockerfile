# syntax=docker/dockerfile:1

FROM golang:1.26.1-alpine AS builder
WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /out/server /app/server

ENV GSERVER_ADDR=localhost:8000

EXPOSE 8000
CMD ["/app/server"]
