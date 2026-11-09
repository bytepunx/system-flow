---
id: TH-0949
title: "Plan for S-0628: a container story's close-out and checks run in its container; needs a cost of delay input (recommend penalty 40 USD a week) and the remedy confirmed"
anchor:
  path: wip/kanban/stories/S-0628-the-close-out-job-and-the-manifest-s-test-commands-run-the-story-worktree-s.md
  item: S-0628
status: open
participants: [planner-S-0628, orchestrator]
created: 2026-11-09T22:00:27Z
updated: 2026-11-09T22:08:43Z
---

# TH-0949 Plan for S-0628: a container story's close-out and checks run in its container; needs a cost of delay input (recommend penalty 40 USD a week) and the remedy confirmed

On wip/kanban/stories/S-0628-the-close-out-job-and-the-manifest-s-test-commands-run-the-story-worktree-s.md.

## Entries

### 2026-11-09T22:00:27Z planner-S-0628
**Two questions for the operator.**

**1. Cost of delay input.** S-0628 has no inputs and no epic, so `flai cod` refuses it.

Recommended: `penalty_per_week: 40` USD on S-0628.

- No instance has been observed, and the analysis estimates no time lost.
- This project's story agents run on `claude-code` on the host. Only `planning.agent` uses `container`. So the gap is open today only for a project that runs story agents in containers, and for this one once it does.
- The risk is real for those: code the agent writes runs unreviewed on the host, as the operator, at every close-out. So I put it below S-0627's 100, a risk that applies today, and above S-0606's 30.
- No revenue rides on it, so a penalty fits.

Alternative: `penalty_per_week: 15`, if you count it as documentation work until story agents run in containers here.

**2. The remedy.** The issue asks for one before it is built. Recommended: **A**.

- **A. Run it inside (planned).** For a story whose agent's harness is `container`, the close-out job runs `flai verify --record-issues` in the story's container, with the story's mounts, so the tiers run there. The commit and the sync rounds run on the host through flai's own code. Git there takes its repository from the main checkout, not from the worktree's `.git` file or `config.worktree`, which the container can write, and runs only the main checkout's hooks. The operator's checks for that story run in its container too. An ADR and the guides say what still runs on the host.
- **B. Document the boundary only.** An ADR and the guides say the close-out, the tiers, and the checks run the worktree's scripts on the host. This is about a third of the work. But the cause in criterion 1 would still occur, so it would need a criterion reworded, which is your call.

## Plan (remedy A)

| Layer | Task | Touches | after |
|-------|------|---------|-------|
| 1 | T-2772 ADR-0272: a container story's verify runs in its container, and what stays on the host | the ADR, `design/adrs/README.md` | none |
| 1 | T-2775 A `docker run` argv for one command in a story's container, with its six story mounts and no login or key | `flai/internal/harness/container.go`, `container_test.go` | none |
| 2 | T-2777 The close-out job for a container story: verify in the container; commit, `--sync-only`, and the rounds on the host. The test reproducing I-0295 fails on today's code | `flai/internal/serve/verifyjob.go`, `verifyjob_test.go`, `flai/internal/serve/container.go` | T-2775 |
| 3 | T-2782 The operator's checks for a container story run in its container, and a cancel stops it | `flai/internal/serve/checks.go`, `checks_test.go` | T-2775, T-2777 |
| 4 | T-2784 `flai-cli.md`, `strategic-agents.md`, the operators' guide and container runbook, and `docs/users/flai.md` | those five | T-2772, T-2777, T-2782 |
| 5 | T-2785 Close I-0295 | the issue, `design/issues/summary.md` | all |

T-2782 waits for T-2777 because it reuses the image choice T-2777 exposes from `containerImageHost`.

Forecast: 2h, raised from `flai forecast`'s 30m. The job reimplements the close-out script's commit and its three sync rounds in Go, and its tests fake docker. Delivery is 2026-11-10T23:00:00Z, a day out, since the story is a draft.

## Assumptions

- What stays on the host by design is recorded, not changed. That covers the merge queue's batch tiers, an acceptance's `fix_command`, and `flai test` or `flai verify` run by hand. Each runs after a decision or by the operator's own hand.
- The image's flai runs `flai verify` in the container. This repository's close-out script builds the worktree's flai, and that build is the worktree's code, which must not run on the host.
- I-0295's Impact names `manifest.go:1570`. That line is `KlusterErrors`; the tiers are `TestTiers`. T-2785's close reason says so.
- ADR-0253 says the close-out runs inside the container, and `docs/users/flai.md` says it cannot (S-0520). ADR-0272 settles this, and T-2784 brings the documents in line.

## Proposed, not created

- A story of its own: `flai/internal/docedit/file.go` checks paths as strings in `FileReach` and then follows symlinks in `readText`, where `hostapi/docs.go` resolves real paths. The analysis lists it as hygiene beside I-0295. It is another cause. Say if you want it drafted.

I hold here until question 1 is answered. If you choose B in question 2, I will rewrite T-2775, T-2777, and T-2782, which I wrote, and say how on this thread.

### 2026-11-09T22:08:43Z orchestrator (recommendation)
**Recommended:** `penalty_per_week: 40` USD on S-0628, and remedy **A**, as the planner recommends.

**Question 1, the cost of delay input.** No instance has been observed, and the gap is open only for story agents on the `container` harness. Where it is open, the agent's code runs unreviewed on the host as the operator at every close-out. So 40 sits between S-0606's 30 and S-0627's 100, which is a risk that applies today.

- No revenue rides on it, so a penalty fits. The `throughput` policy orders ready by forecast, so the value decides no order today.
- Alternative: `penalty_per_week: 15`, if you count it as documentation work until story agents here run in containers.
- Why I recommend rather than set it: you have finalized S-0628 and moved it to ready, so it is no longer a backlog story. `plan_backlog_stories` lets me set only a backlog story's first inputs. Once you confirm, or the recommendation is confirmed after `confirm_after`, I set the input under `set_cost_of_delay`.

**Question 2, the remedy: A.** ADR-0253 already decided that the close-out runs inside the agent's container, against the main checkout's state. A makes the close-out job and the operator's checks do what that ADR says. B would record the opposite and leave criterion 1 unmet.

- Keeping the commit and the sync on the host through flai's own code, with only the main checkout's hooks, keeps the worktree's code off the host. That is the issue's point.
- You finalized S-0628 with A's tasks (T-2772 to T-2785) as written. That already points to A.

**The symlink check in `flai/internal/docedit/file.go`.** Draft it as a story of its own, outside S-0628. It is a different cause, and S-0628 is large already at 2h.

Source: design/adrs/0253-a-container-that-runs-a-story-s-agent-mounts-the-main-checkout-s-wip-and-system.md § Decision
