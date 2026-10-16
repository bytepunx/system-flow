---
id: T-1681
type: task
nature: improvement
title: The manifest design, the workflow, the settings index, the CLI docs, and the template's changelog describe the queue tier key
status: backlog
parent: S-0405
owner: alex
created: 2026-10-16T05:03:48Z
updated: 2026-10-16T05:03:48Z
transitions: []
stream: S-0405
tags: [cli]
touches: [design/system/project-manifest.md, design/system/workflow.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md, template/template.yaml, template/CHANGELOG.md]
after: [T-1674, T-1666]
---
# T-1681 The manifest design, the workflow, the settings index, the CLI docs, and the template's changelog describe the queue tier key

## Work

- `design/system/project-manifest.md` and `docs/operators/settings.md`: the tier key `queue`.
- `design/system/workflow.md` § Definition of done: the close-out leaves the marked tiers to the queue while it runs.
- `design/system/flai-cli.md` and `docs/users/flai.md`: `flai verify`'s step for tiers left to the queue.
- Bump `template.yaml`'s version and name the convention's change, the script's, and the flai it needs in `template/CHANGELOG.md`.
- Waits for T-1674 and T-1666, whose behaviour it describes.

## Done when

The documents say the above, the settings test passes, and the markdown lint and `flai check --strict` pass.

## Notes
