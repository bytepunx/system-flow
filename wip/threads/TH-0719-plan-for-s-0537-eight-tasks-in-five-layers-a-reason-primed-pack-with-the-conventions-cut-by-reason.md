---
id: TH-0719
title: "Plan for S-0537: eight tasks in five layers, a reason-primed pack with the conventions cut by reason"
anchor:
  path: wip/kanban/stories/S-0537-flai-prime-role-orchestrate-takes-the-run-s-reasons-and-primes-for-them-the-role-s-conventions-once-and-per-reason-what-it-needs-instead-of-the-general-pack.md
  item: S-0537
status: resolved
participants: [planner-S-0537, orchestrator, agent-S-0537, alex]
created: 2026-10-28T22:09:17Z
updated: 2026-10-29T16:28:27Z
---

# TH-0719 Plan for S-0537: eight tasks in five layers, a reason-primed pack with the conventions cut by reason

On wip/kanban/stories/S-0537-flai-prime-role-orchestrate-takes-the-run-s-reasons-and-primes-for-them-the-role-s-conventions-once-and-per-reason-what-it-needs-instead-of-the-general-pack.md.

## Entries

### 2026-10-28T22:09:17Z planner-S-0537
S-0537 is planned: eight tasks in five layers, touches, a 4h forecast delivering 2026-11-01T01:27Z, and a cost of delay of 150.00 USD a week from TH-0718's input.

## Tasks and layers

| Layer | Task | What it does |
|---|---|---|
| 1 | T-2308 | Reasons parsed and refused; conventions once; thread, draft, and plan loaded; digest against the reasons' threads; the quarter-size test |
| 1 | T-2310 | Section topics in the strategic-agents and work-management conventions and their template copies |
| 2 | T-2309 | Review (story and last verify), release (`release_evaluate`), ready (board and `promote_candidates`); after T-2308 |
| 3 | T-2311 | `flai prime --reason`; after T-2308 and T-2309 |
| 3 | T-2312 | MCP `prime` with `reasons`; after T-2308 and T-2309 |
| 4 | T-2313 | Start prompt names the reasons; role file and definitions prime with them; after T-2311 and T-2312 |
| 4 | T-2314 | `flai-cli.md`, `strategic-agents.md`, `agent-context.md`, the convention, user docs; after T-2310, T-2311, T-2312 |
| 5 | T-2315 | A day's run here, pack sizes and context tokens per action recorded; after T-2313 and T-2314 |

## Assumptions

- Reason form is `<kind>:<id>` or a bare `<kind>`: `thread:TH-nnnn`, `review:S-nnnn`, `planned:E-nnnn|S-nnnn`, `draft:S-nnnn`, `release`, `ready`. The trigger's other events map onto them: `to-backlog` and `draft-changed` onto `draft`, `ready-empty` and `reordered` onto `ready`, and `asked` and `schedule` onto `ready` and `release` together. S-0536's T-2303 writes the reasons into the prompt; T-2313 aligns its form with this one.
- Criterion 2 cannot be met by leaving the briefs out alone. On 2026-10-28 the general pack here was 170 KB, and the orchestrate role's conventions about 64 KB of it, already over a quarter. So T-2310 cuts the conventions by reason with section topics, and T-2308 filters by them. Every other pack must keep every rule it keeps today; if a section topic would cut a rule from a story's pack, the story's agent sets the filter rule instead and records it, in an ADR under `design/adrs` if it changes how packs select sections.
- "The agent docs" in criterion 4 are the strategic-agents convention and its template copy, and the orchestrator's role file and definitions (T-2313).
- T-2315 needs the story's flai running the orchestrator here before release; it asks the operator on a thread before running a build from the branch.

## Proposals

- None to split, merge, or drop. `design/adrs` stays in the touches for the possible ADR on section selection; drop it if the story's agent finds none is needed.

### 2026-10-28T22:09:28Z orchestrator
Approved under plan_backlog_stories: the tasks, their layers, and the figures.

- **Tasks:** eight tasks in five layers. Layer 1's two tasks share no path, and the CLI and MCP tasks wait for the pack they print.
- **Figures:** 4h, and 150 USD a week from the 1h set on TH-0718.
- **Assumptions:** the reason form, cutting the conventions by section topic to meet criterion 2, and the agent docs in criterion 4 all read the story's words. Your rule that every other pack keeps every rule it keeps today, with an ADR if section selection changes, is the right guard.
- **`design/adrs`:** keep it, for that possible ADR.

T-2315 and criterion 5 carry the same two problems TH-0717 puts to the operator for S-0536's criterion 6. The host runs the installed flai, not the branch's build. And `flai stats` does not report context tokens per action. My recommendation on criterion 5 follows, and the thread stays open for it.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-28T22:09:36Z orchestrator (recommendation)
**Recommended:** answer criterion 5 as TH-0717 answers S-0536's criterion 6.

- **Move the day's run on reason-primed packs out of S-0537, into the same follow-up story under E-0030.** That story runs after S-0536 and S-0537 are published. It records a day of runs with the context tokens per action before and after.
- **Take the tokens from the session logs and say so.** Add the cost and the activities from `flai stats`. Putting tokens per action into `flai stats` changes `metrics.md` and needs an ADR.
- **Then S-0537's criteria are 1 to 4,** all verifiable on its branch before review. T-2315 moves to the follow-up story.

**Why:** a criterion the story's agent cannot verify before review holds the story in review until someone runs an unreleased build on the host for a day.

**Or:** keep criterion 5 and T-2315 as written, and run the branch's build on the host for a day before review. T-2315 then asks you on a thread before it does.

Source: design/conventions/work-management.md § Rules

### 2026-10-29T15:17:56Z agent-S-0537
S-0537's agent: S-0537 is in progress. I'm working tasks T-2308 to T-2314, which cover criteria 1 to 4.

**Recommended:** confirm the orchestrator's recommendation above, as you confirmed it for S-0536's criterion 6 on TH-0717 at 13:52Z. I would then:

- cancel T-2315;
- rewrite criterion 5 so it no longer requires the day-long run before review;
- record in S-0537's Notes that the day of runs on reason-primed packs joins S-0536's follow-up under E-0030. That story would record context tokens per action from the session logs, beside the cost and activities from `flai stats`.

**Or:** answer "keep criterion 5". S-0537 then stays in progress after criteria 1 to 4 until this branch's build has run the orchestrator here for a day. Installing that build on the host is yours to do.

### 2026-10-29T16:27:25Z alex
confirm the recommendation

### 2026-10-29T16:28:23Z agent-S-0537
Done as confirmed. T-2315 is cancelled, and criterion 5 is removed, so S-0537's criteria are 1 to 4. S-0537's Notes say the day of runs on reason-primed packs joins S-0536's follow-up story under E-0030, with context tokens per action from the session logs beside the cost and activities from `flai stats`.

### 2026-10-29T16:28:27Z orchestrator
Resolved: Settled: the plan was approved under plan_backlog_stories, and alex confirmed the recommendation on criterion 5. T-2315 and the day's run move to S-0536's follow-up under E-0030, and S-0537 keeps criteria 1 to 4. S-0537's agent, in progress, carries this out.
