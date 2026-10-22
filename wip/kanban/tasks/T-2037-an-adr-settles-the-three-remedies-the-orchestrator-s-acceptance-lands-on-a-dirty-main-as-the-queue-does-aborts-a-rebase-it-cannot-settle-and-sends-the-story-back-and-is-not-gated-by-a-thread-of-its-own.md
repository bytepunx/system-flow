---
id: T-2037
type: task
nature: improvement
title: "An ADR settles the three remedies: the orchestrator's acceptance lands on a dirty main as the queue does, aborts a rebase it cannot settle and sends the story back, and is not gated by a thread of its own"
status: backlog
parent: S-0467
owner: alex
created: 2026-10-22T20:59:12Z
updated: 2026-10-22T20:59:12Z
transitions: []
stream: S-0467
tags: [flai]
touches: [design/adrs/0186-the-orchestrator-s-acceptance-lands-on-a-dirty-main-as-the-queue-does.md, design/adrs/README.md]
---
# T-2037 An ADR settles the three remedies: the orchestrator's acceptance lands on a dirty main as the queue does, aborts a rebase it cannot settle and sends the story back, and is not gated by a thread of its own

## Work

Write one ADR with `flai adr new`, refining ADR-0093, ADR-0175, and ADR-0178, that records the remedy for each cause I-0150 names. Read I-0150's instances first: TH-0387 (S-0341, a gate the orchestrator raised itself), TH-0389, TH-0391, and TH-0394 (S-0342 and S-0343 held on a dirty main), and TH-0397 (S-0342's acceptance stopped mid-rebase on `design/system/flai-cli.md`, a `claims.shared` path).

The decision, in three parts:

1. Dirty main. `flai accept --by orchestrator` without the merge queue lands as `flai accept --batch` does under ADR-0178: it leaves the main checkout's uncommitted files outside `wip/` as they are, out of the lint fix and out of the acceptance commit, which holds only the acceptance; `git merge --ff-only` moves main and the checkout together; and it refuses only when the story's change would overwrite one of those files, naming the files and the command that frees them, as a blocker the dry run reports too. An acceptance the orchestrator made itself that stopped before its commit (the done transition is `by: orchestrator`, and `flai accept` finds it done and archived but uncommitted) is finished by the orchestrator with `flai accept <S-nnnn> --by orchestrator`; the rule that the operator completes a resumed acceptance keeps to acceptances someone else made.
2. A rebase it cannot settle. When the rebase of the story branch in `flai accept --by orchestrator` stops on a conflict the merge steps cannot settle (the issue files, the ADR index, and a shared markdown file's sections), flai aborts the rebase as the merge queue's batch does, leaves the branch at its head and the worktree clean, sends the story back to in-progress by `flai` with the paths as the reason, logs the conflict as today, and opens one thread on the story naming the paths, so that `flai serve` starts the story's agent as it does for a story the queue sent back, and the agent runs `flai stream sync` and resolves it. An acceptance by hand keeps today's behaviour: the operator resolves in the worktree or aborts.
3. A gate of its own. A thread the orchestrator opened on a story in review, with no entry by anyone else, does not gate the story: the acceptance preview lists it as a note, not a blocker, the merge queue's eligibility leaves it out the same way, and the acceptance resolves it when the story is archived (ADR-0109). A thread anyone else opened or wrote on holds the story until it is resolved. The orchestrator accepts on flai's blockers and the acceptance criteria alone; a concern beyond them, such as a design change the story was not asked for, is raised on a thread and does not hold the story.

Give the ADR its context from I-0150 and the analysis, the decision, its consequences for the orchestrator's prompt, the convention, and the docs, and `topics: [cli, orchestration]`. Add its row to `design/adrs/README.md` as `flai adr new` does. The predicted file name is the story's touch; `flai task done` widens the touches to the name `flai adr new` gives it.

## Done when

- One ADR under `design/adrs/` records the three decisions above with their context and consequences, refines ADR-0093, ADR-0175, and ADR-0178, and is indexed in `design/adrs/README.md`.
- `flai check --strict` passes on the ADR and the index.

## Notes

Layer 1 of the plan: the code tasks wait for this decision. The planner proposed the three remedies from I-0150's instances; the plan's thread on S-0467 names the assumptions.
