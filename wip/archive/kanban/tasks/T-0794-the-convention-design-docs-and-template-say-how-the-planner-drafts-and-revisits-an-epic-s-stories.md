---
id: T-0794
type: task
nature: feature
title: The convention, design, docs, and template say how the planner drafts and revisits an epic's stories
status: done
parent: S-0209
owner: alex
created: 2026-10-04T04:03:02Z
updated: 2026-10-04T04:13:11Z
transitions:
  - to: ready
    at: 2026-10-04T04:03:32Z
    by: agent-S-0209
  - to: in-progress
    at: 2026-10-04T04:10:01Z
    by: agent-S-0209
  - to: done
    at: 2026-10-04T04:13:11Z
    by: agent-S-0209
stream: S-0209
tags: []
touches: [template/root/design/conventions/strategic-agents.md, design/conventions/strategic-agents.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users, template/template.yaml, template/CHANGELOG.md]
after: [T-0791, T-0792, T-0793]
usage:
  source: log
  seconds: 190
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 39
      output: 8893
      cache_read: 3397337
      cache_write: 38604
      cost: 1.1129
---
# T-0794 The convention, design, docs, and template say how the planner drafts and revisits an epic's stories

## Work

- `strategic-agents.md`'s "As the planner", in `template/root/design/conventions/` and copied above the marker in `design/conventions/`: the plan thread names the stories, their order, and the assumptions; a revisit re-enriches each open story and proposes splits, merges, additions, and drops in the epic's thread, creating drafts only for additions; the final summary names the stories created and revisited.
- `design/system/strategic-agents.md`: the prompt, the guard's draft rule, the checked `item_new`, and the activity entry's items.
- `design/system/flai-cli.md` and `docs/users` (`flai-reference.md` regenerated with `make flai-reference`): the guard, `item_new`, and the planner's entry.
- Template release: bump `template/template.yaml` and add the `template/CHANGELOG.md` entry, naming `planner.md` too.
- Waits for T-0791, T-0792, and T-0793: it describes what they built.

## Done when

- [ ] The two convention copies agree above the marker.
- [ ] `make flai-reference` leaves no diff, and the docs test passes.

## Notes

Shares the second layer with the `planner.md` task: their paths do not meet.
