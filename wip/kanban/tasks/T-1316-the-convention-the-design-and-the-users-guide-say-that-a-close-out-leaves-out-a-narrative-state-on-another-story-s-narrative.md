---
id: T-1316
type: task
nature: improvement
title: The convention, the design, and the users' guide say that a close-out leaves out a narrative.state on another story's narrative
status: backlog
parent: S-0323
owner: alex
created: 2026-10-08T00:27:05Z
updated: 2026-10-08T00:27:05Z
transitions: []
stream: S-0323
tags: [flai, template]
touches: [design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-1314]
---
# T-1316 The convention, the design, and the users' guide say that a close-out leaves out a narrative.state on another story's narrative

## Work

Say what T-1314 built wherever S-0280 said it for `item.archive`. It waits for T-1314, so the words describe the behaviour as built.

- `design/conventions/work-management.md`, in the close-out rule beside `item.archive`: a `narrative.state` on another story's narrative is left out too, since only that story's agent writes it and its own close-out stops on it. Make the same change in `template/root/design/conventions/work-management.md`, the baseline copy, and add an entry to `template/CHANGELOG.md`.
- `design/system/continuous-improvement.md`, where it lists what a close-out never records, and `design/system/flai-cli.md`, in the `flai check --story` row: name the rule, the ADR T-1313 wrote, and S-0323.
- `docs/users/flai.md`, under `--story` and `--record-issues`: say it is left out and that an unscoped `flai check` still warns it.
- Bump `updated` on each design and docs file changed.

## Done when

- Each file names the leave-out and links T-1313's ADR where it links ADR-0122.
- The markdown lint passes under `flai test` on the changed paths.
- `flai check --strict` passes.

## Notes
