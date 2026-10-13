---
id: T-1563
type: task
nature: feature
title: "This repository's design records, and its conventions gain, a design-system convention: build UI from the tokens and propose changes on the story's branch"
status: backlog
parent: S-0386
owner: alex
created: 2026-10-13T17:49:17Z
updated: 2026-10-13T17:49:17Z
transitions: []
stream: S-0386
tags: [conventions]
touches: [design/system/conventions.md, design/conventions/design-system.md, design/conventions/README.md]
---
# T-1563 This repository's design records, and its conventions gain, a design-system convention: build UI from the tokens and propose changes on the story's branch

## Work

Record the convention in `design/system/conventions.md` first, then write `design/conventions/design-system.md` in the conventions' shape, with `roles: [story, explore, verify]`, its Rules: build UI from the design system's tokens through the stylesheet `flai design css --write` keeps; add no colour, size, font, radius, or shadow outside them; make a change the UI needs to the design system document on the story's branch, so the operator previews it on the Design page before accepting; record the reason for a design decision in the document's guidance. Add it to `design/conventions/README.md` in its order. It waits for nothing.

## Done when

- The convention and the design say the same, the README lists it, and `flai check --strict` and the markdown lint pass.

## Notes

- Layer 1.
