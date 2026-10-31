---
id: T-2380
type: task
nature: improvement
title: The testing convention, flai-cli.md, the manifest design, and the user docs say which tiers the close-out runs and which the batch repeats
status: backlog
parent: S-0547
owner: alex
created: 2026-10-31T07:40:29Z
updated: 2026-10-31T07:40:29Z
transitions: []
stream: S-0547
tags: [flai]
touches: [design/conventions/code-quality.md, design/conventions/tooling.md, design/system/flai-cli.md, design/system/project-manifest.md, docs/users/flai.md]
after: [T-2375, T-2378]
---
# T-2380 The testing convention, flai-cli.md, the manifest design, and the user docs say which tiers the close-out runs and which the batch repeats

## Work

Criteria 2 and 4. Describe what T-2375 decided and T-2378 built, linking the new ADR:

- `design/system/project-manifest.md`: the `when` key in the tier example and the key list, its rule beside `queue`, the `manifest.tests` errors T-2376 added, and the paragraph on this repository's tiers: `go-test` on `{importers}`, and `integration` and `smoke` run in the close-out when their `when` matches, with `go-test` then covered.
- `design/system/flai-cli.md`: the `flai verify` and close-out text and the merge queue pass, saying which tiers the close-out runs, which it leaves, and that the batch runs every tier again.
- `docs/users/flai.md`: the same for `flai verify`, `flai test`, and the close-out, for users.
- The testing convention: the project additions of `design/conventions/code-quality.md` (the tiers line) and `design/conventions/tooling.md` (the `flai test` line), which say the close-out's go-test runs only the changed packages. The story's declared touch `design/conventions/testing.md` does not exist; these two hold the tiers, as S-0548's plan found.

Waits for T-2375, whose ADR it links, and T-2378, whose behaviour it describes. It runs beside T-2379, which changes none of these files, and describes the `when` patterns T-2379 writes as that task names them.

## Done when

- Each of the five files says which tiers the close-out runs, which it leaves to the batch, and that the batch repeats them, and links the new ADR where it gives the decision.
- `flai check --strict` and the markdown lint pass on the five files.

## Notes

Drafted by planner-S-0547. If T-2379 settles on other `when` patterns than `["flai/**", "template/**"]`, align `project-manifest.md` before the close-out.
