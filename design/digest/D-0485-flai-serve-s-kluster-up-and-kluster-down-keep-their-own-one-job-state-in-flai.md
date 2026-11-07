---
id: D-0485
title: flai serve's kluster.up and kluster.down keep their own one-job state in flai/internal/hostapi/kluster.go and do not take the deploy jobs' lock, .flai-cache/kluster/job.lock, so the two are not joined yet
kind: decision
sources: [S-0578, S-0580]
topics: [kluster, deployment, serve, flai, go, code]
paths: [flai/internal/hostapi/kluster.go, flai/internal/kluster/deploy.go, design/system/kluster.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0485 flai serve's kluster.up and kluster.down keep their own one-job state in flai/internal/hostapi/kluster.go and do not take the deploy jobs' lock, .flai-cache/kluster/job.lock, so the two are not joined yet

Two stories of E-0032 each gave the project's kluster a one-job-at-a-time rule, apart.

- S-0580 built the host API's `kluster.up` and `kluster.down` for the dashboard's Klusters page. flai serve runs those jobs in its own process, one per project at a time, and keeps the running and last job in memory (`flai/internal/hostapi/kluster.go`; `design/system/kluster.md` § The jobs).
- S-0578 built `deploy`, `teardown`, and `promote`, which take `.flai-cache/kluster/job.lock` and name themselves in `job.json`; a job that finds it held is refused with `kluster.ErrBusy` (`design/system/kluster.md` § The lock).

Neither takes the other's: an up or a down started from the dashboard does not hold the job lock, and a deploy does not see an up running. S-0580 said so to S-0578 on MS-0330:

```text
S-0580's `kluster.up` and `kluster.down` keep their own one-job-per-project state in flai serve (`flai/internal/hostapi/kluster.go`). They do not use your job lock yet, so S-0579 or S-0581 may want to join the two.
```

S-0578 recorded it in its `## Decisions` and left the join to S-0579 (flai serve's deploy jobs) or S-0581 (the Deployments page). Whichever story runs deploys from flai serve beside the Klusters page's up and down decides whether to join them.

Sources: S-0578 (wip/archive/agents/S-0578.md), S-0580, MS-0330.
