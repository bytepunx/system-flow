---
id: E-0017
type: epic
nature: improvement
title: Story-loop work whose outcome the repository determines moves from the agent into flai, reached from the CLI, the host channel, and MCP
status: in-progress
owner: alex
created: 2026-10-05T01:35:27Z
updated: 2026-10-06T23:51:15Z
transitions:
  - to: ready
    at: 2026-10-06T22:47:57Z
    by: alex
  - to: in-progress
    at: 2026-10-06T23:51:15Z
    by: agent-S-0273
tags: [cli, mcp]
topics: [automation, mcp, hostapi]
touches: [flai, design/conventions, template]
usage:
  source: sum
  seconds: 6820
  models:
    - model: claude-opus-5-5
      input: 1720
      output: 717287
      cache_read: 94608418
      cache_write: 2721498
      cost: 49.0364
    - model: claude-sonnet-5-5
      input: 44
      output: 10122
      cache_read: 470318
      cache_write: 110688
      cost: 0.4721
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6h
    by: planner-E-0017
    at: 2026-10-06T11:34:11Z
---
# E-0017 Story-loop work whose outcome the repository determines moves from the agent into flai, reached from the CLI, the host channel, and MCP

## Outcome

A model turn carries a judgement. Work whose outcome the repository's state fully determines, closing a task, verifying a branch, ticking a criterion, recording narrative state, opening a story, waiting for an answer, is done by flai in one call that answers structured data the agent acts on, never a log it reads. Each operation is reached the same way from the shell (`flai`), the host channel (`flai serve` to the dashboard), and MCP (`flai mcp`), as acceptance, publishing, and `flai stream sync` already are.

Evidence, from the 108 story runs logged on this host to 2026-10-04 (about 43 hours of agent time, $778, a mean of 80 model turns per story, 11.9 million cache-read tokens per story): each task transition is a run of six to nine single-purpose turns (tick criteria with sed, commit, sync, move, log, touches, check, inbox), about 2,900 flai ceremony turns in all; tests and lint were run by hand 1,449 times, each returning a log for the model to read; the Sonnet verifier cost 88 minutes over 22 stories and in S-0248 ran the close-out three times; 650 turns woke from `wait_for_events` with nothing to do; narratives, task bodies, and criteria were edited by hand 253 times. Turns whose every call is pure ceremony are 9% of all turns, about $70; the mechanical work around them is the larger share.

This epic is the story agent's counterpart of S-0217 under E-0016, which exposes the orchestrator's deterministic operations as commands: the same principle, applied to the loop a story agent runs.

Done when a story agent's run shows no turn that only commits, syncs, moves, logs, ticks, or reads a test log, and the conventions and the harness prompt send it to the commands instead.

## Stories

Drafted with the epic; finalize and order them on the board. Not in the epic because stories exist: S-0261 (prime pack size), S-0266 (verifier exit status), S-0267 (duplicate test tier), S-0268 (prompt: batching, stash), S-0249 (strict check scoped to the story).
- S-0269 One command closes a task: flai task done commits, syncs, moves, logs, widens touches, checks, and answers the inbox
- S-0270 Verification is a module: flai verify runs the tiers the diff selects and answers structured findings, replacing the verifier sub-agent's run
- S-0271 Criteria and narrative state are commands: flai story tick checks a criterion and flai stream state writes Current state and Next steps
- S-0272 An agent with an open question ends instead of waiting: flai serve restarts it on the answer, and wait_for_events keeps a timeout only for an agent with work in hand
- S-0273 flai test runs the project's test and lint tiers for a path or a package and answers pass or the first failures as findings
- S-0274 Opening a story is one call: flai story start moves it to in-progress, opens the stream, primes, and answers the first inbox together
- S-0275 Issue and ADR handling from the worktree is one call each: bump, close, and adr new number from the whole repository and commit on the story branch
- S-0293 flai stats classifies a story run's tool calls, so the ceremony turns E-0017 removes are measured per story and over time

## Notes

Found by the operator's review of the S-0248 agent log and a classification of every main-agent tool call in the logs under `~/.flai/serve/agents/` on 2026-10-04.
- 2026-10-06T22:47:57Z: moved to ready: follows S-0269, which moved to ready
- 2026-10-06T23:51:15Z: moved to in-progress: follows S-0273, which moved to in-progress
