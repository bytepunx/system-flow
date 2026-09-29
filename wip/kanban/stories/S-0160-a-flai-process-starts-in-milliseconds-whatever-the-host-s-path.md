---
id: S-0160
type: story
nature: remediation
title: A flai process starts in milliseconds whatever the host's PATH
status: backlog
parent: E-0012
owner: alex
created: 2026-09-29T07:00:30Z
updated: 2026-09-29T07:00:30Z
transitions: []
tags: [cli]
topics: [server-side, back-end]
touches: [flai/cmd, flai/go.mod, flai/go.sum, design/tech/go-libraries.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0160 A flai process starts in milliseconds whatever the host's PATH

## Goal

A flai process takes 145 ms to start on the operator's host before it does anything, and 8 ms with a short `PATH`. `atotto/clipboard`, which `huh` brings in through `bubbles/textarea`, looks for clipboard programs on `PATH` when its package loads, and the host's `PATH` has 54 entries, 17 of them Windows folders under `/mnt` that are slow to look in. Every write from the dashboard, every read it delegates, and every `flai` an agent or a hook runs pays it. Cause 5 of `design/system/server-performance.md`.

## Acceptance criteria
- [ ] `GODEBUG=inittrace=1 flai version` on the operator's host shows no package init over 5 ms.
- [ ] `flai version` takes under 30 ms there with the host's own `PATH`.
- [ ] Prompts that need a terminal still work, and `design/tech/go-libraries.md` records the change of dependency.

## Tasks

## Notes

Measured by S-0152: `init github.com/atotto/clipboard @3.3 ms, 145 ms clock`. `go mod why`: `flai/cmd` → `charmbracelet/huh` → `bubbles/textarea` → `atotto/clipboard`. Options: use huh fields that do not pull `textarea`, if there is a way; replace the prompt library for the few prompts flai has; or keep huh behind a separate binary or build. Choose in the story, with a line in `design/tech`.
