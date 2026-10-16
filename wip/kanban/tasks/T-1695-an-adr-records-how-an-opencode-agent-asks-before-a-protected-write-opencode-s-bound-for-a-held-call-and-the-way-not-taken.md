---
id: T-1695
type: task
nature: feature
title: An ADR records how an OpenCode agent asks before a protected write, OpenCode's bound for a held call, and the way not taken
status: backlog
parent: S-0408
owner: alex
created: 2026-10-16T05:16:20Z
updated: 2026-10-16T05:16:20Z
transitions: []
stream: S-0408
tags: [cli]
touches: [design/adrs]
---
# T-1695 An ADR records how an OpenCode agent asks before a protected write, OpenCode's bound for a held call, and the way not taken

## Work

- From S-0406's finding on answering an `ask` headless, write the ADR with `flai adr new`, refining ADR-0124 and ADR-0130: hold-and-ask through the plugin's `permission.asked` and `flai guard`, or `--auto` behind `deny` rules with the agent told to ask on a thread.
- Name OpenCode's bound for a held call, below whatever timeout OpenCode puts on a plugin or a permission, as ADR-0124 names Claude Code's four minutes.
- Waits for nothing in this story.

## Done when

The ADR is accepted in `design/adrs`, listed in its index, and names the way chosen, the bound, and the alternative.

## Notes
