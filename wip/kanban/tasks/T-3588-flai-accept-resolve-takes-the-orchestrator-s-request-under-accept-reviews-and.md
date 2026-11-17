---
id: T-3588
type: task
nature: improvement
title: flai accept --resolve takes the orchestrator's request under accept_reviews, and an MCP tool resolve asks for it on the host
status: backlog
parent: S-0845
owner: alex
created: 2026-11-17T12:34:05Z
updated: 2026-11-17T12:34:05Z
transitions: []
stream: S-0845
tags: [flai, orchestrator]
touches: [flai/cmd/accept_resolve.go, flai/cmd/accept_resolve_test.go, flai/internal/mcpserver/resolve.go, flai/internal/mcpserver/resolve_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go]
after: [T-3587]
---
# T-3588 flai accept --resolve takes the orchestrator's request under accept_reviews, and an MCP tool resolve asks for it on the host

## Work

Build what T-3587's ADR decides; it waits for that ADR because the ADR settles who may ask and when. As the plan assumes: in `flai/cmd/accept_resolve.go`, `resolveFlagsRefused` stops refusing `--by orchestrator` and `FLAI_ROLE=orchestrate` while `orchestration.permissions.accept_reviews` is on, and refuses them while it is off, naming the permission; every other refusal ADR-0313 keeps stays. Add an MCP tool `resolve` in `flai/internal/mcpserver/resolve.go`, beside `accept.go`, that runs `flai accept <S-nnnn> --resolve --by orchestrator --json` on the host and answers the job ADR-0313 §1 gives, `{id, job, started, running, log, reason}`, or the refusal. Register it in `folder.go`. Test both in `accept_resolve_test.go` and `resolve_test.go`, and the tool list in `folder_test.go`.

## Done when

- The orchestrator's `--resolve` request starts a resolve job while `accept_reviews` is on and is refused while it is off, each with a test.
- The MCP tool `resolve` answers the job or the refusal, with a test.
- `flai test` passes on the paths changed.

## Notes
