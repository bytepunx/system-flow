---
id: T-0333
type: task
nature: feature
title: A story with an unresolved narrative question is refused into review, and its own inbox entries don't linger once it no longer needs the designer
status: done
parent: S-0089
owner: alex
created: 2026-09-23T00:12:03Z
updated: 2026-09-23T00:12:47Z
transitions:
  - to: ready
    at: 2026-09-23T00:12:26Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T00:12:28Z
    by: system-flow
  - to: done
    at: 2026-09-23T00:12:47Z
    by: system-flow
stream: S-0089
tags: []
---

# T-0333 A story with an unresolved narrative question is refused into review, and its own inbox entries don't linger once it no longer needs the designer

## Work
`inbox.designer` (`flai/internal/hostapi/people.go`) added a "question" entry for every hand-written bullet under every narrative's `## Open questions`, unconditionally — a thread mirrored there is tracked as a thread (answered once the designer's own word is last) and drops out correctly, but a freeform bullet has no such tracking, so it showed forever, including for stories already in review, done, or cancelled, and nothing stopped a story from moving to review with one still open.

Two fixes, both in `flai/`:
- `internal/workitem/narrative.go`: exported `OpenQuestions(body) []string`, moved from `hostapi` (now the one implementation, used by both the rule below and the inbox).
- `internal/workitem/rules.go`'s `Move`: a story moving to `review` is refused if its narrative (`wip/agents/<id>.md`) still has an open question, naming it. No narrative, or none unreadable, is not an error (fails open) — this is the shared path behind `flai move`, the MCP `item_move`, and the dashboard's `item.move` alike.
- `internal/hostapi/people.go`'s `inbox.designer`: a narrative's question entries are only added while its story is still `backlog`, `ready`, or `in-progress` (found by ID in the item list); once it is `review`, `done`, `cancelled`, or archived (not found), they are left out — a second, independent guarantee, for anything already in review before this shipped, or moved there by hand.

## Done when
The three acceptance criteria: a story with an unresolved question refuses `flai move ... review`, naming the question; the question no longer appears in the inbox once the story is in review (or done, or cancelled); answering it (removing the bullet, as the convention already has an agent do by hand) lets the story through and never reappears. Go tests in `internal/workitem` and `internal/hostapi`; the flaiover fixture in `src/lib/server/inbox.test.ts` (which moved a story to review with an unanswered question, unnoticed until now) updated to answer it first, plus a new test for the refusal itself. `go test ./...`, `pnpm run check`/`test:unit`/`lint`, `flai check --strict` clean.

## Notes
Criterion 2 ("cannot be moved to review") is the preventive fix; criterion 1/3 (never shown once in review, never lingers once answered) is enforced independently by the inbox filter, in case a story got there some other way (an older item, a manual edit).
