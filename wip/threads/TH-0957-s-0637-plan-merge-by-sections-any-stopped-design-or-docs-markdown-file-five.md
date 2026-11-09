---
id: TH-0957
title: "S-0637 plan: merge by sections any stopped design or docs markdown file, five tasks in four layers, 1h30m, 175 USD a week"
anchor:
  path: wip/kanban/stories/S-0637-stories-that-each-add-a-section-at-one-spot-in-a-shared-design-document.md
  item: S-0637
status: open
participants: [system-flow]
created: 2026-11-09T22:10:27Z
updated: 2026-11-09T22:10:27Z
---

# TH-0957 S-0637 plan: merge by sections any stopped design or docs markdown file, five tasks in four layers, 1h30m, 175 USD a week

On wip/kanban/stories/S-0637-stories-that-each-add-a-section-at-one-spot-in-a-shared-design-document.md.

## Entries

### 2026-11-09T22:10:27Z system-flow
**Recommendation:** confirm the fix below, or choose an alternative. The tasks are drafted for it.

**Fix (recommended):** widen ADR-0175's section merge from markdown inside `claims.shared` to every markdown file under the manifest's design and docs folders. Issue files and generated indexes keep their own steps, and wip stays out. All three I-0315 instances were settled by hand the same way: keep both sections. The section merge already does that for two sections added at one spot ("one both added appears once per side, ours first"). It skips `design/system/server-performance.md` only because that file is not shared. The gate is one matcher, `sharedMarkdown` in `flai/internal/issues/generated.go`. `issues.FilesIn` passes it to every rebase flai runs: stream sync, task done's sync, `flai accept`, and the merge queue's batch. A section the two stories changed on the same lines still stops the rebase, as it does today.

**Alternatives:**

- Add each such file to `claims.shared`. That list is yours, and a file only shows up as a hub after it conflicts.
- A planning rule that gives each story's section its own file. It fixes nothing for stories already written.

**Tasks and layers:**

| Layer | Task | Touches | After |
|---|---|---|---|
| 1 | T-2811 ADR refining ADR-0175 | `design/adrs/0273-…md`, `design/adrs/README.md` | none |
| 2 | T-2812 widen the matcher, update its unit test | `flai/internal/issues/generated.go`, `generated_test.go` | T-2811 |
| 2 | T-2813 design and user docs | `design/system/workflow.md`, `agent-coordination.md`, `flai-cli.md`, `docs/users/flai.md` | T-2811 |
| 3 | T-2814 real-git tests reproducing I-0315 at acceptance and at sync | `flai/cmd/accept_conflict_test.go`, `stream_sync_test.go` | T-2812 |
| 4 | T-2815 close I-0315 | the issue file, `design/issues/summary.md` | T-2813, T-2814 |

T-2812 and T-2813 share no file and can run together.

**Figures:**

- Forecast: 1h30m, delivery 2026-11-10T01:15Z. `flai forecast` gave 3m from 2 criteria and 0 touches. I raised it for the ADR, the Go change, two real-git tests, and four documents.
- Cost of delay: 175 USD a week, `flai cod`'s figure from your input of 1h10m a cycle. The input counts 2 occurrences; I-0315 now records 3. Raise the input if you want the extra occurrence counted. I left it as it is.

**Assumptions:**

- I found no earlier plan and no thread on S-0637, so this is a first plan.
- The ADR number 0273 is a guess. S-0625 and S-0628 both claim 0272. The story's agent renames the touch to whatever `flai adr new` gives.
- The trial merge keeps reporting these files as a conflict between two stories, as ADR-0175 says for shared files. Making it run the section merge too would cut the conflict messages. I would make that a separate story, not part of this one.
- The ordering stays ADR-0175's: main's section first, then the story's. That gives acceptance order, which is not always the cause order a human picks.
