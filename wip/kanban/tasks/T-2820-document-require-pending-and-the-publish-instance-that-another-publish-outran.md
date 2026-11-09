---
id: T-2820
type: task
nature: improvement
title: Document --require-pending and the publish instance that another publish outran
status: backlog
parent: S-0643
owner: agent
created: 2026-11-09T22:13:08Z
updated: 2026-11-09T22:13:08Z
transitions: []
stream: S-0643
tags: [docs, release]
touches: [docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/system/flai-cli.md, design/system/strategic-agents.md]
after: [T-2818]
---
# T-2820 Document --require-pending and the publish instance that another publish outran

## Work

Document the flag T-2818 adds and what a publish instance does when another publish outran it.

- `docs/users/flai.md`: the `--pending` line in `## Accept and release`, and `### When the orchestrator publishes`.
- `design/system/flai-cli.md`: the Commands row for `flai release --pending`.
- `design/system/strategic-agents.md`: `### Publishing a release (S-0222)`, which names only step 2's stop today.

Regenerate `docs/users/flai-reference.md` and the flag rows of `docs/operators/settings.md` with `make flai-reference`, never by hand.

This task waits for T-2818 because the reference is generated from the help text T-2818 writes. It shares no path with T-2819, so the two run together.

## Done when

- The three documents above describe `--require-pending`, its exit 2, and the release its output names.
- `docs/users/flai-reference.md` and `docs/operators/settings.md` are regenerated, and `TestReferenceIsCurrent` and the `settings-index` tests pass.

## Notes
