---
id: TH-0862
title: Findings 10 and 11 of the 2026-11-07 analysis link I-0284 where I-0292 and I-0293 were meant
anchor:
  path: design/analysis/2026-11-07-all.md
status: resolved
participants: [orchestrator, alex]
created: 2026-11-07T08:37:48Z
updated: 2026-11-10T09:30:10Z
---

# TH-0862 Findings 10 and 11 of the 2026-11-07 analysis link I-0284 where I-0292 and I-0293 were meant

On design/analysis/2026-11-07-all.md.

## Entries

### 2026-11-07T08:37:48Z orchestrator
**Recommended:** fix the two issue links in `design/analysis/2026-11-07-all.md`. Finding 10 should link I-0292 and Finding 11 should link I-0293.

**What is wrong:** the verifier found this while I accepted S-0584. Two lines end with a garbled link:

```text
Finding 10: Issue: [I-0284](../issues/I-0284-...md)0.
Finding 11 (line 212): Issue: [I-0284](../issues/I-0284-...md)2.
```

Neither I-0292 nor I-0293 is named anywhere else in the report. `flai issue list --report latest` still lists both, so S-0584's criteria held, and I accepted it.

**Why it is yours:** only the analyzer writes under `design/analysis`, and I edit no documents. You can fix it by hand or ask the analyzer to.

### 2026-11-10T09:30:10Z alex
Resolved.
