---
id: T-2408
type: task
nature: improvement
title: The work-management and strategic-agents conventions state the trial merge and the after rule, in the tree and the template, with a changelog entry
status: done
parent: S-0551
owner: alex
created: 2026-10-31T07:44:19Z
updated: 2026-10-31T10:06:28Z
transitions:
  - to: ready
    at: 2026-10-31T08:57:56Z
    by: agent-S-0551
  - to: in-progress
    at: 2026-10-31T08:57:59Z
    by: agent-S-0551
  - to: done
    at: 2026-10-31T10:06:28Z
    by: agent-S-0551
stream: S-0551
tags: [flai, conventions, template]
touches: [design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md]
after: [T-2388]
usage:
  source: log
  seconds: 4109
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 70
      output: 22157
      cache_read: 4357054
      cache_write: 159799
      cost: 2.2053
---
# T-2408 The work-management and strategic-agents conventions state the trial merge and the after rule, in the tree and the template, with a changelog entry

## Work

State the ADR's rules in the conventions, baseline first in the template, then the tree's copy:

- `work-management.md`: a ready story that overlaps an in-progress one starts on a clean trial merge and is held only on a conflict; `after` names a story whose output this one reads, with the reason in the Notes; `flai check` advises of one with none.
- `strategic-agents.md`, under "As the planner": the same `after` rule, in place of "only a story whose accepted change it needs on main".
- `template/CHANGELOG.md`: one entry for the convention change and the planner's role file, which the role-file task changes.

It waits for the ADR task, whose wording it carries. It shares no file with the role-file task, so the two run together.

## Done when

- Both conventions, in the tree and the template, state the trial merge and the `after` rule, and the template copies match the tree's baseline.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes

The baseline above each convention's marker changes in the template and the tree together, as `CLAUDE.md` says for convention changes.
