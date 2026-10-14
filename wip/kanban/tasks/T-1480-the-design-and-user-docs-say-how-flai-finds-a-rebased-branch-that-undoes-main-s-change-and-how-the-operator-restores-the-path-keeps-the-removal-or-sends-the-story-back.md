---
id: T-1480
type: task
nature: remediation
title: The design and user docs say how flai finds a rebased branch that undoes main's change, and how the operator restores the path, keeps the removal, or sends the story back
status: backlog
parent: S-0371
owner: alex
created: 2026-10-12T17:16:06Z
updated: 2026-10-14T21:01:32Z
transitions: []
stream: S-0371
tags: [cli, docs]
touches: [design/system/flai-cli.md, design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1478, T-1575]
---
# T-1480 The design and user docs say how flai finds a rebased branch that undoes main's change, and how the operator restores the path, keeps the removal, or sends the story back

## Work

It waits for T-1478 and T-1575, the last code tasks (T-1575 waits for T-1574 and T-1479), so that the docs describe the wording, the flags, and the writes as built, and the reference is generated once from every help.

- `design/system/flai-cli.md`: the `flai stream sync` row says the rebase reads the old merge base, reports `undoes_main` after a clean rebase, and records it, linking the ADR from T-1476; the `flai task done` row's sync step says the same; a new `flai stream restore` row; the `flai accept` row gains `--keep-removal` and `--reason`, the blocker, and the refusal; the `flai hostapi` row's methods gain `stream.restore` and `accept.run`'s `keep_removals`.
- `design/system/workflow.md` § Branches and collisions (ADR-0019): one paragraph on the check, the record, and the three ways out, and the ADR.
- `design/system/strategic-agents.md` § Accepting a story (S-0221): the orchestrator never keeps a removal; an undone path is a blocker it leaves for the operator, on a thread, as step 4 says.
- `docs/users/flai.md` § Story branches, § Closing a task, § Accept and release, and the `flai hostapi` section: the warning, `flai stream restore`, `--keep-removal` with `--reason`, the two host writes, and that the dashboard's review page will resolve it, once its story lands.
- `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md`: regenerated with `make flai-reference`, which writes both; `TestReferenceIsCurrent` fails until it is run.

## Done when

- The six files say what T-1477 to T-1575 built, and the reference and the flag index match the helps.
- The markdown lint and `flai check --strict` pass on them.

## Notes

`docs/users/flai.md`, `flai-reference.md`, and `flai-cli.md` are in the manifest's shared paths, so they hold no other story.
