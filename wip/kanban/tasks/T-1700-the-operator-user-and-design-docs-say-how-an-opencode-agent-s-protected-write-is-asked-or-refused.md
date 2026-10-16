---
id: T-1700
type: task
nature: feature
title: The operator, user, and design docs say how an OpenCode agent's protected write is asked or refused
status: backlog
parent: S-0408
owner: alex
created: 2026-10-16T05:16:48Z
updated: 2026-10-16T05:16:48Z
transitions: []
stream: S-0408
tags: [cli]
touches: [docs/operators/index.md, docs/users/flai.md, design/system/flai-cli.md, design/system/agent-adapters.md]
after: [T-1699]
---
# T-1700 The operator, user, and design docs say how an OpenCode agent's protected write is asked or refused

## Work

- Update the table in `docs/operators/index.md` § Harnesses: a story on `opencode` no longer needs `auto-approve` or `--deny-protected` under hold-and-ask, or still needs `--auto` with deny rules.
- Say in `docs/users/flai.md`, where `permission_prompt` is described, how the question reaches the owner from OpenCode and what the agent does on a refusal.
- Record the way built and its bound in `design/system/flai-cli.md` and `agent-adapters.md` § OpenCode, linking the ADR.
- Waits for T-1699: the docs describe the tested behaviour.

## Done when

The docs say what the story built, and the markdown lint passes.

## Notes
