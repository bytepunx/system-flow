---
id: TH-0207
title: "S-0271 plan: seven tasks in four layers, and one question: is the tick half met by flai criteria tick?"
anchor:
  path: wip/kanban/stories/S-0271-criteria-and-narrative-state-are-commands-flai-story-tick-checks-a-criterion-and-flai-stream-state-writes-current-state-and-next-steps.md
  item: S-0271
status: resolved
participants: [planner-S-0271, alex]
created: 2026-10-06T22:50:53Z
updated: 2026-10-06T22:59:39Z
---

# TH-0207 S-0271 plan: seven tasks in four layers, and one question: is the tick half met by flai criteria tick?

On wip/kanban/stories/S-0271-criteria-and-narrative-state-are-commands-flai-story-tick-checks-a-criterion-and-flai-stream-state-writes-current-state-and-next-steps.md.

## Entries

### 2026-10-06T22:50:53Z planner-S-0271
**Question, needs your answer.** The story still names `flai story tick`, `item_tick` and `item.tick`. ADR-0089 (S-0282) already delivered all three under other names:

- the command `flai criteria tick` and `untick`;
- the MCP tool `criteria_tick`;
- the host method `item.criteria`;
- ticking on the dashboard's item page, through `item.criteria`.

The harness prompt and `work-management.md` already send the agent to them. You resolved TH-0176 without applying my proposal to narrow the story, so I've left its words as you wrote them.

**My recommendation (a):** the existing ticking meets the tick half of criteria 1 and 2. The story's agent records that in the Notes and adds no new names. Reply `a` for this.

**Option (b):** add `flai story tick <id> <n>` (`--untick`) as an alias of `flai criteria tick`. It would refuse a story that is not in progress or review, as criterion 1 says; `flai criteria tick` today refuses only an item that cannot be edited. With (b) I'd add one task: `flai/cmd/criteria.go`, a test, and the reference. MCP and the host channel would keep `criteria_tick` and `item.criteria` under their names. It adds about 5m.

Reply `b` for this, or `c` to also add the MCP and host-channel aliases `item_tick` and `item.tick`.

**The plan.** The tasks are below, in layers. Each task waits for the tasks in the layers above it that it names.

1. **T-1054:** `workitem` replaces the two sections and leaves everything else alone. It refuses a story with no narrative or not in progress or review, and text the lint rejects.
2. **T-1055:** `flai stream state` as text and `--json`, with `-` reading from standard input.
   **T-1056:** the MCP tool `stream_state`, in a new `mcpserver/stream.go`.
   **T-1059:** a `flai check` rule for the sections' shape.
3. **T-1058:** the host method `stream.state`. It waits for T-1055 because it runs the command.
   **T-1061:** `tooling.md`, `work-management.md` and `session-start.md`, their template copies, `template/CHANGELOG.md`, and the harness prompt.
4. **T-1063:** `agent-narrative.md`, `flai-cli.md`, `docs/users/flai.md`, and the regenerated `flai-reference.md` and `docs/operators/settings.md` flag index.

**Figures:**

- **Touches:** I added six paths to the story's touches, so it now has 25, all files and no folder. The six are `narrative_test.go`, `mcpserver/stream.go` and `stream_test.go`, `check/check.go` and `check_test.go`, and `docs/operators/settings.md`.
- **Forecast:** flai's 39m stands, at size 28.
- **Cost of delay:** stays at 16 USD a week, the epic planner's share by turns removed. flai's 77.75 shares by forecast duration.

**Assumptions:**

1. **The check rule's shape:** `## Current state` must not be empty or the placeholder, and `## Next steps` must be a non-empty list. The rule applies only to stories in progress or review. Its severity is chosen so that today's narratives still pass `--strict`. It matches what `scripts/close-out.sh` checks; I left close-out unchanged.
2. **The dashboard:** it does not call `stream.state`. Criterion 2 asks only that it tick, which it already does.
3. **The guard:** it needs no change. `flai stream state` is a flai write like `stream log`.
4. **ADRs:** no ADR is needed. The story adds a command, and ADR-0089 already covers ticking.

I'm waiting on the question above. The figures and the tasks for (a) are already written.

### 2026-10-06T22:56:08Z alex
a

### 2026-10-06T22:56:35Z planner-S-0271
Thanks, I've applied (a). The existing `flai criteria tick`, `criteria_tick`, and `item.criteria` count as the tick half of criteria 1 and 2, and no alias is added. S-0271's `### Planning` notes now record this, and the story's agent will note it beside those criteria when it ticks them.

The tasks, layers and figures are unchanged:

- **Layer 1:** T-1054.
- **Layer 2:** T-1055, T-1056, T-1059.
- **Layer 3:** T-1058, T-1061.
- **Layer 4:** T-1063.
- **Figures:** 39m forecast, 16 USD a week cost of delay.

### 2026-10-06T22:59:39Z alex
Resolved.
