---
id: I-0257
title: A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host
class: impression
status: open
count: 2
cost: 5m
first_reported: 2026-11-02T19:06:40Z
last_reported: 2026-11-07T08:17:36Z
updated: 2026-11-10T08:37:15Z
usage:
  strategic:
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 33
          output: 3103
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13855
          cost: 0.278
---

# I-0257 A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host

## Description
A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host

## Instances

### 2026-11-02T19:06:40Z
Story: S-0509.
ADR-0247 mounts the main checkout read-write for the planner, the orchestrator, the analyzer, and the curator, with `.git/hooks` read-only over it. `.git/config` stays writable, because the orchestrator's `flai accept` deletes a merged branch's section from it. Keys such as `core.fsmonitor`, `core.sshCommand`, or `core.hooksPath` name commands that git runs on the host at its next command there. The guard refuses each role `git config`, but it does not see a write through the shell. That is no worse than a host run, but it is a way back out of the container. Remedy to weigh: mount `.git/config` read-only and move the branch deletion of an acceptance to the host (`flai serve`), or have `flai serve` check `.git/config` for command-naming keys after each strategic container run.

### 2026-11-07T08:17:36Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: found again as the strategic half of the container-to-host boundary; the story half is a new issue.

## Impact
flai/internal/harness/container.go:276-277 mounts the main checkout read-write with .git/hooks alone read-only; ADR-0247 lines 62 and 67 record that .git/config can name core.fsmonitor or core.sshCommand and that a shell write is not seen; no Go code checks hooksPath, fsmonitor, or sshCommand

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0688 remediates this issue, created from it at 2026-11-10T08:37:15Z.
