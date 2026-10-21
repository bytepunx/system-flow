---
id: S-0453
type: story
nature: remediation
title: A thread awaits an agent only when the agent takes part in it or it is on the agent's own story, so a story's agent no longer counts the threads awaiting the operator as its own
status: done
owner: alex
created: 2026-10-20T21:49:07Z
updated: 2026-10-21T07:49:34Z
transitions:
  - to: ready
    at: 2026-10-21T00:36:24Z
    by: alex
  - to: in-progress
    at: 2026-10-21T00:37:12Z
    by: agent-S-0453
  - to: review
    at: 2026-10-21T06:07:24Z
    by: agent-S-0453
  - to: in-progress
    at: 2026-10-21T07:17:03Z
    by: flai
  - to: review
    at: 2026-10-21T07:43:54Z
    by: agent-S-0453
  - to: done
    at: 2026-10-21T07:49:34Z
    by: alex
tags: [flai, template]
topics: [mcp, inbox, orchestration]
touches: [design/adrs/README.md, design/system/agent-narrative.md, design/system/flai-cli.md, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/work.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/messages.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/orchestrator_threads_test.go, flai/internal/mcpserver/story_start_test.go, flai/internal/mcpserver/task_test.go, flai/internal/storystart/start_test.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/cmd/story_start_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/strategic-agents.md, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, design/adrs/0181-a-thread-awaits-an-agent-only-when-the-agent-takes-part-in-it-or-it-is-on-the.md, design/issues/I-0174-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md, design/issues/summary.md]
after: [S-0450]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 6323
  priced_by: harness
  turns:
    - day: 2026-10-21
      ceremony: 5
      hand_edits: 2
      work: 138
  models:
    - model: claude-opus-5-5
      input: 568
      output: 169988
      cache_read: 42667934
      cache_write: 1576246
      cost: 22.1477
  strategic:
    - kind: planner
      seconds: 709
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 152
          output: 42615
          cache_read: 8116741
          cache_write: 324710
          cost: 0.5652
        - model: claude-opus-5-5
          input: 268
          output: 34097
          cache_read: 37861975
          cache_write: 383777
          cost: 11.1721
    - kind: orchestrator
      seconds: 6154
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 305
          output: 5739
          cache_read: 100706223
          cache_write: 76307
          cost: 25.2655
        - model: claude-sonnet-5-5
          input: 35
          output: 153
          cache_read: 780774
          cache_write: 122491
          cost: 0.7392
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: orchestrator
    at: 2026-10-20T22:23:30Z
  value: 50
  by: planner-S-0453
  at: 2026-10-20T22:24:31Z
forecast:
  duration: 51m
  delivery: 2026-10-21T16:15:00Z
  basis: "Its own forecast of 51m; 29th in the pull order with an in-progress limit of 5, behind S-0445, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0428, S-0429, S-0430, S-0431, S-0432, S-0433, S-0434 and S-0435."
  by: flai
  at: 2026-10-21T00:29:28Z
finalized:
  by: alex
  at: 2026-10-20T22:12:55Z
---
# S-0453 A thread awaits an agent only when the agent takes part in it or it is on the agent's own story, so a story's agent no longer counts the threads awaiting the operator as its own

## Goal

`Summarize` in `flai/internal/inbox/inbox.go` marks a thread `awaiting: "you"` whenever its last entry is by anyone but the reader. So every story's agent counts every thread awaiting the operator as its own. Measured on 2026-10-20, a new story agent's inbox was 43,279 bytes, 42,403 of them threads: 31 threads, all `awaiting: "you"`, none on the agent's story. They were planners' questions, the orchestrator's escalations, and flai serve's notices, all waiting on the operator.

S-0450 (ADR-0179) keeps `story_start` and `task_done` under Claude Code's limit by listing only the story's threads and counting the rest in `threads_omitted`. The rule underneath still stands. The standalone `inbox` is near the 50,000-character limit and grows with every thread the operator has not answered. `awaiting_you` and `threads_omitted` count threads the agent should not answer. And `wait_for_work` (`flai/internal/mcpserver/work.go`) answers `thread` to an agent with no story in progress for any new thread awaiting the operator.

This story fixes the rule. The operator chose this on 2026-10-20 over shorter summaries everywhere and a general size cap. ADR-0179's 40,000-byte fit stays as the safety net.

## Acceptance criteria
- [x] An accepted ADR records the rule. A thread awaits an agent only when the agent is among its participants, or the thread is on the agent's own story or one of its tasks, and its last entry is not the agent's. ADR-0090's rule for an opener whose recommendation awaits confirmation stands. Any other open thread is `awaiting: "operator"` when it awaits the operator, else `"other"`. The ADR refines ADR-0179 and ADR-0090.
- [x] `inbox`, `wait_for_work`, `story_start`, and `task_done` apply the rule. A test reproduces the 2026-10-20 measurement: open threads awaiting the operator, none on the reading agent's story. The story's agent then gets no thread listed, `awaiting_you` 0, `threads_omitted` 0, and no `thread` answer from `wait_for_work`, while a thread the operator answered on the agent's story still awaits it.
- [x] The orchestrator still finds the threads awaiting the operator in `inbox` without `all: true`, as its role's default, and its instructions in `flai/internal/harness/harness.go` name them by `awaiting: "operator"`.
- [x] The operator's view is unchanged: the dashboard and `flai thread list` show the same threads awaiting the operator as before.

