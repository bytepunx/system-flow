---
id: T-3464
type: task
nature: improvement
title: This repository's manifest drops queue and when and its workflow file names the suites in their place, with no change to what any step runs
status: backlog
parent: S-0813
owner: alex
created: 2026-11-16T23:50:04Z
updated: 2026-11-16T23:50:16Z
transitions: []
stream: S-0813
tags: [flai]
touches: [system-flow.yaml, design/workflows/board.md, design/conventions/tooling.md]
after: [T-3455, T-3459]
---
# T-3464 This repository's manifest drops queue and when and its workflow file names the suites in their place, with no change to what any step runs

## Work

Criterion 4. This repository runs the baseline today: it has no `design/workflows/board.md`, and its `integration` and `smoke` tiers set `queue: true` and `when: ["flai/**", "template/**"]` in `system-flow.yaml`. It waits for T-3459, so that the file is what flai runs, and T-3455, so that the check's advice gives the lines and goes quiet once they are in.

- Before any change, record what each step selects for three sets of paths: one under `flai/`, one under `template/`, and one only under `scripts/`. The steps are the close-out (`flai verify`), the move to review, the acceptance's gates, and the merge queue's batch. Use `flai test`, `flai verify`, or a test, whichever prints a selection. Write the result in this task's `## Notes`.
- Create `design/workflows/board.md` from the baseline, `flai/internal/workflow/baseline/board.md`, with the lines T-3455's advice gives.
- Remove `queue` and `when` from the two tiers with `flai manifest set`, since `system-flow.yaml`'s `tests:` is flai's to write. If it cannot remove a key, say so in `## Notes` and on the story's thread before you edit the file by hand.
- Record the same selections again, and compare them in `## Notes`. They must be equal. A difference stops the task: raise it on the story's thread rather than accept it.
- Rewrite the `when` sentence in `design/conventions/tooling.md`'s project additions to say the workflow file now decides it. Leave the baseline above the marker alone.

## Done when

- The two selections in `## Notes` are equal for every step and path set.
- `flai check --strict` passes in the story's worktree with no `manifest.tests-queue` finding.
- `flai test --all` passes, or the close-out does.

## Notes
