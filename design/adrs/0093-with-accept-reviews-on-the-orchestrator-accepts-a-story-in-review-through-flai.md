---
id: ADR-0093
title: "With accept_reviews on, the orchestrator accepts a story in review through flai accept --by orchestrator, only when the verifier passed at the branch head, every criterion is ticked, the diff stays within the touches, no thread is open, and its evidence covers every criterion"
status: accepted
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0032]
topics: [cli, dashboard, orchestration]
---

# ADR-0093 With accept_reviews on, the orchestrator accepts a story in review through flai accept --by orchestrator, only when the verifier passed at the branch head, every criterion is ticked, the diff stays within the touches, no thread is open, and its evidence covers every criterion

## Context

[ADR-0032](0032-accepting-a-story-merges-it-publishing-is-a-deliberate-batched-step-over.md) made acceptance one flow, `flai accept`, which merges a story's branch, moves it to done, archives it, and commits; the operator runs it, or an agent on the operator's word. `item_move` refuses every agent a move to done. S-0218 gave the orchestrator the permission `orchestration.permissions.accept_reviews`, off by default, and its guard lets the orchestrator run `flai accept` while it is on, but nothing says when the orchestrator may accept, what it must check first, or what it leaves behind for the operator to read.

Nothing records a verifier's verdict today: a story's close-out prints its result and stores nothing, and a story can change after the run that passed. The planner's plan for S-0221 (TH-0134) recommended that the orchestrator run its own verifier and hand flai the commit it verified; the operator resolved the plan without objection.

## Decision

With `orchestration.permissions.accept_reviews` on, the orchestrator accepts a story in review through `flai accept <S-nnnn> --by orchestrator`, and flai accepts it only when four conditions hold and the evidence covers every criterion.

1. **Who.** `--by orchestrator` is the orchestrator's acceptance, whoever runs it. flai refuses it unless `accept_reviews` is on in the project's manifest, naming the permission, and the guard refuses it as well. Under `FLAI_ROLE=orchestrate`, `flai accept` with any other `--by`, or without one, is refused by the guard and by flai. The done transition is recorded `by: orchestrator`. `item_move` still refuses every agent a move to done; its refusal to the orchestrator names `flai accept --by orchestrator` and the permission.
2. **The four conditions.** Each is a blocker of the acceptance preview when it fails, beside the blockers the preview already reports, and `flai accept` refuses on any blocker before anything is merged:
   - the verifier's run passed at the story branch's head: `--verified <commit>` names the commit the verifier checked, and flai refuses when it is missing or is not the head;
   - every acceptance criterion is ticked;
   - every file the branch changes is under the story's touches, read as a claim reads them;
   - no thread on the story or one of its tasks is open.
3. **The evidence.** `--evidence <file>`, or `-` for standard input, is markdown the orchestrator writes from its verifier's report: a `Verdict:` line, and one list item per acceptance criterion, `- <n>: <files>`, naming the files changed on the branch that meet criterion `<n>`. flai refuses an orchestrator's acceptance without evidence, and one whose evidence leaves a criterion with no file the branch changed: a story whose criteria the orchestrator cannot check against the diff is never accepted.
4. **What is recorded.** flai writes the evidence, with the verified commit, under `### Accepted by the orchestrator` in the story's `## Notes`, in the acceptance commit, and returns it in `--json`. The orchestrator logs the acceptance with `activity_log`, naming the story and the commit, as it logs every action. The dashboard's review and story pages name who accepted, and link an orchestrator's evidence.
5. **When it does not accept.** The orchestrator leaves the story in review and opens a thread on it saying what is missing: each blocker, and each criterion it could not check against the diff. It logs that decision.

## Consequences

- The operator can let accepted work flow without watching each review, and still read, on each story the orchestrator accepted, what was verified at which commit and which files meet which criterion.
- Acceptance stays one flow. The orchestrator's acceptance merges, archives, and commits exactly as the operator's; it never publishes, which is `publish`'s and S-0222's.
- A story with a thread open, an unticked criterion, or a file outside its touches waits for the operator, or for its agent to fix it, however the orchestrator judges it.
- A commit added after the verifier's run makes the acceptance refuse, so the orchestrator verifies again.
- The verdict is the orchestrator's word: flai checks that the commit verified is the head and that the evidence covers every criterion, not that the verifier ran.

## Alternatives considered

- **The close-out records a pass on the story, which flai checks.** A recorded verdict that flai reads would not rest on the orchestrator's word, but it widens the story into the close-out script and the template, and a pass recorded by the story's own agent is the run the operator would otherwise review. Declined on TH-0134.
- **The orchestrator moves the story to done with `item_move`.** It would skip the preview and leave no evidence. `item_move` stays the operator's no.
- **Evidence in the activity log alone.** The log is one line per action; the Notes keep the evidence with the story, in its archive, where the operator reads it.
