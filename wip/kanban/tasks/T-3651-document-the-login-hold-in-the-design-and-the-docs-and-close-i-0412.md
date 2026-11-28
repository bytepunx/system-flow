---
id: T-3651
type: task
nature: remediation
title: Document the login hold in the design and the docs, and close I-0412
status: backlog
parent: S-0855
owner: alex
created: 2026-11-28T00:00:10Z
updated: 2026-11-28T00:00:10Z
transitions: []
stream: S-0855
tags: [flai]
touches: [design/system/flai-cli.md, design/system/strategic-agents.md, docs/operators/index.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0412-flai-serve-restarts-the-orchestrator-every-minute-while-the-harness-s-login-has.md, design/issues/summary.md]
after: [T-3649, T-3650]
---
# T-3651 Document the login hold in the design and the docs, and close I-0412

## Work

Waits for T-3649 and T-3650, since it describes the hold, the thread, and the banner as they were built.

- `design/system/flai-cli.md`, the host's budget state: the login kind, what sets it, its backoff and probe, its notice and thread, and what clears it. Cite the ADR.
- `design/system/strategic-agents.md`: the orchestrator's start table and its run-end table. A run ended on a failed login waits for the budget state, not `orchestrateRetry`. Do the same in the curator's and the analyzer's rows where they name the usage limit.
- `docs/operators/index.md`, "When a provider is out of budget", and its notices: what the operator sees when the login fails and how to renew it.
- `docs/users/flai.md` and `docs/users/flaiover.md`: the banner's words and the thread.
- Close I-0412 with `flai issue close I-0412 --reason "<what fixed it>"`, naming the ADR and the hold. That writes the issue's file and `design/issues/summary.md`.

## Done when

- Each document says what the code does, and `flai check --strict` and the markdown lint pass on the story.
- I-0412 is closed with its reason, and the story's second criterion can be ticked.

## Notes
