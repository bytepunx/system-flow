---
id: T-1695
type: task
nature: feature
title: An ADR records how an OpenCode agent asks before a protected write, OpenCode's bound for a held call, and the way not taken
status: done
parent: S-0408
owner: alex
created: 2026-10-16T05:16:20Z
updated: 2026-11-07T21:50:20Z
transitions:
  - to: ready
    at: 2026-11-07T21:48:54Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T21:48:56Z
    by: system-flow
  - to: done
    at: 2026-11-07T21:50:20Z
    by: agent-S-0408
stream: S-0408
tags: [cli]
touches: [design/adrs]
usage:
  source: log
  seconds: 84
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 3701
      cache_read: 1476940
      cache_write: 60687
      cost: 0.7279
---
# T-1695 An ADR records how an OpenCode agent asks before a protected write, OpenCode's bound for a held call, and the way not taken

## Work

- From S-0406's finding on answering an `ask` headless, write the ADR with `flai adr new`, refining ADR-0124 and ADR-0130: hold-and-ask through the plugin's `permission.asked` and `flai guard`, or `--auto` behind `deny` rules with the agent told to ask on a thread.
- Name OpenCode's bound for a held call, below whatever timeout OpenCode puts on a plugin or a permission, as ADR-0124 names Claude Code's four minutes.
- Waits for nothing in this story.

## Done when

The ADR is accepted in `design/adrs`, listed in its index, and names the way chosen, the bound, and the alternative.

## Notes
