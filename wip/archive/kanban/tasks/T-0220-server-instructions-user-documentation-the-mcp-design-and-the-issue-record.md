---
id: T-0220
type: task
nature: feature
title: Server instructions, user documentation, the MCP design, and the issue record
status: done
parent: S-0061
owner: alex
created: 2026-09-19T08:15:26Z
updated: 2026-09-19T08:19:17Z
transitions:
  - to: ready
    at: 2026-09-19T08:19:17Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T08:19:17Z
    by: system-flow
  - to: done
    at: 2026-09-19T08:19:17Z
    by: system-flow
stream: S-0061
tags: []
---

# T-0220 Server instructions, user documentation, the MCP design, and the issue record

## Work
The server instructions and the `inbox` and `wait_for_events` tool descriptions say what a first look returns and what `changes_omitted` means. `docs/users/flai.md` and the MCP part of `design/system/flai-cli.md` say the same. Record the issue on this branch, as the story's notes ask: the unbounded first look, found on 2026-09-19, a defect in S-0058. The narrative of S-0053 also noted that `inbox` reports an agent's own block and unblock back to it, because those intervals carry no author; say in this story's open questions whether that belongs here or in a story of its own, without widening this one.

## Done when
- The instructions, descriptions, user docs, and design agree with the behaviour
- The issue is recorded on the branch
- The criteria are ticked and `flai check --strict` is clean

## Notes
