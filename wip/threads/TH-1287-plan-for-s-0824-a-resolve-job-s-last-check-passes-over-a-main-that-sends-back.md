---
id: TH-1287
title: "Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0824-a-rebase-and-resolve-job-runs-the-whole-close-out-again-each-time-main-moves.md
  item: S-0824
status: open
participants: [planner-S-0824, orchestrator]
created: 2026-11-17T08:35:52Z
updated: 2026-11-17T08:37:02Z
---

# TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers

On wip/kanban/stories/S-0824-a-rebase-and-resolve-job-runs-the-whole-close-out-again-each-time-main-moves.md.

## Entries

### 2026-11-17T08:35:52Z planner-S-0824
S-0824 is planned with six draft tasks in four layers. It assumes the first of I-0398's options; please confirm it, or name another, before T-3551's ADR is written.

The recommended option: in a rebase and resolve job, the close-out's last check (`flai verify --sync-only`) passes when main moved but the acceptance's `verify_behind_main` would raise no blocker for the head it verified. That means no path main brought in overlaps the story's change set by ADR-0303's rule, or the overlapping paths send back no tier. The acceptance still judges anything main brings in later. Another round runs only when the gate would block, and its round line says why. In I-0398's case, S-0742 and S-0753 would have cost a round each only if their paths overlapped S-0817's change set.

The tasks and their layers:

| Layer | Task | What |
|---|---|---|
| 1 | T-3551 | ADR refining ADR-0313 and ADR-0156: the condition, the flag or variable, the round fields, and the alternatives |
| 2 | T-3552 (after T-3551) | `flai verify --sync-only` with a flag judges a moved main by the gate's rule; the gate moves out of `preview` so both share it |
| 2 | T-3553 (after T-3551) | `closeout.Status`, the job's record, and `resolve.status` carry the round, the rounds, and each round's why; fixes `stepAfter` missing the behind-main sync |
| 3 | T-3554 (after T-3552, T-3553) | The resolve job sets the variable; both close-out scripts and the container loop pass the flag; tests that reproduce I-0398 |
| 3 | T-3555 (after T-3553) | The review page shows the round and why |
| 4 | T-3556 (after T-3554, T-3555) | `flai-cli.md`, `flaiover-dashboard.md`, `docs/users/flai.md`, `docs/users/flaiover.md`; close I-0398 |

The assumptions:

- Only a resolve job passes over a moved main. A close-out job, and the agent's close-out before review, keep ADR-0156's rule, since a story moving to review must contain main.
- Overlap-based reuse in later rounds (option 2) is left for T-3551 to weigh, not built. A round that overlaps `flai/**` selects the integration tier either way.
- One order for acceptances and resolve jobs (option 3) is out: ADR-0156 rejected holding acceptances.
- No `after`. S-0808 changes what the gate judges, but S-0824 builds on the gate as main has it. S-0822 and S-0823 share `resolvejob.go`, `closeout.go`, and the review page, which a trial merge judges.
- The story keeps one folder touch, `design/adrs/`, because the ADR's file name comes from its title.

The figures:

- Forecast: 2h30m, raised from flai's 1h2m to match S-0804's 2h17m across the same packages.
- Cost of delay: 150 USD a week, from the 1h per cycle input, unchanged.

I was asked to plan again from a thread called "asked". No such thread is on S-0824 or I-0398, and the story had no tasks or earlier plan, so I planned it from scratch.

### 2026-11-17T08:37:02Z orchestrator (recommendation)
**Recommended: the first option, as the planner wrote it.** In a rebase and resolve job only, the close-out's last check passes over a main that moved, as long as `verify_behind_main` would raise no blocker at the head it verified. Another round runs only when that gate would block, and the round line says why. T-3551's ADR records this, refining ADR-0313 and ADR-0156 §3.

**Why:**

- The acceptance already runs this gate on whatever main brings in later (ADR-0291, ADR-0303). Running the whole close-out again for every move of main repeats that judgement. In I-0398 it cost three rounds and 33 minutes.
- ADR-0156 rejected the third option, holding acceptances while a close-out runs, because "That delays the operator's acceptances for minutes".
- The close-out a story's agent runs before review keeps ADR-0156's rule, so a story still enters review containing main.

**Why I am asking you:** ADR-0156 §3 says main's commits "are not passed over". Passing over them in a resolve job is a new decision, and it changes what a passing verdict means there.

**Alternative:** the second option, ADR-0303's overlap rule applied to every round in place of each tier main's paths select. It saves less: a round that overlaps `flai/**` still runs the integration tier.

I leave S-0824 a draft until you answer, because the tasks follow from the option.

Source: design/adrs/0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md § Alternatives considered
