---
id: ADR-0253
title: "A container that runs a story's agent mounts the main checkout's wip and system-flow.yaml read-only and its .flai-cache/verify read-write, so flai verify and the close-out run inside it against the main checkout's state"
status: accepted
date: 2026-11-06
supersedes: []
superseded_by: []
refines: [ADR-0203, ADR-0246]
topics: [harness, docker]
---

# ADR-0253 A container that runs a story's agent mounts the main checkout's wip and system-flow.yaml read-only and its .flai-cache/verify read-write, so flai verify and the close-out run inside it against the main checkout's state

## Context

The close-out's first step is `flai verify`, run through `scripts/flai.sh`. Inside the agent container it stopped with exit 127 before `flai verify` began ([I-0218](../issues/I-0218-scripts-close-out-sh-stops-at-its-first-step-inside-the-agent-container-since-scripts-flai-sh-builds-flai-with-go-which-the-image-lacks.md), S-0505): `scripts/flai.sh` built flai with `go`, which the image lacked. S-0518 made `scripts/flai.sh` run the worktree's `bin/flai` or the installed flai when `go` is absent, and S-0510's layer carries the tiers' tools. What remained is the main checkout. [ADR-0203](0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md) mounted the worktree alone, and [ADR-0246](0246-a-container-that-runs-a-story-s-agent-mounts-the-main-checkout-s-git-read-only.md) added the main checkout's `.git` read-only and the story's registration read-write. With no main checkout's manifest beside that `.git`, S-0521 made `flai verify` refuse, naming the host's MCP tools, in place of taking the worktree for the project. So the close-out inside a container still stops at its first step, with S-0521's refusal in place of exit 127.

`flai verify S-nnnn` reads from the main checkout:

| What | Where it is read |
|------|------------------|
| `system-flow.yaml`, by whose presence flai knows the main checkout is there, and the manifest the queue check reads | `flai/internal/workitem/repo.go`, `mainCheckout`; `flai/cmd/verify.go` |
| `wip/`: the narrative's `## Current state` and `## Next steps`, the board, the items, the threads, and the messages that `flai check --strict` scoped to the story reads under `WipDir()` | `flai/internal/verify/story.go`; `flai/internal/check/check.go` |
| the markdownlint configuration at the main checkout's root, which lints `wip/` | `flai/internal/check/check.go`, `markdown` |
| `.git`: the base branch, the story's diff, and the refs | mounted by ADR-0246 |

It writes one thing there: its report, under `<main>/.flai-cache/verify/` (`flai/internal/verify/story.go` and `record.go`). `flai verify --last`, the close-out's reuse of tiers, the close-out job's verdict, and the orchestrator's acceptance read it on the host.

## Decision

A container that runs a story's agent mounts the main checkout's `wip/` and `system-flow.yaml` read-only and its `.flai-cache/verify/` read-write, each at its host path, beside the worktree, the story's cache folder, the login, and ADR-0246's two git mounts.

| Mount | Mode | Why |
|-------|------|-----|
| `<main>/wip` | read-only | the narrative, the board, the items, the threads, and the messages that `flai verify` checks |
| `<main>/system-flow.yaml` | read-only | by it flai knows the main checkout is there, so `flai verify` runs rather than refusing |
| `<main>/.flai-cache/verify` | read-write | the report `flai verify` writes, where the host reads it |

So `flai verify` and the close-out run inside the container against the same narrative, board, and branch as on the host, and record their report where the host reads it.

1. **The verify folder is the operator's.** The container adapter and `scripts/agent-container.sh` create `<main>/.flai-cache/verify` on the host before `docker run` when it is absent, since Docker would create a missing bind-mount source as root and the container's `--user` could then not write the report.
2. **The lint configuration falls back to the worktree's.** `flai check` lints `wip/` with the main checkout's markdownlint configuration, and with the worktree's when the main checkout's root holds none. Inside the container the root holds only what is mounted, so the worktree's tracked configuration lints `wip/`, and no configuration file needs a mount of its own. On the host the two are the same file once the story is accepted.
3. **A test holds the list.** A `flai verify` test runs over a fixture whose main checkout holds only `.git`, `wip/`, and `system-flow.yaml`, so a new read of anything else from the main checkout fails a test, rather than stopping a close-out in a container.

Nothing else of the main checkout is mounted: not `.flai-cache` as a whole, which holds the host configuration, the MCP token, the gateway key files, and other stories' agent folders; not `bin/` or the tools, which are the project's layer's (S-0510); not `design/`, `docs/`, or `scripts/`, which the worktree has.

## Consequences

- The close-out inside a container runs past `flai verify`'s own steps, rebase, sync, narrative, and check, to the tiers. The tiers then need the project's layer (S-0510).
- Commits stay the host's (ADR-0246). The close-out's commit step stops inside a container when `--record-issues` recorded something or anything else is uncommitted; the agent commits through `task_done` and asks for the close-out from the host.
- The container's flai configuration (S-0521's `FLAI_CONFIG` under the story's agent folder) turns on no merge-queue host action, so `flai verify` inside it runs the tiers marked `queue: true` itself.
- The report's owner is the container's `--user`, the operator's uid, as on the host.
- The container's agent can read every work item, thread, and narrative of the project, not only its own story's. It could already reach them through the MCP server; it cannot write them, since the mount is read-only.
- The container adapter's mounts grow from five to eight, and so do `scripts/agent-container.sh`'s and its runbook's mount table.

## Alternatives considered

- **The main checkout read-only whole, with `.flai-cache` masked.** One mount serves every read, but it brings the operator's untracked files into the container, and a mask that misses a secret under `.flai-cache` leaks it. Not taken.
- **Keep S-0521's refusal and run the close-out of a container's story on the host, through the MCP tool `verify`.** No mount is needed, but the tiers run on the host, which takes away what S-0510's layer is for: a close-out inside the boundary. Kept as the path while the mounts are absent, not as the design.
- **Read the worktree's own tracked `wip/` in place of the main checkout's.** It needs no mount, but it is the board as the story's branch last saw it: S-0505 measured 2,725 items there against the main checkout's 2,738. It is the silent wrong root S-0521 refused. Not taken.
