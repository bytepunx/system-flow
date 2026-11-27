---
id: I-0390
title: "Design sentences the code contradicts: ci.md's race detector on macOS, markdown.md's yaml in flaiover, metrics.md's WIP series and unnamed keys, flai-cli.md's missing dashboard logs row, the analyzer's MultiEdit"
class: defect
status: open
count: 2
cost: 15m
first_reported: 2026-11-16T06:16:30Z
last_reported: 2026-11-27T22:47:23Z
updated: 2026-11-27T22:47:23Z
usage:
  strategic:
    - kind: analyzer
      seconds: 185
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 102
          output: 12117
          cache_read: 740537
          cache_write: 63820
          cost: 2.0683
        - model: claude-haiku-5-5
          input: 36
          output: 14843
          cache_read: 2398931
          cache_write: 47318
          cost: 0.1557
        - model: claude-opus-5-5
          input: 29
          output: 4665
          cache_read: 1264948
          cache_write: 28114
          cost: 0.487
---

# I-0390 Design sentences the code contradicts: ci.md's race detector on macOS, markdown.md's yaml in flaiover, metrics.md's WIP series and unnamed keys, flai-cli.md's missing dashboard logs row, the analyzer's MultiEdit

## Description
Design sentences the code contradicts: ci.md's race detector on macOS, markdown.md's yaml in flaiover, metrics.md's WIP series and unnamed keys, flai-cli.md's missing dashboard logs row, the analyzer's MultiEdit

## Instances

### 2026-11-16T06:16:30Z
Report: design/analysis/2026-11-16-all.md.
Analysis 2026-11-16: design/tech/ci.md:37 says the race detector runs on Linux and macOS, .github/workflows/flai.yml:38-41 runs it on ubuntu only; design/tech/markdown.md:13 lists yaml 2.9 in flaiover, which flaiover/package.json does not name; design/system/metrics.md:42 defines WIP as a series where metrics.go:179 emits one count, and § Aggregates names no JSON key for summary, throughput, cfd, burnup, aging, or groups; design/system/flai-cli.md § Commands has no row for flai dashboard logs (flai/cmd/dashboard.go:1038); design/system/strategic-agents.md:940 lists the analyzer's file tools without MultiEdit (flai/internal/guard/claudecode.go:39).

### 2026-11-27T22:47:23Z
Report: design/analysis/2026-11-27-all.md.
Analysis of 2026-11-27 adds nine sentences: two the code has outgrown (project-manifest.md's "Nothing acts on pull-request yet", strategic-agents.md's "No page shows the curator's runs yet"), two manifest keys, one config key and two stored entries, four command rows, and three flaiover versions the lockfile has moved past. The first five instances stand. I-0401 no longer reproduces at ea7d124c6: TestPublishFormulaUnderJudgement passes and design/workflows/formulas holds no publish formula.

## Impact
project-manifest.md line 379 "Nothing acts on pull-request yet" against flai/cmd/accept_remote.go acceptThroughRemote (ADR-0170); strategic-agents.md line 1135 "No page shows the curator's runs yet" against flaiover/src/routes/workflow/[role]/+page.svelte; manifest keys checks and dashboard.bind (flai/internal/manifest/manifest.go) absent from project-manifest.md; config key dashboard.no_restart and stored entries components and mcp absent from flai-cli.md § Configuration; flai-cli.md rows omit flai serve components max-servers, flai formula run --ready, MCP tools subprojects and subproject_edit, flai kluster deployments --with-limit; design/tech/sveltekit.md 2.63/5.56/1.60 against pnpm-lock.yaml 2.70.3/5.57.0/1.63.0.

## Remediation

Found by the analysis in [2026-11-16-all.md](../analysis/2026-11-16-all.md).

Found by the analysis in [2026-11-27-all.md](../analysis/2026-11-27-all.md).
