FROM golang:1.15.2-alpine3.12 AS builder

WORKDIR /rmqbe

COPY . .
RUN go build ./cmd/rmqbe

FROM alpine:3.12.0

ENV MONGODB_URI "mongodb://mongo:27017"

COPY --from=builder /rmqbe/rmqbe /

ENTRYPOINT [ "/rmqbe" ]
