---
id: T-2722
type: task
nature: remediation
title: Close I-0249 with the fix named
status: done
parent: S-0616
owner: alex
created: 2026-11-09T10:21:05Z
updated: 2026-11-09T10:32:16Z
transitions:
  - to: ready
    at: 2026-11-09T10:32:05Z
    by: agent-S-0616
  - to: in-progress
    at: 2026-11-09T10:32:06Z
    by: agent-S-0616
  - to: done
    at: 2026-11-09T10:32:16Z
    by: agent-S-0616
stream: S-0616
tags: [issues]
touches: [design/issues/I-0249-testrenderprototypetemplate-refuses-any-double-curly-brace-in-a-rendered.md, design/issues/summary.md]
after: [T-2721]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 48
      cache_read: 536327
      cache_write: 949
      cost: 0.1955
---
# T-2722 Close I-0249 with the fix named

## Work

Run `flai issue close I-0249 --reason "<reason>"` in the story's worktree. The reason should say three things:

- `TestRenderPrototypeTemplate` now looks for `{{` only in files rendered from a `.tmpl` source and in files that differ from their verbatim source;
- the S-0424 and S-0554 workarounds are gone;
- the formulas README shows `{{name}}` again.

Check that `design/issues/summary.md` no longer lists I-0249.

This task waits for T-2721: the reason cites the README shown passing on the real template, which T-2721 does.

## Done when

- I-0249's status is closed, and the reason names the fix.
- `design/issues/summary.md` no longer lists it.
- `flai check --strict` reports nothing new on the two files.

## Notes
