---
id: T-1712
type: task
nature: feature
title: An ADR records whether an OpenCode agent's status comes from opencode run --format json or from opencode serve's event stream, and what each costs on a busy host
status: backlog
parent: S-0411
owner: alex
created: 2026-10-16T05:19:36Z
updated: 2026-10-16T05:19:36Z
transitions: []
stream: S-0411
tags: [cli]
touches: [design/adrs]
---
# T-1712 An ADR records whether an OpenCode agent's status comes from opencode run --format json or from opencode serve's event stream, and what each costs on a busy host

## Work

- From S-0406's finding on `--attach`, write the ADR with `flai adr new`, refining ADR-0130: the status of an OpenCode agent comes from `opencode run --format json`, as Claude Code's does from stream-json, or from a running `opencode serve`'s event stream with each story's run attached to it.
- Weigh what each costs with several stories on one host: MCP start-up per run, one server's life cycle under `flai serve`, and what the dashboard can show from each.
- Waits for nothing in this story.

## Done when

The ADR is accepted in `design/adrs`, listed in its index, and names the source, the cost on a busy host, and the alternative.

## Notes
