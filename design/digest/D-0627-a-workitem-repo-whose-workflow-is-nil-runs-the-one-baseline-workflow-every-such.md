---
id: D-0627
title: A workitem Repo whose Workflow is nil runs the one baseline workflow every such Repo shares read-only, and a test opens a repository with another workflow through rules_test.go's repoWithWorkflow
kind: decision
sources: [S-0737]
topics: [flai, go, code]
paths: [flai/internal/workitem/repo.go, flai/internal/workitem/rules_test.go, flai/internal/workflow/baseline.go]
date: 2026-11-14
status: active
absorbed_by: ""
asked_again: []
cited: [S-0738, TH-1206, TH-1207, TH-1208, TH-1209, TH-1210, TH-1211, TH-1212, TH-1213, TH-1214, TH-1215, TH-1218, TH-1219, TH-1220]
---

# D-0627 A workitem Repo whose Workflow is nil runs the one baseline workflow every such Repo shares read-only, and a test opens a repository with another workflow through rules_test.go's repoWithWorkflow

S-0737 made the board's states and transition map data, read from a workflow (ADR-0301). Its narrative's Decisions record how a `Repo` gets its workflow:

> Repo.Workflow, nil for the baseline, is the hook that lets a test open a repository with another workflow. rules_test.go's repoWithWorkflow builds one for later stories' tests.
  
> The baseline is parsed at BaselinePath, flai/internal/workflow/baseline/board.md, so that a problem names the embedded file. Baseline() hands out a copy. workitem keeps one baseline read-only in a package variable, baselineFlow, and every Repo without a Workflow of its own shares it.

So a later E-0034 story that needs a board other than the baseline in a test calls `repoWithWorkflow(t, text)` in `flai/internal/workitem/rules_test.go` rather than writing a workflow file, and never changes `baselineFlow`: a code path that must alter a workflow takes its own copy from `workflow.Baseline()`. `design/system/flai-cli.md` says only that workitem reads the workflow a `Repo` runs, the baseline unless a test gives it another.

Source: S-0737's narrative, `wip/archive/agents/S-0737.md`, § Decisions.
