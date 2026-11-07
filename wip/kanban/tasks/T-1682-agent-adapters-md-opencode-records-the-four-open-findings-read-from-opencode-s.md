---
id: T-1682
type: task
nature: research
title: "agent-adapters.md § OpenCode records the four open findings, read from OpenCode's source and documentation at a pinned release"
status: in-progress
parent: S-0406
owner: alex
created: 2026-10-16T05:13:59Z
updated: 2026-11-07T09:40:11Z
transitions:
  - to: ready
    at: 2026-11-07T09:40:10Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T09:40:11Z
    by: system-flow
stream: S-0406
tags: [cli]
touches: [design/system/agent-adapters.md]
usage:
  source: log
  seconds: 1638
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 172
      output: 985
      cache_read: 10404179
      cache_write: 204379
      cost: 3.8035
---
# T-1682 agent-adapters.md § OpenCode records the four open findings, read from OpenCode's source and documentation at a pinned release

## Work

- Choose one OpenCode release to pin, the newest tagged one when the task starts, and name it in `design/system/agent-adapters.md` § OpenCode as the release the adapter and its fixtures follow. Read its public source and documentation; install and run nothing (TH-0532, D-0506).
- Record under § OpenCode, for that release:
  - whether a plugin's `tool.execute.before` input can tell a sub-agent's call (one the Task tool started) from the primary agent's, and by which field
  - whether an `ask` in headless `opencode run` can be answered from outside, through a plugin's `permission.asked` or the server API, and within what bound
  - the fields `step_finish` carries in `--format json`: tokens by kind, `cost`, the model, and any session or generation ID a gateway's spend log could join on
  - whether `--attach` to a running `opencode serve` is the better way to run several stories on one host, against one `opencode run` per story
- Name for each the source file or documentation page it was read from, by link to the release's tag, not to a line. Say that each was read, not run against a binary, and replace each **to check** the finding resolves.
- Waits for nothing: it needs only the public source.

## Done when

Each of the four findings is recorded with the pinned release and the source it was read from, and the section has no **to check** left on them.

## Notes
