---
id: TH-1304
title: "Plan for S-0837: eight tasks in four layers retire the publish formula and make release_publish the orchestrator's one publish"
anchor:
  path: wip/kanban/stories/S-0837-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md
  item: S-0837
status: resolved
participants: [planner-S-0837, orchestrator]
created: 2026-11-17T12:39:55Z
updated: 2026-11-17T12:41:03Z
---

# TH-1304 Plan for S-0837: eight tasks in four layers retire the publish formula and make release_publish the orchestrator's one publish

On wip/kanban/stories/S-0837-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md.

## Entries

### 2026-11-17T12:39:55Z planner-S-0837
S-0837 has eight draft tasks in four layers, 46 file touches, a forecast of 1h38m, and a cost of delay of 30.00 USD a week.

## Layers

| Layer | Task | What |
|-------|------|------|
| 1 | T-3600 | Accept ADR-0316 as superseding ADR-0240 and ADR-0273, set their `superseded_by`, and open a thread on S-0811 for its criterion 2 |
| 2 | T-3601 | `release_publish` fetches, evaluates again, publishes only what is pending, answers `covered_by`, and drops the formula refusal; the board-before-and-after test |
| 2 | T-3602 | `flai formula run` refuses the orchestrator every formula; the `push` gate in `flai release --pending` stays |
| 2 | T-3603 | The orchestrator's role file, all three copies, publishes with `release_publish` only; definition regenerated |
| 2 | T-3604 | `flai upgrade` removes an unchanged retired template file and reports a changed one |
| 3 | T-3605 | Delete the template's `publish.md`, both `scripts/fetch.sh`, and `publish_judgement_test.go`; fix the READMEs (after T-3604) |
| 4 | T-3606 | Both conventions, in both copies, and the template's changelog |
| 4 | T-3607 | The designs and the user, operator, and contributor guides |

Every layer-2 task waits for T-3600, because the story asks that the ADR be accepted, or changed, first. The layer-2 tasks share no file. T-3606 and T-3607 wait for the code and role tasks they describe.

## Assumptions

- **Supersede, not refine.** ADR-0316 lists ADR-0240 and ADR-0273 under `refines`. Criterion 1 asks that they name it as what supersedes them. T-3600 moves them to `supersedes` before acceptance, and writes `superseded_by` on the two by hand, since `flai adr accept` does not set it.
- **`--judged` stays only if `release_publish` uses it.** T-3600 asks you, and T-3601 builds what you say.
- **The dashboard's Publish is unchanged.** T-3601 gives the orchestrator's path its own arguments (`--require-pending`, exit 2 mapped) if `publishArgs` is shared.
- **This repository's `publish.md` is already gone.** The stopgap commit 4052fc603 removed it, so only the template's copy is deleted.
- **S-0811 is not edited here.** It is finalized and ready. T-3600 asks on a thread on S-0811 for its criterion to be rewritten. S-0811 is not named in `after`: neither story reads the other's output, and a trial merge judges their overlap on `flai/cmd/release.go` and the formula.

## Proposals

None to split, merge, or drop.

### 2026-11-17T12:40:48Z orchestrator
Approved. The plan fits S-0837.

- **The tasks cover the eight criteria.** Criteria 1 and 8 are T-3600. Criteria 2 and 3 are T-3601. Criterion 4 is T-3601 and T-3602. Criterion 5 is T-3604 and T-3605. Criterion 6 is T-3603. Criterion 7 is T-3606 and T-3607.
- **The layers hold.** Layer 2 waits for T-3600, because the goal puts the ADR's acceptance first. Layer 2's tasks share no file. T-3605 waits for T-3604's retirement path. The docs come last, so they describe what was built.
- **The assumptions stand.** Moving ADR-0240 and ADR-0273 from `refines` to `supersedes` follows criterion 1, and setting `superseded_by` is the one edit the decisions convention allows on an accepted ADR. Asking on S-0811 rather than editing it follows criterion 8. Asking the operator about `--judged` at T-3600 is the story's agent's to do.
- **The figures stand.** The forecast is 1h38m, flai's figure. The value is 30.00 USD a week, from the input set on TH-1292.

Resolving the thread. I am finalizing S-0837 under `finalize_drafts`.

Source: design/adrs/0316-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md

### 2026-11-17T12:41:03Z orchestrator
Resolved: The orchestrator approved S-0837's plan: eight tasks in four layers, a 1h38m forecast, and a value of 30.00 USD a week.
