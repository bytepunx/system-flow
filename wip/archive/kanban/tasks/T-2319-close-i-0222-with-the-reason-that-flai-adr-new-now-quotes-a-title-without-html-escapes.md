---
id: T-2319
type: task
nature: remediation
title: Close I-0222 with the reason that flai adr new now quotes a title without HTML escapes
status: done
parent: S-0538
owner: alex
created: 2026-10-29T02:58:03Z
updated: 2026-10-29T03:09:57Z
transitions:
  - to: ready
    at: 2026-10-29T03:09:24Z
    by: system-flow
  - to: in-progress
    at: 2026-10-29T03:09:25Z
    by: system-flow
  - to: done
    at: 2026-10-29T03:09:57Z
    by: agent-S-0538
stream: S-0538
tags: [flai]
touches: [design/issues/I-0222-flai-adr-topics-rewrites-an-adr-s-front-matter-with-and-in-its-title-escaped-as-and.md, design/issues/summary.md]
after: [T-2318]
usage:
  source: log
  seconds: 32
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 15
      output: 2994
      cache_read: 1099664
      cache_write: 29255
      cost: 0.5139
---
# T-2319 Close I-0222 with the reason that flai adr new now quotes a title without HTML escapes

## Work

- In the story's worktree, run `flai issue close I-0222 --reason "<reason>"` from the shell, the reason saying that the escapes came from `flai adr new`'s quoting with `json.Marshal`, not from `flai adr topics`, and that the quoting now writes the less-than, greater-than, and ampersand characters as themselves, with the test that reproduces it named.
- Give the reason through the shell, not a JSON tool argument: a JSON string decodes a backslash-u escape to its character, which is how the issue's own title lost the escapes it describes.
- Check that `design/issues/summary.md` no longer lists I-0222.

Waits for T-2318: the reason names the fix and its test, which exist only once T-2318 is done.

## Done when

- I-0222's status is closed, with the reason above.
- `flai check --strict` reports nothing on `design/issues`.

## Notes

`flai issue close` writes the issue file and regenerates the summary; both are this task's touches.
