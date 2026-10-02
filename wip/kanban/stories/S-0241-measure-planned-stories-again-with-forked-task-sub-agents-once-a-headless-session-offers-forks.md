---
id: S-0241
type: story
nature: experiment
title: Measure planned stories again with forked task sub-agents once a headless session offers forks
status: backlog
owner: arobson
created: 2026-10-02T17:14:07Z
updated: 2026-10-02T17:14:07Z
transitions: []
tags: [template]
topics: [conventions]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks

## Goal

S-0176 measured a story's agent handing its tasks to task sub-agents, and found that each sub-agent started fresh and read the code again, because the headless `claude -p` session flai serve starts refused `subagent_type: fork`. A fork inherits the parent's conversation and prompt cache, so it would start primed. S-0230 probed again on 2026-10-02 (Claude Code 2.1.286) and forks were still refused (TH-0069). When a headless session offers forks, repeat S-0176's measurement with forked task sub-agents and record whether the sub-agents' cost falls.

## Acceptance criteria

- [ ] A headless session started as flai serve starts one is shown to offer `subagent_type: fork`, with the Claude Code version that does
- [ ] The measurement S-0176 made (`design/system/agent-context.md` § Tasks in parallel) is repeated with forked task sub-agents, and `design/system/agent-context.md` records the result beside the first
- [ ] The experiment's results document is written under `design/experiments`, with a recommendation

## Tasks

## Notes

Split from S-0230's fourth criterion on TH-0069. Do not move it to ready until a headless session offers forks: it cannot be worked before then.
