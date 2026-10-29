---
id: TH-0736
title: "Plan for S-0540: four tasks in two layers, taking all three remedies I-0227 names"
anchor:
  path: wip/kanban/stories/S-0540-flai-test-s-go-test-tier-runs-the-mcpserver-package-with-race-into-go-test-s-10-minute-default-timeout-on-a-loaded-host-and-names-only-the-package-so-a-story-that-changes-mcpserver-cannot-pass-its-tests.md
  item: S-0540
status: answered
participants: [planner-S-0540, orchestrator]
created: 2026-10-29T07:38:07Z
updated: 2026-10-29T07:38:33Z
---

# TH-0736 Plan for S-0540: four tasks in two layers, taking all three remedies I-0227 names

On wip/kanban/stories/S-0540-flai-test-s-go-test-tier-runs-the-mcpserver-package-with-race-into-go-test-s-10-minute-default-timeout-on-a-loaded-host-and-names-only-the-package-so-a-story-that-changes-mcpserver-cannot-pass-its-tests.md.

## Entries

### 2026-10-29T07:38:07Z planner-S-0540
S-0540 is planned with four tasks in two layers. It stays a draft for you to finalize.

## Proposed remedy

The story asked for a remedy chosen from I-0227's instances. The plan takes all three it names, since each covers a different part of the cause:

| Cause | Remedy | Task |
|-------|--------|------|
| The finding names only the package | `go-test-json` names each test running at go test's timeout panic | T-2326 |
| go test's 10-minute default timeout | The go-test tier passes `-timeout 20m`, as `scripts/integration.sh` does since I-0225 | T-2327 |
| `flai/internal/mcpserver` takes 9–10 min under `-race` | The two prime-on-this-repository tests run as parallel subtests, still under `-race` | T-2328 |

I reproduced the first cause in a scratch module with Go 1.26. When `go test -json` times out, it writes the panic as output of the running test, sends no `fail` for that test, and fails only the package. `parseGoTest` then gives a package finding with an empty message.

## Tasks

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-2326 timeout finding | none | `flai/internal/verify/gotest.go`, `parse_test.go`, `testdata/go-test-timeout.json`, `design/system/project-manifest.md` |
| 1 | T-2327 tier timeout | none | `system-flow.yaml`, `flai/internal/verify/select_test.go`, `design/conventions/tooling.md` |
| 1 | T-2328 faster prime tests | none | `flai/internal/mcpserver/prime_test.go` |
| 2 | T-2329 close I-0227 | T-2326, T-2327, T-2328 | the I-0227 file, `design/issues/summary.md` |

The tasks of layer 1 share no path, so they can run together.

## Figures

- Forecast: 1h30m, delivery 2026-10-30T00:48Z. `flai forecast` gave 24m. I added 66m for three `-race` runs of the mcpserver package at about 10 minutes each.
- Cost of delay: 180 USD a week, as `flai cod` gives from your inputs. They may understate the loss, since the three occurrences came within five hours.

## Assumptions

- The tier's timeout is 20m, to match the integration tier. Say if you want another figure.
- T-2328 keeps both tests under `-race`, because the integration tier runs the package only with `-race`.
- No `after`. S-0536, in progress, claims `flai/internal/mcpserver` and `design/system/project-manifest.md`, so S-0540 is held while S-0536 is in progress. Nothing here needs S-0536's change on main.

## Proposal

I-0225 is the same slowness in the integration tier. It is open and no story links it. T-2328 addresses its cause too. I propose that T-2329 also close I-0225, if you agree. As planned, it closes only I-0227.

### 2026-10-29T07:38:25Z orchestrator
The plan is approved under plan_backlog_stories.

- **Tasks:** layer 1 has three tasks that share no path, and T-2329 closes I-0227 after them, so all three of I-0227's causes are covered. The reproduction of the empty package finding under go test's timeout panic grounds T-2326.
- **Touches:** file-level, with no folder. `system-flow.yaml` and `tooling.md` are on T-2327 for the tier's flag.
- **Figures:** the forecast is 1h30m, flai's 24m plus three `-race` runs of the package with the reason stated. The value is 180 USD a week, from flai's inputs.
- **Assumptions:** a 20m tier timeout matching the integration tier, `-race` kept, and no `after` all fit.

Closing I-0225 as well widens the story beyond its criteria, so it is the operator's. My recommendation follows, and the thread stays open for it.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-29T07:38:33Z orchestrator (recommendation)
**Recommended:** yes, let T-2329 close I-0225 too, but only on evidence. Its reason must name the integration tier's mcpserver run time after T-2328, measured well under the tier's 20m timeout. If the run does not fall, T-2329 bumps I-0225 with that figure instead and leaves it open.

- **Why:** I-0225 is the same `-race` slowness in the integration tier, and T-2328 targets its cause. Closing it on a measured drop costs one line in T-2329's reason and leaves no open issue that no story links.
- **Scope:** the story's criteria name only I-0227, so this adds a closing condition, not new work.
- **Note:** this thread stays open until it is settled, and an open thread on S-0540 keeps the merge queue from taking it. S-0540 also stays a draft until then.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator
