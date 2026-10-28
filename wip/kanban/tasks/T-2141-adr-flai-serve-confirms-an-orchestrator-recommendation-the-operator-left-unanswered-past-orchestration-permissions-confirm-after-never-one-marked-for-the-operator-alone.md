---
id: T-2141
type: task
nature: improvement
title: "ADR: flai serve confirms an orchestrator recommendation the operator left unanswered past orchestration.permissions.confirm_after, never one marked for the operator alone"
status: done
parent: S-0502
owner: alex
created: 2026-10-26T04:42:52Z
updated: 2026-10-28T12:47:20Z
transitions:
  - to: ready
    at: 2026-10-28T12:46:33Z
    by: agent-S-0502
  - to: in-progress
    at: 2026-10-28T12:46:34Z
    by: agent-S-0502
  - to: done
    at: 2026-10-28T12:47:20Z
    by: agent-S-0502
stream: S-0502
tags: [flai]
touches: [design/adrs/0206-flai-serve-confirms-a-recommendation-the-operator-left-unanswered-past-confirm-after.md, design/adrs/README.md, design/adrs/0206-flai-serve-confirms-an-orchestrator-recommendation-the-operator-left-unanswered.md]
usage:
  source: log
  seconds: 46
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 3544
      cache_read: 724764
      cache_write: 22442
      cost: 0.3517
---
# T-2141 ADR: flai serve confirms an orchestrator recommendation the operator left unanswered past orchestration.permissions.confirm_after, never one marked for the operator alone

## Work

Write the ADR with `flai adr new`, refining ADR-0090 (who makes a pending recommendation the answer), ADR-0091 and ADR-0081 (what `flai stats` reports of the waits), with topics `cli`, `orchestration`, and `analysis`. It records the decision the operator confirmed on TH-0598 and I-0195 names, settled as follows; where this task's agent finds a reason to decide otherwise, it asks on the plan's thread before writing.

1. **The permission.** `orchestration.permissions.confirm_after`, a duration such as `48h`, off when unset or zero, declared beside the other permissions and shown in the settings catalog as a duration (`KindDuration`), with a meaning and a risk as each permission has.
2. **Who confirms, and when.** `flai serve`, per project, on its per-minute look, while the project's `orchestrate` host action is on. It confirms a pending recommendation (ADR-0090's definition) when all of these hold: its author is `orchestrator`, the orchestrator's run's author name (ADR-0087); it has waited `confirm_after` since its `at`; no entry by the project's owner (the manifest's `owner`) came after it, whoever opened the thread, since `PendingRecommendation` skips the opener's entries and the operator may have opened the thread; the thread was opened by an agent `flai serve` starts (`agent-S-nnnn`, `planner-S-nnnn`, `planner-E-nnnn`), not by `flai serve` itself (upgrade asks, restart asks, forge failures, Claude Code checks) and not by the permission prompt (a title matching the digest's `permissionTitle`, `Allow … ?`), which keeps `.claude/` writes and the host changes flai asks about for the operator; and it is not marked for the operator alone.
3. **The operator-alone mark.** The orchestrator marks a recommendation on money (a cost of delay input, an estimate, spend) or a host change (an install, a key, a host action) as the operator's alone: `thread_reply` with `operator_only: true`, `flai thread reply --recommend --operator-only`. flai can tell neither class from the text, and the orchestrator already classes them under `answer_threads`. The ADR fixes the mark's form in the file so that an older flai reads the entry as ADR-0090 reads a recommendation: recommended, a one-line paragraph just above the `Source:` line, `Awaits: the operator alone`, read into the entry's `operator_only` and kept in its text, as the source line is.
4. **The confirming entry.** By `flai serve`, through `threads.Confirm`'s path: `Confirmed the recommendation of <time> <author> after <duration> with no answer from the operator (orchestration.permissions.confirm_after).`, citing the recommendation's source, status `answered`, so the asking agent is started again as on the operator's confirmation. The line keeps the `Confirmed the recommendation of` prefix so an older flai's metrics still see a confirmation, and ends with a period so the digest's `confirmed` pattern still matches; the digest harvests no answer from it, since it reads the owner's entries alone. The operator may still reply differently afterwards, as on any answered thread.
5. **The log.** Each confirmation is logged in `wip/agents/orchestrator.md` through `serve.LogActivity`, kind `orchestrator`, with no run, so no seconds and no cost, naming the thread, its story, and the duration waited, and in the host's journal as `flai serve`'s other actions are.
6. **The metrics.** `design/system/metrics.md` § Waiting: who ended a wait gains `timeout`, the entry by `flai serve` that confirms a recommendation by this rule; `waiting.weeks[].threads` gains `timeout` with `count` and `total_seconds` beside `orchestrator` and `confirmed`, and `confirmed` keeps only the operator's confirmations; `waiting.longest[].awaited` reads `flai serve` for such a wait.
7. **Alternatives considered.** The orchestrator confirming its own recommendations (rejected on TH-0598: confirming is never the recommender's, and the rule stays in one place the guard and the stats see); flai classing money by the thread's words (brittle); leaving every planner's thread out (the instances, TH-0583 and TH-0588, were planners' threads); a blanket timer on every thread awaiting the operator (most are money or host questions).

`flai adr new` adds the index row to `design/adrs/README.md`. The ADR's number is predicted as 0206 in this task's touches (0204 is the next free on main, S-0504 may take it, and S-0498 has taken 0205 on its branch); `flai task done` widens the touch to the name `flai adr new` gives.

## Done when

- The ADR exists under `design/adrs/` with `refines: [ADR-0090, ADR-0091]` and the seven points above as its decision, and its row is in `design/adrs/README.md`.
- `flai check --strict` scoped to the story is clean and the markdown lint passes on the ADR.

## Notes

Layer 1, waits for nothing. Every other task waits for it, since the manifest key, the mark, the confirming entry, and the metrics field take their names from it.
