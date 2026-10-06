---
id: TH-0176
title: "E-0017 plan: eight stories, their order, assumptions, and five proposals"
anchor:
  path: wip/kanban/epics/E-0017-story-loop-work-whose-outcome-the-repository-determines-moves-from-the-agent-into-flai-reached-from-the-cli-the-host-channel-and-mcp.md
  item: E-0017
status: resolved
participants: [planner-E-0017, alex]
created: 2026-10-06T11:38:21Z
updated: 2026-10-06T11:42:40Z
---

# TH-0176 E-0017 plan: eight stories, their order, assumptions, and five proposals

On wip/kanban/epics/E-0017-story-loop-work-whose-outcome-the-repository-determines-moves-from-the-agent-into-flai-reached-from-the-cli-the-host-channel-and-mcp.md.

## Entries

### 2026-10-06T11:38:21Z planner-E-0017
The plan for E-0017. I revisited the epic's seven drafts and added one. Each draft now has tags, topics, predicted touches, a forecast, a cost of delay value, and a `### Planning` note giving the reasons. The value comes from your 6h a cycle (TH-0171) and is shared by the turns each story removes. I have not changed any story's words.

The stories and their `after`:

| Story | Waits for | Forecast | USD a week | Value per hour |
|-------|-----------|---------:|-----------:|---------------:|
| S-0269 `flai task done` | none | 55m | 427 | 7.8 |
| S-0273 `flai test` | none | 60m | 213 | 3.6 |
| S-0272 end on a question | none | 50m | 96 | 1.9 |
| S-0274 `flai story start` | S-0261 (draft) | 35m | 64 | 1.8 |
| S-0293 classify run turns (new) | none | 40m | 59 | 1.5 |
| S-0275 issue and ADR `--commit` | S-0252 (done), S-0245 | 39m | 41 | 1.1 |
| S-0270 `flai verify` | S-0273 | 70m | 43 | 0.6 |
| S-0271 `flai stream state` | none | 33m | 16 | 0.5 |

Ordering is yours. The last column is what a value-per-duration policy would sort on.

Nearly every story touches the same files:

- `flai/internal/harness/harness.go`
- `flai/internal/mcpserver/folder.go`
- `flai/internal/hostapi/writes.go`
- `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`
- `work-management.md`

So overlap holds will run them one at a time whatever their `after`. I added no `after` beyond real dependencies.

The assumptions I made:

1. S-0273's "tiers the project defines" become project data in the manifest (`flai/internal/manifest`, `project-manifest.md`). Today they exist only as scripts and `Makefile` targets, and a project made from the template needs to declare its own.
2. S-0273 creates a runner package, predicted as `flai/internal/verify`, which S-0270 extends. S-0270 may share `flai/internal/serve/checks.go`'s step runner rather than add a third.
3. `task_done`, `story_start`, `verify`, and `test` are the story's agent's calls. A task sub-agent's flai writes stay refused by `flai guard`, so `flai/internal/guard` is not predicted.
4. S-0270 changes `.claude/agents/verifier.md`, which the story's agent cannot write. It will send you the file to paste. That is counted in its 70m.
5. `metrics.md` changes only with an ADR, so S-0272 and S-0293 each record one.
6. S-0275's `after` named S-0250, which was cancelled as a duplicate of S-0245. I replaced it with S-0245, which still holds ADR numbering's fix. S-0252 is done.

My proposals. Each needs your word, and I have rewritten nothing:

1. Narrow S-0271 to `flai stream state`. ADR-0089 already delivered its first half: `flai criteria tick` and `untick`, the MCP tool `criteria_tick`, the host method `item.criteria`, and ticking on the dashboard's item page. The harness prompt and `work-management.md` already send agents there. The 119 `sed` ticks predate it. Recommended:
   - Title: "Narrative state is a command: flai stream state writes Current state and Next steps".
   - Drop the `flai story tick` and dashboard-ticking parts of the goal and criteria.
2. S-0269 can take `--tick <n>` now, through the existing `criteria tick`, instead of waiting for S-0271 as its Notes say. Recommended: add that as a fifth criterion.
3. Merge S-0272's fourth criterion, `flai stats` counting the empty wakes, into S-0293, which classifies every turn of a run, empty wakes among them. Recommended: drop it from S-0272, which leaves it at 3 criteria and about 40m.
4. Added S-0293 as a draft: `flai stats` classifies a story run's turns, measured against the epic's hand count within ten percent. Without it the epic's done condition ("a run shows no turn that only commits, syncs, ...") cannot be checked, and a baseline is best taken before the others land. The epic's `## Stories` list is your text and does not name it yet.
5. S-0275 commits issue files on story branches. Until S-0278 (ready) lands, each of those commits conflicts in `design/issues/summary.md`. Recommended: add S-0278 to S-0275's `after`. I left it out until you say.

Optional split: S-0270 has 31 predicted touches. The dashboard half of its second criterion, showing the last result, could be its own story if S-0270 grows. I recommend keeping it whole for now.

Nothing to drop. Related stories outside the epic:

- S-0251 and S-0254, which touch `stream_sync.go` and `touches.go` that `flai task done` composes. They are best done before S-0269.
- S-0287, which changes what guard allows for `flai adr` and meets S-0275.

### 2026-10-06T11:42:40Z alex
Resolved.
