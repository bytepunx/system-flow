---
id: S-0285
type: story
nature: improvement
title: A story's agent waits for its sub-agents with wait_for_events, which cannot see them finish, so each wait runs to its timeout
status: ready
owner: alex
created: 2026-10-06T06:03:25Z
updated: 2026-10-06T07:01:55Z
transitions:
  - to: ready
    at: 2026-10-06T06:04:29Z
    by: alex
tags: [flai, template]
topics: [cli, conventions]
touches: [flai/internal/harness, flai/internal/mcpserver, flai/internal/guard, flai/cmd/guard.go, design/conventions/delegation.md, template, ".claude/settings.json", design/system/flai-cli.md, design/system/agent-context.md, docs/users, docs/operators, design/issues/I-0083-a-story-s-agent-waits-for-its-sub-agents-with-wait-for-events-which-cannot-see-them-finish-so-each-wait-runs-to-its-timeout.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h18m
    by: flai
    at: 2026-10-06T06:03:25Z
finalized:
  by: alex
  at: 2026-10-06T06:04:29Z
---
# S-0285 A story's agent waits for its sub-agents with wait_for_events, which cannot see them finish, so each wait runs to its timeout

## Goal

This story remediates [I-0083](../../../design/issues/I-0083-a-story-s-agent-waits-for-its-sub-agents-with-wait-for-events-which-cannot-see-them-finish-so-each-wait-runs-to-its-timeout.md), "A story's agent waits for its sub-agents with wait_for_events, which cannot see them finish, so each wait runs to its timeout". The issue recommends this solution:

Two changes, both needed. The wording T-0869 added did not hold in two stories, so better wording is not enough without a check in flai.

1. **Say how to wait.** The start prompt in `flai/internal/harness`, `delegation.md` here and in the template, and the `wait_for_events` tool description say how the agent waits for a sub-agent when it has nothing else to do, in a way that keeps the session alive however long the sub-agent runs. Ending the turn is not that way: it works while the sub-agent finishes within ten minutes, as S-0219's did, and loses the session and the sub-agent's work past that, as in S-0220 (I-0084). The candidate is a launch with `run_in_background: false`, which returns the result as the tool's result; test it against the Claude Code that flai serve runs, 2.1.290 when this was recorded, with a sub-agent that runs longer than ten minutes, and with a layer of several launched in one message. They say when `wait_for_events` is right: only for a thread awaiting the designer. The sentence no longer depends on the agent having chosen the background.
2. **Make flai catch it.** `flai guard` already runs as a `PreToolUse` hook on the story agent's `mcp__flai__` calls, so flai sees the `wait_for_events` call when it is made. Either of these keeps a mistaken wait short:
   - A Claude Code hook on a sub-agent's stop (`SubagentStop`) tells flai, and a `wait_for_events` held by that session returns at once, saying which sub-agent finished. Confirm first what that hook's input carries.
   - `flai guard` refuses a `wait_for_events` call made while a sub-agent of the session is running and no thread on the story awaits the designer, and its refusal says how to wait instead.

The measure, on a story flai serve works after the change is installed: no wait on a sub-agent outlasts the sub-agent by more than a few seconds.

## Acceptance criteria
- [ ] The story's narrative records what was tested against the Claude Code that flai serve runs: whether a launch with `run_in_background: false` returns the sub-agent's result as the tool's result when the sub-agent runs longer than ten minutes and when several are launched in one message, how long a `claude -p` session that ends its turn with a sub-agent out is kept, and what a sub-agent stop hook's input carries
- [ ] flai serve's start prompt, `delegation.md` here and in the template, and the `wait_for_events` tool description say how a story's agent waits for a sub-agent, in the way the test found to work, and say that `wait_for_events` is for a thread awaiting the designer; a test holds the prompt and the tool description to it
- [ ] flai keeps a story's agent from waiting past its sub-agent's finish with `wait_for_events`: a held call returns within seconds of the sub-agent stopping and names it, or the call is refused with what to do instead; a test covers the way chosen
- [ ] A `wait_for_events` held for a thread awaiting the designer while a sub-agent runs still returns on the thread's answer, with a test
- [ ] The design, the user and operator docs, and the template, with a template release, describe the change and any hook it adds to `.claude/settings.json`
- [ ] The way the prompt names keeps the session alive for a sub-agent that runs longer than ten minutes, so that the first remediation of I-0084 is made; I-0084 stays open for its second, which is the operator's to decide
- [ ] I-0083 is closed with `flai issue close I-0083 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0083. time_lost_per_cycle 1h18m: 39m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-06T06:02:33Z, 0 days before this story; under one cycle counts as one).

Written on the operator's word (alex, 2026-10-06), from the examination of S-0282's run; the operator asked for it at the top of ready.

The evidence is in the Claude Code transcripts of the agents of S-0218, S-0219, S-0220, and S-0282 on this host, under `~/.claude/projects/`, and their narratives. S-0219's run waited by ending its turn and lost nothing, because its sub-agents were short.

Corrected on 2026-10-06, at about 07:00Z: the goal first named ending the turn as the way to wait. S-0220 showed that Claude Code ends the process ten minutes after the turn ends, sub-agent or not (I-0084), so the goal and the first criterion now ask for a way that holds for a long sub-agent.

The story's own agent runs on the installed flai, whose prompt is the old one, so it cannot measure the fix on itself. The measure in the goal is the operator's to read on a later story.

A write under `.claude/` is refused today with no thread, which S-0283 is to fix. If this story changes `.claude/settings.json`, ask the operator on a thread to paste it, as S-0219 did with TH-0161.
