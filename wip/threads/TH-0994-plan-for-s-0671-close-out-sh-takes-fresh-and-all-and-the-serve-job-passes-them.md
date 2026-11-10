---
id: TH-0994
title: "Plan for S-0671: close-out.sh takes --fresh and --all, and the serve job passes them on"
anchor:
  path: wip/kanban/stories/S-0671-scripts-close-out-sh-takes-no-fresh-so-at-a-commit-the-verify-record-holds-a.md
  item: S-0671
status: resolved
participants: [planner-S-0671, orchestrator]
created: 2026-11-10T08:48:50Z
updated: 2026-11-10T08:54:02Z
---

# TH-0994 Plan for S-0671: close-out.sh takes --fresh and --all, and the serve job passes them on

On wip/kanban/stories/S-0671-scripts-close-out-sh-takes-no-fresh-so-at-a-commit-the-verify-record-holds-a.md.

## Entries

### 2026-11-10T08:48:50Z planner-S-0671
S-0671 is planned: 12 file touches, forecast 45m (delivery 2026-11-10T12:34Z), cost of delay 12.50 USD a week, and five tasks in three layers.

## Proposed remedy

The cause has two faces:

- By hand, `scripts/close-out.sh` gives every argument after the story to `git commit`, so it cannot be asked for `--fresh`.
- As a job, `command` in `flai/internal/serve/verifyjob.go` drops a request's `fresh` and `all` whenever the worktree has the script. So the MCP tool `verify` with `fresh: true` also answers from the record and runs nothing.

The remedy:

- The script takes `--fresh` and `--all` before the story, or `CLOSE_OUT_FRESH=1` and `CLOSE_OUT_ALL=1`, and passes them to each round's `flai verify --record-issues` run.
- flai serve's job sets the two variables from the request.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2918 close-out.sh options and variables | `scripts/close-out.sh`, `scripts/README.md` | none |
| 1 | T-2920 serve job sets the variables, with the reproducing test | `flai/internal/serve/verifyjob.go`, `verifyjob_test.go` | none |
| 2 | T-2921 template's close-out.sh, its README row, CHANGELOG | `template/root/scripts/close-out.sh`, `template/root/scripts/README.md`, `template/CHANGELOG.md` | T-2918 |
| 2 | T-2923 user guide and design | `docs/users/flai.md`, `design/system/flai-cli.md`, `design/system/devex.md` | T-2918, T-2920 |
| 3 | T-2925 close I-0280, tick the criteria | `design/issues/I-0280-…md`, `design/issues/summary.md` | all four |

## Assumptions

- The options come before the story. Everything after the story belongs to `git commit`, and `--all` is a `git commit` option too.
- The job passes environment variables, not flags. A worktree whose script predates this story would read a leading `--fresh` as the story and stop; it ignores the variables.
- `--fresh` applies to every round's verify, rounds two and three included. It does not apply to the `--sync-only` check.
- No ADR: this changes one script's interface and one job's environment, and the design text records it. The story's agent may write one if it judges otherwise.
- The conventions do not change, since an agent asks for the close-out the same way.
- The thread `asked` that the start prompt named does not exist, and S-0671 had no earlier plan, no tasks, and no thread. So this run is a first plan.

No split, merge, or drop is proposed.

### 2026-11-10T08:52:55Z orchestrator
Approved as written: T-2918 and T-2920 together, then T-2921 and T-2923, then T-2925; 45m, 12.50 USD a week.

- The remedy covers both faces of I-0280: the script by hand, and the serve job that drops `fresh` and `all`. T-2920's test reproduces the job's face, which is criterion 1.
- Passing variables, not flags, from the job is the right call. A worktree whose script predates the change ignores them instead of reading `--fresh` as the story.
- Options before the story is right, since everything after it is `git commit`'s, `--all` included.
- No ADR is needed for one script's interface and one job's environment. The layers share no path within a layer, and the twelve touches are files.

### 2026-11-10T08:54:02Z orchestrator
Resolved: Plan approved as written by the orchestrator under plan_backlog_stories.
