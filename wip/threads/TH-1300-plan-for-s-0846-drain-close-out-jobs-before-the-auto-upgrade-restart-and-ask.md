---
id: TH-1300
title: "Plan for S-0846: drain close-out jobs before the auto-upgrade restart, and ask again what a restart still interrupts"
anchor:
  path: wip/kanban/stories/S-0846-flai-serve-s-auto-upgrade-restart-interrupts-the-close-out-jobs-it-was-running.md
  item: S-0846
status: resolved
participants: [planner-S-0846, orchestrator]
created: 2026-11-17T12:34:24Z
updated: 2026-11-17T12:38:32Z
---

# TH-1300 Plan for S-0846: drain close-out jobs before the auto-upgrade restart, and ask again what a restart still interrupts

On wip/kanban/stories/S-0846-flai-serve-s-auto-upgrade-restart-interrupts-the-close-out-jobs-it-was-running.md.

## Entries

### 2026-11-17T12:34:24Z planner-S-0846
S-0846 is planned: four tasks in four layers, touches, forecast, and cost of delay written.

## Cause

`flai/internal/serve/upgrade.go` `look` waits only for dashboard writes before it asks the host to upgrade (ADR-0158 §5). The restart stops serve, which stops each close-out job's process group. The next serve's `verifyPass.sweep` ends the job `interrupted` and asks nothing again. Rebase and resolve jobs share that pass.

## Proposed fix (decided in T-3583's ADR before code)

1. Drain: ask for no upgrade while a close-out or rebase and resolve job runs, and take no new job while an upgrade is due.
2. Bound: 45 minutes from the first due ask, then upgrade anyway. The longest job on record is the 33-minute rebase and resolve of S-0824.
3. Ask again once: at start, serve writes a new request for a job the restart still interrupted, but not for one that was already asked again.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3583 ADR and `design/system/flai-cli.md` | flai-cli.md (the ADR's file is added by task done) | none |
| 2 | T-3584 drain in the upgrade and the verify pass, with the reproducing tests | upgrade.go, upgrade_test.go, verifyjob.go, verifyjob_test.go, serve.go | T-3583 |
| 3 | T-3585 ask again at the sweep | verifyjob.go, verifyjob_test.go, resolvejob.go, resolvejob_test.go | T-3584 (same file) |
| 4 | T-3586 operator and user docs, `flai issue close I-0388` | docs/operators/index.md, docs/users/flai.md, I-0388, summary.md | T-3585 |

## Figures

- Forecast: 1h30m, up from flai's 26m, delivery 2026-11-17T22:50Z. The raise covers an ADR before code, the concurrency tests, and a close-out of a flai/ change.
- Cost of delay: 300 USD/week, from the inputs carried over from I-0388.

## Assumptions

- A fixed 45-minute bound, not a new manifest setting.
- Kluster, merge-queue batch, and image jobs are out of scope. If the operator wants them drained too, I would add a follow-up story.
- A request made during the drain gets a pending answer, so the asker's ten-second wait does not withdraw it and tell the agent serve is not running. The ADR settles the exact form.

Nothing here needs an answer before the story is pulled. Answer only to change the bound or the scope.

### 2026-11-17T12:37:56Z orchestrator
Approved. The plan fits S-0846.

- T-3583 proposes the drain, its 45-minute bound, and the ask-again rule in an ADR before any code, as the goal asks.
- T-3584 and T-3585 build the drain and the ask-again, each with tests that reproduce I-0388. That meets criterion 1.
- T-3586 documents it and closes I-0388. That meets criterion 2.
- The touches are 12 files and no folder. Each task's touches are inside them. The ADR's own file is added by `flai task done` once `flai adr new` numbers it.
- T-3585 waits for T-3584 because both change `verifyjob.go`, so no layer has two tasks on one path.
- The forecast of 1h30m and the value of 300 USD a week are flai's figures, with the planner's stated reason for the forecast.

The fixed bound with no manifest key, and leaving the kluster, merge-queue, and image jobs out, are proposals. T-3583 puts them to the operator in its ADR, so this approval does not decide them.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-17T12:38:32Z orchestrator
Resolved: The orchestrator approved S-0846's plan under plan_backlog_stories.
