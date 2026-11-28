---
id: TH-1334
title: "Plan for S-0849: a container's kluster --job asks through its agent cache folder, five tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0849-flai-kluster-deploy-job-teardown-job-and-flai-checks-run-job-from-a-story-s.md
  item: S-0849
status: open
participants: [planner-S-0849, orchestrator]
created: 2026-11-28T00:05:34Z
updated: 2026-11-28T00:06:52Z
---

# TH-1334 Plan for S-0849: a container's kluster --job asks through its agent cache folder, five tasks in four layers

On wip/kanban/stories/S-0849-flai-kluster-deploy-job-teardown-job-and-flai-checks-run-job-from-a-story-s.md.

## Entries

### 2026-11-28T00:05:34Z planner-S-0849
The plan for S-0849 is written: five draft tasks in four layers, plus the story's touches, forecast, and cost of delay. One decision is yours: which fix to build.

## Proposed fix (recommended)

In a story's container, where the main checkout's `.flai-cache/kluster` is absent, the three `--job` commands write their request in `.flai-cache/agents/S-nnnn/kluster-requests/`. The container already mounts that folder read-write (ADR-0319), so no mount is added. flai serve's request pass also takes requests from each story's folder. It refuses a request whose story is not the folder's. It writes the answer and the job state there, so `--wait` works from the container. The folder sits beside the copy, `kluster/`, not inside it, because a teardown that passes removes the copy whole. A new ADR records this, refining ADR-0292 §1 and ADR-0319.

Alternative: in a container, `--job` refuses at once and names the MCP tool. That is one small task plus the docs, but the documented command stays broken where agents now run. If you choose it, T-3655 records it instead, and I will replan T-3656 to T-3658.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-3655 Write the ADR | `design/adrs/README.md` (the ADR file itself through the story's `design/adrs` folder touch) |
| 2 | T-3656 Folder-based request functions and `AgentRequestsDir` in the kluster package, after T-3655 | `flai/internal/kluster/request.go`, `request_test.go` |
| 3 | T-3657 flai serve takes and answers in each story's folder, after T-3656 | `flai/internal/serve/kluster_requests.go`, `kluster_requests_test.go` |
| 3 | T-3658 The commands ask through the folder in a container, with the test that reproduces I-0405, and the help and reference, after T-3656 | `flai/cmd/kluster.go`, `flai/cmd/checks.go`, `flai/cmd/kluster_container_test.go`, `docs/users/flai-reference.md` |
| 4 | T-3659 `kluster.md` and `docs/users/flai.md`, then close I-0405, after T-3657 and T-3658 | `design/system/kluster.md`, `docs/users/flai.md`, I-0405, `design/issues/summary.md` |

T-3657 and T-3658 share no path, so they run together.

## Figures

- Forecast: 1h30m, raised from flai's 33m. S-0834, the closest story, took 85m. Delivery 2026-11-28T08:45:00Z.
- Cost of delay: 25 USD a week, from your input of 10m per cycle. Kept as is.

## Assumptions

- The root-based functions in `request.go` stay as wrappers, so `flai/internal/mcpserver/kluster.go` does not change. The MCP tools keep asking under the main checkout's folder on the host.
- The test that reproduces the defect answers through the folder itself, as the ADR-0319 container tests fake a container, since this repository runs no kluster or container in tests.
- No `after` on another story. S-0835 overlaps on `kluster_requests_test.go`, `kluster.md`, and `docs/users/flai.md`; a trial merge judges that.

## What you need to do

Confirm the recommended fix, or choose the alternative. Then finalize the draft.

### 2026-11-28T00:06:52Z orchestrator (recommendation)
**Recommended:** build the planner's proposed fix, a requests folder in the story's agent cache folder, and approve the plan as it stands.

**Why this fix:**

- ADR-0319 gave the agent cache folder its kluster use and mounts it read-write, with no new mount. Its Consequences name this exact gap: the three `--job` commands write under `.flai-cache/kluster/requests/`, which is not mounted. The fix closes that gap the same way ADR-0319 closed I-0361.
- Putting the folder beside `kluster/` and not inside it is right: ADR-0319 §1 says a teardown that passes removes the copy whole.
- flai serve refuses a request whose story is not the folder's, so a container still asks only for its own namespace. ADR-0292's boundary holds.
- The alternative, a refusal that names the MCP tool, leaves the documented command broken in the containers where agents now run.

**The plan fits the story:**

- T-3658 holds the test that reproduces I-0405 (criterion 1). T-3659 closes I-0405 (criterion 2).
- T-3657 and T-3658 share no path. The layers follow the real waits.
- The 1h30m forecast states its reason, S-0834's 85m. The 25 USD a week value comes from your 10m input.

**Why I recommend rather than answer:** the fix is a new ADR that refines ADR-0292 §1 and ADR-0319. That decision is not recorded yet, so it is yours.

**What happens next:** once this is settled, I finalize S-0849 under `finalize_drafts`. `flai promote --drafts` already lists it complete. Until then it stays a draft.

Source: design/adrs/0319-flai-copies-a-story-namespace-s-record-step-logs-state-and-kubeconfig-into-the.md § Consequences
