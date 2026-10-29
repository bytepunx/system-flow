---
id: D-0426
title: "flai/internal/context cannot import flai/internal/verify, since verify imports context through check, so a pack that needs a story's last flai verify result reads .flai-cache/verify/<S-nnnn>.json itself"
kind: decision
sources: [S-0537, T-2309]
topics: [flai, go, code, cli, orchestration]
paths: [flai/internal/context/reasons_figures.go, flai/internal/context, flai/internal/verify, flai/internal/check]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0426 flai/internal/context cannot import flai/internal/verify, since verify imports context through check, so a pack that needs a story's last flai verify result reads .flai-cache/verify/<S-nnnn>.json itself

S-0537's T-2309 built the orchestrator's `review` reason, which loads a story in review with its last `flai verify` result. It found that `flai/internal/context` cannot call `verify.LastReport`: the `verify` package imports `context` through `check`, so the call would be an import cycle.

S-0537's narrative records the decision: "`verify` imports `context` through `check`, so the review loader reads `.flai-cache/verify/<S>.json` itself rather than calling `verify.LastReport`." `loadReview` in `flai/internal/context/reasons_figures.go` reads the report from the project's cache dir (`.flai-cache/verify/<story>.json`, as `flai verify --last --json` prints it) and says so in its comment. When no file is there, the pack states that the story has no stored result.

The same seam came up for the `ready` reason's text. `flai board` and `flai promote --candidates` render their text in `cmd` with unexported functions, so `context` copies their format and leaves out the board's merge-queue line. Release figures are `release.EvaluateRepo`'s JSON, which the MCP tool `release_evaluate` gives too.

A later change that has `context` read another package's stored result should check this import direction first. It may need to read the stored file, or move the reader to a package below both.

Source: S-0537 (wip/archive/agents/S-0537.md, Decisions, "T-2309's seams").
