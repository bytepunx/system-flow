---
id: I-0398
title: A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted
class: efficiency
status: open
count: 1
cost: 22m
first_reported: 2026-11-17T07:13:20Z
last_reported: 2026-11-17T07:13:20Z
updated: 2026-11-17T07:13:20Z
---

# I-0398 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted

## Description
A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted

## Instances

### 2026-11-17T07:13:20Z
Job S-0817-20261117T042837Z ran three rounds of the close-out. Its rebase was clean in 3 seconds. Round 1 passed at 04:39:33Z, round 2 at about 04:49:41Z, round 3 at 05:02:04Z; between them the acceptances of S-0742 and S-0753 moved main, each bringing in 36 or 37 paths that select every tier, so each round ran integration again, 10 to 11 minutes. The page showed the job running for 33 minutes. Had main moved a third time the job would have ended at its round limit with nothing to show.

## Impact
22 minutes of a 33 minute job were rounds two and three; with several stories in review each acceptance restarts every other story's resolve job, so review throughput falls as the review queue grows

- time_lost_per_cycle: 1h

## Remediation

Decide, in an ADR refining ADR-0313 and ADR-0156, what a rebase and resolve job owes when main moves under it, and build it. Options to weigh: the job passes at the base its last round ran on and leaves the rest to the acceptance's own `verify_behind_main` gate (ADR-0291, ADR-0303), which sends back only the tiers the new paths overlap; the later rounds apply ADR-0303's overlap rule in place of every tier main's paths select; or flai serve holds acceptances and resolve jobs in one order, so that main does not move under a job it is about to accept. ADR-0312's acceptance suites (S-0808) bear on the same cost. The job's record and the review page say which round the job is on and why it started another.

Story S-0824 remediates this issue, created from it at 2026-11-17T07:13:20Z.
