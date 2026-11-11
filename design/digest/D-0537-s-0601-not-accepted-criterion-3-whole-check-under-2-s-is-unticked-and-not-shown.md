---
id: D-0537
title: "S-0601 not accepted: criterion 3 (whole check under 2 s) is unticked and not shown: I reworded criterion 3 and ticked it"
kind: answer
sources: [TH-0910, S-0601]
topics: [server-side, back-end, flai, cli, go, code, all]
paths: [flai/internal/check/check.go, flai/internal/components/props.go, flai/internal/check/concurrent_test.go, flai/internal/components/propscache.go, flai/internal/components/propscache_test.go, flai/internal/components/props.mjs, design/system/server-performance.md, design/system/component-library.md, design/issues/I-0312-a-whole-flai-check-takes-3-3-to-4-9-s-on-a-loaded-host-after-s-0601-not-under-2.md, design/issues/summary.md]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: [S-0677, TH-0961]
---

# D-0537 S-0601 not accepted: criterion 3 (whole check under 2 s) is unticked and not shown: I reworded criterion 3 and ticked it

The operator, alex, answered TH-0910 at 2026-11-09T08:19:31Z.

What was asked:

```text
Recommended: reword criterion 3 to what the story showed, or accept it as it stands with criterion 3 left to I-0312. If you reword it, I accept on my next run.
```

The answer:

```text
I reworded criterion 3 and ticked it
```

What the rewording was, from S-0601's `## Decisions` (`wip/archive/agents/S-0601.md`):

```text
2026-11-09, the operator's decision: criterion 3 no longer asks for under 2 s. It was changed to under 5 s, then to halving a warm whole check against main on the same host, and ticked. S-0601's 3.3 to 4.9 s against main's 7.9 to 11.7 s, measured taking turns, meets it.
```

So S-0601's criterion 3 asks for a warm whole `flai check --strict` at half of main's on the same host, not under 2 s. The 2 s target stays open as I-0312, which names the item list, `workItems`, and the digest rule as the remaining critical path. The orchestrator accepted S-0601 on 2026-11-09. This is the rule D-0532 records for S-0603: a performance criterion out of reach on a loaded host is reworded to the relative fall measured against main, and the absolute target goes to an issue.

Sources: TH-0910, S-0601.
