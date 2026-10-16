---
id: T-1689
type: task
nature: feature
title: The operator, user, and design docs, the reference, and the template describe the opencode harness and what a story on it has and lacks
status: backlog
parent: S-0406
owner: alex
created: 2026-10-16T05:14:43Z
updated: 2026-10-16T05:14:43Z
transitions: []
stream: S-0406
tags: [cli, template]
touches: [docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, template/root/opencode.json, template/CHANGELOG.md]
after: [T-1684, T-1688]
---
# T-1689 The operator, user, and design docs, the reference, and the template describe the opencode harness and what a story on it has and lacks

## Work

- Add § Harnesses to `docs/operators/index.md`: installing OpenCode on the host, `flai serve agent harness opencode`, `--deny-protected` or `auto-approve`, and a table of what a story on `opencode` has and lacks today (no guard, no ask, no roles, no measured cost).
- Name the harness in `docs/users/flai.md` § Who works a story: its agent, and in `design/system/flai-cli.md` where the adapters are listed; regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Add `template/root/opencode.json`: a skeleton with `$schema` and no MCP entry and no key, since flai adds the entry at start. The template keeps `CLAUDE.md`; say so, and the new check warning, in `template/CHANGELOG.md`.
- Waits for T-1684 (the check's wording) and T-1688 (the tested behaviour the docs describe).

## Done when

The docs and the template say what the story built, the reference is regenerated, and the markdown lint and the settings index test pass.

## Notes
