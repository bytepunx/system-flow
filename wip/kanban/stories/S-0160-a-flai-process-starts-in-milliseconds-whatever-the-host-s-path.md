---
id: S-0160
type: story
nature: remediation
title: A flai process starts in milliseconds whatever the host's PATH
status: in-progress
parent: E-0012
owner: alex
created: 2026-09-29T07:00:30Z
updated: 2026-09-29T19:38:50Z
transitions:
  - to: ready
    at: 2026-09-29T19:18:48Z
    by: alex
  - to: in-progress
    at: 2026-09-29T19:33:57Z
    by: agent-S-0160
tags: [cli]
topics: [server-side, back-end]
touches: [flai/cmd, flai/internal/prompt, flai/go.mod, flai/go.sum, design/tech/go-libraries.md, design/tech/README.md, design/adrs, design/system/server-performance.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 392
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 96
      output: 613
      cache_read: 4001253
      cache_write: 107282
      cost: 1.6334
---
# S-0160 A flai process starts in milliseconds whatever the host's PATH

## Goal

A flai process takes 145 ms to start on the operator's host before it does anything, and 8 ms with a short `PATH`. `atotto/clipboard`, which `huh` brings in through `bubbles/textarea`, looks for clipboard programs on `PATH` when its package loads, and the host's `PATH` has 54 entries, 17 of them Windows folders under `/mnt` that are slow to look in. Every write from the dashboard, every read it delegates, and every `flai` an agent or a hook runs pays it. Cause 5 of `design/system/server-performance.md`.

## Acceptance criteria
- [x] `GODEBUG=inittrace=1 flai version` on the operator's host shows no package init over 5 ms.
- [x] `flai version` takes under 30 ms there with the host's own `PATH`.
- [x] Prompts that need a terminal still work, and `design/tech/go-libraries.md` records the change of dependency.

## Tasks
- T-0569 flai asks its prompts through a line-based prompt package and no longer depends on huh
- T-0570 Measure flai version's start on this host and record it in the design

## Notes

Measured by S-0152: `init github.com/atotto/clipboard @3.3 ms, 145 ms clock`. `go mod why`: `flai/cmd` → `charmbracelet/huh` → `bubbles/textarea` → `atotto/clipboard`. Options: use huh fields that do not pull `textarea`, if there is a way; replace the prompt library for the few prompts flai has; or keep huh behind a separate binary or build. Choose in the story, with a line in `design/tech`.

Chosen: flai asks through `flai/internal/prompt`, its own line-based package on the standard library, and `huh` is gone with 24 modules it brought ([ADR-0052](../../../design/adrs/0052-flai-asks-its-questions-a-line-at-a-time-with-its-own-prompt-package-not-with.md), refining ADR-0006). huh imports `textarea` whatever fields are used; a `replace` of the clipboard module breaks `go install`; a build tag or second binary keeps 25 modules for six questions. Given up: arrow-key selection and form styling.

Verified on 2026-09-29 with the host's own `PATH` (54 entries, 17 under `/mnt`): `flai version` median 5.8 ms over 30 runs (main: 163 ms); the largest package init 0.6 ms; `flai new` and `flai move E-0001 cancelled` answered at a pseudo-terminal, a bad answer asked again.
