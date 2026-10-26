FROM golang:1.27.2-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0
COPY ./server/go.mod ./server/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY ./server .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
	go build -ldflags="-s -w" -o ./bin/server ./cmd/demo

FROM alpine:latest
WORKDIR /app
ARG EXPOSE_PORT
EXPOSE $EXPOSE_PORT
COPY --from=builder /app/bin/server ./bin/server
CMD ["./bin/server"]
