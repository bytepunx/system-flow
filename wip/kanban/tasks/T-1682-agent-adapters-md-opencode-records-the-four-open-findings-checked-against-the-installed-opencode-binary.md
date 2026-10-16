---
id: T-1682
type: task
nature: research
title: "agent-adapters.md § OpenCode records the four open findings, checked against the installed opencode binary"
status: backlog
parent: S-0406
owner: alex
created: 2026-10-16T05:13:59Z
updated: 2026-10-16T05:13:59Z
transitions: []
stream: S-0406
tags: [cli]
touches: [design/system/agent-adapters.md]
---
# T-1682 agent-adapters.md § OpenCode records the four open findings, checked against the installed opencode binary

## Work

- With the `opencode` the operator installed on the host, check and record under `design/system/agent-adapters.md` § OpenCode, for that version:
  - whether a plugin's `tool.execute.before` input can tell a sub-agent's call (one the Task tool started) from the primary agent's, and by which field
  - whether an `ask` in headless `opencode run` can be answered from outside, through a plugin's `permission.asked` or the server API, and within what bound
  - the fields `step_finish` carries in `--format json`: tokens by kind, `cost`, the model, and any session or generation ID a gateway's spend log could join on
  - whether `--attach` to a running `opencode serve` is the better way to run several stories on one host, against one `opencode run` per story
- Say how each was checked: the command, a short plugin, or the request. Replace each **to check** the finding resolves.
- Waits for nothing: it needs only the binary.

## Done when

Each of the four findings is recorded with the version and how it was checked, and the section has no **to check** left on them.

## Notes
