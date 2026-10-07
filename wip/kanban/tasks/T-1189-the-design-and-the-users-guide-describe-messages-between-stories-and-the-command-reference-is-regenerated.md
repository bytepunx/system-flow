---
id: T-1189
type: task
nature: feature
title: The design and the users' guide describe messages between stories, and the command reference is regenerated
status: backlog
parent: S-0330
owner: alex
created: 2026-10-07T20:14:26Z
updated: 2026-10-07T20:14:26Z
transitions: []
stream: S-0330
tags: [flai]
touches: [design/system/agent-coordination.md, design/system/repository-layout.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1187, T-1188]
---
# T-1189 The design and the users' guide describe messages between stories, and the command reference is regenerated

## Work

Document what T-1187 and T-1188 built. It waits for both, so the docs say what the commands do.

- `design/system/agent-coordination.md`: a section on messages between stories, linking the ADR.
- `design/system/repository-layout.md`: the `wip/messages/` folder, if the ADR keeps messages there.
- `design/system/flai-cli.md` § Commands and `docs/users/flai.md`: `flai message` and its subcommands.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

## Done when

- Each document names the commands and links the ADR.
- The markdown lint passes on the changed files.
- `flai check --strict` passes.

## Notes
