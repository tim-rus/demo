# System design map

## Destination

`docs/system_design.md`: a wire-ready runtime architecture for the service-center ticket-routing simulator. Service responsibilities, protocol edges, NATS queues, SurrealDB ownership, observability wiring, and dev/prod deployment. No payload schemas, retry/idempotency semantics, agent internals, or pricing rules.

## Notes

- Domain: synthetic customer requests routed through a distributed service-center simulator. Docs are in English.
- Source brief: attached `.prompts/0-system_design.md`; treat as revisable. Never scan or edit `.prompts`.
- Tracker: local Markdown under `.scratch/system-design/`. Every HITL session calls `/grilling` and `/domain-modeling`.
- Agreed bar: wire-ready runtime, not contract-level; runtime-only scope; token/cost tracked only where it changes wiring.
- Session rules from `AGENTS.md`: new branch per session, atomic commits, `/caveman` voice, `/caveman-commit` messages.

## Decisions so far

<!-- one line per resolved ticket: gist + link. Gist only; detail stays in the ticket. -->
- ConnectRPC/MCP paths per edge, serving constraints, and brief-vs-diagram mismatch: [protocol capabilities facts](./issues/03-protocol-capabilities.md).

## Not yet specified

- Failure-injection path: the UI can initiate failures; owner and transport stay unclear until queue lifecycle and orchestrator/runner split settle.
- UI trace and load read path: which queue or store feeds the UI depends on data ownership and queue lifecycle.
- Admin PWA scope: only "admin interface" is given; capabilities that change runtime wiring may need a decision later.
- Runner state persistence: the brief shows no runner datastore; revisit if the control-plane split exposes a need.

## Out of scope

- Message, API, and database schemas; retry and idempotency semantics, deferred by the runtime-only scoping round.
- Agent graph, prompts, model selection, RAG retrieval strategy, and classifier/worker taxonomy.
- Token pricing rules and prompt-accounting internals.
- Code-level implementation spec and runbook-style operational procedures.
