---
id: D-0544
title: A flaiover component test that times out under load is fixed by settling once after a burst of harness messages, not after each, and not by raising its timeout (I-0319)
kind: cause
sources: [I-0319, S-0640]
topics: [testing, flaiover, sveltekit]
paths: [flaiover/src/routes/design/components/components.svelte.test.ts]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0544 A flaiover component test that times out under load is fixed by settling once after a burst of harness messages, not after each, and not by raising its timeout (I-0319)

I-0319 was closed on 2026-11-10 by S-0640 (T-2805). The props-panel test in `components.svelte.test.ts`, "lists the calls of the component's function props, and clears them on another pick", hit vitest's 5,000 ms timeout on a loaded host and stopped the flaiover tier.

The cause, from I-0319's Remediation: "the test settled after each of the 55 harness messages it posted, about 1,240 setTimeout(0) turns in one test, so a loaded host pushed it past 5,000 ms."

The fix, from S-0640's Decisions (`wip/archive/agents/S-0640.md`): "The fix posts the 55-call burst with a new `post` helper that does not settle, then settles once. `harness` is now `post` plus `settle`. Raising the timeout was rejected: it hides the cost." The burst assertion now checks all 50 entries, newest first, which is the check that a burst lands whole.

The effect: "Before: 2,247, 2,447, and 2,174 ms at a load of 15 to 17 on 24 cores. After: 674, 672, and 682 ms at a load of 12 to 19." No deterministic test reproduces the timeout, since it needs a load of 50 or more, so the before-and-after timings were the evidence.

A later flaiover test that times out under load is looked at for per-step settles or waits it can batch before its timeout is touched. A sibling, `design.svelte.test.ts`, timing out while mounting the design page, is a different cause, I-0330 (S-0655).

Sources: I-0319, S-0640.
