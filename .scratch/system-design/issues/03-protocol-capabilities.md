# Research protocol capabilities

Type: research
Status: open
Parent: ../map.md

## Question

Research only; report facts, make no design decision. What protocol paths are technically available for browser-to-API, API-to-RAG, API-to-scheduler, runner-to-RAG, and runner-to-scheduler given Node 26.13.1 (`@connectrpc/connect-node`, `@modelcontextprotocol/node`), Go 1.27.4 (ConnectRPC), Python 3.14.10 (FastAPI, FastMCP), and lighttpd-served React/Vue clients? Surface the contradiction between brief prose (RAG mcp/rest, scheduler mcp/connectrpc) and both diagrams (API connects by ConnectRPC to RAG and scheduler). Note any constraint on serving MCP from Node and Python runtimes.
