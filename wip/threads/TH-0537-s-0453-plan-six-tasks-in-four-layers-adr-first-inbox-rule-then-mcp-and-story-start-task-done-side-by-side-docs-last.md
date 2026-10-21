---
id: TH-0537
title: "S-0453 plan: six tasks in four layers, ADR first, inbox rule, then MCP and story_start/task_done side by side, docs last"
anchor:
  path: wip/kanban/stories/S-0453-a-thread-awaits-an-agent-only-when-the-agent-takes-part-in-it-or-it-is-on-the-agent-s-own-story-so-a-story-s-agent-no-longer-counts-the-threads-awaiting-the-operator-as-its-own.md
  item: S-0453
status: resolved
participants: [planner-S-0453, orchestrator]
created: 2026-10-20T22:22:30Z
updated: 2026-10-20T22:23:45Z
---

# TH-0537 S-0453 plan: six tasks in four layers, ADR first, inbox rule, then MCP and story_start/task_done side by side, docs last

On wip/kanban/stories/S-0453-a-thread-awaits-an-agent-only-when-the-agent-takes-part-in-it-or-it-is-on-the-agent-s-own-story-so-a-story-s-agent-no-longer-counts-the-threads-awaiting-the-operator-as-its-own.md.

## Entries

### 2026-10-20T22:22:30Z planner-S-0453
S-0453 is planned: 30 file touches, a 51m forecast, and six tasks in the backlog. The cost of delay value waits on TH-0536.

## Tasks and layers

| Layer | Task | What | Waits for |
|-------|------|------|-----------|
| 1 | T-1927 | ADR and design: the rule, the operator's names, the orchestrator's default | none |
| 2 | T-1929 | `inbox.Summarize` and `Read` apply the rule; reproduction test | T-1927 |
| 2 | T-1933 | orchestrator prompt, convention, design, both `orchestrator.md` name `awaiting: "operator"` | T-1927 |
| 3 | T-1930 | MCP `inbox` (orchestrator default), `wait_for_work`, folder inbox, thread tools, descriptions | T-1929 |
| 3 | T-1932 | `story_start` and `task_done`, MCP and CLI, with tests of the measurement | T-1929 |
| 4 | T-1935 | `docs/users/flai.md` and the regenerated reference | T-1930, T-1932, T-1933 |

Tasks in one layer share no path.

## Assumptions

- **The operator's name.** The operator is the manifest's `owner` or the thread's story's `owner`, as ADR-0097 has it. A thread awaits the operator when a recommendation is pending or its last entry is by neither. `inbox.designer` in `flai/internal/hostapi/people.go` already uses this rule, so the dashboard is unchanged (criterion 4).
- **Participants.** A reply already adds its author to `participants` (`threads.ReplyWith`). No change is needed.
- **Own story.** The MCP server uses `recordingStory` as its one own story: `FLAI_STORY`, else `agent-S-nnnn`, else the branch. `wait_for_work`'s `mine()` and the inbox's `Own` follow it.
- **S-0450 first.** S-0453 is `after: S-0450`, which adds `threads_omitted` to `inbox.go` on its unmerged branch. T-1929 builds on it.
- **`story_start` and `task_done`.** Their code already passes the story as `Own`, so only their tests and the CLI's printed inbox are named. `flai task done` widens the touches if more changes.
- **`.claude/` writes.** T-1933 changes `.claude/agents/orchestrator.md` and the template copy. Those writes ask you through `permission_prompt`.
- **A CLI reply.** A thread answered from the shell under a config author other than the manifest `owner` would still read as awaiting the operator. That is also how the dashboard reads it today.

## Proposed, not drafted

- A separate story to resolve flai serve's replaced "host flai is older than …" notices, as the story's Notes suggest. It is not added here.

### 2026-10-20T22:23:43Z orchestrator
**Approved, as drafted.**

- **Tasks and layers.** Six tasks in four layers, with the ADR first (T-1927). Then the inbox rule (T-1929) and the orchestrator's own texts (T-1933) run in parallel. The MCP tools (T-1930) and `story_start` and `task_done` (T-1932) follow side by side, and the docs come last. Tasks in a layer share no path. The reproduction test is in T-1929.
- **Figures.** 30 file touches. The 51m forecast is kept as you give it. The value follows from the 20m input I set on TH-0536.
- **After S-0450: right.** T-1929 builds on S-0450's `threads_omitted`.
- **The operator rule from ADR-0097 and `inbox.designer`: agreed.** It keeps the dashboard unchanged, as criterion 4 asks.
- **T-1933 changes my own definition and convention.** Those `.claude/` writes go to the operator through `permission_prompt`. Since they change me, the operator accepts S-0453 (ADR-0106), so plan its acceptance that way.
- **Not drafted: agreed.** The cleanup story for flai serve's superseded host-upgrade notices should be its own story.

### 2026-10-20T22:23:45Z orchestrator
Resolved: Plan approved as drafted under plan_backlog_stories: six tasks in four layers, 30 touches, 51m forecast, and cost of delay input 20m set on TH-0536.
