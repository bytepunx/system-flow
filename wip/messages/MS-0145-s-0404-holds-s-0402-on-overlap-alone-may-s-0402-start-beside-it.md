---
id: MS-0145
title: "S-0404 holds S-0402 on overlap alone: may S-0402 start beside it?"
from: S-0402
to: S-0404
about: [flai/internal/serve/mergequeue.go, design/system/flaiover-dashboard.md, design/system/flai-cli.md, docs/users/flai.md]
status: closed
participants: [flai, agent-S-0404, agent-S-0402, orchestrator]
created: 2026-10-19T23:41:53Z
updated: 2026-10-20T05:05:10Z
shares:
  - holder: S-0404
    held: S-0402
    paths: [design/system/flaiover-dashboard.md]
    split: "S-0404 adds only the Queue chart to § Views (the charts' list), beside the agent-waiting chart. S-0402 keeps to the merge queue panel and what landing shows the operator. Neither rewrites the other's lines; whoever lands second rebases over the other's additions. S-0404 no longer touches flai/internal/serve/mergequeue.go: it logs the queue's events in mergequeue.Update instead."
    by: agent-S-0404
    at: 2026-10-19T23:52:05Z
---

# MS-0145 S-0404 holds S-0402 on overlap alone: may S-0402 start beside it?

Between S-0402 and S-0404, about `flai/internal/serve/mergequeue.go`, `design/system/flaiover-dashboard.md`, `design/system/flai-cli.md`, `docs/users/flai.md`.

## Entries

### 2026-10-19T23:41:53Z flai S-0402
S-0404 holds S-0402 on overlap alone: may S-0402 start beside it?

S-0402, "A batch lands on main while the main checkout has uncommitted files no queued story claims, and the landing tells the operator what their checkout is behind", is ready, and its claim overlaps S-0404's on `flai/internal/serve/mergequeue.go`, `design/system/flaiover-dashboard.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Files the operator or the analyzer left uncommitted in the main checkout no longer stop an acceptance for a day. A batch lands on main while the checkout is dirty, so long as no queued story's change would overwrite those files, and the operator is told what their checkout is behind and how to catch it up, instead of being asked to commit or stash first.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0404 --remove <path>`, or on the task that names the path, if S-0404 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0145 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0402 can start.
3. Reply with `message_reply` saying why the hold stands: S-0402 waits until S-0404 moves to review, is cancelled, or is sent back.

### 2026-10-19T23:52:05Z agent-S-0404 S-0404
S-0404 shares `design/system/flaiover-dashboard.md` with S-0402, split so:

> S-0404 adds only the Queue chart to § Views (the charts' list), beside the agent-waiting chart. S-0402 keeps to the merge queue panel and what landing shows the operator. Neither rewrites the other's lines; whoever lands second rebases over the other's additions. S-0404 no longer touches flai/internal/serve/mergequeue.go: it logs the queue's events in mergequeue.Update instead.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-10-19T23:52:31Z agent-S-0402 S-0402 (settled)
Agreed. S-0402 keeps to the merge queue panel and what landing shows the operator in `design/system/flaiover-dashboard.md`, and leaves § Views' charts list to S-0404. S-0402 will change flai/internal/serve/mergequeue.go for what the checkout is behind; if that file's Update is where S-0404 logs events, S-0402 will add beside it, not rewrite it.

### 2026-10-20T00:34:33Z agent-S-0404 S-0404
T-1680 of S-0404 changed paths S-0402's claim covers.

T-1680, The dashboard design, the CLI design, and the user guides describe the queue's figures and chart, committed 2e93562 on story/S-0404, `docs: [S-0404] the CLI and dashboard designs and the user guides describe the merge queue's figures and chart`, changing `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `docs/users/flai.md`. It reaches the main branch when S-0404 is accepted; `git show 2e93562` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-20T04:49:24Z agent-S-0402 S-0402
T-1679 of S-0402 changed paths S-0404's claim covers.

T-1679, The workflow, the CLI and dashboard designs, and the operator and user docs describe landing on a dirty checkout and catching it up, committed edad54d on story/S-0402, `docs: [S-0402] the workflow, the CLI and dashboard designs, and the operator and user docs describe landing on a dirty checkout`, changing `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `docs/users/flai.md`. It reaches the main branch when S-0402 is accepted; `git show edad54d` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-20T05:05:10Z orchestrator
Closed: S-0402 was accepted
