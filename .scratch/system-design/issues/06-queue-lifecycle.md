# Decide queue lifecycle

Type: grilling
Status: open
Parent: ../map.md

## Question

What is the end-to-end lifecycle across `queue/generation`, `queue/requests`, `queue/messages`, and `queue/reports`? Decide who publishes and consumes each queue, what each stage means, how generation becomes requests, and where reports and messages end up on the UI read path. No payload schemas and no retry semantics.
