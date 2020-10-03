FROM golang:1.15-alpine AS build
WORKDIR /rmqbe
COPY . .
RUN go build ./cmd/rmqbe

FROM alpine

ENV MONGODB_URI "mongodb://mongo:27017"

COPY --from=build /rmqbe/rmqbe /
ENTRYPOINT [ "/rmqbe" ]
