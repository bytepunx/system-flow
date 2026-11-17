---
id: T-3450
type: task
nature: feature
title: An ADR maps a tier's queue and when onto the workflow file's verify and verify_at.batch, and says what of ADR-0174 and ADR-0233 it supersedes
status: backlog
parent: S-0813
owner: alex
created: 2026-11-16T23:49:05Z
updated: 2026-11-16T23:49:05Z
transitions: []
stream: S-0813
tags: [flai]
touches: [design/adrs/README.md, design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md, design/adrs/0233-while-the-merge-queue-can-batch-the-close-out-runs-a-queue-tier-whose-when.md]
---
# T-3450 An ADR maps a tier's queue and when onto the workflow file's verify and verify_at.batch, and says what of ADR-0174 and ADR-0233 it supersedes

## Work

Criterion 1. ADR-0312 §8 retires `queue` and `when` once a project's file names suites for `batch`, and leaves the mapping to this story. Write it with `flai adr new` before any code, since every later task reads it.

- Say what each key becomes in the workflow file:
  - `queue: true` with no `when`: the tier's suite is left out of the `verify:` of the move to review and stays in `verify_at.batch`.
  - `queue: true` with `when`: the tier's suite is in the move to review's `verify:` only for stories whose claim or changed paths match `when`. ADR-0312's `verify:` has no paths of its own. Decide whether a `verify:` entry gains a `when` of its own (`{suite, when}`), or whether the tier's `paths` stand in for `when` and the difference is accepted. This repository's `smoke` selects `scripts/**/*.sh` while its `when` is `flai/**` and `template/**`, and `integration` selects `scripts/integration.sh`. Only the first choice leaves what a close-out runs here unchanged, as criterion 4 asks.
  - A tier with neither key: in every suite list, as today.
- Say what reads the keys while a project still sets them: the effective workflow derived from them, until the project's file names suites for `batch`. From then on the keys are ignored, and `flai check` advises their removal.
- Say how a project's file comes to name the lines. `flai check` advises them, and `flai upgrade` writes them only into a workflow file it adds to a project that has none. A project's own file is never written, as ADR-0307 says, and S-0806's criterion 4 says too.
- Supersede, or refine, the parts of ADR-0174 and ADR-0233 that the file now decides, and set `superseded_by` on each one it supersedes whole. Use `flai adr topics` for topics, not a hand edit.

No `after`: the first layer.

## Done when

- The ADR is written, proposed, and indexed in `design/adrs/README.md`. It names the mapping of each key, the precedence while a project still sets them, what `flai check` and `flai upgrade` do, and what of ADR-0174 and ADR-0233 it supersedes or refines.
- `flai check --strict` and the markdown lint pass on the ADR files.

## Notes
