---
id: ADR-0052
title: "flai asks its questions a line at a time with its own prompt package, not with huh"
status: accepted
date: 2026-09-29
supersedes: []
superseded_by: []
refines: [ADR-0006]
topics: [cli, go]
---

# ADR-0052 flai asks its questions a line at a time with its own prompt package, not with huh

## Context

ADR-0006 chose Charm's huh and lipgloss for flai's prompts and tables. lipgloss was never used directly; huh asked flai's few questions: the template's variables and layout in `new` and `import`, whether to apply an import and move its folders, where a loose markdown file belongs, how to settle an upgrade conflict, and whether to cancel what is open under an item.

huh imports `bubbles/textarea`, which imports `atotto/clipboard`, and that package looks for clipboard programs on `PATH` when it loads. Every flai process pays it, prompt or none: 145 ms on the operator's host, whose `PATH` has 54 entries, 17 of them Windows folders under `/mnt`, against 8 ms with a short `PATH` (cause 5 of [server performance](../system/server-performance.md)). `flai serve` starts flai for writes and delegated reads, and agents and hooks run it constantly.

## Decision

flai asks its questions through `flai/internal/prompt`, a package of its own on the standard library: a value with a default and a check, yes or no with a default, or one of a few numbered options, each written as a line and answered with a line. huh and lipgloss are no longer dependencies. The rest of ADR-0006 stands.

## Consequences

- A flai process starts without searching `PATH`, and 25 modules leave `go.mod`.
- Questions are plain lines: no arrow-key selection, no form styling. Every one still has a flag or `--yes` that answers it without a terminal.
- Tests answer prompts by giving the app a prompter over a string, rather than a confirmation hook.

## Alternatives considered

- huh with fields that do not use `textarea`: huh imports it whatever fields are used, so the cost stays.
- A `replace` of `atotto/clipboard` with an empty module: `go install` refuses a module with `replace` directives.
- huh behind a build tag or a second binary: keeps 25 modules and a second build for six questions.
