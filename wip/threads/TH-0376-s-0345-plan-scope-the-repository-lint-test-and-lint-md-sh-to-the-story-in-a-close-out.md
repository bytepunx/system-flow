---
id: TH-0376
title: "S-0345 plan: scope the repository lint test and lint-md.sh to the story in a close-out"
anchor:
  path: wip/kanban/stories/S-0345-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md
  item: S-0345
status: resolved
participants: [planner-S-0345, alex]
created: 2026-10-08T08:41:16Z
updated: 2026-10-08T08:47:06Z
---

# TH-0376 S-0345 plan: scope the repository lint test and lint-md.sh to the story in a close-out

On wip/kanban/stories/S-0345-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md.

## Entries

### 2026-10-08T08:41:16Z planner-S-0345
The plan for S-0345 is written: 4 tasks in 2 layers, with touches, a forecast of 30m, and a cost of delay of 180 USD a week. The reasons are under `### Planning` in the story's Notes.

## Proposed fix

A close-out runs two markdown checks over the whole repository. Each will scope itself to the story when `CLOSE_OUT_STORY` is set, as `TestMonorepoIsClean` and `scripts/check.sh` already do (S-0249, ADR-0085). After the fix, a finding in a `wip/` file the story does not change is logged as a note and does not fail. Without the variable, as in `make lint-md` and CI, every file still counts. Main's `wip/` stays linted there.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1363 `TestRepositoryLintsClean` scopes itself to the story, with a test that reproduces I-0117 | `flai/internal/mdlint/mdlint_test.go`, `flai/internal/mdlint/repo_test.go` | none |
| 1 | T-1364 `lint-md.sh` with no files leaves out main's `wip/` and lints only the `wip/` files the story changes | `scripts/lint-md.sh` | none |
| 2 | T-1365 Document both in `scripts/README.md` and the `flai verify` row of `flai-cli.md` | `scripts/README.md`, `design/system/flai-cli.md` | T-1363, T-1364 |
| 2 | T-1366 Close I-0117 with the reason | the I-0117 file, `design/issues/summary.md` | T-1363, T-1364 |

The two tasks in each layer share no path, so each pair can run together.

## Assumptions

1. **Smoke's lint is in scope.** The four instances failed only in integration. Smoke runs `lint-md.sh` over the whole repository right after integration, so it would fail next on the same line, as it did under I-0027.
2. **The test moves to an external test package.** `workitem` imports `mdlint`, so the internal test package cannot import `storygit` without an import cycle.
3. **`scripts/smoke.sh` is not changed.** S-0340 is changing it now. `CLOSE_OUT_STORY` already reaches `lint-md.sh` through the environment.
4. **No new ADR.** The fix applies ADR-0085's scoping to two more checks. The `flai verify` row already says the tests that check the repository read `CLOSE_OUT_STORY`.

## Overlap

`scripts/README.md` is also in S-0340's touches. S-0345 is held while S-0340 is in progress. To avoid the hold, T-1365 could leave the README out and document the change only in `flai-cli.md`, which is a shared path. I recommend keeping the README in, because it is where `lint-md.sh` is described.

## Proposal (not drafted)

Projects made from the template have the same cause in a broader form. `template/root/system-flow.yaml.tmpl`'s markdown tier runs `template/root/scripts/lint-md.sh` over every file on every close-out, `wip/` included. I would make that a story of its own rather than grow this one.

### 2026-10-08T08:47:06Z alex
Resolved.
