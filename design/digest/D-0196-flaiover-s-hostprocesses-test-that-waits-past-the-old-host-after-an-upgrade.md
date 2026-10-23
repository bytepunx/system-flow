---
id: D-0196
title: flaiover's HostProcesses test that waits past the old host after an upgrade fails under machine load
kind: cause
sources: [I-0053]
topics: [testing, flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src/lib/components/HostProcesses.svelte.test.ts, design/issues/I-0053-flaiover-s-hostprocesses-test-that-waits-past-the-old-host-after-an-upgrade-fails-under-machine-load.md, design/issues/summary.md, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md]
date: 2026-10-06
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0196 flaiover's HostProcesses test that waits past the old host after an upgrade fails under machine load

I-0053 was closed on 2026-10-06.

The cause:

```text
flaiover's HostProcesses test that waits past the old host after an upgrade fails under machine load
```

The fix:

```text
S-0264: HostProcesses' reconnect tests wait for their outcome (vi.waitFor, 4 s) with a 3 s give-up instead of asserting after a fixed 40 ms against a 50 ms give-up; a new test that answers the old host's poll after 80 ms fails at the old timings and passes now, and the file passed 6 of 6 runs while the Go tiers loaded the machine (load 10-13).
```

Source: I-0053.
