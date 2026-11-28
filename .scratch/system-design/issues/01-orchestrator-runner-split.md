# Decide orchestrator and runner split

Type: grilling
Status: open
Parent: ../map.md

## Question

What is the boundary between orchestrator control plane and runner execution? The brief makes runner a ConnectRPC client of orchestrator, sends runner to RAG and scheduler over HTTP, and shows no runner datastore. Decide who creates and tracks a unit of work, who calls RAG and scheduler, what runner reports back to orchestrator, and whether runner keeps any state. Do not define payload schemas or retry semantics.
