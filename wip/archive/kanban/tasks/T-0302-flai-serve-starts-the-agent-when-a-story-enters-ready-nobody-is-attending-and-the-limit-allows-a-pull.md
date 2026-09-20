---
id: T-0302
type: task
nature: feature
title: flai serve starts the agent when a story enters ready, nobody is attending, and the limit allows a pull
status: done
parent: S-0079
owner: alex
created: 2026-09-20T15:47:55Z
updated: 2026-09-20T15:53:48Z
transitions:
  - to: ready
    at: 2026-09-20T15:53:48Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T15:53:48Z
    by: system-flow
  - to: done
    at: 2026-09-20T15:53:48Z
    by: system-flow
stream: S-0079
tags: []
---
# T-0302 flai serve starts the agent when a story enters ready, nobody is attending, and the limit allows a pull

## Work
On a change under wip, flai serve works out which stories entered ready. It starts the command when the action is enabled, a command is set, the in-progress limit leaves room, no agent it started for the project is still running, and nobody is attending: no MCP cursor and no narrative written in the last minutes. One at a time per project; when its agent ends and ready stories remain, the next is started. Environment: FLAI_AGENT, FLAI_STORY, FLAI_SESSION; run in the project's directory, detached, output to a log. Starts and failures are journalled and kept in a state file.

## Done when
- Tests with a stub command: started once for a story entering ready, not when attended, not over the limit, not twice, the next after the first ends, a command that cannot start reported
- The rule for attending is written in the design
- The Go tests and lint pass

## Notes
