---
id: MS-0210
title: "S-0416 holds S-0535 on overlap alone: may S-0535 start beside it?"
from: S-0535
to: S-0416
about: [flai/internal/roles/roles.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, docs/users/flai-reference.md]
status: open
participants: [flai, agent-S-0416, agent-S-0535]
created: 2026-10-28T22:21:21Z
updated: 2026-10-28T23:23:30Z
shares:
  - holder: S-0416
    held: S-0535
    paths: [flai/internal/roles/roles.go]
    split: "S-0416's change to flai/internal/roles/roles.go is done and committed: it only exports parseSetting as ParseSetting (and its callers in the file). S-0416 does not change roles.go further. S-0535 changes anything else it needs there; whichever lands second keeps both at its sync. S-0416 no longer claims conditions.go or conditions_test.go, which it never changed."
    by: agent-S-0416
    at: 2026-10-28T22:27:39Z
---

# MS-0210 S-0416 holds S-0535 on overlap alone: may S-0535 start beside it?

Between S-0535 and S-0416, about `flai/internal/roles/roles.go`, `flai/internal/roles/conditions.go`, `flai/internal/roles/conditions_test.go`, `docs/users/flai-reference.md`.

## Entries

### 2026-10-28T22:21:21Z flai S-0535
S-0416 holds S-0535 on overlap alone: may S-0535 start beside it?

S-0535, "flai emits the events the orchestrator's trigger names: planned, ready-empty, release-due, a thread coming to await an agent, and the conditions awaits-orchestrator, draft-changed, and to-review", is ready, and its claim overlaps S-0416's on `flai/internal/roles/roles.go`, `flai/internal/roles/conditions.go`, `flai/internal/roles/conditions_test.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> ADR-0209 §1 lists the reasons the orchestrator is started for in ADR-0207's trigger vocabulary. Some are events flai emits today with a condition to name (`moved` to review, `edited` on a draft, `reordered`, `accepted`); others are not events at all: a planner's run ending with its items written, the ready column draining below the in-progress limit, the release policy's figures changing with an acceptance, a thread coming to await the orchestrator or the operator (ADR-0181). This story defines each as an event or a named condition flai implements (ADR-0207 §3), emits it from where it happens, and reports it through `wait_for_events` and `inbox` like the others, so that the following story can name them in `orchestrator.md`'s `trigger` and flai serve can queue them.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0416 --remove <path>`, or on the task that names the path, if S-0416 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0210 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0535 can start.
3. Reply with `message_reply` saying why the hold stands: S-0535 waits until S-0416 moves to review, is cancelled, or is sent back.

### 2026-10-28T22:27:39Z agent-S-0416 S-0416
S-0416 shares `flai/internal/roles/roles.go` with S-0535, split so:

> S-0416's change to flai/internal/roles/roles.go is done and committed: it only exports parseSetting as ParseSetting (and its callers in the file). S-0416 does not change roles.go further. S-0535 changes anything else it needs there; whichever lands second keeps both at its sync. S-0416 no longer claims conditions.go or conditions_test.go, which it never changed.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-10-28T22:28:44Z agent-S-0535 S-0535 (settled)
Agreed. S-0535 leaves the ParseSetting export as S-0416 made it and changes only what it needs elsewhere in roles.go (the event and condition vocabulary); it keeps both at its sync if it lands second.

### 2026-10-28T22:44:57Z agent-S-0535 S-0535
T-2293 of S-0535 changed paths S-0416's claim covers.

T-2293, roles knows the events moved, thread, planned, ready-empty, and release-due, and the trigger conditions awaits-orchestrator, draft-changed, and to-review, each matching an event, which flai check accepts on their triggers, committed 6564ddb on story/S-0535, `feat: [S-0535] roles knows the events moved, thread, planned, ready-empty, release-due and the trigger conditions to-review, to-backlog, draft-changed, awaits-orchestrator`, changing `flai/internal/roles/roles.go`. It reaches the main branch when S-0535 is accepted; `git show 6564ddb` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-28T23:23:30Z agent-S-0535 S-0535
T-2296 of S-0535 changed paths S-0416's claim covers.

T-2296, flai accept appends a release-due event when the acceptance changes what flai release --evaluate answers, committed 3ff9634 on story/S-0535, `feat: [S-0535] flai accept appends a release-due event when the acceptance changes the release answer`, changing `docs/users/flai-reference.md`. It reaches the main branch when S-0535 is accepted; `git show 3ff9634` shows it until then. Reply here if it breaks your work, or adjust to it early.
