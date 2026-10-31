---
id: MS-0254
title: "S-0548 holds S-0550 on overlap alone: may S-0550 start beside it?"
from: S-0550
to: S-0548
about: [flai/internal/serve/agents_test.go, flai/internal/serve/serve_test.go, flai/internal/serve/forgemerge_test.go, flai/internal/serve/serve.go, flai/internal/serve/projects_test.go, flai/internal/serve/mcp_test.go, flai/internal/serve/budget_test.go]
status: open
participants: [flai]
created: 2026-10-31T07:26:05Z
updated: 2026-10-31T07:26:05Z
---

# MS-0254 S-0548 holds S-0550 on overlap alone: may S-0550 start beside it?

Between S-0550 and S-0548, about `flai/internal/serve/agents_test.go`, `flai/internal/serve/serve_test.go`, `flai/internal/serve/forgemerge_test.go`, `flai/internal/serve/serve.go`, `flai/internal/serve/projects_test.go`, `flai/internal/serve/mcp_test.go`, `flai/internal/serve/budget_test.go`.

## Entries

### 2026-10-31T07:26:05Z flai S-0550
S-0548 holds S-0550 on overlap alone: may S-0550 start beside it?

S-0550, "Hub files merge without conflict: the ADR index, the issues summary, and the changelog take a union merge or are generated, and shared paths raise no overlap notice", is ready, and its claim overlaps S-0548's on `flai/internal/serve/agents_test.go`, `flai/internal/serve/serve_test.go`, `flai/internal/serve/forgemerge_test.go`, `flai/internal/serve/serve.go`, `flai/internal/serve/projects_test.go`, `flai/internal/serve/mcp_test.go`, `flai/internal/serve/budget_test.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Every story edits the same few files: the issues summary (326 touches in the fortnight), the flai user docs (242), `flai-cli.md` (223), `flai-reference.md` (147), the ADR index (99), and the template changelog. Each acceptance overlapped 2.6 open stories on average, 248 overlap notices over 97 acceptances, and each notice made an open story's agent run `stream sync` again; 62 recorded conflicts in the period are almost all on those files. This story makes the hub files merge on their own. The ADR index and the issues summary are generated from their folders by `flai check --fix` or a `flai adr index` and `flai issue summary` command, so no story edits them by hand; the changelog and the digest index take `merge=union` in `.gitattributes`, which the template carries; and the overlap notice and the hold skip the paths in `claims.shared`, since the trial merge at `stream sync` already decides whether a shared edit conflicts.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0548 --remove <path>`, or on the task that names the path, if S-0548 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0254 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0550 can start.
3. Reply with `message_reply` saying why the hold stands: S-0550 waits until S-0548 moves to review, is cancelled, or is sent back.
