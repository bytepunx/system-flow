---
id: TH-1274
title: Should a move out of review stop a running rebase and resolve job?
anchor:
  path: wip/kanban/stories/S-0822-a-rebase-and-resolve-job-that-passes-records-stayed-in-review-without-reading.md
  item: S-0822
status: open
participants: [planner-S-0822]
created: 2026-11-17T08:18:55Z
updated: 2026-11-17T08:18:55Z
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
