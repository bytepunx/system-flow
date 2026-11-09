---
id: D-0530
title: "flai accept's preflight markdown lint takes the preview's DirtyFiles, which keep git's quoting of a file name with a tab, newline, or double quote, so such a file gets no lint at preflight: accepted as an edge"
kind: decision
sources: [S-0609, T-2658]
topics: [flai, cli, markdown]
paths: [flai/cmd/accept_markdown.go, flai/internal/preview/accept.go, flai/cmd/accept.go]
date: 2026-11-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0530 flai accept's preflight markdown lint takes the preview's DirtyFiles, which keep git's quoting of a file name with a tab, newline, or double quote, so such a file gets no lint at preflight: accepted as an edge

S-0609 made `flai accept` lint each set of files once. Its preflight lint over the main checkout's uncommitted files now takes the preview's lists, read once from `git status --porcelain` (`design/system/server-performance.md`, cause 19).

Those lists do not unescape a name git quotes. S-0609's narrative records the gap and the choice: "T-2658: the preflight lint takes the preview's `DirtyFiles` lists, which do not unescape a name git quotes (tab, newline, `"`); such a file would get no lint at preflight. Accepted as an edge no markdown file here has" (wip/archive/agents/S-0609.md, Decisions).

A story that finds a markdown file missed by the preflight lint, or that adds file names with such characters, fixes the unescaping in the preview's lists rather than working around it in the lint.

Source: S-0609, T-2658.
