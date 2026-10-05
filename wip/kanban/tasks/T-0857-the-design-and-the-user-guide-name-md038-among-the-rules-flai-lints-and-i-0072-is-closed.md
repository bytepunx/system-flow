---
id: T-0857
type: task
nature: remediation
title: The design and the user guide name MD038 among the rules flai lints, and I-0072 is closed
status: backlog
parent: S-0262
owner: alex
created: 2026-10-05T03:12:58Z
updated: 2026-10-05T03:12:58Z
transitions: []
stream: S-0262
tags: [flai]
touches: [design/system/flai-cli.md, docs/users/flai.md, design/issues]
after: [T-0856]
---
# T-0857 The design and the user guide name MD038 among the rules flai lints, and I-0072 is closed

## Work

- In `design/system/flai-cli.md`, add MD038, spaces inside a code span (S-0262), beside MD007 where `flai check` lists the rules ADR-0061 names and the ones added since, and on the `mdlint/` line of the layout.
- In `docs/users/flai.md`, add code spans to the list of what flai's lint checks under `flai check`.
- Close I-0072 with `flai issue close I-0072 --reason` saying that mdlint now reports MD038, from the story's worktree, and regenerate `design/issues/summary.md`.
- ADR-0061 is accepted and stays as it is, as S-0240 left it for MD007.
- Waits for T-0856: it describes the rule T-0856 adds, and closes the issue only once the rule is in.

## Done when

- `flai-cli.md` and `docs/users/flai.md` name MD038 or code spans among the rules flai lints
- I-0072 is closed with a reason naming MD038, and `summary.md` is current
- The markdown lint and `flai check --strict` pass on the changed files

## Notes

S-0240's T-0710 did the same for MD007.
