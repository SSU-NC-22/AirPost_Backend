# logic-core — Kafka consumer + Elasticsearch indexer + email rules.
FROM golang:1.24-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/logic-core .

FROM alpine:3.20
RUN apk add --no-cache curl ca-certificates
COPY --from=build /out/logic-core /usr/local/bin/logic-core
EXPOSE 8084
ENTRYPOINT ["logic-core"]
