---
id: T-1189
type: task
nature: feature
title: The design and the users' guide describe messages between stories, and the command reference is regenerated
status: done
parent: S-0330
owner: alex
created: 2026-10-07T20:14:26Z
updated: 2026-10-07T20:54:08Z
transitions:
  - to: ready
    at: 2026-10-07T20:50:21Z
    by: agent-S-0330
  - to: in-progress
    at: 2026-10-07T20:50:21Z
    by: agent-S-0330
  - to: done
    at: 2026-10-07T20:54:08Z
    by: agent-S-0330
stream: S-0330
tags: [flai]
touches: [design/system/agent-coordination.md, design/system/repository-layout.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/users/conventions.md, template/root/wip/messages/README.md]
after: [T-1187, T-1188]
usage:
  source: log
  seconds: 227
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 76
      output: 33465
      cache_read: 4426209
      cache_write: 145102
      cost: 2.4279
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
