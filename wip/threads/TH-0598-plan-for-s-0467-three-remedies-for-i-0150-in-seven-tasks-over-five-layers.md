---
id: TH-0598
title: "Plan for S-0467: three remedies for I-0150 in seven tasks over five layers"
anchor:
  path: wip/kanban/stories/S-0467-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md
  item: S-0467
status: answered
participants: [planner-S-0467, orchestrator]
created: 2026-10-22T21:02:46Z
updated: 2026-10-22T21:03:00Z
---

# TH-0598 Plan for S-0467: three remedies for I-0150 in seven tasks over five layers

On wip/kanban/stories/S-0467-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md.

## Entries

### 2026-10-22T21:02:46Z planner-S-0467
Plan by planner-S-0467. The story is finalized, so its words stand; I added tags `flai` and `template` (flai check requires a component tag), topic `orchestration`, 34 file touches, a forecast of 2h (flai's 54m raised, reason under `### Planning`), and the cost of delay value 14400 USD a week from the operator's 96h input.

## The proposed solution

I-0150's three causes, from TH-0387, TH-0389, TH-0391, TH-0394, and TH-0397, and what fixes each:

1. Dirty main. `flai accept --by orchestrator` lands as the merge queue's landing does (ADR-0178): the main checkout's uncommitted files outside `wip/` stay as they are, out of the commit, and only a change that would overwrite one is refused, as a blocker the dry run lists. The orchestrator finishes an acceptance of its own that stopped before its commit (S-0338's case); a resumed acceptance someone else made stays the operator's.
2. A rebase it cannot settle. After ADR-0175's section merge, a stop `flai accept --by orchestrator` still cannot settle is aborted by flai, as the queue's batch does, and the story is sent back to in-progress by `flai` with the paths, with one thread, and `flai serve` starts its agent as for a story the queue sent back. The operator's acceptance by hand keeps today's behaviour.
3. A gate of its own. A thread the orchestrator opened that nobody else wrote on does not block its acceptance or the queue's eligibility; the acceptance resolves it at archive (ADR-0109). A thread anyone else opened or wrote on holds the story until resolved. The prompt and convention say: accept on flai's blockers and the criteria alone; a concern beyond them goes on a thread that does not hold the story.

## Tasks and layers

| Layer | Tasks | Touches |
|-------|-------|---------|
| 1 | T-2037 ADR | `design/adrs/` (new ADR, README) |
| 2 | T-2039 dirty main; T-2042 own thread | accept.go, preview, batch-dirty; orchestrator.go, eligible.go, threads.go |
| 3 | T-2047 abort and send back (after T-2039: both change `flai/cmd/accept.go`) | branch.go, storygit/sync.go, serve/agents.go |
| 4 | T-2049 prompt, agent definition, convention, template; T-2052 design and user docs | harness.go, `.claude/agents`, conventions; strategic-agents.md, workflow.md, flai-cli.md, docs/users/flai.md |
| 5 | T-2053 close I-0150, tick the criteria | I-0150, summary.md |

## Assumptions

- One ADR covers all three decisions, refining ADR-0093, ADR-0175, and ADR-0178; its file name is predicted as 0186 and `flai task done` widens the touch to the real name.
- The abort-and-send-back applies to the orchestrator's acceptance only; an acceptance by hand still leaves the stop in the worktree for the operator, and the batch is unchanged.
- "Its own gate" is a thread with every entry the orchestrator's. One reply by the operator or an agent makes it hold again, so the operator can still keep a story in review by writing on it.
- flai aborts the rebase itself, so `flai guard` need not grant the orchestrator `git rebase --abort`.
- No new flag, so `docs/users/flai-reference.md` is untouched.

Nothing to split, merge, or drop. If you would rather the abort apply to every acceptance, or the own-thread rule be stricter (the orchestrator must resolve its thread first), say so and the story's agent adjusts T-2047 or T-2042.

### 2026-10-22T21:02:55Z orchestrator
Approved under plan_backlog_stories, as a plan: seven tasks in five layers, with one ADR first. T-2047 waits on T-2039 because both change `flai/cmd/accept.go`. The 34 file touches stand, as does the 2h forecast raised from flai's 54m.

The three remedies are decisions not yet recorded, and they widen what my own acceptances may do. They are the operator's, so a recommendation follows and the thread stays open. Until it is answered, the story's agent can draft the ADR (T-2037) and should build nothing past it.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-22T21:03:00Z orchestrator (recommendation)
This story changes the rules on my own acceptances, so weigh my view with that in mind. Recommended for the operator: take remedies 1 and 2 as proposed, and take remedy 3 only in the stricter form.

1. **Dirty main: take it.** It does what the merge queue already does under ADR-0178. Uncommitted files stay out of the commit, and only a change that would overwrite one is refused. It widens nothing the queue does not already do.
2. **A rebase it cannot settle: take it.** Sending the story back with the paths starts its agent, which is the party that can resolve the conflict. Today the story sits in review until you step in.
3. **My own thread not gating: take the stricter alternative.** The rule would be: I must resolve my own thread before the acceptance, and a thread anyone else wrote on still holds the story. Under the proposed form, a concern I raise could never stop a story. TH-0589 is the case in point. It held S-0483 because the operator's instruction on TH-0578 had not been carried out, and that hold is what brought it to you. A thread I can resolve myself still ends the "gate it raised itself" wait, and it keeps a concern visible as a decision rather than letting it lapse at archive.

The value is 14,400 USD a week, from your 96h input, which puts this story near the top of anything ordered by value.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md
