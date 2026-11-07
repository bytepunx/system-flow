---
id: T-2533
type: task
nature: remediation
title: "The four I-0247 tests pass together under -race -count=10 with GOMAXPROCS=2 beside a CPU load, and I-0247 is closed with what fixed each"
status: done
parent: S-0569
owner: alex
created: 2026-11-07T03:48:10Z
updated: 2026-11-07T05:27:26Z
transitions:
  - to: ready
    at: 2026-11-07T05:23:18Z
    by: agent-S-0569
  - to: in-progress
    at: 2026-11-07T05:23:20Z
    by: agent-S-0569
  - to: done
    at: 2026-11-07T05:27:26Z
    by: agent-S-0569
stream: S-0569
tags: [flai, tests]
touches: [design/issues/I-0247-testservesaystheflaiisolderthantheproject-and.md, design/issues/summary.md]
after: [T-2522, T-2525, T-2528, T-2529]
usage:
  source: log
  seconds: 246
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 5651
      cache_read: 2669296
      cache_write: 89777
      cost: 1.1473
---
# T-2533 The four I-0247 tests pass together under -race -count=10 with GOMAXPROCS=2 beside a CPU load, and I-0247 is closed with what fixed each

## Work

Waits for T-2522, T-2525, T-2528, and T-2529: it runs the four tests they fix together and closes the issue with what each found. Run the four tests at once, `GOMAXPROCS=2 go test -race -count=10 -run '<the four names>' ./internal/serve ./internal/mcpserver ./cmd`, beside a CPU load like the close-outs met (load average 40 or more), and record the load and the result in the narrative. Then run the integration tier once with `flai test integration`, or `scripts/integration.sh`, and record its time.

Close I-0247 with `flai issue close I-0247 --reason`, naming for each test the cause its task found and its fix, and the stress run that passed. Tick story criterion 1 and 2 once each is verified.

## Done when

- [ ] The four tests pass ten times in a row together under `GOMAXPROCS=2 -race` beside a load of 40 or more, recorded in the narrative
- [ ] The integration tier passes
- [ ] I-0247 is closed with its reason, and `design/issues/summary.md` lists it closed
- [ ] Story criteria 1 and 2 are ticked

## Notes
