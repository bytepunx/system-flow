---
id: T-0791
type: task
nature: feature
title: The planner's prompt asks for a plan thread, revisit proposals, and a summary naming the stories, and the guard holds its stories to drafts
status: done
parent: S-0209
owner: alex
created: 2026-10-04T04:02:44Z
updated: 2026-10-04T04:08:06Z
transitions:
  - to: ready
    at: 2026-10-04T04:03:31Z
    by: agent-S-0209
  - to: in-progress
    at: 2026-10-04T04:03:32Z
    by: agent-S-0209
  - to: done
    at: 2026-10-04T04:08:06Z
    by: agent-S-0209
stream: S-0209
tags: []
touches: [flai/internal/harness, flai/internal/guard, flai/cmd/guard.go]
usage:
  source: log
  seconds: 274
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 59
      output: 266
      cache_read: 1996600
      cache_write: 81583
      cost: 0.8477
---
# T-0791 The planner's prompt asks for a plan thread, revisit proposals, and a summary naming the stories, and the guard holds its stories to drafts

## Work

- `planPrompt` (`flai/internal/harness/harness.go`), for an epic: create each story with `draft` true in the backlog; open one thread on the epic summarising the plan: the stories, their order (`after`), and the assumptions made; when the epic has stories, revisit each one not done or cancelled, re-enrich it, and propose in that thread each story it would split, merge, add, or drop, creating drafts only for additions, never cancelling or rewriting a finalized story's words without asking; end with one line naming the stories created and the stories revisited.
- `flai guard` in a planner session (`flai/internal/guard`): refuse `item_new` of a story without `draft: true`, and `flai story new` without `--draft`, saying that the planner's stories are drafts for the operator to finalize. `Event.ToolInput` reads `type` and `draft`.
- `flai guard`'s help in `flai/cmd/guard.go` says so.
- Waits for nothing: the first layer.

## Done when

- [ ] `harness_test.go` pins the epic prompt's thread, revisit, draft, and summary instructions.
- [ ] `guard_test.go` covers a planner's story with and without draft over MCP and the CLI, and a task or an epic, which pass.
- [ ] `go test ./internal/harness ./internal/guard ./cmd -run Guard` passes.

## Notes

The guard is where the planner's rules are enforced (S-0208); a story it writes that is not a draft could be moved to ready by anyone without the operator finalizing it.
