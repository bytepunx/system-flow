---
id: TH-0777
title: "S-0550 plan: eight tasks in four layers; cost of delay inputs needed"
anchor:
  path: wip/kanban/stories/S-0550-hub-files-merge-without-conflict-the-adr-index-the-issues-summary-and-the.md
  item: S-0550
status: open
participants: [planner-S-0550, orchestrator]
created: 2026-10-31T07:46:01Z
updated: 2026-10-31T07:53:27Z
---

# TH-0777 S-0550 plan: eight tasks in four layers; cost of delay inputs needed

On wip/kanban/stories/S-0550-hub-files-merge-without-conflict-the-adr-index-the-issues-summary-and-the.md.

## Entries

### 2026-10-31T07:46:01Z planner-S-0550
**Question: cost of delay inputs.** S-0550 and E-0031 have no cost of delay inputs, so `flai cod S-0550` refuses and I wrote no value.

My recommendation, on S-0550:

| Input | Value | Why |
|---|---|---|
| `time_lost_per_cycle` | 10m | 248 notices over 97 acceptances is about 2.6 extra `stream sync` runs per acceptance, at 2 to 4 min of agent time each, plus the 62 hub-file conflicts worked by hand |
| `penalty_per_week` | none | no spend comes from these files |
| `revenue_per_week` | none | |

You can set them with `flai edit S-0550 --time-lost-per-cycle 10m`, or reply with other figures and I write the value.

**Plan.** There were no tasks before this run. I wrote eight.

| Layer | Task | Criterion | Waits for |
|---|---|---|---|
| 1 | T-2390 ADR: generated ADR index, union-merged changelogs, acceptance notice skips shared paths | 1–3 | none |
| 1 | T-2394 `flai check` compares the issues summary with its folder; `flai check --fix` | 1 | none |
| 1 | T-2397 acceptance notice skips `claims.shared`; serve test with two `docs/users/flai.md` stories | 3 | none |
| 1 | T-2399 `merge=union` in both `.gitattributes`; `flai release` names a `## Unreleased` section; `git.md` | 2 | none |
| 1 | T-2402 before figures for 2026-10-15 to 29 and how to measure after | 4 | none |
| 2 | T-2405 ADR index render, `flai adr index`, adr new/accept/supersede regenerate | 1 | T-2390 (decision, `README.md`) |
| 3 | T-2409 `adr.index` stale check and fix; rebase regenerates the index | 1 | T-2394 (`check.go`, `fix.go`), T-2405 (render) |
| 4 | T-2413 flai-cli.md, flai.md, flai-reference.md, workflow.md, continuous-improvement.md, decisions.md, changelog | 1–3 | T-2409, T-2397, T-2399, T-2402 |

Forecast: 5h, up from `flai forecast`'s 20m. Delivery at 18:30Z, counted from S-0548's move to review. S-0550's `### Planning` gives the reasons.

**What the code already does.** These findings change the scope.

- The ready-story hold already skips shared paths (ADR-0096 §2). Only the acceptance notice reports them. T-2397 adds the test the criterion asks for.
- The issues summary is already generated, and `flai issue summary` exists. What is missing is a check against its folder and `--fix`.
- The ADR index is edited row by row. ADR-0157 rejected regenerating it, so T-2390 supersedes it in part.

**Assumptions.**

- **Digest index.** The digest has no index file (`design/digest/README.md`), so criterion 2's digest clause has nothing to mark. T-2390's ADR records this, and no task adds one.
- **Changelog headings.** `merge=union` alone breaks the changelog. Two stories that each write a `## x.y.z` heading leave two sibling headings, which MD024 refuses (I-0107). T-2399 has stories write under `## Unreleased`, and `flai release` names that section with the version.
- **Template copy.** `template/CHANGELOG.md` sits outside `template/root`. T-2399 adds `template/.gitattributes` only if `flai template push` publishes `template/` as the repository root.
- **Trial merge.** It is not verified that `git merge-tree --write-tree` honours a `merge=union` attribute. The planner cannot run git writes. T-2399 tests it.
- **Notice scope.** "Sends no overlap notice" means the acceptance notice. `flai task done`'s message about a task's commit still names shared paths.
- **Where the counts live.** The overlap notices and conflicts are in `.flai-cache/overlaps.jsonl` and `.flai-cache/conflicts.jsonl`, not the serve journal. Criterion 4 can be read only a fortnight after acceptance, so it stays unticked at close-out. T-2402 records the before figures and the method.

**Proposals.**

- `flai/internal/sync` and `flai/internal/claims` in S-0550's touches name no package that exists. I kept them as declared. The story's agent may remove them.
- Criterion 2's digest clause could be dropped from the story's words. That is your call.

**Defect.** My start prompt said to read the thread `asked`. No such thread exists: the installed flai reads the trigger `asked` as a thread name. TH-0771 reported the same.

### 2026-10-31T07:53:27Z orchestrator (recommendation)
**Recommendation:** set `time_lost_per_cycle: 10m` on S-0550, as the planner recommends:

```bash
flai edit S-0550 --time-lost-per-cycle 10m
```

Then drop "and the digest index" from criterion 2. The digest has no index file.

## The figure

- **The goal is time lost.** It costs about 2.6 extra `stream sync` runs per acceptance, plus the hub-file conflicts worked by hand. That is agent time, not revenue or a penalty.
- **The story gives no reason for an alternative.** No spend comes from these files, so no penalty applies.
- **Put it on the story, not on E-0031**, as on TH-0773 to TH-0775. Inputs on the epic would stop flai taking inputs on any of its stories.

## Why this needs you

I tried to set it, and flai refused: "S-0550 is ready, not in the backlog: the orchestrator gives cost of delay inputs to a story in the backlog alone". You moved S-0550 to ready at 2026-10-31T07:24:05Z. A confirmation of this entry does not write the input. Running the command does.

## The plan: approved

- **The tasks fit the criteria.**
  - Criterion 1: T-2394, T-2405, and T-2409.
  - Criterion 2: T-2399.
  - Criterion 3: T-2397.
  - Criterion 4: T-2402.
  - T-2390's ADR comes first because ADR-0157 rejected regenerating the index. T-2413 carries the docs.
- **The assumptions hold.** Stories write under `## Unreleased`, because MD024 would refuse two sibling version headings. T-2399 tests that `merge-tree` honours `merge=union`. Criterion 4 stays unticked at close-out, with T-2402's before figures and method.
- **The forecast and touches fit.** The 5h forecast over flai's 20m states its reason. ADR-0096 §3 narrows the declared `flai/internal/serve` and `template` folders to the files the tasks name.

## The proposals

- **Touches `flai/internal/sync` and `flai/internal/claims`:** the story's agent removes them with `flai touches`. They name no package. Today they also overlap the same two folders in S-0551's touches.
- **Criterion 2's digest clause:** drop it, as above. No digest index exists to mark `merge=union`, and the criterion's words are yours. Left as it is, the clause cannot be ticked honestly at close-out.
- **The `asked` trigger defect:** TH-0771 already reports it. Nothing more is needed here.

The thread stays open for the cost of delay input and the criterion edit.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
