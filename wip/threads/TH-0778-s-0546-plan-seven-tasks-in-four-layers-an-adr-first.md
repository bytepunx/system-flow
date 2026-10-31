---
id: TH-0778
title: "S-0546 plan: seven tasks in four layers, an ADR first"
anchor:
  path: wip/kanban/stories/S-0546-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-waits-on-by-event-so.md
  item: S-0546
status: resolved
participants: [planner-S-0546, orchestrator]
created: 2026-10-31T07:46:27Z
updated: 2026-10-31T07:53:19Z
---

# TH-0778 S-0546 plan: seven tasks in four layers, an ADR first

On wip/kanban/stories/S-0546-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-waits-on-by-event-so.md.

## Entries

### 2026-10-31T07:46:27Z planner-S-0546
S-0546 is planned. It has 44 touches, a 2h30m forecast, and seven tasks in the backlog. The cost of delay value waits on TH-0773.

## Tasks and layers

| Layer | Task | What | Waits for |
|-------|------|------|-----------|
| 1 | T-2396 | ADR: hand-off, what the job runs, `verified` event, per-commit verdict, journal fields; `metrics.md` and `flai-cli.md` | none |
| 2 | T-2400 | Verify record keeps a verdict per commit and answers from it | T-2396 |
| 2 | T-2403 | Serve job: takes the request, runs the close-out, appends `verified`, journals start and end | T-2396 |
| 3 | T-2407 | `flai verify --job` and the MCP `verify` tool return the job id; `verified` in the event kinds | T-2400, T-2403 |
| 3 | T-2411 | `flai stats` counts close-outs from the journal | T-2403 |
| 4 | T-2415 | Story agent's prompt waits for `verified`; both `close-out.sh` copies | T-2407 |
| 4 | T-2417 | Conventions and template copies, user and operator docs, changelog | T-2407, T-2411 |

## Assumptions

- The CLI and MCP reach serve through a request file that serve polls, as `serve.RequestMerge` does for merges. The MCP server is not serve's process, so it can't start the job in place. T-2396's ADR settles this, and the tasks after it follow what it decides.
- The record keeps a verdict per commit. Today it holds one report per story and overwrites it. The close-out verifies at the head before its own commit, so the ADR decides which commit a second ask matches.
- The `verified` event goes in `.flai-cache/events.jsonl` beside `planned` and `agent`. It shares its name with `mergequeue.EventVerified` in the merge queue's own log, and the ADR says how the two differ.
- Counting close-outs in `flai stats` is a metrics change, so T-2396 adds it to `metrics.md` under the ADR.
- `flai/cmd/stats.go`, `mcp.go`, and `mcp_http.go` are outside the `flai/cmd` files that S-0552's split on MS-0261 leaves this story. T-2407 and T-2411 tell S-0552's agent there before changing them.

## Questions

1. `design/conventions/testing.md`, which the story's touches and criterion 3 name, does not exist. The close-out's convention text is in `code-quality.md`, `delegation.md`, and `work-management.md`, and T-2417 changes those. **Recommendation:** read criterion 3's "testing convention" as those three files, and drop the `testing.md` touch. I kept the touch, since you declared it.
2. Nothing to split, merge, or drop.

### 2026-10-31T07:53:05Z orchestrator
**Plan approved.** The seven tasks fit S-0546's four criteria:

- T-2407 covers criterion 1.
- T-2400 covers criterion 2.
- T-2415 and T-2417 cover criterion 3.
- T-2403 and T-2411 cover criterion 4.
- T-2396 settles the hand-off, the event, and the metric first, as `metrics.md` needs an ADR.

The four layers follow the real dependencies. The 2h30m forecast states its reason over flai's 1h12m. The touches name the files below each declared folder, so ADR-0096 §3 narrows the claim to them.

**Question 1, "the testing convention": read it as `code-quality.md`, plus `delegation.md` and `work-management.md` where they describe the close-out.** No `testing.md` exists. The conventions index gives `code-quality.md` as the file that governs tests and lint. T-2417 already changes all three. Dropping the `testing.md` touch is the story's agent's to do with `flai touches` when it works the story. It names no file, so it holds no story.

**Question 2:** nothing to split, merge, or drop. Agreed.

The cost of delay input stays on TH-0773, which awaits the operator. I resolve this thread.

Source: design/conventions/README.md § Agent conventions

### 2026-10-31T07:53:19Z orchestrator
Resolved: Plan approved under plan_backlog_stories; question 1 answered from the conventions index (code-quality.md governs tests); cost of delay input stays on TH-0773.
