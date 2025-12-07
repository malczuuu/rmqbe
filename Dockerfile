FROM golang:1.24-alpine AS builder

WORKDIR /build

RUN apk add --no-cache curl git
RUN curl -sL https://taskfile.dev/install.sh | sh
RUN ./bin/task --version

COPY . .
RUN ./bin/task build-prod

FROM alpine:3.19

ENV MONGODB_URI "mongodb://mongo:27017"

RUN apk --no-cache add ca-certificates dumb-init

WORKDIR /rmqbe

COPY --from=builder /build/dist/rmqbe /rmqbe

EXPOSE 8000

# use dumb-init to handle signals properly
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/rmqbe/rmqbe"]
