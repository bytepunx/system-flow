---
id: T-1918
type: task
nature: improvement
title: Document the embedded inbox and its fit in the user docs, the generated reference, and the session-start convention
status: done
parent: S-0450
owner: alex
created: 2026-10-20T19:05:45Z
updated: 2026-10-20T19:54:26Z
transitions:
  - to: ready
    at: 2026-10-20T19:52:14Z
    by: agent-S-0450
  - to: in-progress
    at: 2026-10-20T19:52:15Z
    by: agent-S-0450
  - to: done
    at: 2026-10-20T19:54:26Z
    by: agent-S-0450
stream: S-0450
tags: [flai, template]
touches: [docs/users/flai.md, docs/users/flai-reference.md, design/conventions/session-start.md, template/root/design/conventions/session-start.md]
after: [T-1916, T-1917]
usage:
  source: log
  seconds: 131
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 6278
      cache_read: 1185176
      cache_write: 59491
      cost: 0.7247
---
# T-1918 Document the embedded inbox and its fit in the user docs, the generated reference, and the session-start convention

## Work

- `docs/users/flai.md`: say in "Starting a story", "Closing a task", and "Serving agents over MCP" that the inbox `story_start` and `task_done` carry lists the story's threads, counts the others, and fits 40,000 bytes, and that `inbox` lists everything.
- `docs/users/flai-reference.md`: regenerate it with `make flai-reference` if T-1916 or T-1917 changed the commands' help; never edit it by hand.
- `design/conventions/session-start.md` and `template/root/design/conventions/session-start.md`: the sentence that `story_start` answers its first part, or only its header when that part does not fit, must stay true; change it in both, in the same words, only if the fit changed what it says. This is the convention text the template ships, so the change lands in the project's copy and the template's together.

It waits for T-1916 and T-1917, because it documents what they built and the help they changed.

## Done when

- The user docs describe the answers as built, and the reference matches `flai --help` output.
- The two session-start copies agree with each other and with the built behaviour.
- `flai check --strict` and the markdown lint pass.

## Notes
