---
id: T-2076
type: task
nature: remediation
title: The design and the user docs say a checked edit refuses only what it introduces inside its story's scope, and warns on the rest
status: done
parent: S-0486
owner: alex
created: 2026-10-23T17:50:58Z
updated: 2026-10-24T00:11:19Z
transitions:
  - to: ready
    at: 2026-10-24T00:00:08Z
    by: agent-S-0486
  - to: in-progress
    at: 2026-10-24T00:00:09Z
    by: agent-S-0486
  - to: done
    at: 2026-10-24T00:11:18Z
    by: agent-S-0486
stream: S-0486
tags: [flai]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/edit.go, flai/cmd/criteria.go, design/adrs/0197-a-checked-edit-of-a-work-item-refuses-only-a-new-finding-inside-its-story-s.md, docs/users/flaiover.md]
after: [T-2075]
usage:
  source: log
  seconds: 669
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 48
      output: 14064
      cache_read: 2759521
      cache_write: 88458
      cost: 1.3923
---
# T-2076 The design and the user docs say a checked edit refuses only what it introduces inside its story's scope, and warns on the rest

## Work

Describe the behaviour T-2075 builds, citing the ADR of T-2074:

- `design/system/flai-cli.md`, the Commands table: the `flai edit` row says the check runs with the change in place and refuses only a finding the edit introduces inside the item's story's scope (`check.ScopeToStories`, the story itself or a task's parent), returns a new finding outside it, or one the scope leaves out, as a warning or not at all, and refuses a new finding anywhere but on a written path for an epic; the `flai criteria` row says `tick` and `untick` are refused on the same terms, so another story's unwritten narrative or an untracked file of another story no longer refuses a tick (I-0182).
- `flai/cmd/edit.go` and `flai/cmd/criteria.go`: the long help that `make flai-reference` renders into `docs/users/flai-reference.md` § `flai edit` and § `flai criteria`, which now say "a change flai check finds fault with" without the scope. Regenerate the reference with `make flai-reference` and commit the generated file with the sources.
- `docs/users/flai.md` § Ticking acceptance criteria and the `criteria_tick` row of the MCP tools table, where the sentence on the checks `flai edit` makes is finished with the scope.

Waits for T-2075 so the words describe what the code does.

## Done when

The three documents and the two help texts say what the ADR decides and name it; `make flai-reference` leaves `docs/users/flai-reference.md` unchanged after the commit; the markdown lint and `flai check --strict` report nothing on them.

## Notes

`design/system/flai-cli.md`, `docs/users/flai.md`, and `docs/users/flai-reference.md` are touches of S-0429, in progress, and S-0468, in review; the lines this task changes are the `flai edit` and `flai criteria` rows, which neither story changes.
