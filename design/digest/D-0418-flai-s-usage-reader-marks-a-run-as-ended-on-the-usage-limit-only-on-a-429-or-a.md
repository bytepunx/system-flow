---
id: D-0418
title: flai's usage reader marks a run as ended on the usage limit only on a 429 or a rejecting rate_limit_event, so another API error, such as a 529, is not taken for a usage limit
kind: decision
sources: [S-0514, T-2198]
topics: [planning, digest, flai, go]
paths: [flai/internal/usage/claudecode.go, flai/internal/usage/event.go, flai/internal/serve/limit.go]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0418 flai's usage reader marks a run as ended on the usage limit only on a 429 or a rejecting rate_limit_event, so another API error, such as a 529, is not taken for a usage limit

S-0514 taught flai's usage reader to mark a run that ended on the harness's usage limit (ADR-0214). Its narrative records the rule for what counts as one, under T-2198:

> a Limit needs a 429 or a rejecting `rate_limit_event`, so another API error, such as a 529, is not taken for a usage limit. The reset time travels on the End, from the rejecting event before it.

ADR-0214 and `design/system/strategic-agents.md` say what a run on the limit does, but not that other API errors are left out. A 529 (overloaded) is transient and is not a limit to hold starts until a reset. A run that ends on one is recorded as any failed run is. The same narrative, under T-2200, says a limit outranks the exit code. S-0534's host-wide budget state builds on this reader (ADR-0209 §5) and keeps the same rule.

Sources: S-0514 (`wip/archive/agents/S-0514.md`, Decisions), T-2198.
