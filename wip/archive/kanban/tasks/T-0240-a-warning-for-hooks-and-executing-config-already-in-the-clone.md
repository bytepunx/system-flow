---
id: T-0240
type: task
nature: feature
title: A warning for hooks and executing config already in the clone
status: done
parent: S-0064
owner: alex
created: 2026-09-19T10:18:10Z
updated: 2026-09-19T10:24:52Z
transitions:
  - to: ready
    at: 2026-09-19T10:24:52Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T10:24:52Z
    by: system-flow
  - to: done
    at: 2026-09-19T10:24:52Z
    by: system-flow
stream: S-0064
tags: []
---

# T-0240 A warning for hooks and executing config already in the clone

## Work
A warning when the clone already holds what the fix prevents from now on: an executable hook that is not a `.sample`, and config keys that make git run a command or send data elsewhere (`core.hooksPath`, `core.fsmonitor`, `core.sshCommand`, `core.editor`, `core.pager`, `credential.helper`, `alias.*` starting with `!`, `include.path`, `includeIf`, `filter.*`, `diff.*.textconv` and `command`, `merge.*.driver`, `url.*.insteadOf` and `pushInsteadOf`, `gpg.program`, `sequence.editor`). It is said by `flai dashboard` at start, which is where the risk is taken, and listed by `flai dashboard status`; not by `flai check`, which runs in CI where none of this exists and whose strict gate a developer's legitimate hook would fail. Tests.

## Done when
- The warning names each finding and is tested
- A clean clone says nothing

## Notes
