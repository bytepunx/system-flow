---
id: T-1302
type: task
nature: improvement
title: The convention, the design, and the users' guide say that a close-out leaves out a markdown finding on another open story's narrative
status: backlog
parent: S-0318
owner: alex
created: 2026-10-08T00:10:44Z
updated: 2026-10-08T00:10:44Z
transitions: []
stream: S-0318
tags: [flai, template]
touches: [design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-1301]
---
# T-1302 The convention, the design, and the users' guide say that a close-out leaves out a markdown finding on another open story's narrative

## Work

Say what T-1301 built wherever the close-out's recording is described. It waits for T-1301, so that the words match the code. It shares no path with T-1303, so the two can run together.

- `design/conventions/work-management.md` and its template copy, `template/root/design/conventions/work-management.md`: the close-out rule names `wip.overlap` and `item.archive`. Add that a markdown finding on another open story's narrative is left out and recorded in no issue, because that story's own close-out finds it. Change the baseline in both files alike.
- `template/CHANGELOG.md`: an entry for the template convention's change, as S-0280's entry reads.
- `design/system/continuous-improvement.md`, where it says flai records the findings a close-out meets outside the story (ADR-0085). Name the new ADR beside ADR-0122.
- `design/system/flai-cli.md`, the `flai check` row.
- `docs/users/flai.md`, the `--story` and `--record-issues` paragraphs of `flai check`.

## Done when

- Each file above says what is left out, and links the ADR T-1300 wrote.
- `flai check --strict` passes, and the markdown lint passes on the changed files.

## Notes
