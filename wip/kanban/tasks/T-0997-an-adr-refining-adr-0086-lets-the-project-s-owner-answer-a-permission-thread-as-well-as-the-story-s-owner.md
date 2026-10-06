---
id: T-0997
type: task
nature: remediation
title: An ADR refining ADR-0086 lets the project's owner answer a permission thread as well as the story's owner
status: backlog
parent: S-0284
owner: alex
created: 2026-10-06T06:24:05Z
updated: 2026-10-06T06:24:05Z
transitions: []
stream: S-0284
tags: [flai]
touches: [design/adrs]
---
# T-0997 An ADR refining ADR-0086 lets the project's owner answer a permission thread as well as the story's owner

## Work

I-0081's one instance: S-0218's owner is `arobson`, the operator replied `allow` on TH-0158 as `alex`, the `owner` in `system-flow.yaml` and the name the dashboard writes as (`owner` in `flai/internal/hostapi/writes.go`). `awaitAnswer` in `flai/internal/mcpserver/permission.go` takes an answer from the story's owner only, so the write waited until Claude Code's 1800 s MCP idle timeout ended it.

Propose the fix on S-0284's plan thread before building it, as the story's goal asks. The recommended fix is that an answer may come from the story's owner or the project's owner (`repo.Manifest.Owner`), never from the asking agent, with entries by anyone else, other agents included, still not answers. When both names are empty, anyone but the agent answers, as now.

Once the operator agrees, write the ADR with `flai adr new`. It refines ADR-0086's decision 2: who may answer and why. Set its `refines` to ADR-0086 and its topics to `cli`. Leave ADR-0086 as it is: it is accepted. Add the row to `design/adrs/README.md` if `flai adr new` does not.

## Done when

- the operator has agreed to the fix on S-0284's plan thread, or chosen another, and the ADR records what they chose
- the ADR is accepted, refines ADR-0086, and names I-0081 as its context
- `flai check --strict` is clean on `design/adrs`

## Notes
