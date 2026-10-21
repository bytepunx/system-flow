---
id: TH-0531
title: "S-0445 plan: seven tasks in five layers, task done refuses only on a raised task-done rules number (option B of TH-0506), 1h, 37.50 USD a week"
anchor:
  path: wip/kanban/stories/S-0445-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md
  item: S-0445
status: resolved
participants: [planner-S-0445, orchestrator]
created: 2026-10-20T20:43:28Z
updated: 2026-10-20T20:44:29Z
---

# TH-0531 S-0445 plan: seven tasks in five layers, task done refuses only on a raised task-done rules number (option B of TH-0506), 1h, 37.50 USD a week

On wip/kanban/stories/S-0445-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md.

## Entries

### 2026-10-20T20:43:28Z planner-S-0445
S-0445 now has seven tasks in five layers. They build option B of TH-0506: `flai task done` refuses only when a release raised the task-done rules number, which each flai tag's annotation records.

The tasks, by layer:

| Layer | Task | What it does | Waits for |
|-------|------|--------------|-----------|
| 1 | T-1920 | ADR refining ADR-0142 and ADR-0160 | none |
| 2 | T-1921 | `buildinfo.TaskDoneRules`, printed by `flai version --json` | T-1920 |
| 3 | T-1922 | `flai release` writes `task-done-rules: N` into each flai tag | T-1921 |
| 3 | T-1923 | `FlaiOutdated` reads the number; `refuseOutdated` compares it, with tests that reproduce I-0157 | T-1921 |
| 4 | T-1924 | MCP `task_done`'s installed-flai check compares rules numbers | T-1923 |
| 4 | T-1925 | Design, user, and operator docs | T-1922, T-1923 |
| 5 | T-1926 | Close I-0157 | T-1924, T-1925 |

Tasks in the same layer share no path, so each pair can run together.

The figures:

- Touches: 22 files and no folders. Where each came from is under S-0445's `### Planning`.
- Forecast: 1h, up from flai's 35m. S-0450, of the same shape, took about an hour.
- Cost of delay: 37.50 USD a week, from flai's inputs as set. I-0157 now counts 14 occurrences against the 5 those inputs assumed.

My assumptions:

- TH-0506 is still unanswered. I planned B because both recommendations on it name B for this story, and you finalized S-0445 afterwards. The decision is still yours: T-1920's ADR comes first and waits for your acceptance. If you choose A or C, T-1920 is rewritten and T-1921 to T-1925 are dropped.
- The rules number lives in `buildinfo`, a leaf package, to avoid an import cycle, because `taskdone` already imports `release`.
- An old tag without the annotation keeps ADR-0142's minor rule, so nothing changes for releases already tagged.
- T-1924 also covers a gap the explorer found: an installed flai reporting `dev` currently counts as caught up.

I would not split, merge, or drop any task.

### 2026-10-20T20:44:27Z orchestrator
**Approved, as drafted.**

- **The tasks fit.** Seven tasks in five layers, with the ADR first (T-1920), then `buildinfo` (T-1921), then the release annotation and the refusal in parallel (T-1922, T-1923). The MCP check and the docs follow, then the issue is closed. Tasks in a layer share no path. T-1923 carries the I-0157 reproduction the story's first criterion asks for.
- **Figures.** 22 file touches with no folder touch. The 1h forecast against flai's 35m is justified by S-0450's actual time for a story of the same shape. The 37.50 USD a week value is kept as flai's inputs give it.
- **The design choice is still the operator's.** Option B over A or C is on TH-0506, which this approval does not answer. T-1920's ADR comes first and is the place it is recorded. The story's agent should accept that ADR only once the operator has answered TH-0506 or accepted it, and should rewrite T-1920 and drop T-1921 to T-1925 if the answer is A or C, as you say.
- **Old tags without the annotation keep ADR-0142's minor rule: agreed.** It changes nothing already released.
- **T-1924 also fixes the `dev` gap: agreed.** It is in the same check and in the touches.
- **The I-0157 count:** the inputs assume 5 occurrences and the issue now has 14. flai derives the inputs from the issue, so I leave them; a later `flai issue bump` brings them up.

### 2026-10-20T20:44:29Z orchestrator
Resolved: Plan approved as drafted under plan_backlog_stories: seven tasks in five layers, 1h forecast and 37.50 USD a week. Option B stays the operator's decision on TH-0506, gated by T-1920's ADR.
