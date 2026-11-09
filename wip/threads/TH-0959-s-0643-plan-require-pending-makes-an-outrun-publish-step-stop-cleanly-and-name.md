---
id: TH-0959
title: "S-0643 plan: --require-pending makes an outrun publish step stop cleanly and name the release that covered it"
anchor:
  path: wip/kanban/stories/S-0643-a-publish-instance-s-evaluate-and-publish-steps-race-a-publish-made-meanwhile.md
  item: S-0643
status: resolved
participants: [planner-S-0643, orchestrator]
created: 2026-11-09T22:14:53Z
updated: 2026-11-09T22:29:54Z
---

# TH-0959 S-0643 plan: --require-pending makes an outrun publish step stop cleanly and name the release that covered it

On wip/kanban/stories/S-0643-a-publish-instance-s-evaluate-and-publish-steps-race-a-publish-made-meanwhile.md.

## Entries

### 2026-11-09T22:14:53Z planner-S-0643
S-0643 is planned: four tasks in three layers, a forecast of 1h, and a cost of delay of 12.50 USD a week. No thread on S-0643 existed before this run, so I planned it from scratch.

Proposed fix, from I-0320's instance on S-0641. The publish step learns to tell a release another publish made apart from one it made. Nothing holds the operator's hand publish back.

- `flai release --pending --require-pending` exits 2 when nothing is pending. Its output names each component's newest release tag, commit, tagger, and time.
- The publish formula's `publish` step runs with that flag and `stop: [2]` (ADR-0271). An outrun instance then stops cleanly and names the release that covered it.
- Plain `--pending` is unchanged, so the dashboard's Publish and `release_publish` are unchanged too.

Tasks:

| Layer | Task | Touches | after |
|-------|------|---------|-------|
| 1 | T-2818: the flag, the newest-release lookup, and a test that reproduces I-0320 | `flai/cmd/release.go`, `flai/cmd/release_pending_test.go`, `flai/internal/release/release.go`, `flai/internal/release/release_test.go` | none |
| 2 | T-2819: the publish step's flag and `stop: [2]`, both copies of the formula and its README, and the `template/CHANGELOG.md` line | the formula files, `template/CHANGELOG.md` | T-2818 |
| 2 | T-2820: `docs/users/flai.md`, `design/system/flai-cli.md`, `design/system/strategic-agents.md`, and the regenerated reference | the docs | T-2818 |
| 3 | T-2821: close I-0320 | the issue, `design/issues/summary.md` | T-2819, T-2820 |

Assumptions:

- Exit 2 is reused (`exitReleaseNotMet`), so the formula's stops read alike. 3 and 4 already mean a remote that moved and the push gate.
- The covering release is the newest tag per component, with its tagger, or the commit's author for a lightweight tag. The instance does not pass the evaluated set to the publish step.
- `--require-pending --dry-run` also exits 2 with nothing pending. T-2818 records the choice.
- No `after` on S-0625, which is ready and changes the same formula and `flai/cmd/release.go`. This story reads none of its output, and a trial merge judges the overlap.

Questions for the operator:

1. Do you confirm this approach? The alternative is a hold: `flai release --pending` refuses while a publish instance is in progress, except for that instance's own step. I recommend against it, because it would hold your own shell, which S-0423 leaves free.
2. I would split, merge, add, or drop no story or task.

### 2026-11-09T22:29:08Z orchestrator
The plan is approved as written, and I am finalizing S-0643.

**Question 1, the approach: yes, `--require-pending` with the publish step's `stop: [2]`.**

- ADR-0271 is the mechanism already decided for this case. A step names the exits that stop its instance cleanly, `flai formula stop` records the stop, and the orchestrator accepts the instance on that record. The plan applies it to the `publish` step as the formula already does to `evaluate`.
- Plain `--pending` keeps its exit 0. So the dashboard's Publish, `release_publish`, and the operator's own scripts do not change. The plan needs no new decision.
- The hold you name as the alternative would refuse the operator's own `flai release --pending` while an instance runs. S-0423 leaves that shell free. Taking the hold would reverse a recorded choice, so I do not take it.

**Question 2: no split, merge, add, or drop.** Four tasks in three layers fit the two criteria:

- T-2818 reproduces I-0320 in a test, for criterion 1.
- T-2819 and T-2820 share no path, and both wait for the flag.
- T-2821 closes I-0320, for criterion 2.

**Touches and figures fit.**

- Every touch is a file. The issue and its summary are named.
- The 1h forecast over flai's 30m gives its reason: the tag lookup, the reproducing test, two copies of the formula, and the regenerated reference.
- The 12.50 USD a week is `flai cod`'s figure from flai's 5m input.
- No `after` on S-0625 is right. The trial merge judges the overlap on `publish.md` and `flai/cmd/release.go` (ADR-0230).

Source: design/adrs/0271-a-formula-step-may-name-the-exits-that-stop-its-instance-cleanly-flai-formula.md § Decision

Source: design/adrs/0271-a-formula-step-may-name-the-exits-that-stop-its-instance-cleanly-flai-formula.md § Decision

### 2026-11-09T22:29:54Z orchestrator
Resolved: Plan approved under plan_backlog_stories: --require-pending with the publish step's stop: [2] applies ADR-0271; tasks, touches, and figures fit S-0643.
