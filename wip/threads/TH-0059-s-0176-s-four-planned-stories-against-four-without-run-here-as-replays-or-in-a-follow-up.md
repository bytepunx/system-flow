---
id: TH-0059
title: "S-0176's four planned stories against four without: run here as replays, or in a follow-up?"
anchor:
  path: wip/kanban/stories/S-0176-the-story-s-agent-plans-which-tasks-can-run-in-parallel-and-works-them-with-sub-agents.md
  item: S-0176
status: resolved
participants: [agent-S-0176, alex]
created: 2026-10-01T11:40:21Z
updated: 2026-10-01T13:20:17Z
---

# TH-0059 S-0176's four planned stories against four without: run here as replays, or in a follow-up?

On wip/kanban/stories/S-0176-the-story-s-agent-plans-which-tasks-can-run-in-parallel-and-works-them-with-sub-agents.md.

## Entries

### 2026-10-01T11:40:21Z agent-S-0176
The sixth criterion asks for at least four stories run with the plan and four comparable ones without. A story gets the plan only through the prompt and conventions of the flai that `flai serve` runs, which is the installed one. S-0176 is an `experiment`: per git.md it stays on its branch and acceptance refuses it, so no story `flai serve` starts will run with the plan while this story is open.

Recommended, (a): measure it here. This story's own run is the first planned run. Three more are replays: archived stories with three or more tasks, at least two of them independent, re-run with `claude -p` (same model and effort) from their base commit in a scratch worktree, with this branch's flai, prompt, and conventions. Their original logged runs are the four without the plan, comparable by construction (same story, same start). Review defects are compared as a fresh verifier's findings on each diff. Each replay calls a paid service and should cost about what its original did, roughly 5 to 12 US dollars by the runs agent-context.md measures, so about 20 to 40 dollars in all. I will list the three stories and their original cost in this thread before starting any.

(b): split, as S-0188 followed S-0175 (TH-0042). This story builds everything, measures its own run and one replay, and records a provisional recommendation; a follow-up story measures four stories that `flai serve` runs once you install this branch's flai on the host.

I am building the task `after` field, the plan's display, the prompt and conventions, and the worktree evidence meanwhile; none of that waits on this answer.

### 2026-10-01T11:44:00Z alex
a

### 2026-10-01T11:52:32Z agent-S-0176
Going with (a). The three replays, each with three or more tasks of which at least two touch nothing in common, all run by `flai serve` on 2026-10-01 with the claude-code prompt that delegates (S-0175):

| Story | Tasks (independent) | Original time to review | Original cost (USD) |
|-------|---------------------|-------------------------|---------------------|
| S-0186 | 3 (3: flaiover page, API, docs) | 11m 21s | 6.04 |
| S-0185 | 3 (2: T-0645 and T-0646 share flai/cmd/new.go) | 22m 26s | 7.83 |
| S-0182 | 4 (3: serve/agents, serve/start and harness, docs; T-0652 investigates first) | 47m 18s | 11.13 |

Expected replay spend: about 25 USD, the originals' sum. The fourth planned run is this story; its unplanned comparison is S-0189 (6 tasks, improvement, flai and template, 25.82 USD, 53m 30s).

Each replay runs in a scratch clone at the story's base commit, with this branch's flai, prompt, and conventions, and the story put back to `ready` there, so nothing it does reaches this project's board or main. One confound I will record rather than remove: the replays also carry S-0189's prompt (cheaper sub-agents, one verifier run), which the originals did not. I start the replays once layer 2 of the plan is in, so they run with the plan's display too.

### 2026-10-01T11:56:42Z alex
ok

### 2026-10-01T13:20:17Z agent-S-0176
Resolved: Option (a): S-0186, S-0185, and S-0182 replayed with the plan in scratch clones against their originals; results and the adapt recommendation in design/system/agent-context.md § Tasks in parallel.
