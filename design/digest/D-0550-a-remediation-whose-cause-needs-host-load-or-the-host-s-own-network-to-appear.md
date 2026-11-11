---
id: D-0550
title: A remediation whose cause needs host load or the host's own network to appear meets a criterion's 'test that reproduces it' with a fake that stands in for the host or with before-and-after timings, and says so against the criterion
kind: decision
sources: [S-0640, I-0319, S-0645, I-0321, S-0655, I-0330, S-0650, I-0325, S-0649, I-0324, S-0652, I-0329, I-0331, S-0684, T-2915, I-0234, S-0703, T-3047, I-0267]
topics: [testing, go, flaiover]
paths: [design/conventions/code-quality.md, flaiover/src/lib/components, flai/cmd/dashboard_porthold_test.go, flaiover/src/routes/design/design.svelte.test.ts, flai/cmd/components_test.go, flai/internal/channel/channel_test.go, flai/internal/channel/channeltest/dashboard.go, flai/internal/serve/candidates_test.go]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: [S-0669, T-3072]
---

# D-0550 A remediation whose cause needs host load or the host's own network to appear meets a criterion's 'test that reproduces it' with a fake that stands in for the host or with before-and-after timings, and says so against the criterion

`code-quality.md` asks that a fixed bug get a test that reproduces it, and a remediation story's first criterion usually asks for one "where one fits". Remediations accepted on 2026-11-10 found no deterministic test could reproduce the load or the host itself, and said what stood in.

S-0640 fixed I-0319, a flaiover props-panel test that timed out at 5,000 ms under load. From its Decisions (`wip/archive/agents/S-0640.md`):

```text
Criterion 1 asks for a test that reproduces the cause "where one fits". The timeout needs a host load of 50 or more, so no deterministic test reproduces it. The before-and-after timings are the evidence.
```

Its timings were 2,247, 2,447, and 2,174 ms before and 674, 672, and 682 ms after, at similar loads, measured as D-0545 says. S-0655, which fixed I-0330, said the same (`wip/archive/agents/S-0655.md`, Decisions): "Criterion 1 asks for a test that reproduces the cause "where one fits". None fits: the cause is host load, which a unit test cannot make reliably." (D-0552)

S-0645 fixed I-0321, `TestTCPAnswersAListenerAndNotAClosedOne`, which failed when WSL2's mirrored networking let a Windows listener answer a port free in Linux (D-0444). From its Decisions (`wip/archive/agents/S-0645.md`):

```text
Criterion 1's "test that reproduces it" does not fit: no test can make the host's relay listen. The accepting fake stands in for it.
```

Where the race itself, not the load, can be widened by a seam, a deterministic test does fit. S-0650 reproduced I-0325 with `FLAI_FAKE_COMPONENTS_VITE_DELAY`, a delay before the fake vite listens (D-0551), and S-0649 reproduced I-0324 with `lockedBuffer.holdAnswered`, which holds back the log line the test waits for (D-0543). Each failed under the old code and passes with the fix.

So a later remediation whose cause depends on the host's load or network first looks for a seam that widens the race without load; else it gives a fake that behaves as the host did, through a seam such as S-0645's `dialFunc`, or measured timings before and after. It states which against the criterion, rather than leave the criterion unmet. Where a load can be applied by hand, D-0492 says how to reproduce a Go test's failure under it.

S-0652 found seams of the same kind for I-0329 and I-0331: a 1 s start delay for the fake vite, and a 500 ms listen delay, each failing the old test on every run (D-0551), and `FLAI_SLOW_REQUEST=1ns` to force the slow-request log that load had caused (D-0558).

Where the timings themselves are the criterion and the host's load moves between runs, the baseline is main measured alternately with the branch under the same load. S-0684, which sped up `flai/internal/mcpserver`'s tests for I-0234, found the load at 53 to 156 against the issue's 40 to 57, so it "timed main's package alternately with the branch under the same load, as the baseline" (`wip/archive/agents/S-0684.md`, Decisions), and recorded each run with its load in T-2915's Notes and I-0234's close reason. It ran those timed `-race` runs itself, not through a sub-agent, as a measurement the story asked for rather than test feedback (D-0571).

S-0703, which fixed I-0267, took the seam route and, for its repeated runs, the host as it was rather than an added load. From its Decisions (`wip/archive/agents/S-0703.md`):

```text
T-3047 ran its repeated runs at the host's own load of about 4, not under synthetic load: other agents' close-outs share this host, and the reproduction test proves the mechanism without load by shortening the ping.
```

Its reproduction test pings every 300 ms and failed 3 of 3 runs against the old test dashboard (D-0597); the four affected tests then passed 20 of 20 runs with `GOMAXPROCS=1` and `-race`. So where a seam proves the mechanism, the repeated runs need no load added on a host other agents share.

Sources: S-0640 (`wip/archive/agents/S-0640.md`, Decisions); I-0319; S-0645 (`wip/archive/agents/S-0645.md`, Decisions); I-0321; S-0655 and I-0330; S-0650 and I-0325; S-0649 and I-0324; S-0652, I-0329, and I-0331; S-0684 (T-2915) and I-0234; S-0703 (`wip/archive/agents/S-0703.md`, Decisions), T-3047, and I-0267.
