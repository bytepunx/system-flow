---
id: I-0295
title: The close-out job and the manifest's test commands run the story worktree's scripts on the host, so a container harness does not keep an agent's code off the host
class: impression
status: open
count: 1
first_reported: 2026-11-07T08:17:37Z
last_reported: 2026-11-07T08:17:37Z
updated: 2026-11-07T23:40:34Z
usage:
  strategic:
    - kind: analyzer
      seconds: 41
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
---

# I-0295 The close-out job and the manifest's test commands run the story worktree's scripts on the host, so a container harness does not keep an agent's code off the host

## Description
The close-out job and the manifest's test commands run the story worktree's scripts on the host, so a container harness does not keep an agent's code off the host

## Instances

### 2026-11-07T08:17:37Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: by design the project's tests run somewhere; the risk is the boundary being read wider than it is. Hygiene beside it: docedit/file.go:51-81 refuses paths by string and os.Stat follows symlinks, where hostapi/docs.go:145-149 resolves real paths.

## Impact
flai/internal/serve/verifyjob.go:290-291 runs sh scripts/close-out.sh <story> from the worktree inside flai serve (ADR-0229); flai/internal/serve/checks.go:155 and :204 run system-flow.yaml tests[].command and checks on the host; the worktree is the container's read-write mount (ADR-0253); .claude/settings.json is in the same mount (container.go:276-277); nothing in strategic-agents.md § In a container or docs/operators says the host is not protected from the worktree's scripts

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0628 remediates this issue, created from it at 2026-11-07T23:40:34Z.
