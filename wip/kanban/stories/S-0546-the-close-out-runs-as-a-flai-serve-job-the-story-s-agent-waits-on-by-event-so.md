---
id: S-0546
type: story
nature: improvement
title: The close-out runs as a flai serve job the story's agent waits on by event, so it never dies at the harness's tool timeout and is never re-run from the start
status: ready
parent: E-0031
owner: alex
created: 2026-10-30T01:51:56Z
updated: 2026-10-31T08:20:13Z
transitions:
  - to: ready
    at: 2026-10-31T07:24:17Z
    by: alex
tags: [flai]
topics: [testing, mcp, cli]
touches: [flai/internal/serve, flai/internal/mcpserver, flai/cmd, scripts/close-out.sh, design/system/flai-cli.md, design/conventions/testing.md, docs/users/flai.md, design/adrs/0224-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-waits-on-by-event.md, design/adrs/README.md, design/system/metrics.md, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/serve.go, flai/internal/inbox/events.go, flai/internal/inbox/events_test.go, flai/internal/hostapi/writes.go, flai/cmd/serve_actions.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, flai/internal/mcpserver/server.go, flai/internal/metrics/closeouts.go, flai/internal/metrics/closeouts_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/cmd/stats.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, template/root/scripts/close-out.sh, docs/users/flai-reference.md, docs/operators/index.md, design/conventions/code-quality.md, design/conventions/delegation.md, design/conventions/work-management.md, template/root/design/conventions/code-quality.md, template/root/design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 130
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 42
          output: 300
          cache_read: 1358341
          cache_write: 142271
          cost: 0.3395
forecast:
  duration: 2h30m
  delivery: 2026-10-31T11:43:00Z
  basis: "Its own forecast of 2h30m; 1st in the pull order with an in-progress limit of 5, behind S-0548, S-0550 and S-0552."
  by: flai
  at: 2026-10-31T08:20:13Z
---
# S-0546 The close-out runs as a flai serve job the story's agent waits on by event, so it never dies at the harness's tool timeout and is never re-run from the start

## Goal

`scripts/close-out.sh` runs `flai verify` inside one Bash tool call of the agent's turn. Claude Code ends a tool call at 600 s, and 64 of the fortnight's 267 close-outs died there, 146 ended with no summary line, and only 88 reported every step passing; the agent then ran the whole close-out again from the start, every tier included, p50 109 s and p85 the full 600 s, for 18.5 hours of the fortnight's 113 session-hours. This story makes the close-out a job: the agent asks flai serve, through the `verify` MCP tool or `flai verify --job`, to run the close-out in the story's worktree, gets a job id at once, and waits for its `verified` event with `wait_for_events`, so no tool call outlives its timeout and a close-out is never cut short. The job's tiers, output, and verdict are recorded where `flai verify` records them today, and a close-out that is asked for again at the same commit answers from that record instead of running.

## Acceptance criteria
- [ ] `verify` (MCP) and `flai verify --job` start the close-out as a serve job in the story's worktree and return its id without waiting; `wait_for_events` reports `verified` with the job's verdict when it ends, and `flai verify` without `--job` behaves as today for a session without serve
- [ ] A close-out asked for again at a commit the record already holds a verdict for answers from the record, saying so, and runs nothing
- [ ] `scripts/close-out.sh`, the testing convention, `flai-cli.md`, and the user docs say the close-out is a job the agent waits on, and the story agent's prompt tells it to wait rather than call the script in a turn
- [ ] The serve journal records each job's start, end, duration, and verdict, and `flai stats` counts close-outs from it

## Tasks

The planner drafts them.
- T-2396 An ADR settles the close-out job: how the CLI and MCP reach flai serve, what the job runs, the verified event, the per-commit verdict, and the journal's fields, and metrics.md defines the close-out count
- T-2400 The verify record keeps a verdict per commit, and a close-out asked for again at a commit it holds answers from it, saying so, and runs nothing
- T-2403 flai serve takes a close-out request, runs the close-out in the story's worktree as a job, appends a verified event with its verdict, and journals its start, end, duration, and verdict
- T-2407 flai verify --job and the MCP verify tool ask flai serve for the close-out job and return its id at once, and wait_for_events lists verified among its kinds
- T-2411 flai stats counts close-outs from the serve journal, by verdict, with their durations, as metrics.md defines
- T-2415 The story agent's prompt asks for the close-out job and waits for verified, and scripts/close-out.sh and the template's copy say they run as the job, not in a turn
- T-2417 The conventions, their template copies, the user and operator docs, and the changelog say the close-out is a serve job the agent waits on

## Notes

Evidence: the fortnight's story agent logs, counted in the operator's conversation on 2026-10-29. The harness timeout is the harness's, not a setting flai can raise, so the fix is to leave the turn.

### Planning

Touches, by source:

- Declared, kept: `flai/internal/serve`, `flai/internal/mcpserver`, `flai/cmd`, `scripts/close-out.sh`, `design/system/flai-cli.md`, `design/conventions/testing.md`, `docs/users/flai.md`. The three folders stay because the story declared them; the tasks name the files in each, and flai narrows each folder in the claim to those files (ADR-0096). `design/conventions/testing.md` does not exist: the close-out's convention text is in `code-quality.md`, `delegation.md`, and `work-management.md`, which the layout adds; the plan's thread asks whether criterion 3's "testing convention" means those.
- Design: the ADR the job needs (predicted number 0224; the agent renames it to the number `flai adr new` gives) and `design/adrs/README.md`; `design/system/metrics.md`, since counting close-outs in `flai stats` is a metrics change and needs the ADR.
- Layout, from the explorer's read of the code: the verify record (`flai/internal/verify/story.go`, keyed by story and overwritten each run, so a per-commit verdict is new); the serve job beside the merge hand-off (`flai/internal/serve/verifyjob.go`, started from `serve.go` as `newMergePass` is); the event kind (`flai/internal/inbox/events.go`); the journal entry (`flai/internal/hostapi/writes.go`'s `Entry`, written by `flai/cmd/serve_actions.go`); the CLI and MCP entry points (`flai/cmd/verify.go`, `flai/internal/mcpserver/verify.go`, `server.go`'s event-kind description, `flai/cmd/mcp.go` and `mcp_http.go` wiring); the stats (`flai/internal/statsread/statsread.go`, `flai/internal/metrics`, `flai/cmd/stats.go`); the story agent's prompt (`flai/internal/harness/harness.go`'s close-out line); the template's copies of the script and the conventions, and `template/CHANGELOG.md`.
- Co-change: `docs/users/flai-reference.md` (22 percent, regenerated from the verify help) and `docs/operators/index.md` (13 percent, the journal's fields) from `flai touches suggest`; the rest of its list is hub files this story has no reason to change.
- Shared with others: S-0552's split on MS-0261 leaves S-0546 `flai/cmd/verify.go` and the job's new files; `flai/cmd/stats.go` is in S-0552's half, and this story needs one line there, so whichever lands second rebases over the other.

Figures:

- Forecast: `flai forecast` gave 1h12m, which is 89 s per unit of size over 108 done large improvement stories, times size 48 (4 criteria, 44 touches). Raised to 2h30m: seven tasks in four layers, including a new concurrent serve job behind an ADR. Also, each task's tests and the close-out run the serve and mcpserver packages, which take 8 to 10 minutes under `-race` (I-0225) until S-0548 lands. flai's replan keeps the 2h30m and sets the delivery from the pull order.
- Cost of delay: no value yet. Neither S-0546 nor E-0031 has inputs, and they are the operator's. TH-0773 asks for them, recommending `time_lost_per_cycle: 9h`, from the goal's 18.5 lost hours a fortnight.

Plan: T-2396 first; T-2400 and T-2403 together; T-2407 and T-2411 together; T-2415 and T-2417 together. Each task's `## Work` says what it waits for and why.
