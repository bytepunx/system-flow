---
id: T-1150
type: task
nature: improvement
title: The convention, the design, and the users' guide say that wip.overlap compares claims and that a close-out records no overlap
status: backlog
parent: S-0279
owner: alex
created: 2026-10-07T01:21:13Z
updated: 2026-10-07T01:21:13Z
transitions: []
stream: S-0279
tags: [flai, template]
touches: [design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/workflow.md, design/system/agent-coordination.md, design/system/continuous-improvement.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1145]
---
# T-1150 The convention, the design, and the users' guide say that wip.overlap compares claims and that a close-out records no overlap

## Work

Bring every document that describes `wip.overlap` or the close-out's notes in line with what T-1143 and T-1145 built, linking T-1141's ADR. It waits for T-1145, the last code task, so that it describes the behaviour as shipped and so that the reference is generated from the final help. It shares no file with T-1149, so the two run together.

- `design/conventions/work-management.md` and its template copy, `template/root/design/conventions/work-management.md`: the close-out rule says that "any `wip.overlap`" is printed as a note and recorded in an issue. Say instead that only an overlap naming the story is printed, and that none is recorded. Edit the baseline in both copies in the same words, and add a line to `template/CHANGELOG.md`.
- `design/system/flai-cli.md`: the `flai check` and `flai touches` rows.
- `design/system/workflow.md`: the open-story bullet under Branches and collisions.
- `design/system/agent-coordination.md`: the overlap-test row.
- `design/system/continuous-improvement.md`: the bullet on what a close-out records.
- `docs/users/flai.md`: the `flai check --story` paragraph.
- `docs/users/flai-reference.md`: regenerate it with `make flai-reference` from T-1145's help. Do not edit it by hand.

## Done when

- No document says that a close-out records a `wip.overlap`, or that `wip.overlap` compares raw touches. A search for `wip.overlap` under `design/` and `docs/` finds only text that matches the ADR.
- The convention's two copies match, and `template/CHANGELOG.md` names the change.
- `make flai-reference` leaves `docs/users/flai-reference.md` unchanged after the commit.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes
