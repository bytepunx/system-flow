---
id: D-0314
title: "I-0149's fix: one host upgrade, plus docs and a thread that names a published release: Recommendation: run `flai host upgrade` once, and reply \"upgraded\" here"
kind: answer
sources: [TH-0461]
topics: [release, host, mcp, flai, cli, go, code, all]
paths: [docs/operators/index.md, docs/operators/runbooks/update.md, design/system/flai-cli.md, design/issues/I-0149-the-host-s-installed-flai-stays-behind-the-project-s-releases-after-i-0147-closed-so-the-mcp-task-done-still-refuses.md, design/issues/summary.md, design/adrs/0161-flai-serve-s-upgrade-thread-names-a-flai-release-only-once-it-is-published-and.md, design/adrs/README.md, flai/cmd/serve.go, flai/cmd/serve_actions.go, flai/internal/serve/serve.go, flai/internal/serve/upgrade.go, flai/internal/serve/upgrade_test.go, docs/operators/runbooks/migrate.md]
date: 2026-10-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0314 I-0149's fix: one host upgrade, plus docs and a thread that names a published release: Recommendation: run `flai host upgrade` once, and reply "upgraded" here

The operator, alex, answered TH-0461 at 2026-10-15T22:53:57Z.

What was asked:

```text
**Cause.** The host runs flai 1.48.0 (`~/.flai/bin/flai`, `flai host` PID 26855, `flai serve` 26878). The fixes for this exist but shipped later: ADR-0158 (serve checks every minute, then upgrades or asks on a thread) in 1.50.1, and ADR-0160 (an outdated MCP `task_done` closes the task with the installed flai) in 1.51.1. A 1.48.0 host has neither, so nothing upgrades it and nothing asks you. No code change on main can reach a binary that is already running. flai 1.52.0 is published now.
```

The answer:

```text
**Recommendation:** run `flai host upgrade` once, and reply "upgraded" here. Keep 2b in S-0391. Choose `auto-upgrade` on or off as on TH-0460. Both threads ask for the same upgrade, so one answer settles both.

## Why keep 2b

- **It is the same cause.** With `auto-upgrade` off, ADR-0158's thread can name a tag whose binaries are not published yet, about five minutes early. Then `flai host upgrade --version X` refuses, and the host stays behind: that is I-0149 again, through the very path meant to prevent it.
- **It is small and tested.** One check, naming the newest published release, recorded as a new ADR refining ADR-0158, with a test.
- **It belongs to criterion 1,** "the cause I-0149 describes no longer occurs", better than to a separate story.

## Why this needs you

- The upgrade changes your machine. No agent may do it (`safety.md`).
- 2b adds code and an ADR beyond the planner's no-code plan on TH-0460. Its agent asked whether you want it.

Source: design/adrs/0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md § Decision
```

Source: TH-0461.
