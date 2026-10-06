---
id: T-1003
type: task
nature: remediation
title: Close I-0082 saying what fixed it
status: done
parent: S-0283
owner: alex
created: 2026-10-06T06:27:04Z
updated: 2026-10-06T10:10:36Z
transitions:
  - to: ready
    at: 2026-10-06T10:10:11Z
    by: agent-S-0283
  - to: in-progress
    at: 2026-10-06T10:10:11Z
    by: agent-S-0283
  - to: done
    at: 2026-10-06T10:10:36Z
    by: agent-S-0283
stream: S-0283
tags: [issues]
touches: [design/issues/I-0082-flai-s-permission-prompt-returns-a-result-claude-code-calls-invalid-so-a-claude-write-and-a-tmp-write-are-refused-at-once-with-no-thread.md, design/issues/summary.md]
after: [T-1001]
usage:
  source: log
  seconds: 25
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1734
      cache_read: 302092
      cache_write: 9638
      cost: 0.1607
---
# T-1003 Close I-0082 saying what fixed it

## Work

It waits for T-1001, because the reason names T-1001's fix and its test. It shares no path with T-1002, so the two run together.

In the story's worktree, run `flai issue close I-0082 --reason` with what fixed it: the cause T-1001 confirmed, that `permission_prompt` now answers with one text block and no structured content, and the test that reproduces it. Say too that the fix reaches running agents once a release with it is installed, since the MCP server agents use is the installed `flai`.

## Done when

- I-0082's status is closed with that reason, and `design/issues/summary.md` lists it as closed
- the story's second acceptance criterion is ticked with `flai criteria tick`

## Notes
