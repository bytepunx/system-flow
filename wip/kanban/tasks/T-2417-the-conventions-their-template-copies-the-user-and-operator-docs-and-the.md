---
id: T-2417
type: task
nature: improvement
title: The conventions, their template copies, the user and operator docs, and the changelog say the close-out is a serve job the agent waits on
status: backlog
parent: S-0546
owner: alex
created: 2026-10-31T07:44:57Z
updated: 2026-10-31T07:45:30Z
transitions: []
stream: S-0546
tags: [flai]
touches: [design/conventions/code-quality.md, design/conventions/delegation.md, design/conventions/work-management.md, template/root/design/conventions/code-quality.md, template/root/design/conventions/delegation.md, template/root/design/conventions/work-management.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, template/CHANGELOG.md]
after: [T-2407, T-2411]
---
# T-2417 The conventions, their template copies, the user and operator docs, and the changelog say the close-out is a serve job the agent waits on

## Work

Where the conventions tell the story's agent to run `scripts/close-out.sh` itself, say instead that it asks for the close-out job and waits for `verified`. Change the files in `design/conventions/` first, then their copies in `template/root/design/conventions/`, in the same task. The places are `work-management.md`'s close-out rule, `delegation.md`'s close-out lines, and `code-quality.md`'s "run the whole suite yourself through the close-out". `design/conventions/testing.md`, which the story names, does not exist (see the plan's thread). In `docs/users/flai.md`, under "Verify a story before review", describe `flai verify --job`, the `verified` event, the answer from the record, and the close-out counts in `flai stats`. Regenerate `docs/users/flai-reference.md` with `make flai-reference`. Give the journal's new fields in `docs/operators/index.md`. Add the template's changes to `template/CHANGELOG.md`.

Waits for T-2407 and T-2411, whose flag, tool, and counts it describes. Shares no path with T-2415, so the two run together.

## Done when

- [ ] No convention or doc tells an agent to run the close-out script in a turn
- [ ] `make flai-reference` leaves `docs/users/flai-reference.md` unchanged after the commit
- [ ] `flai check --strict` and the markdown lint pass
- [ ] Story criterion 3 (the conventions and docs)

## Notes
