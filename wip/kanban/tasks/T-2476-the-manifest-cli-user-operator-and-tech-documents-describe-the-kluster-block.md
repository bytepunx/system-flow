---
id: T-2476
type: task
nature: feature
title: The manifest, CLI, user, operator, and tech documents describe the kluster block, the three commands, and the kluster host action
status: backlog
parent: S-0577
owner: alex
created: 2026-11-06T21:43:04Z
updated: 2026-11-06T21:43:04Z
transitions: []
stream: S-0577
tags: [kluster, docs]
touches: [design/system/project-manifest.md, design/system/flai-cli.md, design/system/kluster.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/operators/index.md, design/tech/kluster.md, design/tech/README.md]
after: [T-2473, T-2475]
---
# T-2476 The manifest, CLI, user, operator, and tech documents describe the kluster block, the three commands, and the kluster host action

## Work

It waits for T-2473 and T-2475, whose host action and commands it documents as built.

- `design/system/project-manifest.md`: the `kluster:` block, key by key, with an example, linking `kluster.md` and the ADR.
- `design/system/flai-cli.md`: `flai kluster up`, `down`, and `status`.
- `design/system/kluster.md`: fill in its section on the commands with what S-0577 built.
- `docs/users/flai.md`: the three commands, with examples. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- `docs/operators/settings.md`: a row for each `kluster:` key, so the settings-index test passes.
- `docs/operators/index.md`: a short entry for the `kluster` host action, off by default, turned on with `flai serve enable kluster` or the Features page. S-0579 adds what flai serve does with it.
- `design/tech/kluster.md`, new, with front matter: kluster, the version tested, where flai calls it, why, and how to install it. Index it in `design/tech/README.md`.

Bump `updated` on every document changed.

## Done when

- Each document above says what the story built, and the reference is regenerated.
- The settings-index test and `flai check --strict` pass.
- The markdown lint is clean on every path.

## Notes
