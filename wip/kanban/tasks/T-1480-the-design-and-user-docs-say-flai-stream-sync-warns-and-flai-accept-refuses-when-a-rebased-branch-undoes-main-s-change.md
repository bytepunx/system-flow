---
id: T-1480
type: task
nature: remediation
title: The design and user docs say flai stream sync warns and flai accept refuses when a rebased branch undoes main's change
status: backlog
parent: S-0371
owner: alex
created: 2026-10-12T17:16:06Z
updated: 2026-10-12T17:16:06Z
transitions: []
stream: S-0371
tags: [cli, docs]
touches: [design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1478, T-1479]
---
# T-1480 The design and user docs say flai stream sync warns and flai accept refuses when a rebased branch undoes main's change

## Work

It waits for T-1478 and T-1479, so that the docs describe the wording and the refusal as built, and the reference is generated once from both helps.

- `design/system/flai-cli.md`: the `flai stream sync` row says the rebase reads the old merge base and reports `undoes_main` after a clean rebase, linking the ADR from T-1476; the `flai task done` row's sync step and the `flai accept` row say the same, the second that it refuses.
- `design/system/workflow.md` § Branches and collisions (ADR-0019): one sentence on the check and the ADR.
- `docs/users/flai.md` § Story branches, § Closing a task, and the acceptance section: the warning, what to do about it, and the refusal.
- `docs/users/flai-reference.md`: regenerated with `make flai-reference`.

## Done when

- The four files say what T-1478 and T-1479 built, and the reference matches the helps.
- The markdown lint and `flai check --strict` pass on them.

## Notes

All four but `workflow.md` are in the manifest's shared paths, so they hold no other story.
