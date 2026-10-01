---
id: T-0671
type: task
nature: improvement
title: The explorer and verifier run on a cheaper model
status: done
parent: S-0189
owner: arobson
created: 2026-10-01T10:44:01Z
updated: 2026-10-01T10:51:58Z
transitions:
  - to: ready
    at: 2026-10-01T10:44:11Z
    by: agent-S-0189
  - to: in-progress
    at: 2026-10-01T10:44:11Z
    by: agent-S-0189
  - to: done
    at: 2026-10-01T10:51:58Z
    by: agent-S-0189
blocked:
  - from: 2026-10-01T10:45:06Z
    until: 2026-10-01T10:51:58Z
    reason: "Claude Code refuses this session's edits to .claude/agents/; asked the designer to make them (TH-0057)"
stream: S-0189
tags: []
touches: [template/root/.claude/agents/, ".claude/agents/", design/system/agent-context.md]
usage:
  source: log
  seconds: 438
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 116
      output: 23789
      cache_read: 6283438
      cache_write: 116356
      cost: 2.7807
---
# T-0671 The explorer and verifier run on a cheaper model

## Work

Set `model` in the template's `.claude/agents/explorer.md` and `verifier.md` to a model cheaper than the story's agent's, copy them to `.claude/agents/`, and say in `design/system/agent-context.md § Sub-agents` which model each runs and why.

## Done when

- Both definitions, in `template/root/.claude/agents/` and `.claude/agents/`, name the model, and the copies are identical.
- `agent-context.md` says which and why, and no longer says they name no model.

## Notes
