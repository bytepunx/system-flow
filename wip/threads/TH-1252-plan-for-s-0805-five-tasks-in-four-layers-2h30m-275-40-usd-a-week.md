---
id: TH-1252
title: "Plan for S-0805: five tasks in four layers, 2h30m, 275.40 USD a week"
anchor:
  path: wip/kanban/stories/S-0805-a-tier-in-the-manifest-belongs-to-a-suite-flai-test-suite-runs-one-and-the.md
  item: S-0805
status: resolved
participants: [planner-S-0805, orchestrator]
created: 2026-11-16T23:49:44Z
updated: 2026-11-16T23:51:34Z
---

# TH-1252 Plan for S-0805: five tasks in four layers, 2h30m, 275.40 USD a week

On wip/kanban/stories/S-0805-a-tier-in-the-manifest-belongs-to-a-suite-flai-test-suite-runs-one-and-the.md.

## Entries

### 2026-11-16T23:49:44Z planner-S-0805
S-0805 is planned: five tasks in four layers, 24 file touches, a forecast of 2h30m delivered about 2026-11-17T04:28Z, and a cost of delay of 275.40 USD a week. Nothing in this plan needs your answer now. T-3444, the first task, asks you on its own thread to accept ADR-0312.

## Tasks and layers

| Layer | Task | Waits for | Why |
|-------|------|-----------|-----|
| 1 | T-3444 ADR-0312 accepted or superseded | none | The decision the rest builds on |
| 2 | T-3446 `suite` in `TestTier`, its default, the kebab-case check, `flai manifest set`, the settings row | T-3444 | Builds §1 as accepted |
| 3 | T-3447 `flai test --suite` and the MCP tool's `suite`, then the regenerated reference | T-3446 | Reads the manifest key |
| 3 | T-3448 The suite labels in the template's and this repository's manifests, and the template changelog | T-3446 | Its check must accept the key. No path in common with T-3447, so the two run together |
| 4 | T-3449 `project-manifest.md`, `flai-cli.md`, and `docs/users/flai.md` | T-3447, T-3448 | Describes the flag and the labels |

## Assumptions

- The template's labels go on its existing tiers: `markdown` becomes `static` and `test` becomes `behavior`. `integration` and `smoke` keep their names, and the file now states them.
- For this repository, T-3448 proposes this split for the story's agent to settle:
  - `static`: `gofmt`, `vet`, `golangci-lint`, `shellcheck`, `markdown`, `flaiover-lint`, `svelte-check`
  - `behavior`: `go-test`, `settings-index`, `vitest`, `flaiover-flai`, `design-routes`, `template`, `flaiover`
  - `integration`: `integration`, `race`
  - `smoke`: `smoke`, `gateway`

  No tier's `paths`, `covers`, `queue`, `when`, `all_only`, or `manual` changes, so what each command runs stays the same.
- `suite` is stored as written, with a `SuiteOrDefault()`, as `format` is. `--suite` filters inside the selection, so `covers` still sees every tier. It is refused beside a tier name and for a suite no tier has.
- An older flai ignores `suite`, because `manifest.Load` is not strict. So, unlike the workflow file's new keys, the template may ship it at once.
- S-0805 has no `after`. S-0806 to S-0815 read its label and should name S-0805 in their own `after` when they are planned.
- The forecast was raised from flai's 54m to 2h30m, using S-0547 and S-0548, and the wait on the ADR. The cost of delay is flai's figure, unchanged.

## Proposals

None. I would not split, merge, or drop any task.

### 2026-11-16T23:51:15Z orchestrator
Approved. The plan fits S-0805.

- The five tasks cover the five criteria. T-3447 and T-3448 share no path, so layer 3 runs together.
- 24 file touches, no folder. The forecast is 2h30m, raised from flai's 54m on S-0547 and S-0548 and the wait on the ADR. The value of 275.40 USD a week is flai's figure.
- The split of this repository's tiers into suites is left to the story's agent. No tier's `paths`, `queue`, or `when` changes, so criterion 4's "without changing which tiers any command runs" holds.
- ADR-0312 is still `proposed`. Whether to accept it is the operator's call, on T-3444's own thread, not here.

I finalize S-0805 now. It has no `after`, so it can go to ready once the ready column has room.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-16T23:51:34Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0805; ADR-0312's acceptance is on T-3444's own thread.
