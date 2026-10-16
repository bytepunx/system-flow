---
id: T-1694
type: task
nature: feature
title: The operator, user, and design docs say what flai guard holds on OpenCode, and agent-adapters.md records the sub-agent finding's consequence
status: backlog
parent: S-0407
owner: alex
created: 2026-10-16T05:15:52Z
updated: 2026-10-16T05:15:52Z
transitions: []
stream: S-0407
tags: [cli]
touches: [docs/operators/index.md, docs/users/flai.md, design/system/flai-cli.md, design/system/agent-adapters.md]
after: [T-1693]
---
# T-1694 The operator, user, and design docs say what flai guard holds on OpenCode, and agent-adapters.md records the sub-agent finding's consequence

## Work

- Update the table in `docs/operators/index.md` § Harnesses and the guard's section of `docs/users/flai.md`: the plugin, what it refuses, and whether roles may run.
- In `design/system/flai-cli.md`, list `opencode` beside `claude-code` among the guard's readers.
- In `design/system/agent-adapters.md` § OpenCode, record the guard as built, and, where `Guard` stayed false, why and what OpenCode would need to change.
- Waits for T-1693: the docs describe the tested behaviour.

## Done when

The docs say what the story built, and the markdown lint passes.

## Notes
