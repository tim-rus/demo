# Decide deployment topology

Type: grilling
Status: open
Parent: ../map.md

## Question

How do `compose.yaml`, `compose.dev.yaml`, and `taskfile.yaml` realize the two deployment modes? Decide which containers run in dev versus prod, where Vector runs in each mode and how it reads infra-container logs plus `./logs/*` in dev, and how taskfile starts services with one log file per service. Brief constraints: dev runs infra and SurrealDB in containers with services local; prod runs everything in containers.
