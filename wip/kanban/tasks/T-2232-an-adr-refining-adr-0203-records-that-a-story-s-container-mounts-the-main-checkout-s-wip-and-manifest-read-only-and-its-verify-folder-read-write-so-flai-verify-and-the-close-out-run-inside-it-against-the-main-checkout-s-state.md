---
id: T-2232
type: task
nature: remediation
title: An ADR refining ADR-0203 records that a story's container mounts the main checkout's wip and manifest read-only and its verify folder read-write, so flai verify and the close-out run inside it against the main checkout's state
status: backlog
parent: S-0520
owner: alex
created: 2026-10-28T13:57:07Z
updated: 2026-10-28T13:57:07Z
transitions: []
stream: S-0520
tags: [cli, docker]
touches: [design/adrs/README.md]
---
# T-2232 An ADR refining ADR-0203 records that a story's container mounts the main checkout's wip and manifest read-only and its verify folder read-write, so flai verify and the close-out run inside it against the main checkout's state

## Work

Write the ADR with `flai adr new`, refining [ADR-0203](../../../design/adrs/0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md) and the ADR S-0517 writes for the two git mounts, so that `design/adrs/README.md` gains its row. Context: I-0218's instance, and what `flai verify` reads from the main checkout, from the code: the narrative's `## Current state` and `## Next steps` (`flai/internal/verify/story.go`, `narrative`), the board, the items, the threads, and the messages that `flai check --strict` scoped to the story reads under `WipDir()` (`flai/internal/check/check.go`), `system-flow.yaml`, by which `workitem.Open` knows the main checkout is there (`flai/internal/workitem/repo.go`, `mainCheckout`) and the queue check reads the manifest (`flai/cmd/verify.go`), the markdownlint configuration (`check.go`, `markdown`), and the `.git` S-0517 mounts; and what it writes: its report under `<main>/.flai-cache/verify/<S-nnnn>.json` (`story.go`, `SaveReport`), which `flai verify --last`, the close-out's reuse of tiers, and the orchestrator's acceptance read on the host. Decision: a story's container mounts, beside the worktree, the story's cache folder, the login, and S-0517's two git mounts, the main checkout's `wip/` and `system-flow.yaml` read-only and its `.flai-cache/verify/` read-write, each at its host path, so that `flai verify` inside the container checks the narrative, the board, and the branch the host would and records its report where the host reads it; the markdown lint configuration is read from the worktree when the main checkout's cannot be, so no configuration file needs a mount of its own; and a test keeps `flai verify` to those reads, so a read of anything else from the main checkout is a test failure, not a silent stop in a container. Say what stays out: `.flai-cache` as a whole, which holds the config, the MCP token, the LiteLLM key file, and other stories' agent folders; `bin/` and the tools, which are the project's layer's (S-0510); and the commit step, which S-0517's ADR keeps the host's, so the close-out inside a container stops at `commit` when `--record-issues` recorded something, and the agent commits through `task_done`. Alternatives: the main checkout read-only whole with `.flai-cache` masked (one mount, every read served, but the operator's untracked files in the container too); the refusal kept and the close-out in a container run as the host's MCP `verify` (no mount, but the tiers run on the host, which defeats S-0510's layer); the worktree's own `wip/` read in place of main's (what S-0505 measured as 2,725 items against 2,738, the silent wrong root S-0521 refused). Consequences: three more mounts in the adapter and the script; the report's owner is the container's `--user`, the operator's uid; the container's own config (S-0521's `FLAI_CONFIG` under the agent folder) carries no merge-queue host action, so `flai verify` inside it runs the `queue: true` tiers itself. Accept it with `flai adr accept`.

## Done when

The ADR exists under `design/adrs/` with status accepted, its title names the three mounts and that `flai verify` and the close-out run inside the container, its context cites I-0218 and the reads above with the files they are in, its decision lists the mounts with their modes, the lint fallback, and the guarding test, its alternatives name the three above with why each is not taken, and `design/adrs/README.md` lists it; `flai check --strict` and the markdown lint pass on both.

## Notes

The ADR goes first because the mounts' list is what the other tasks build and test. S-0517's ADR number is not known at planning; the story's agent reads it from `design/adrs/README.md` once S-0517 is on main.
