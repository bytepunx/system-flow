---
id: D-0550
title: A remediation whose cause needs host load or the host's own network to appear meets a criterion's 'test that reproduces it' with a fake that stands in for the host or with before-and-after timings, and says so against the criterion
kind: decision
sources: [S-0640, I-0319, S-0645, I-0321]
topics: [testing, go, flaiover]
paths: [design/conventions/code-quality.md, flaiover/src/lib/components, flai/cmd/dashboard_porthold_test.go]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0550 A remediation whose cause needs host load or the host's own network to appear meets a criterion's 'test that reproduces it' with a fake that stands in for the host or with before-and-after timings, and says so against the criterion

`code-quality.md` asks that a fixed bug get a test that reproduces it, and a remediation story's first criterion usually asks for one "where one fits". Two remediations accepted on 2026-11-10 found no deterministic test could, and said what stood in.

S-0640 fixed I-0319, a flaiover props-panel test that timed out at 5,000 ms under load. From its Decisions (`wip/archive/agents/S-0640.md`):

```text
Criterion 1 asks for a test that reproduces the cause "where one fits". The timeout needs a host load of 50 or more, so no deterministic test reproduces it. The before-and-after timings are the evidence.
```

Its timings were 2,247, 2,447, and 2,174 ms before and 674, 672, and 682 ms after, at similar loads, measured as D-0545 says.

S-0645 fixed I-0321, `TestTCPAnswersAListenerAndNotAClosedOne`, which failed when WSL2's mirrored networking let a Windows listener answer a port free in Linux (D-0444). From its Decisions (`wip/archive/agents/S-0645.md`):

```text
Criterion 1's "test that reproduces it" does not fit: no test can make the host's relay listen. The accepting fake stands in for it.
```

So a later remediation whose cause depends on the host's load or network gives as its reproduction a fake that behaves as the host did, through a seam such as S-0645's `dialFunc`, or else measured timings before and after. It states which against the criterion, rather than leave the criterion unmet. Where a load can be applied by hand, D-0492 says how to reproduce a Go test's failure under it.

Sources: S-0640 (`wip/archive/agents/S-0640.md`, Decisions); I-0319; S-0645 (`wip/archive/agents/S-0645.md`, Decisions); I-0321.
