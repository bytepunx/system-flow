---
id: I-0417
title: config.Save writes the host configuration world-readable while its notify URL may carry a secret
class: impression
status: open
count: 1
cost: 5m
first_reported: 2026-11-27T22:47:02Z
last_reported: 2026-11-27T22:47:02Z
updated: 2026-11-27T22:47:02Z
usage:
  strategic:
    - kind: analyzer
      seconds: 85
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 55
          output: 5532
          cache_read: 391176
          cache_write: 29513
          cost: 0.9652
        - model: claude-opus-5-5
          input: 29
          output: 4664
          cache_read: 1264948
          cache_write: 28114
          cost: 0.487
---

# I-0417 config.Save writes the host configuration world-readable while its notify URL may carry a secret

## Description
config.Save writes the host configuration world-readable while its notify URL may carry a secret

## Instances

### 2026-11-27T22:47:02Z
Report: design/analysis/2026-11-27-all.md.
Read in the code today; no occurrence observed. On a shared host a local user other than the operator can read the reminder channel's URL, a webhook address that may hold its secret.

## Impact
flai/internal/config/config.go:1280 config.Save: os.Chmod(tmpName, 0o644) before the rename over ~/.flai/config.json; the Notify type at lines 260-265 says of its URL "It may carry a secret in its path or query, so flai logs its host". The dashboard token, the kubeconfig, and the agent configuration are written 0600. safety.md: secrets never go into logs or reports; use the project's secret mechanism. No instance observed.

## Remediation

Found by the analysis in [2026-11-27-all.md](../analysis/2026-11-27-all.md).
