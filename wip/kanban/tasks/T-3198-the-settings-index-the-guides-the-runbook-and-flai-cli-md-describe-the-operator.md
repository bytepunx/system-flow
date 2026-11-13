---
id: T-3198
type: task
nature: feature
title: The settings index, the guides, the runbook, and flai-cli.md describe the operator's build command and the switch that turns the automatic image job off
status: backlog
parent: S-0759
owner: alex
created: 2026-11-12T16:05:58Z
updated: 2026-11-12T16:05:58Z
transitions: []
stream: S-0759
tags: [cli]
touches: [docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/runbooks/agent-container.md, design/system/flai-cli.md]
after: [T-3191, T-3194]
---
# T-3198 The settings index, the guides, the runbook, and flai-cli.md describe the operator's build command and the switch that turns the automatic image job off

## Work

- `docs/operators/settings.md`: a row under § flai configuration for each new key, near `agent.harnesses.<name>.images.<name>` (line 61), and the flags under § Flags (line 365 on). `TestSettingsIndex` fails on a key or flag with no row.
- `docs/users/flai.md` § In a container: both settings.
- `docs/users/flai-reference.md`: `make flai-reference`, never by hand.
- `docs/operators/runbooks/agent-container.md`: a section, "Your own images", on naming a build command, its variables, its log, and turning the automatic job off, beside `--image`.
- `design/system/flai-cli.md`: mark the substitution built in the section S-0733 wrote, and bump `updated`.

It waits for T-3191 and T-3194, whose flags and behaviour it describes.

## Done when

- The five files describe both settings as built.
- `TestSettingsIndex`, `TestFlagIndexIsCurrent`, `TestReferenceIsCurrent`, `flai check --strict`, and the markdown lint pass.

## Notes
