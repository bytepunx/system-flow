---
id: T-1066
type: task
nature: remediation
title: The flai docs and the CLI design say flai new, import, and upgrade take the newest template tag and when upgrade asks
status: backlog
parent: S-0301
owner: alex
created: 2026-10-06T22:50:45Z
updated: 2026-10-06T22:50:45Z
transitions: []
stream: S-0301
tags: [docs, template]
touches: [docs/users/flai.md, docs/users/flai-reference.md, docs/contributors/template.md, design/system/flai-cli.md]
after: [T-1064, T-1065]
---
# T-1066 The flai docs and the CLI design say flai new, import, and upgrade take the newest template tag and when upgrade asks

## Work

User-facing behaviour changes, so the docs change in the same story.

- `docs/users/flai.md`, "Upgrade to a newer template": with no `--ref`, upgrade takes the newest version tag; `--ref` wins and is recorded; the question asked when `system-flow.yaml`, the lock and the newest tag disagree, with its choices; and the refusal without a terminal and the `--ref` it names. Say the same for `flai new` and `flai import` where the guide creates a project.
- `docs/contributors/template.md`, "Fork it and point flai at the fork" and "Versioning": a branch ref is fetched again before use, so drop the warning that it goes stale. A config ref left empty means the newest `v<version>` tag, so tag each release with `flai template push --tag`.
- `design/system/flai-cli.md`: the rows for `flai new` and `flai upgrade` in the command table.
- Regenerate `docs/users/flai-reference.md` with `scripts/flai-reference.sh` from the help text T-1064 and T-1065 changed. Never edit it by hand.
- This task waits for T-1064 and T-1065, because it documents and regenerates what they ship.

## Done when

- [ ] The three hand-written documents describe the default ref, `--ref` precedence, the question and its choices, and the run without a terminal, as the ADR from T-1062 says.
- [ ] `docs/users/flai-reference.md` matches `scripts/flai-reference.sh`'s output, and `flai/cmd/reference_test.go` passes.
- [ ] `flai check --strict` and the markdown lint pass.

## Notes

Drafted by the planner.
