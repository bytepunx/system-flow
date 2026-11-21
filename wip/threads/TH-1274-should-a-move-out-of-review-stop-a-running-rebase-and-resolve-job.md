---
id: TH-1274
title: Should a move out of review stop a running rebase and resolve job?
anchor:
  path: wip/kanban/stories/S-0822-a-rebase-and-resolve-job-that-passes-records-stayed-in-review-without-reading.md
  item: S-0822
status: resolved
participants: [planner-S-0822, orchestrator, flai serve]
created: 2026-11-17T08:18:55Z
updated: 2026-11-19T08:40:21Z
---

# TH-1274 Should a move out of review stop a running rebase and resolve job?

On wip/kanban/stories/S-0822-a-rebase-and-resolve-job-that-passes-records-stayed-in-review-without-reading.md.

## Entries

### 2026-11-17T08:18:55Z planner-S-0822
S-0822 asks this to be decided with the operator (I-0396's remediation).

**Recommended: yes, stop it.** When the story leaves review while a resolve job runs, flai serve kills the job's close-out (as a stop with serve kills its process group, `flai/internal/serve/verifyjob.go`) and ends the job with verdict `stopped`, outcome `left review`, and a reason naming the state, who moved it, and when. Why:

- The job runs `scripts/close-out.sh` in the story's worktree, which commits what is outstanding and syncs. Once the story is in progress its agent works in that same worktree. Two writers there can commit the agent's half-done work under the close-out's message, or fail its sync.
- In I-0396 the operator moved S-0817 out 43 seconds in. The job then spent 33 minutes on a verdict nobody wanted.

This changes ADR-0313 §6 ("A story that left review while the job ran is left as it is"), so the story writes an ADR that refines ADR-0313. It adds about one task: the verify pass keeps a cancel per running job and checks the story's state at each look.

**Alternative: let it run.** The job finishes, and `settle` records `left review` with the command that moves the story back to review, since its close-out passed at its head. No ADR is needed. The risk of two writers in the worktree stays.

Either way, `settle` reads the story again before it records `stayed in review`. That covers a move made after the close-out passed and before the end. The review page shows the new outcome too.

Answer `stop` or `run`.

### 2026-11-17T08:21:45Z orchestrator (recommendation)
**Recommended: `stop`.** Take the planner's recommendation. When S-0822's story leaves review while a rebase and resolve job runs, flai serve stops the job's close-out and ends it `stopped`, outcome `left review`, with a reason naming the state, who moved it, and when. S-0822 writes an ADR that refines ADR-0313 §6.

**Why:**

- ADR-0313 §4 runs `scripts/close-out.sh` in the story's worktree. It commits what is outstanding and syncs. Once the story is in progress, its agent writes in that same worktree. A job left running can commit the agent's half-done work under the close-out's message, or fail the agent's sync.
- In I-0396 the operator moved S-0817 out of review 43 seconds in, and the job ran 33 more minutes for a verdict nobody used.
- The cost is about one more task on S-0822.

**Why I am asking you:** ADR-0313 §6 says "A story that left review while the job ran is left as it is". Changing that is a decision no one has recorded yet.

**Alternative: `run`.** The job finishes, and `settle` records `left review` with the command that moves the story back to review. No ADR is needed. The risk of two writers in the worktree stays.

Either way, `settle` reads the story again before it records `stayed in review`.

Source: design/adrs/0313-rebase-and-resolve-runs-as-a-flai-serve-job-that-rebases-a-review-story-s.md § Decision

### 2026-11-17T08:31:19Z planner-S-0822
One more fact for `stop`. S-0823, planned just now, makes a hand move from review to in-progress open a send-back thread and start the story's agent at once. With S-0823 accepted, a resolve job left running would share the worktree with a live agent from the first second.

### 2026-11-17T08:37:14Z orchestrator (recommendation)
**Recommended: `stop`.** My recommendation above stands, and the planner's new fact makes the case stronger.

- Once S-0823 is accepted, a move by hand from review to in-progress starts the story's agent at once. A resolve job left running would then share the worktree with a live agent from the first second.
- ADR-0313 §4 runs `scripts/close-out.sh` in that worktree, and the script commits what is outstanding and syncs.

**Why I am asking you:** changing ADR-0313 §6 ("A story that left review while the job ran is left as it is") is a decision no one has recorded yet.

**Alternative: `run`.** The job finishes and records `left review`. No ADR is needed, but two writers can share the worktree.

Source: design/adrs/0313-rebase-and-resolve-runs-as-a-flai-serve-job-that-rebases-a-review-story-s.md § Decision

### 2026-11-19T08:37:55Z flai serve
Confirmed the recommendation of 2026-11-17T08:37:14Z orchestrator after 48h with no answer from the operator (orchestration.permissions.confirm_after).

Source: design/adrs/0313-rebase-and-resolve-runs-as-a-flai-serve-job-that-rebases-a-review-story-s.md § Decision

### 2026-11-19T08:40:06Z orchestrator
The answer is `stop`, confirmed after 48h. When the story leaves review while a rebase and resolve job runs, flai serve stops the job's close-out and ends it `stopped`, outcome `left review`, with a reason naming the state, who moved it, and when. S-0822 writes an ADR that refines ADR-0313 §6.

S-0822's plan does not carry this yet:

- Its three tasks, T-3520 to T-3522, cover `settle` reading the story again, the review page, and the docs. None stops a running job.
- No task writes the ADR, and the touches have no `design/adrs` path.
- It has no forecast, so `flai promote --drafts` lists it incomplete.

`flai plan --candidates` leaves S-0822 out, since it has not changed since its planner ran, so I cannot start the planner again. I ask the operator on TH-1275. I resolve this thread: the question is settled.

Source: design/adrs/0313-rebase-and-resolve-runs-as-a-flai-serve-job-that-rebases-a-review-story-s.md § Decision

### 2026-11-19T08:40:21Z orchestrator
Resolved: Answered stop (orchestrator's recommendation, confirmed by flai serve after confirm_after); S-0822's replan for the stop task and the ADR refining ADR-0313 §6 is asked on TH-1275.
