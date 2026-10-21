---
id: T-1973
type: task
nature: remediation
title: flai check --story takes several stories, a finding outside every one of them a note
status: done
parent: S-0464
owner: alex
created: 2026-10-21T07:30:02Z
updated: 2026-10-21T07:49:09Z
transitions:
  - to: ready
    at: 2026-10-21T07:32:20Z
    by: agent-S-0464
  - to: in-progress
    at: 2026-10-21T07:32:23Z
    by: agent-S-0464
  - to: done
    at: 2026-10-21T07:49:09Z
    by: agent-S-0464
stream: S-0464
tags: [flai]
touches: [flai/cmd/check.go, flai/cmd/check_test.go, flai/internal/check/story.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, docs/users/flai-reference.md]
---
# T-1973 flai check --story takes several stories, a finding outside every one of them a note

## Work

A merge queue batch holds up to `merge_queue.batch_max` stories, and its repository check must count only what they bring, as a close-out's counts only what its one story brings (ADR-0085). Today `flai check --story` names one story.

- Let `--story` take several stories, comma-separated or repeated (`S-0001,S-0002`). A finding is inside when it is inside any story named, and a note (`Outside`) only when it is outside every one. Each rule that passes over findings outside the story (`item.archive`, `board.wip-limit`, `issues.no-story`, `threads.archived`, `components.examples`, the narrative rules, `wip.overlap`) treats the set the same way.
- `check.RunScoped` and `ScopeToStory` in `flai/internal/check/story.go` and `scope.go` take the list; one story keeps today's behaviour exactly.
- `--record-issues` with more than one story is refused with a message saying it records the findings outside one story's close-out; a batch never records issues, since its worktree is never committed and a dirty one would stop the bisection's checkouts.
- The summary line names the stories (`outside S-0001, S-0002`).
- Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

It waits for nothing: layer 1, beside the merge queue task, with no path in common.

## Done when

- A test in `flai/internal/check/scope_test.go` scopes a fixture with a finding inside the second of two stories and one outside both, and gets the first counted and the second a note.
- A test in `flai/cmd/check_test.go` runs `flai check --strict --story S-0001,S-0002` on a fixture whose only finding is outside both and exits 0, runs it unscoped and exits 1, and is refused `--record-issues` with two stories.
- `flai test` passes on the changed paths.

## Notes
