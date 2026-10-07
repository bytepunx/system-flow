---
id: T-1167
type: task
nature: improvement
title: The convention, the design, and the users' guide say that a close-out records no item.archive
status: backlog
parent: S-0280
owner: alex
created: 2026-10-07T15:04:22Z
updated: 2026-10-07T15:04:22Z
transitions: []
stream: S-0280
tags: [flai, template]
touches: [design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-1165]
---
# T-1167 The convention, the design, and the users' guide say that a close-out records no item.archive

## Work

Say what T-1165 built where the close-out's recording is described. It waits for T-1165, so that the words match the code.

- The close-out rule in `design/conventions/work-management.md` and its template copy, `template/root/design/conventions/work-management.md`: it names `wip.overlap`; add that an `item.archive` is recorded in no issue and left out, and that the operator archives with `flai archive` in the main checkout.
- `template/CHANGELOG.md`: an entry for the template convention's change.
- `design/system/continuous-improvement.md`, where it says flai records the findings a close-out meets outside the story (ADR-0085).
- `design/system/flai-cli.md`, the `flai check` row.
- `docs/users/flai.md`, the `--story` and `--record-issues` paragraphs of `flai check`.

Link the ADR T-1164 wrote from each design and docs change. It shares no path with T-1166 and can run beside it.

## Done when

- Each file above says that a close-out records no `item.archive` and leaves out one that does not name the story, and links the ADR.
- The two copies of the convention say the same.
- The markdown lint and `flai check --strict` pass.

## Notes
