---
id: ADR-0230
title: "A ready story whose claim overlaps a story in progress is held only when a trial merge of that story's branch at its base conflicts, and after names a story whose output it reads, with the reason in its Notes"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0046, ADR-0134]
---

# ADR-0230 A ready story whose claim overlaps a story in progress is held only when a trial merge of that story's branch at its base conflicts, and after names a story whose output it reads, with the reason in its Notes

## Context

[ADR-0046](0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md) holds a ready story while its claim overlaps the claim of an open story, and while a story it names in `after:` is not done. [ADR-0096](0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md) narrowed the overlap: only a story in progress holds, an overlap inside `claims.shared` never holds, and a story's tasks narrow its folder touches. [ADR-0134](0134-flai-serve-asks-the-agent-of-a-story-in-progress-about-each-ready-story-it.md) lets the holding agent share the paths. [ADR-0121](0121-flai-tells-two-stories-of-a-trial-merge-conflict-or-a-grown-overlap-with-a.md) tells two stories of a conflict that the trial merge at `flai stream sync` finds.

A claim says that a story may change a path. It does not say that two changes conflict. Each story works in its own worktree, and two stories that change one file conflict only where their changes meet. S-0551 measured the fortnight to 2026-10-30, with an in-progress limit of 5:

| Stories in progress at once | Share of the fortnight |
|-----------------------------|------------------------|
| 0 | 42 percent |
| 1 | 29 percent |
| 2 | 26 percent |
| 3 | 2 percent |
| 4 or 5 | never |

The mean while anything ran was 1.6. The journal recorded 124 overlap holds and 113 `after` holds, against 62 real conflicts. The overlap holds were on `flai/internal/check/check.go`, the strategic-agents design and convention, and the manifest and dashboard designs. The `after` holds were chains the planner made, S-0417 to S-0422 each waiting on the one before, as [I-0152](../issues/I-0152-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md) noted. The ready wait at p85 was 5.4 hours.

People who share one remote merge optimistically and resolve the rare conflict. ADR-0046 declined that as design 3, because a conflict found late wastes agent work and the operator sees nothing before work starts. `git merge-tree --write-tree`, in git 2.38 and later, tells whether two branches merge without touching any worktree, and `flai stream sync` already runs it.

## Decision

**A ready story whose claim overlaps the claim of a story in progress is held only when a trial merge of the in-progress story's branch at the ready story's base conflicts. `after:` names a story whose output this story reads, with the reason in the story's `## Notes`, and `flai check` advises of an `after` entry the Notes do not name.** This refines ADR-0046 and ADR-0134. The rest of each stands, and so does ADR-0096.

### The trial merge

- **When it runs.** The hold works out the overlapping pairs as today: the narrowed claims, the `/`-prefix rule, and `claims.shared` and the shares of ADR-0134, which drop a pair before any trial. For each pair left, the holds ask a trial-merge judge.
- **The base.** The ready story's base is its own branch, `story/S-nnnn`, when it has one, as a story sent back has. Otherwise it is the main branch.
- **The trial.** `git merge-tree --write-tree` of the in-progress story's branch at that base. The conflicting paths leave out the generated and merged files that the trial merge at sync leaves out.
- **Clean.** The pair holds nothing.
- **Conflict.** The pair holds, with the code `overlap`. The reason names the conflicting paths, and `flai serve`'s journal names them when it holds the story. For example: `held (overlap): a trial merge of story/S-0128 at main conflicts on flai/cmd/serve.go; starts when S-0128 moves to review, is cancelled, or is sent back, or when the trial is clean`.
- **Cannot run.** The pair holds as today, and the reason says why: the holding story has no branch yet, or git is older than 2.38. A story whose agent `flai serve` started in the same look has no branch yet, so one look still never starts two stories that overlap.
- **The judge.** It caches each verdict by the two branches' head commits. With no judge, the holds give today's verdict.

### Where it applies

`wait_for_work`, `flai story start`, `flai serve`'s launcher, `inbox`, `flai board`, and the orchestrator's promotion candidates all read the holds with the judge, so they agree. The `no-touches` hold is unchanged. flai serve's ask of ADR-0134 goes on for a story held on overlap, which is now one whose trial conflicts or could not run. The replay of past holds in `flai stats`, `items[].held_seconds` and `claims`, has no judge: it stays a record of the holds the claims give.

### After

- `after:` is honoured as ADR-0046 says: a ready story waits until each story it names is done.
- It names a story whose output this story reads: a file, a command, a schema, or a decision that the named story writes and this one builds on. It is never an order of preference, which the pull order sets, nor a guard against a conflict, which the trial merge judges.
- The story's `## Notes` give the reason, naming the story by its ID.
- `flai check` gives the advisory warning `wip.after-reason` on a story not done or cancelled whose `after:` names a story its Notes do not name. `--strict` passes over it, as it does over `wip.after-hub`. A task's `after` is left out: a task says why it waits in its `## Work`.
- The planner sets `after` on a story only so, with the reason in the Notes.

## Consequences

- An overlap whose branches merge no longer holds. Stories that touch one file in different places start together, so more lanes of the in-progress limit are used.
- A clean trial at start says only that the in-progress story's change lands cleanly on the ready story's base. The two stories' own changes can still meet. The trial merge at sync, the conversation of ADR-0121, and the notice at acceptance catch that, as before, after both stories have started.
- A hold now costs one `git merge-tree` per overlapping pair at each look whose heads changed, which takes milliseconds.
- On a host with git older than 2.38 every overlap holds as before, and the reason says so.
- `items[].held_seconds` and `claims.weeks[].held_seconds` count an overlap that a clean trial let start as held while the story stayed in ready. That story starts at the next look, so the excess is small.
- `after` chains need a reason to stand. `wip.after-reason` points at each one without, and the open chain S-0417 to S-0422 is reviewed against it.

## Alternatives considered

- **Optimistic only, never hold.** Declined, as in ADR-0046: a trial that conflicts already shows a conflict before work starts.
- **A trial merge of both stories' branches.** The ready story has no branch, and so no change, until it starts.
- **Hold on a conflict the claim alone predicts, as today.** This is what S-0551 measured: 124 holds against 62 real conflicts.
- **Drop `after` to stop the chains.** A story that reads another's output does need it on main first.
- **Refuse an `after` with no reason.** An advisory warns without stopping a close-out or an operator's edit, as `wip.after-hub` does.
