---
id: T-2558
type: task
nature: remediation
title: "The tests I-0246 names pass together under -race -count=10 with GOMAXPROCS=2 beside a CPU load, and I-0246 is closed with what fixed each"
status: done
parent: S-0560
owner: alex
created: 2026-11-07T08:20:34Z
updated: 2026-11-07T09:53:04Z
transitions:
  - to: ready
    at: 2026-11-07T09:50:43Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T09:50:46Z
    by: system-flow
  - to: done
    at: 2026-11-07T09:53:04Z
    by: agent-S-0560
stream: S-0560
tags: [flai, tests]
touches: [design/issues/I-0246-testaheldwaitfailsonatruncatedstory-and-testservestartsadefinedroleonaneventando.md, design/issues/summary.md]
after: [T-2556, T-2557]
usage:
  source: log
  seconds: 137
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 24
      output: 95
      cache_read: 3211280
      cache_write: 12688
      cost: 1.155
---
# T-2558 The tests I-0246 names pass together under -race -count=10 with GOMAXPROCS=2 beside a CPU load, and I-0246 is closed with what fixed each

## Work

Waits for T-2556 and T-2557, and through T-2557 for T-2555: it runs what they fix together and closes the issue with what each found.

Run `TestAHeldWaitFailsOnATruncatedStory`, `TestServeStartsADefinedRoleOnAnEventAndOnASchedule`, `TestWaitForWorkPassesOverAHeldStory`, and `TestARepositoryIsOfferedImportedAndThenServed` at once, with `GOMAXPROCS=2 go test -race -count=10 -run '<the four names>' ./internal/serve ./internal/mcpserver`, from `flai/`, beside a CPU load of 40 or more, the load the instances met. Record the load and the result in the narrative. S-0569 fixed the first two tests' own races (T-2522, T-2528) after I-0246's last instance, so this run is also the first evidence on them since. Then run the integration tier once with `flai test integration`, or `scripts/integration.sh`, and record its time.

Close I-0246 with `flai issue close I-0246 --reason`, naming for each test its cause and its fix, S-0569's where that is the fix, the stress run that passed, and any issue T-2556 recorded. Leave I-0273 open: it is flai serve's own window, which S-0569's test waits out, and no test here depends on it. If a test still fails here, do not close the issue: reopen the task it belongs to and say what the run found. Tick story criteria 1 and 2 once each is verified.

## Done when

- [ ] The four tests pass ten times in a row together under `GOMAXPROCS=2 -race` beside a load of 40 or more, recorded in the narrative
- [ ] The integration tier passes
- [ ] I-0246 is closed with its reason, and `design/issues/summary.md` lists it closed
- [ ] Story criteria 1 and 2 are ticked

## Notes
