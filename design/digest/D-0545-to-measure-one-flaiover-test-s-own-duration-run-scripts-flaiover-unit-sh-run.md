---
id: D-0545
title: "To measure one flaiover test's own duration, run scripts/flaiover-unit.sh run --reporter=json -t on that test and read its duration, since flai test reports only a tier's time; pass checks still go through flai test"
kind: decision
sources: [S-0640, T-2805, S-0655, T-2838]
topics: [testing, flaiover]
paths: [scripts/flaiover-unit.sh]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0545 To measure one flaiover test's own duration, run scripts/flaiover-unit.sh run --reporter=json -t on that test and read its duration, since flai test reports only a tier's time; pass checks still go through flai test

S-0640's T-2805 had to show how long one vitest test took before and after its fix. From S-0640's Decisions (`wip/archive/agents/S-0640.md`):

```text
To measure one test's own duration, T-2805 ran `scripts/flaiover-unit.sh run --reporter=json -t` on that test and read only its `duration`. `flai test` gives a tier's duration, not a test's. This is a measurement the task asks for, not a pass check. Pass checks went through `flai test`.
```

S-0655 did the same for I-0330, citing S-0640 (`wip/archive/agents/S-0655.md`, Decisions): "A test's own duration is measured as S-0640 measured it: `scripts/flaiover-unit.sh run --reporter=json -t "<name>"`, reading only that test's `duration`. This is a measurement, not a pass check, and pass checks go through `flai test`." It went on to time the test phase by phase, which told it the time was the page's own work (D-0552).

So a later story that must time a single flaiover test uses the same command and reads the test's `duration` from the JSON report, several runs each side with the host load noted. It is never a substitute for `flai test` deciding pass or fail.

Sources: S-0640 (T-2805); S-0655 (T-2838).
