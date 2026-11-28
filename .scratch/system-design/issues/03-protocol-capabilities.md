# Research protocol capabilities

Type: research
Status: resolved
Parent: ../map.md

## Question

Research only; report facts, make no design decision. What protocol paths are technically available for browser-to-API, API-to-RAG, API-to-scheduler, runner-to-RAG, and runner-to-scheduler given Node 26.13.1 (`@connectrpc/connect-node`, `@modelcontextprotocol/node`), Go 1.27.4 (ConnectRPC), Python 3.14.10 (FastAPI, FastMCP), and lighttpd-served React/Vue clients? Surface the contradiction between brief prose (RAG mcp/rest, scheduler mcp/connectrpc) and both diagrams (API connects by ConnectRPC to RAG and scheduler). Note any constraint on serving MCP from Node and Python runtimes.

## Answer

Research only. No design decision recorded here. Source facts current as of 2026-11-28.

### Runtime facts

- Node 26.13.1 satisfies listed library floors: `@connectrpc/connect-node` 2.3.0 declares `node >=22`; `@modelcontextprotocol/node` 2.1.2 declares `node >=20`. [S1][S2]
- `@modelcontextprotocol/node` is a thin Node integration layer for `@modelcontextprotocol/server`. It exports `NodeStreamableHTTPServerTransport`, `toNodeHandler`, and Host/Origin guards. It is a server adapter, not an MCP client package. Its peer dependency is `@modelcontextprotocol/server` 2.4.0. [S2][S3]
- `@modelcontextprotocol/client` 2.4.0 is the package with `StreamableHTTPClientTransport`, `StdioClientTransport`, and `SSEClientTransport`. It is not in the listed Node package set. [S4]
- connect-go servers accept Connect, gRPC, and gRPC-Web by default. connect-go clients use Connect by default and switch with `WithGRPC()` / `WithGRPCWeb()`. [S5]
- Python 3.14.10 is a stable Python release. FastMCP docs require Python 3.11+; FastMCP serves MCP over stdio or HTTP. FastAPI is an HTTP API framework with JSON request/response handling and OpenAPI docs. [S6][S7][S8]
- Connect for Python exists in beta, requires Python 3.11+, and generates ASGI/WSGI servers plus clients. It is not in the listed Python package set. [S9]
- lighttpd `mod_proxy` can proxy HTTP to backends, load-balance across backends, remap host/path headers, append `Forwarded`, and enable WebSocket `Upgrade`. [S10]

### Protocol facts

- Connect protocol supports unary, client streaming, server streaming, and bidirectional streaming with Protobuf or JSON. Bidirectional streaming needs HTTP/2; other RPC types also work on HTTP/1.1. Connect uses no HTTP trailers. Unary uses `application/proto` or `application/json`; streaming uses `application/connect+proto` or `application/connect+json`. [S11]
- Connect Web can call Connect via `createConnectTransport()` using `fetch`, or gRPC-Web via `createGrpcWebTransport()`. It targets modern browsers with `fetch` and works with React/Vue. Browser streaming is server-streaming only: browser vendors do not agree on streaming request bodies. [S12][S13]
- Cross-origin Connect needs CORS preflight allowing `GET,POST` and headers `Content-Type, Connect-Protocol-Version, Connect-Timeout-Ms, X-User-Agent`. gRPC-Web needs `POST` and `Content-Type,Grpc-Timeout,X-Grpc-Web,X-User-Agent`, plus exposed `Grpc-Status, Grpc-Message, Grpc-Status-Details-Bin`. Custom headers such as `Authorization` must be allowed explicitly. [S14]
- MCP messages are JSON-RPC. Standard transports are stdio and Streamable HTTP; custom transports MAY be implemented. stdio means the client launches the server as a subprocess and uses newline-delimited JSON on stdin/stdout, so it is a local process edge, not a network service edge. [S15][S16]
- Streamable HTTP exposes one POST endpoint. Each JSON-RPC message is its own POST; responses are JSON or SSE scoped to that request. Revision 2026-07-28 removed the GET stream endpoint and protocol-level sessions. Servers MUST validate `Origin`; local servers SHOULD bind `127.0.0.1`; servers SHOULD implement auth. [S17]
- TypeScript SDK v2 serves HTTP with `createMcpHandler(factory)` and mounts via `toNodeHandler(handler)`. The handler validates no Host, Origin, or token; callers must mount those checks. Plain `node:http` examples bind `127.0.0.1` and compose `localhostHostValidation` / `localhostOriginValidation`. [S18]
- FastMCP HTTP deployment offers direct `mcp.run(transport="http")` or an ASGI app via `mcp.http_app()` under Uvicorn. Mounting under FastAPI/Starlette requires passing the MCP app lifespan. Multi-worker needs `stateless_http=True`. Host/Origin guard is opt-in; default `host_origin_protection=False`. Browser-direct MCP clients need CORS headers and exposed `mcp-session-id`; CORS is separate from Host/Origin validation. [S7]
- FastMCP Streamable HTTP can emit SSE. FastMCP docs warn reverse proxies must disable buffering and raise timeouts for SSE; their example uses nginx, not lighttpd. lighttpd mod_proxy docs list generic proxy and WebSocket upgrade support, not MCP/SSE-specific behavior. [S7][S10]

