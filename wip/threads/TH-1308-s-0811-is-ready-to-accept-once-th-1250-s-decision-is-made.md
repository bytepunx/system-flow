---
id: TH-1308
title: S-0811 is ready to accept once TH-1250's decision is made
anchor:
  path: wip/kanban/stories/S-0811-the-release-s-publish-runs-the-suites-verify-at-publish-lists-at-the-main.md
  item: S-0811
status: resolved
participants: [orchestrator]
created: 2026-11-18T09:09:44Z
updated: 2026-11-20T08:18:48Z
---

# TH-1308 S-0811 is ready to accept once TH-1250's decision is made

On wip/kanban/stories/S-0811-the-release-s-publish-runs-the-suites-verify-at-publish-lists-at-the-main.md.

## Entries

### 2026-11-18T09:09:44Z orchestrator
**Recommended:** confirm my recommendation on TH-1250, the failure path. When a `verify_at.publish` suite fails, the formula's publish step fails and its instance blocks and asks. Then resolve TH-1250, and I accept S-0811 on my next run.

**What stops the acceptance:** one blocker, `thread_open` on TH-1250. That thread holds a decision no one has recorded yet, so it is yours to settle, not mine to resolve.

**What passed at 5541dbfcd920795b455a584798d9a31543eadc6b, the head of `story/S-0811`:**

- `flai verify S-0811 --last`: every step passed, integration and smoke included.
- The verifier found all four criteria met:
  1. `flai/cmd/release.go`, `flai/cmd/release_publish_suites.go`, `flai/cmd/release_publish_suites_test.go`, `flai/internal/hostapi/writes_outcome.go`, `flai/internal/hostapi/writes_test.go`, `flai/internal/hostapi/writes_actions.go`, `docs/users/flai-reference.md`
  2. `template/root/design/workflows/formulas/publish.md`, `flai/cmd/publish_formula_suites_test.go`, `template/CHANGELOG.md`, on the failure-path reading
  3. `flai/cmd/release.go`, `flai/cmd/release_publish_suites.go`, `flai/cmd/release_publish_suites_test.go`
  4. `design/system/pushing-from-the-board.md`, `design/system/workflow.md`, `docs/operators/runbooks/change-the-board.md`, `docs/users/flai.md`
- The accept dry run names no other blocker.

**If you want the clean stop instead:** say so on TH-1250. S-0811 goes back to in-progress so that T-3443 can add exit 5 to the step's `stop`.

If no one answers TH-1250, `flai serve` confirms my recommendation after `confirm_after` (48h).

### 2026-11-20T08:18:48Z orchestrator
Resolved: TH-1250's decision is made: flai serve confirmed the failure path after 48h, and TH-1250 is resolved. S-0811 now waits only on a new close-out over the current main.
