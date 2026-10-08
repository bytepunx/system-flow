---
id: T-1301
type: task
nature: improvement
title: A check scoped to a story leaves out a markdown finding on another open story's narrative, so a close-out records none
status: done
parent: S-0318
owner: alex
created: 2026-10-08T00:10:38Z
updated: 2026-10-08T04:43:31Z
transitions:
  - to: ready
    at: 2026-10-08T04:35:38Z
    by: agent-S-0318
  - to: in-progress
    at: 2026-10-08T04:35:38Z
    by: agent-S-0318
  - to: done
    at: 2026-10-08T04:43:31Z
    by: agent-S-0318
stream: S-0318
tags: [flai]
touches: [flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md]
after: [T-1300]
usage:
  source: log
  seconds: 473
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 22
      output: 5866
      cache_read: 1702897
      cache_write: 37138
      cost: 0.7188
---
# T-1301 A check scoped to a story leaves out a markdown finding on another open story's narrative, so a close-out records none

## Work

Build the remedy T-1300's ADR decides. It waits for T-1300, because the ADR settles what is left out and what is still recorded.

- In `ScopeToStory` (`flai/internal/check/scope.go`), drop a finding whose rule starts with `markdown.` when its path is the narrative of another story that is not done or cancelled. Drop it beside the `item.archive` and `wip.overlap` cases, with `res.drop`, which takes back the counts. Tell a narrative by its path, `repo.NarrativePath(id)` for a story the repository holds, not by matching the file name alone, and read the story's state from its item.
- Keep every other `markdown.*` finding outside the story as it is now, an outside note. That covers a thread, a task, a narrative under `wip/archive/agents`, and the narrative of a done or cancelled story.
- In `recordOutside` (`flai/cmd/check.go`), change the comment that names what never reaches it. In the command's help, where it names `item.archive`, say that a markdown finding on another open story's narrative is left out too. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Reproduce I-0096 in tests:
  - `flai/internal/check/scope_test.go`: another story in progress whose narrative holds a code span with a trailing space, beside the story checked. The scoped result holds no `markdown.MD038`, and its counts do not include it.
  - The same fixture with that story done or cancelled, or the finding on a thread: the finding is kept as an outside note.
  - `flai/cmd/check_test.go`: `flai check --story S-nnnn --record-issues` with the first fixture records no issue for `markdown.MD038`.
- An unscoped `flai check` still warns the `markdown.MD038`; show it in a test.

## Done when

- The tests above pass under `flai test` on the changed paths, and the full Go tests pass.
- `docs/users/flai-reference.md` matches the help, as `make flai-reference` writes it.
- `flai check --strict` passes.

## Notes