### Edge facts

- **Browser-to-API:** Available with ConnectRPC Web Connect transport (`fetch`) or gRPC-Web transport. Cross-origin calls need the CORS rules above. Client streaming and bidirectional streaming from browser are not available across browsers; server streaming is. lighttpd-served clients are a stated premise; `mod_proxy` can proxy API paths. If the API exposes MCP instead, browser JavaScript can reach Streamable HTTP only with explicit CORS plus Host/Origin handling; REST over plain HTTP is also an ordinary browser `fetch` path. [S10][S12][S13][S14][S7][S8]
- **API-to-RAG:** With Node API, ConnectRPC client paths are Connect, gRPC, or gRPC-Web via `@connectrpc/connect-node`. With FastAPI RAG, a plain HTTP/JSON path is available. With FastMCP RAG, network MCP is Streamable HTTP; Node needs `@modelcontextprotocol/client`, not `@modelcontextprotocol/node`, to call it. stdio is available only when the caller launches RAG as a child process on the same host. Python FastMCP can also call MCP servers via its own client over in-memory, stdio, or HTTP. [S1][S2][S4][S5][S7][S8][S15][S16][S17][S19]
- **API-to-scheduler:** A Go connect-go scheduler accepts Connect, gRPC, and gRPC-Web without extra configuration. A Node API client can speak all three via `@connectrpc/connect-node`; protocol switching is by `Content-Type`. ConnectRPC Go docs do not supply MCP. MCP would require a separate MCP server/client implementation or a custom transport. [S1][S5][S11][S15][S20]
- **Runner-to-RAG:** If runner is Node, same client paths as API-to-RAG. If runner is Python, FastMCP client supports in-memory, stdio, and HTTP MCP; FastAPI supports ordinary HTTP APIs. Connect for Python is official but beta and not listed in the given Python stack. [S7][S8][S9][S15][S16][S17][S19]
- **Runner-to-scheduler:** If runner is Node, same Connect/gRPC/gRPC-Web client paths as API-to-scheduler. If runner is Python, the listed FastAPI/FastMCP libraries document HTTP APIs and MCP transports; they do not document ConnectRPC. Connect for Python is an available beta library if added. Go scheduler side already accepts all three ConnectRPC protocols by default. [S5][S7][S8][S9][S11]

### Brief-vs-diagram contradiction

- Brief prose names RAG as `mcp/rest` and scheduler as `mcp/connectrpc`.
- Both diagrams label API-to-RAG and API-to-scheduler as ConnectRPC.
- Overlap exists only for scheduler ConnectRPC. RAG prose lists MCP/REST while diagrams list only ConnectRPC. Scheduler prose lists MCP and ConnectRPC while diagrams list only ConnectRPC. No selection is made in this ticket.

### Sources

- [S1] npm registry, `@connectrpc/connect-node` 2.3.0 metadata: https://registry.npmjs.org/@connectrpc/connect-node
- [S2] npm registry, `@modelcontextprotocol/node` 2.1.2 metadata and README: https://registry.npmjs.org/@modelcontextprotocol/node
- [S3] npm registry, `@modelcontextprotocol/server` 2.4.0 metadata and README: https://registry.npmjs.org/@modelcontextprotocol/server
- [S4] npm registry, `@modelcontextprotocol/client` 2.4.0 metadata and README; TypeScript SDK client guide: https://registry.npmjs.org/@modelcontextprotocol/client and https://ts.sdk.modelcontextprotocol.io/v2/clients/connect
- [S5] ConnectRPC Go getting started: https://connectrpc.com/docs/go/getting-started/
- [S6] Python 3.14.10 documentation: https://docs.python.org/3.14/
- [S7] FastMCP HTTP deployment: https://gofastmcp.com/docs/deployment/http
- [S8] FastAPI documentation: https://fastapi.tiangolo.com/
- [S9] ConnectRPC Python getting started: https://connectrpc.com/docs/python/getting-started/
- [S10] lighttpd `mod_proxy`: https://redmine.lighttpd.net/projects/lighttpd/wiki/Docs_ModProxy
- [S11] Connect protocol: https://connectrpc.com/docs/protocol/
- [S12] ConnectRPC Web getting started: https://connectrpc.com/docs/web/getting-started/
- [S13] ConnectRPC Web choosing a protocol; supported browsers and frameworks: https://connectrpc.com/docs/web/choosing-a-protocol/ and https://connectrpc.com/docs/web/supported-browsers-and-frameworks/
- [S14] ConnectRPC CORS: https://connectrpc.com/docs/cors/
- [S15] MCP transports concepts: https://modelcontextprotocol.io/docs/concepts/transports
- [S16] MCP spec 2026-07-28 stdio: https://modelcontextprotocol.io/specification/2026-07-28/basic/transports/stdio.md
- [S17] MCP spec 2026-07-28 Streamable HTTP: https://modelcontextprotocol.io/specification/2026-07-28/basic/transports/streamable-http.md
- [S18] TypeScript SDK v2 HTTP serving: https://ts.sdk.modelcontextprotocol.io/v2/serving/http
- [S19] FastMCP client: https://gofastmcp.com/docs/clients/client
- [S20] ConnectRPC multi-protocol support: https://connectrpc.com/docs/multi-protocol/
