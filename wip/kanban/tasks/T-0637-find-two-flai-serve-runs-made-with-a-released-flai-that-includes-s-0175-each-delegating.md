---
id: T-0637
type: task
nature: research
title: Find two flai serve runs made with a released flai that includes S-0175, each delegating
status: done
parent: S-0188
owner: arobson
created: 2026-10-01T08:31:38Z
updated: 2026-10-01T09:01:26Z
transitions:
  - to: ready
    at: 2026-10-01T08:31:52Z
    by: agent-S-0188
  - to: in-progress
    at: 2026-10-01T08:31:52Z
    by: agent-S-0188
  - to: done
    at: 2026-10-01T09:01:26Z
    by: agent-S-0188
stream: S-0188
tags: []
touches: [design/system/agent-context.md]
usage:
  source: log
  seconds: 531
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 63
      output: 8583
      cache_read: 2059935
      cache_write: 49562
      cost: 0.9804
---
# T-0637 Find two flai serve runs made with a released flai that includes S-0175, each delegating

## Work

Confirm the installed flai includes S-0175 (`flai version` against the release tag that carries it), then pick two story runs in `~/.flai/serve/agents/` that started after it was installed and whose logs show an `Agent` call to the explorer or the verifier.

## Done when

Two runs are named, with the flai version each ran under and the sub-agents each started.

## Notes
- flai 1.26.4 (commit 861c9c3, tag `flai/v1.26.4`, built 2026-10-01T08:35:11Z) includes S-0175 (31265ff is its ancestor). `flai serve` restarted on it at 08:39:23Z.
- S-0185, started 08:48:00Z, delegated to the `verifier`; S-0184, started 08:54:08Z, to the `explorer`. Both `claude -p` command lines carry the delegation paragraph ("hand noisy work to sub-agents").
- S-0183, started 08:32:10Z under 1.26.3 with the definitions and `delegation.md` on main but not the prompt, also delegated once, to the `verifier`: a third point, convention only.
