# application — REST API + MySQL. Multi-stage: build a static binary, then run
# it on a tiny base so the image stays small and has no toolchain at runtime.
FROM golang:1.24-alpine AS build
WORKDIR /src

# Cache module downloads separately from the source so code edits don't redo them.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO off -> a fully static binary that runs on the scratch-like alpine runtime.
RUN CGO_ENABLED=0 go build -o /out/application .

FROM alpine:3.20
# curl is used by the compose healthcheck to probe the service.
RUN apk add --no-cache curl ca-certificates
COPY --from=build /out/application /usr/local/bin/application
EXPOSE 8081
ENTRYPOINT ["application"]
