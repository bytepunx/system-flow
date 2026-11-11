---
id: T-3111
type: task
nature: improvement
title: A brief's JSON carries the count of its other reasons in place of the also list
status: backlog
parent: S-0699
owner: alex
created: 2026-11-11T09:04:16Z
updated: 2026-11-11T09:04:16Z
transitions: []
stream: S-0699
tags: [flai, cli, go]
touches: [flai/internal/context/items.go, flai/internal/context/context.go, flai/internal/context/design_test.go, design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md]
---
# T-3111 A brief's JSON carries the count of its other reasons in place of the also list

## Work

The operator confirmed on TH-1023 (Q2 yes) that a brief's JSON carries the count of its other reasons in place of its `also` list. The `also` lists of briefs took about 130 KB of S-0425's and S-0586's packs, 39% of their brief bytes (I-0261). Every brief still loads: Q1 was no, so ADR-0241 stands and this task cuts no brief.

`Selection.Items` (`flai/internal/context/items.go`) sets `Also: b.also` on each design and ADR brief. For an item with step `briefed`, leave `Also` empty and set a new count field in its place, such as `More int` encoded as `more,omitempty`, holding `len(b.also)`. An item that loads, step `named`, `ranked`, `digest`, or `reason`, keeps its full `also` list, and so does a briefed document that ends up loaded whole, which gains the brief's reasons in `also` as now.

`Item.briefReason` (`flai/internal/context/context.go`) prints `<first reason>, and n more` from `len(it.Also)`. Compute it from the count instead, so the text pack prints the same words. Find every other reader of a brief's `Also` with `grep` before you change the field; flaiover reads no pack.

Update the documents that describe the field:

- `design/system/agent-context.md`, under `## Fitting the pack to a budget`: the line "the JSON keeps every reason" becomes that a brief's JSON keeps the count and a loaded item keeps the list.
- `design/system/flai-cli.md`: the pack's JSON `items` field list gains the count field.
- `docs/users/flai.md`, under `flai prime`: the sentence on how a brief counts its later reasons says the JSON carries the count too.

This task waits for no other task. It shares no file with T-2992, so the two run together in the first layer. T-2996 waits for it, since it measures what it saved and also edits `design/system/agent-context.md`.

## Done when

- A test in `flai/internal/context/design_test.go` gives a document three reasons that brief it, and its item then has no `also` and a count of 2, while a loaded item with three reasons keeps its two later ones in `also`.
- `TestPrintedBriefsNameTheirFirstReasonAndCountTheRest` passes with its expected text unchanged, built from the count.
- `design/system/agent-context.md`, `design/system/flai-cli.md`, and `docs/users/flai.md` describe the count, and `updated` is bumped on each.
- `flai test` passes on the changed files, and `flai check --strict` and the markdown lint pass on them.

## Notes
