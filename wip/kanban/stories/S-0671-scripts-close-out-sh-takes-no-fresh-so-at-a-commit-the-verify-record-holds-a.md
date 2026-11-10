---
id: S-0671
type: story
nature: improvement
title: scripts/close-out.sh takes no --fresh, so at a commit the verify record holds a pass for it runs no tier, and a second close-out in another place cannot check the same commit again
status: backlog
owner: alex
created: 2026-11-10T08:34:48Z
updated: 2026-11-10T08:48:36Z
transitions: []
tags: [flai, template, close-out, verify, scripts]
topics: [code, review]
touches: [scripts/close-out.sh, scripts/README.md, template/root/scripts/close-out.sh, template/root/scripts/README.md, template/CHANGELOG.md, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, docs/users/flai.md, design/system/flai-cli.md, design/system/devex.md, design/issues/I-0280-scripts-close-out-sh-takes-no-fresh-so-at-a-commit-the-verify-record-holds-a.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-10T08:34:48Z
  value: 12.5
  by: planner-S-0671
  at: 2026-11-10T08:47:59Z
forecast:
  duration: 45m
  delivery: 2026-11-10T12:34:00Z
  basis: "flai forecast's 24m (101 s per unit of size over 151 done improvement stories, size 14), raised to 45m for five tasks across a shell script, a Go function and its serve test, the template copy, and four documents, and a close-out that runs the template, integration, and smoke tiers; delivery moved by the same 21m."
  by: planner-S-0671
  at: 2026-11-10T08:47:59Z
---
# S-0671 scripts/close-out.sh takes no --fresh, so at a commit the verify record holds a pass for it runs no tier, and a second close-out in another place cannot check the same commit again

## Goal

This story remediates [I-0280](../../../design/issues/I-0280-scripts-close-out-sh-takes-no-fresh-so-at-a-commit-the-verify-record-holds-a.md), "scripts/close-out.sh takes no --fresh, so at a commit the verify record holds a pass for it runs no tier, and a second close-out in another place cannot check the same commit again". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0280 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0280 is closed with `flai issue close I-0280 --reason` saying what fixed it

## Tasks
- T-2918 scripts/close-out.sh takes --fresh and --all before the story, or CLOSE_OUT_FRESH and CLOSE_OUT_ALL, and passes them to each round's flai verify
- T-2920 flai serve's close-out job sets CLOSE_OUT_FRESH and CLOSE_OUT_ALL for the worktree's script when the request asks for fresh or all
- T-2921 The template's scripts/close-out.sh takes --fresh and --all and reads CLOSE_OUT_FRESH and CLOSE_OUT_ALL as this repository's does
- T-2923 The user guide and the design say the close-out takes --fresh and --all and that the job passes them on
- T-2925 Close I-0280 with the reason that the close-out takes --fresh and the job passes it on

## Notes

Cost of delay inputs set by flai from I-0280. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T06:58:50Z, 3.1 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from the issue's instance and the code. The cause has two faces. By hand, `scripts/close-out.sh` passes `--record-issues` alone to `flai verify` and hands every argument after the story to `git commit`, so nobody can ask it for `--fresh` (the instance on S-0565's T-2452). As a job, `command` in `flai/internal/serve/verifyjob.go` drops a request's `fresh` and `all` whenever the worktree has the script ("The script takes only git commit options after the story, so all and fresh do not reach it"). So the MCP tool `verify` with `fresh: true` starts a job whose script answers from the record and runs nothing. The plan:

- `scripts/close-out.sh` takes `--fresh` and `--all` before the story, not after it, since everything after the story is `git commit`'s and `--all` is a `git commit` option too. It also reads `CLOSE_OUT_FRESH=1` and `CLOSE_OUT_ALL=1` from its environment. It passes `--fresh` and `--all` to each round's `flai verify --record-issues` run, and neither to the `--sync-only` check.
- flai serve's close-out job sets `CLOSE_OUT_FRESH=1` and `CLOSE_OUT_ALL=1` in the script's environment when the request asks for `fresh` or `all`. It passes no flags, because a worktree whose script predates this story would take a leading `--fresh` as the story and stop. An older script ignores the variables and runs as it does today.
- The template's `close-out.sh` takes the same change, so projects made from the template get it.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `scripts/close-out.sh` | design: the issue names it |
| `scripts/README.md` | co-change (6 of 17 commits): its row gives the script's usage |
| `template/root/scripts/close-out.sh` | co-change (10 of 17) |
| `template/root/scripts/README.md` | co-change (5 of 17) |
| `template/CHANGELOG.md` | co-change (3 of 17): a template change gets an entry |
| `flai/internal/serve/verifyjob.go` | layout: `command` and `closeOutEnv` build the job's run |
| `flai/internal/serve/verifyjob_test.go` | co-change (3 of 17): the reproducing test |
| `docs/users/flai.md` | layout: the `flai verify` section says the close-out runs without `--fresh`, and what the job runs |
| `design/system/flai-cli.md` | layout: the close-out job's "What it runs" bullet |
| `design/system/devex.md` | co-change (2 of 17): its Close-out row gives the script's usage |
| `design/issues/I-0280-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes them |

Left out of the co-change list: `design/conventions/work-management.md` and its template copy, since an agent still asks for the close-out the same way; `flai/cmd/verify_test.go` and `system-flow.yaml`, since `flai verify` and the tiers do not change. No folder touch. The tags `flai` and `template` name the two components it delivers to.

Figures. `flai forecast` gave 24m: 101 s per unit of size over 151 done improvement stories, times size 14. The duration is raised to 45m, since the story has five tasks across a shell script, a Go function and its serve test, the template copy, and four documents, and its close-out runs the template, integration, and smoke tiers. The delivery moves by the same 21m, to 12:34Z. The cost of delay value, 12.50 USD a week, is `flai cod`'s, from the operator's input of 5m lost per 168h cycle at 150 USD an hour, and stands.

No `after`: the story reads no other story's output.
