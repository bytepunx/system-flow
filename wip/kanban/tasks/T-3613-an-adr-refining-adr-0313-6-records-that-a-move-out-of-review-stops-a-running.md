---
id: T-3613
type: task
nature: remediation
title: "An ADR refining ADR-0313 §6 records that a move out of review stops a running rebase and resolve job"
status: backlog
parent: S-0822
owner: alex
created: 2026-11-27T23:56:10Z
updated: 2026-11-27T23:56:10Z
transitions: []
stream: S-0822
tags: [adr, flai]
touches: [design/adrs/README.md]
---
# T-3613 An ADR refining ADR-0313 §6 records that a move out of review stops a running rebase and resolve job

## Work

Write the decision TH-1274 settled as `stop`, with `flai adr new`, as an ADR that refines [ADR-0313](../../../design/adrs/0313-rebase-and-resolve-runs-as-a-flai-serve-job-that-rebases-a-review-story-s.md). ADR-0313 §6 says "A story that left review while the job ran is left as it is"; the new ADR replaces that sentence and leaves the rest of ADR-0313 standing.

Decide and write:

- **When.** flai serve's verify pass reads, at each look, the state of each story whose rebase and resolve job runs. A story no longer in review, in any state, stops its job.
- **The stop.** The job's context is cancelled, which stops its process group, as a stop with flai serve does. A stop during the rebase step aborts a rebase the job started, so the branch is as it was, as ADR-0313 §3 does for a conflict. A rebase it found waiting is left as it found it.
- **The record.** The job ends with verdict `stopped`, `stopped_at` the step it was at (`rebase` or `close-out`), outcome `left review`, and a reason that names the story's state, who moved it, and when. It sends nothing back, opens no thread, and starts no agent. Its end is journalled and announces a `verified` event, as every job's end does.
- **The race.** A move made after the close-out ended and before `settle` ran is the case T-3520 records: outcome `left review`, with the command that moves the story back to review after a pass.
- **Why.** The job runs `scripts/close-out.sh` in the story's worktree, which commits what is outstanding and syncs, while the story's agent works in that worktree once the story is in progress. In I-0396 the job ran 33 minutes after the operator moved S-0817 out 43 seconds in.
- **A close-out job.** Say whether a close-out job, not a resolve job, is stopped the same way. Recommend leaving it as it is: its agent asked for it from the worktree and waits on its verdict.

Add the ADR to `design/adrs/README.md`. Leave ADR-0313's text as it is, since it is accepted.

This task waits for no other: TH-1274 already settled the decision, and the code tasks read the record's fields from this ADR.

## Done when

- The ADR is accepted, refines ADR-0313, and names when the job stops, how, and what its record holds.
- `design/adrs/README.md` lists it.
- `flai check --strict` and the markdown lint pass with `flai test` on the changed paths.

## Notes
