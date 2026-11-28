use /wayfinder

**goal**: ./docs/system_design.md describing overall design & architecture of a system

**main components**:
- ui - user interface: react(ts) / router / zustand / tailwindcss / shadcn / lighttpd
- admin - admin's interface (pwa): vue(ts) / router / pinia / tailwindcss / shadcn / lighttpd
- api - api service: node(26.13.1) / @connectrpc/connect-node / pino
- generator - requests generator service: go(1.27.4) / slog
- orchestrator - agents orchestration service: go(1.27.4) / connectrpc (server) / slog
- runner - agents (classifiers and workers) runner: go(1.27.4) / smallnest/langgraphgo / connectrpc (client) / slog
- rag - knowledge base service (mcp/rest): python(3.14.10) / uv / FastAPI / FastMCP / structlog
- scheduler - field service scheduling *service* (mcp/connectrpc): nodejs(26.13.1) / @modelcontextprotocol/node / connectrpc / pino

**datastore**: surrealdb (v4)

**transport**:
- ConnectRPC: service-to-service comms
- REST: mcp, infra

**infrastructure:**
- message broker: NATS/JetStream
- logs & metrics collector: Datadog Vector
- logs storage: Victoria Logs
- metrics storage: Victoria Metrics 

**deploy**:
- dev mode: infra, surrealdb runs in docker containers, services run locally
- prod: everything runs in docker containers 

**runner**: taskfile w/ logs collector to ./logs/{service_name}.log

**logging**:
- dev mode: collect from infra containers and ./logs/* to vector
- prod mode: collect from containers to vector
*fields*: service, timestamp, level, mode (prod/dev), message, .*

**interconnections:**
```mermaid
flowchart LR
  subgraph Clients["Clients"]
    UI["UI"]
    Admin["Admin"]
  end

  subgraph Services["Services"]
    API["API"]
    Generator["generator"]
    Orchestrator["orchestrator"]
    Scheduler["scheduler"]
    RAG["rag"]
    Runner["runner"]
  end

  subgraph Broker["Message Broker"]
    QGeneration["queue/generation"]
    QMessages["queue/messages"]
    QReports["queue/reports"]
    QRequests["queue/requests"]
  end

  subgraph Stores["SurrealDB"]
    SurrealAPI["surrealdb/api"]
    SurrealGenerator["surrealdb/generator"]
    SurrealOrchestrator["surrealdb/orchestrator"]
    SurrealRAG["surrealdb/rag"]
    SurrealScheduler["surrealdb/scheduler"]
  end

  subgraph Obs["Observability"]
    Vector["vector"]
    VictoriaLogs["victoria logs"]
    VictoriaMetrics["victoria metrics"]
  end

  UI <-->|connectrpc| API
  Admin <-->|connectrpc| API

  API -->|pub| QGeneration
  API -.->|sub| QMessages
  API -.->|sub| QReports
  Generator -.->|sub| QGeneration
  Generator -->|pub| QRequests

  API <-->|connectrpc| Orchestrator
  API <-->|connectrpc| Scheduler
  API <-->|connectrpc| RAG
  Runner <-->|connectrpc| Orchestrator
  Runner -->|http| RAG
  Runner -->|http| Scheduler

  API -->|query/ws| SurrealAPI
  Generator -->|query/ws| SurrealGenerator
  Orchestrator -->|query/ws| SurrealOrchestrator
  RAG -->|query/ws| SurrealRAG
  Scheduler -->|query/ws| SurrealScheduler

  API -->|logs/http| Vector
  API -->|metrics/http| Vector
  Generator -->|logs/http| Vector
  Generator -->|metrics/http| Vector
  Orchestrator -->|logs/http| Vector
  Orchestrator -->|metrics/http| Vector
  Runner -->|logs/http| Vector
  Runner -->|metrics/http| Vector
  RAG -->|logs/http| Vector
  RAG -->|metrics/http| Vector
  Scheduler -->|logs/http| Vector
  Scheduler -->|metrics/http| Vector

  Vector -->|logs/http| VictoriaLogs
  Vector -->|metrics/http| VictoriaMetrics
  ```

**docker containers map**:
```mermaid
flowchart LR
  subgraph Edge["Edge / UI containers"]
    UI["ui<br/>React + lighttpd"]
    Admin["admin<br/>Vue PWA + lighttpd"]
  end

  subgraph Services["Service containers"]
    API["api<br/>Node + ConnectRPC"]
    Generator["generator<br/>Go"]
    Orchestrator["orchestrator<br/>Go + ConnectRPC server"]
    Runner["runner<br/>Go + ConnectRPC client"]
    RAG["rag<br/>Python FastAPI/FastMCP"]
    Scheduler["scheduler<br/>Node MCP/ConnectRPC"]
  end

  subgraph Data["Data containers"]
    NATS["nats<br/>NATS/JetStream"]
    Surreal["surrealdb<br/>SurrealDB v4"]
  end

  subgraph Obs["Observability containers"]
    Vector["vector<br/>Datadog Vector"]
    VictoriaLogs["victoria_logs<br/>Victoria Logs"]
    VictoriaMetrics["victoria_metrics<br/>Victoria Metrics"]
  end

  UI <-->|connectrpc| API
  Admin <-->|connectrpc| API

  API -->|pub queue/generation| NATS
  API -->|sub queue/messages| NATS
  API -->|sub queue/reports| NATS
  Generator -->|sub queue/generation| NATS
  Generator -->|pub queue/requests| NATS

  API <-->|connectrpc| Orchestrator
  API <-->|connectrpc| Scheduler
  API <-->|connectrpc| RAG
  Runner <-->|connectrpc| Orchestrator
  Runner -->|http| RAG
  Runner -->|http| Scheduler

  API -->|ws query| Surreal
  Generator -->|ws query| Surreal
  Orchestrator -->|ws query| Surreal
  RAG -->|ws query| Surreal
  Scheduler -->|ws query| Surreal

  API -->|logs/metrics http| Vector
  Generator -->|logs/metrics http| Vector
  Orchestrator -->|logs/metrics http| Vector
  Runner -->|logs/metrics http| Vector
  RAG -->|logs/metrics http| Vector
  Scheduler -->|logs/metrics http| Vector

  Vector -->|logs http| VictoriaLogs
  Vector -->|metrics http| VictoriaMetrics
```

**project structure**:
```
  ./
    /configs
    /data
    /logs
    /docker
    /proto
    /migrations
    /admin
    /ui
    /api
    /generator
    /orchestrator
    /runner
    /rag
    /scheduler
    /.env.example
    /.env.dev.example
    /compose.yaml
    /compose.dev.yaml
    /taskfile.yaml
    /buf.*.gen.yaml
    /.gitignore

    /docs
    /AGENTS.md
    /README.md
    
    /.prompts
    /.scratch
```
