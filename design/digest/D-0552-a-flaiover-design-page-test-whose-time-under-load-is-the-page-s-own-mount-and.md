---
id: D-0552
title: A flaiover design-page test whose time under load is the page's own mount and re-renders gets a 15 s timeout on its describe block, since fewer settles would save almost nothing (I-0330)
kind: cause
sources: [I-0330, S-0655, T-2838]
topics: [testing, flaiover, sveltekit, code]
paths: [flaiover/src/routes/design/design.svelte.test.ts, design/issues/I-0330-a-flaiover-design-page-test-takes-up-the-draft-another-tab-writes-times-out-at.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0552 A flaiover design-page test whose time under load is the page's own mount and re-renders gets a 15 s timeout on its describe block, since fewer settles would save almost nothing (I-0330)

I-0330 was closed on 2026-11-10 by S-0655 (T-2838). The design-page test "takes up the draft another tab writes", in `flaiover/src/routes/design/design.svelte.test.ts`, timed out at vitest's 5,000 ms under load and stopped the flaiover tier.

The cause, measured phase by phase: the test's time is the page's own mount and its re-renders into and out of Edit mode, 71 to 622 ms per settle, all in each settle's first timer turn. The other 19 turns take about 1 ms each. From S-0655's Decisions (`wip/archive/agents/S-0655.md`): "fewer settles would save about 60 ms of about 1,000. Stubbing the preview, the sheet, or the editor one at a time showed no component dominating, and the test asserts on the editors."

The fix gives the S-0383 `describe` block a 15 s timeout. "The timeout covers the whole block because every test in it mounts the page the same way (128 to 871 ms each in a whole-file run). 15 s is three times the default and about five times the 3.2 s worst case seen at a load of 58. It follows `flaiover/src/lib/server/search.test.ts`, whose `describe` has a timeout of its own."

This is the other branch of D-0544's rule. I-0319's props-panel test spent its time in per-message settles it could batch, so its timeout was not raised. Here the time is the page's own work, which the test cannot cut, so the timeout is raised for the block. A later flaiover test that times out under load is timed phase by phase first, as D-0545 says, to tell which case it is.

Sources: I-0330; S-0655 (`wip/archive/agents/S-0655.md`, Decisions); T-2838.