## Tasks
- T-1927 Record in an ADR and the design that a thread awaits an agent only when the agent takes part in it or it is on the agent's own story, and is awaiting operator or other otherwise
- T-1929 Make the inbox package's Summarize and Read apply the rule: awaiting you only on the agent's own story or a thread it takes part in, else operator or other
- T-1930 Make the MCP inbox, wait_for_work, the folder inbox, and the thread tools apply the rule, with the orchestrator listing the threads awaiting the operator by default
- T-1932 Make story_start and task_done, MCP and CLI, carry only the threads awaiting the story's agent by the rule, with tests of the 2026-10-20 measurement
- T-1933 Name the threads the orchestrator answers by awaiting operator in its instructions, its agent definition, and the strategic-agents convention
- T-1935 Document in the user docs and the generated reference what inbox and wait_for_work report as awaiting you, operator, and other

## Notes

Decided with the operator on 2026-10-20: the orchestrator sees the threads awaiting the operator by default, as its role, rather than by asking for `all: true`, so that its instructions barely change.

Open for the planner:
- How flai knows the operator's name, to tell "awaiting the operator" from "awaiting another agent". The threads' operator entries are by `alex`, the story's `owner`.
- Whether a reply adds its author to a thread's `participants`.

Seen alongside, and not this story's: 17 of the 31 open threads were flai serve's "host flai is older than …" notices, each replaced by a newer one and never resolved, 18 KB of the 42 KB. Resolving a replaced notice when a newer one opens would be a separate story.

### Planning

Answers to the two open questions, taken into T-1927's ADR:

- The operator is the manifest's `owner` or the thread's story's `owner`, the names ADR-0097 takes an answer from. A thread awaits the operator when a recommendation is pending or its last entry is by neither. `inbox.designer` in `flai/internal/hostapi/people.go` already applies this rule, so the dashboard keeps its view (criterion 4).
- A reply already adds its author to `participants`: `threads.ReplyWith`, `Confirm`, and `Resolve` do. No change is needed.

Touches, all files; no folder touch is kept. The story declared none, so `flai touches suggest` ran from four seeds the explorer found.

| Touch | Source |
|-------|--------|
| `flai/internal/inbox/inbox.go`, `flai/internal/mcpserver/work.go`, `flai/internal/mcpserver/server.go`, `flai/internal/harness/harness.go` | design: the goal and criteria name them; the seeds |
| `flai/internal/harness/harness_test.go`, `flai/internal/mcpserver/folder.go`, `flai/internal/mcpserver/server_test.go`, `flai/internal/mcpserver/folder_test.go`, `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`, `template/CHANGELOG.md` | co-change with the seeds |
| `flai/internal/inbox/inbox_test.go`, `flai/internal/mcpserver/messages.go`, `flai/internal/mcpserver/work_test.go`, `flai/internal/mcpserver/orchestrator_threads_test.go` | layout: the other callers of `Summarize` and the tests of the rule |
| `flai/internal/mcpserver/story_start_test.go`, `flai/internal/mcpserver/task_test.go`, `flai/internal/storystart/start_test.go`, `flai/internal/taskdone/taskdone_test.go`, `flai/cmd/task_done.go`, `flai/cmd/task_done_test.go`, `flai/cmd/story_start_test.go` | layout: criterion 2's `story_start` and `task_done`, and the printed inbox |
| `design/adrs/README.md`, `design/system/agent-narrative.md` | design: criterion 1's ADR and § What an agent is told through MCP |
| `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md`, `design/system/strategic-agents.md`, `.claude/agents/orchestrator.md`, `template/root/.claude/agents/orchestrator.md` | layout: every place that tells the orchestrator which threads await the operator (criterion 3) |

Left out of the co-change list: `flai/internal/serve/agents.go`, whose `awaitsAnswer` is its own rule for starting agents again; `flai/cmd/mcp.go`, `mcp_http.go`, and `cursor.go`, which do not read `awaiting`; `work-management.md` and the settings docs, which name no thread rule. `flai/internal/storystart/start.go` and `flai/internal/taskdone/taskdone.go` are left out too: they already call `inbox.Read` with the story as `Own`. The ADR's file is named by `flai adr new`; `flai task done` widens the touches with it.

Forecast 51m, delivery 2026-10-21T14:04:00Z, as `flai forecast` gives them after the touches were set (9m before, at size 4 with no touches). Kept: the median comes from 30 done large remediation stories on the same model, and size 34 counts the six tasks' 30 files. S-0450, of the same shape, was forecast 37m at size 25.

Cost of delay 50 USD a week, as `flai cod` gives it from 20m lost per 168h cycle at 150 USD an hour. The input was missing; TH-0536 asked for it, and the orchestrator set the recommended 20m. Kept: S-0450 has 12m for the part of this cause it fixes, and what is left reaches every agent's standalone `inbox` and `wait_for_work`, more sessions than I-0167 counted.
- 2026-10-21T07:17:03Z: moved to in-progress: merge queue batch 20261021-071126 fails with S-0453 at 28b5c7d30c60: smoke failed: smoke: smoke: template render and check
