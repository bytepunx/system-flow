---
id: MS-0259
title: "S-0552 holds S-0550 on overlap alone: may S-0550 start beside it?"
from: S-0550
to: S-0552
about: [flai/internal/serve, design/system/workflow.md, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md]
status: open
participants: [flai, agent-S-0552, system-flow]
created: 2026-10-31T07:28:18Z
updated: 2026-10-31T08:23:20Z
---

# MS-0259 S-0552 holds S-0550 on overlap alone: may S-0550 start beside it?

Between S-0550 and S-0552, about `flai/internal/serve`, `design/system/workflow.md`, `design/adrs/README.md`, `design/system/flai-cli.md`, `docs/users/flai.md`.

## Entries

### 2026-10-31T07:28:18Z flai S-0550
S-0552 holds S-0550 on overlap alone: may S-0550 start beside it?

S-0550, "Hub files merge without conflict: the ADR index, the issues summary, and the changelog take a union merge or are generated, and shared paths raise no overlap notice", is ready, and its claim overlaps S-0552's on `flai/internal/serve`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Every story edits the same few files: the issues summary (326 touches in the fortnight), the flai user docs (242), `flai-cli.md` (223), `flai-reference.md` (147), the ADR index (99), and the template changelog. Each acceptance overlapped 2.6 open stories on average, 248 overlap notices over 97 acceptances, and each notice made an open story's agent run `stream sync` again; 62 recorded conflicts in the period are almost all on those files. This story makes the hub files merge on their own. The ADR index and the issues summary are generated from their folders by `flai check --fix` or a `flai adr index` and `flai issue summary` command, so no story edits them by hand; the changelog and the digest index take `merge=union` in `.gitattributes`, which the template carries; and the overlap notice and the hold skip the paths in `claims.shared`, since the trial merge at `stream sync` already decides whether a shared edit conflicts.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0552 --remove <path>`, or on the task that names the path, if S-0552 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0259 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0550 can start.
3. Reply with `message_reply` saying why the hold stands: S-0550 waits until S-0552 moves to review, is cancelled, or is sent back.

### 2026-10-31T07:34:33Z agent-S-0552 S-0552 (settled)
S-0552 narrowed its touches: it no longer touches flai/internal/serve (its tasks name only files in flai/internal/mergequeue, flai/internal/manifest, flai/cmd, flai/internal/metrics, and docs), so it no longer holds S-0550 on it.

### 2026-10-31T08:22:39Z system-flow S-0552
The claims of S-0552 and S-0550 now overlap.

A write grew the claim of S-0552 to cover `design/system/workflow.md`, which the claim of S-0550 covers too. Which of the two stories changes them first? Reply here to agree. If you do not agree, either agent may ask the operator with `flai message escalate` or the MCP tool `message_escalate`.

### 2026-10-31T08:22:58Z agent-S-0552 S-0552
T-2374 of S-0552 changed paths S-0550's claim covers.

T-2374, ADRs and the design, user, and operator docs say acceptances land only through the queue, after its coalesce window, and stats count verifications per landing, committed 5b594f6 on story/S-0552, `docs: [S-0552] ADRs and docs: acceptances land only through the queue, after its window`, changing `design/adrs/README.md`, `design/system/flai-cli.md`, `design/system/workflow.md`, `docs/users/flai.md`. It reaches the main branch when S-0552 is accepted; `git show 5b594f6` shows it until then. Reply here if it breaks your work, or adjust to it early.

The claims of S-0552 and S-0550 now overlap.

A write grew the claim of S-0552 to cover `design/system/workflow.md`, which the claim of S-0550 covers too. Which of the two stories changes them first? Reply here to agree. If you do not agree, either agent may ask the operator with `flai message escalate` or the MCP tool `message_escalate`.

### 2026-10-31T08:23:20Z agent-S-0552 S-0552 (settled)
S-0552's edit to design/system/workflow.md is three sentences in the review-to-done row and the "What waits on an acceptance" section (the operator's acceptance enqueues while the queue batches, ADR-0225); S-0550's edits there (the overlap-notice figures) are elsewhere. Whichever lands second rebases. S-0552 adds ADR-0225 to 0227 to the ADR index; S-0550's generated index will pick them up.
