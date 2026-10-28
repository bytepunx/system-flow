---
id: D-0073
title: A flaiover test that runs about twenty flai processes hit vitest's 5 s timeout on a loaded host
kind: cause
sources: [I-0042, I-0185]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, testing, code, all]
paths: [flaiover/src, flaiover/src/lib/server/search.test.ts, flaiover/src/lib/server/writes.test.ts, flai/cmd, flai/internal/serve, flai/internal/hostapi, flai/internal/harness, flai/internal/config, flai/internal/manifest, flai/internal/itemedit, design/adrs, design/system, design/issues, docs/operators, docs/users]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0073 A flaiover test that runs about twenty flai processes hit vitest's 5 s timeout on a loaded host

A flaiover test that runs flai processes against the monorepo can pass alone and time out at vitest's 5 s default while other stories' close-outs load the host. The fix both times was a longer timeout on that test or block, with the measured figures in a comment, not a change to flai.

- I-0042, closed 2026-10-01: a test that ran about twenty flai processes hit the 5 s timeout. "Fixed by S-0116 (7510387): the test has a 30 s timeout."
- I-0185, closed 2026-10-26: `search.test.ts` 'finds a body phrase across wip and design' timed out. Each monorepo test runs a one-shot `flai hostapi` that builds the whole repository's search index: "1.8 to 2.0 s wall and 1.9 s of CPU for 4,845 markdown files, and 4.8 to 6.1 s each with 24 at once on a 24-core host". S-0501 gave the 'adrs on the monorepo' block a 30 s timeout, "as writes.test.ts did in S-0116". No unit test reproduces it, since it depends on host load. Whether the one-shot build should be cheaper was asked on TH-0654.

Sources: I-0042, I-0185.
