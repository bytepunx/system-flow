---
id: T-1537
type: task
nature: feature
title: The flai design and user docs and the reference describe the design system document, flai design, design.tokens, and design.get
status: backlog
parent: S-0381
owner: alex
created: 2026-10-13T17:46:50Z
updated: 2026-10-13T17:46:50Z
transitions: []
stream: S-0381
tags: [cli]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1534, T-1535, T-1536]
---
# T-1537 The flai design and user docs and the reference describe the design system document, flai design, design.tokens, and design.get

## Work

Describe in `design/system/flai-cli.md` the design system document's format, the `designsystem` package, `flai design init`, `show`, and `css`, the rule `design.tokens`, and the host method `design.get`; in `docs/users/flai.md` how a user starts and checks a design system; and regenerate `docs/users/flai-reference.md` with `make flai-reference`. It waits for T-1534, T-1535, and T-1536, whose behaviour it describes.

## Done when

- The three documents describe what the story built, and the markdown lint passes on them.

## Notes

- Layer 4.
