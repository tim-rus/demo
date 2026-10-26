# AI GENERATED
Пока сгенерировал аишкой, чтобы не возиться с описанием конфигов. Потом перепишу.


# demo

A small full-stack starter built around [ConnectRPC](https://connectrpc.com/): a Go server exposes ping RPCs, a React web UI calls one of them over Connect-Web, and a Node.js client (the future Telegram bot) calls the other over Connect-Node.

The value of the repo is the wiring: protobuf-first APIs, generated Go and TypeScript clients, hot-reload dev tasks, and a Docker Compose setup for all three services.

## Architecture

```
┌──────────────────────┐      Connect-Web (HTTP)      ┌──────────────────────────┐
│  ui (React + Vite)   │ ───────────────────────────> │  server (Go, net/http)   │
│  :UI_HOST_PORT       │      /api.PingService/Ping   │  :SERVER_DOCKER_PORT     │
└──────────────────────┘                              │                          │
                                                      │  apigen → api.PingService│
┌──────────────────────┐      Connect (HTTP/1.1)      │  tggen  → tg.PingService │
│  tg (Node client)    │ ───────────────────────────> │                          │
│  no public port      │      /tg.PingService/Ping    └──────────────────────────┘
└──────────────────────┘
```

Both `Ping` RPCs return `response: "pong"` and a `google.protobuf.Timestamp` of the current time.

## Components

| Path | Description |
| --- | --- |
| `proto/api/` | API protobuf: `api.PingService` with a unary `Ping` RPC. Consumed by the UI. |
| `proto/tg/` | Telegram/bot protobuf: `tg.PingService` with the same `Ping` shape. Consumed by `tg`. |
| `server/` | Go server. Serves both Connect services on one `net/http` mux, with CORS and graceful shutdown. |
| `tg/` | Node.js + TypeScript bot scaffold. Pings the `tg` service after 5 seconds and logs the result. |
| `ui/` | React 19 + Vite + Tailwind CSS 4 + shadcn/ui (Base UI) SPA. The **Go!** button calls the `api` ping. |
| `docker/` | Multi-stage Dockerfiles for the server, the Node client, and the UI (lighttpd static host). |
| `taskfile.yaml` | Dev, codegen, Compose, and log shortcuts. |
| `buf.gen.api.yaml`, `buf.gen.tg.yaml` | Buf codegen templates for the two proto sets. |

## Project layout

```
.
├── compose.yaml
├── taskfile.yaml
├── buf.gen.api.yaml
├── buf.gen.tg.yaml
├── docker/
│   ├── server.Dockerfile
│   ├── tg.Dockerfile
│   └── ui.Dockerfile
├── proto/
│   ├── api/api_ping.proto
│   └── tg/tg_ping.proto
├── server/
│   ├── cmd/demo/
│   │   ├── main.go          # wiring, HTTP server, graceful shutdown
│   │   ├── apisvc.go        # api.PingService implementation
│   │   ├── tgsvc.go         # tg.PingService implementation
│   │   └── cors.go          # allow-all CORS middleware
│   └── interal/proto/       # generated Go + Connect handlers
├── tg/
│   ├── src/
│   │   ├── index.ts
│   │   ├── rpc/ping.ts
│   │   ├── lib/
│   │   └── proto/           # generated TypeScript
│   ├── package.json
│   └── tsconfig.json
└── ui/
    ├── src/
    │   ├── pages/Home.tsx
    │   ├── rpc/ping.ts
    │   ├── routes.ts
    │   ├── proto/           # generated TypeScript
    │   └── ...
    ├── lighttpd.conf
    └── vite.config.ts
```

## Requirements

- **Go** 1.27.2+ (matches `server/go.mod` and `docker/server.Dockerfile`)
- **Node.js** 24+ with npm
- **Docker** with Compose v2 (optional, for containerized runs)
- **Task** ([go-task](https://taskfile.dev/)) for the shortcuts in `taskfile.yaml`
- **Buf** plus the protoc plugins, only when regenerating protobuf code

## Configuration

The app reads a root `.env`. It is gitignored, so create your own:

```dotenv
# development
SERVER_MODE=http
SERVER_DEV_HOST=localhost
SERVER_DEV_PORT=9000
UI_DEV_PORT=3000

# Docker Compose
SERVER_HOST=localhost
SERVER_HOST_PORT=9000
SERVER_DOCKER_PORT=8080
UI_HOST_PORT=3000
UI_DOCKER_PORT=80

# Telegram bot (reserved)
TG_TOKEN=
```

| Variable | Used by | Purpose |
| --- | --- | --- |
| `SERVER_MODE` | dev clients, Compose UI build | URL scheme: `http` or `https`. |
| `SERVER_DEV_HOST` | `task run:tg`, `task run:ui` | Host the dev clients connect to. |
| `SERVER_DEV_PORT` | `task run:server`, dev clients | Go server port in development. |
| `SERVER_HOST` | Compose UI build | Host baked into the browser-facing API URL. |
| `SERVER_HOST_PORT` | Compose UI build, port mapping | Published server port; should be reachable from the browser. |
| `SERVER_DOCKER_PORT` | Compose server | Server port inside the container and on the Compose network. |
| `UI_DEV_PORT` | Vite dev server | UI dev port (defaults to `3000`). |
| `UI_HOST_PORT` | Compose UI | Published UI port. |
| `UI_DOCKER_PORT` | Compose UI | UI container port. lighttpd listens on **80**, so use `80`. |
| `TG_TOKEN` | `tg` container | Reserved for the Telegram bot token; not used by the scaffold yet. |

`Task` loads the root `.env` through `dotenv: ['.env']`, and Docker Compose picks it up automatically. The Go server itself only reads `HOST` and `PORT`; `task run:server` sets both from `SERVER_DEV_PORT`.

## Quick start (local)

```bash
# 1. install dependencies
(cd server && go mod download)
npm install --prefix tg
npm install --prefix ui

# 2. create the root .env (see Configuration)

# 3. run server + tg + ui together
task run
```

`task run` starts all three as parallel task dependencies. By default:

- UI dev server: <http://localhost:3000>
- API: `http://localhost:SERVER_DEV_PORT`
- `tg` logs its ping result about 5 seconds after starting.

You can also start each process separately:

```bash
task run:server
task run:tg
task run:ui
```

Open the UI and click **Go!** — it calls `/api.PingService/Ping` and shows `pong: <timestamp>`.

## Quick start (Docker Compose)

```bash
task compose        # docker compose up --build
```

Then open `http://localhost:${UI_HOST_PORT}` (default `3000`).

- `server` is published on `${SERVER_HOST_PORT}`.
- `ui` is served by lighttpd inside the container on port 80; `${UI_HOST_PORT}:${UI_DOCKER_PORT}` publishes it.
- `tg` has no published port and reaches the server at `http://server:${SERVER_DOCKER_PORT}` on the `demo_net` network.

The UI build bakes `VITE_SERVER_BASE_URL=${SERVER_MODE}://${SERVER_HOST}:${SERVER_HOST_PORT}` into the static bundle, so changing any of those three values requires a rebuild.

Useful commands:

```bash
task compose:server     # rebuild and run only the server
task compose:tg         # rebuild and run the tg client
task compose:ui         # rebuild and run only the UI
task compose:stop       # stop the stack
task logs               # follow all logs (100 lines of backlog)
task logs:server
task logs:tg
task logs:ui
```

## Generating protobuf code

Generated code is committed, so you only need this when a `.proto` file changes.

```bash
task buf       # regenerate both api and tg
task buf:api   # only proto/api
task buf:tg    # only proto/tg
```

| Proto input | Generated Go | Generated TypeScript |
| --- | --- | --- |
| `proto/api/api_ping.proto` | `server/interal/proto/apigen/` | `ui/src/proto/` |
| `proto/tg/tg_ping.proto` | `server/interal/proto/tggen/` | `tg/src/proto/` |

The templates expect these plugins on `PATH`:

```bash
# Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

# TypeScript
npm install -g @bufbuild/protoc-gen-es @connectrpc/protoc-gen-connect-es

# Buf CLI, if it is not installed yet
go install github.com/bufbuild/buf/cmd/buf@latest
```

Each template also generates Connect-ES service files (`*_connect.ts`). The clients currently import the service descriptors from the protobuf-es v2 `*_pb.ts` files instead.

## API

Both services use the Connect protocol, so they also work with gRPC and gRPC-Web clients. The server accepts HTTP/1.1 and cleartext HTTP/2 (h2c).

| Procedure | Endpoint |
| --- | --- |
| `api.PingService.Ping` | `POST /api.PingService/Ping` |
| `tg.PingService.Ping` | `POST /tg.PingService/Ping` |

Schema (identical in both packages):

```proto
service PingService {
  rpc Ping(PingRequest) returns (PingResponse) {}
}

message PingRequest {
  string msg = 1;
}

message PingResponse {
  string response = 1;
  google.protobuf.Timestamp timestamp = 2;
}
```

Example with `curl` (Connect protocol, JSON):

```bash
curl -sS \
  -H 'Content-Type: application/json' \
  -d '{"msg":"hello"}' \
  http://localhost:9000/api.PingService/Ping
```

## Notes

- The Go module path is `demo`, and the generated package directory is spelled `interal` (not `internal`). Imports, buf templates, and generated code all use that spelling; keep it consistent.
- CORS is allow-all and h2c is enabled for convenience. Tighten both before exposing the server publicly.
- `tg` is only a Connect client today. It does not connect to Telegram yet, and `TG_TOKEN` is unused.
- The UI container is a static lighttpd site with SPA fallback to `index.html`. Set `UI_DOCKER_PORT=80` unless you also change `lighttpd.conf`.
