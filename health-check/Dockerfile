# health-check — live drone-tracking WebSocket (HTTP :8083, WS listen :8085).
FROM golang:1.24-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/health-check .

FROM alpine:3.20
RUN apk add --no-cache curl ca-certificates
COPY --from=build /out/health-check /usr/local/bin/health-check
EXPOSE 8083 8085
ENTRYPOINT ["health-check"]
