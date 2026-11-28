# Decide data ownership

Type: grilling
Status: open
Blocked by: 01
Parent: ../map.md

## Question

Which service owns each SurrealDB namespace (`api`, `generator`, `orchestrator`, `rag`, `scheduler`), and what runtime state does each keep? The brief shows runner with no namespace and API as the only service querying SurrealAPI. Decide whether those stay true and whether API persists data it consumes from NATS. Do not design schemas.
