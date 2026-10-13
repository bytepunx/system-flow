---
id: T-1562
type: task
nature: feature
title: flai prime --story briefs the design system and its convention for a story that reaches a sub-project it governs, and leaves both out where there is none
status: backlog
parent: S-0386
owner: alex
created: 2026-10-13T17:49:14Z
updated: 2026-10-13T17:49:14Z
transitions: []
stream: S-0386
tags: [cli]
touches: [flai/internal/context/design.go, flai/internal/context/design_test.go]
---
# T-1562 flai prime --story briefs the design system and its convention for a story that reaches a sub-project it governs, and leaves both out where there is none

## Work

In `flai/internal/context/design.go`, a story whose topics, tags, or touches reach a sub-project the design system document's `governs` names, or whose topics include `design`, gets the document briefed and the convention `design-system.md` in its pack; a project with no design system document gets neither, as before. Tests cover a story on a governed sub-project, one on another, one with topic `design`, and a project without the document. It waits for nothing; it reads the `governs` key T-1561 writes.

## Done when

- `go test ./internal/context/` passes with the four cases.

## Notes

- Layer 1. If the convention's selection lives outside `design.go`, widen the touches to that file.
