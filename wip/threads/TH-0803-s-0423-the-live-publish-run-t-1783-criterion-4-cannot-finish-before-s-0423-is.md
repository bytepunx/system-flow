---
id: TH-0803
title: "S-0423: the live publish run (T-1783, criterion 4) cannot finish before S-0423 is accepted; move it to a follow-up story?"
anchor:
  path: wip/kanban/stories/S-0423-the-template-ships-a-baseline-publish-formula-that-fetches-evaluates-orchestration-release-and-publishes-what-is-pending-the-orchestrator-starts-it-in-place-of-calling-release-publish-and-it-is-run-once-here.md
  item: S-0423
status: resolved
participants: [agent-S-0423, orchestrator, alex]
created: 2026-11-02T04:12:14Z
updated: 2026-11-02T04:15:15Z
---

# TH-0803 S-0423: the live publish run (T-1783, criterion 4) cannot finish before S-0423 is accepted; move it to a follow-up story?

On wip/kanban/stories/S-0423-the-template-ships-a-baseline-publish-formula-that-fetches-evaluates-orchestration-release-and-publishes-what-is-pending-the-orchestrator-starts-it-in-place-of-calling-release-publish-and-it-is-run-once-here.md.

## Entries

### 2026-11-02T04:12:14Z agent-S-0423
T-1783 and criterion 4 ask for the `publish` formula to be run once here, with its instance accepted, before S-0423 goes to review. That can't happen while S-0423 is still open.

## Why

- The formula's step 2, `flai release --evaluate --require-met`, and the `push` gate in `flai release --pending` exist only on `story/S-0423`. The instance's agent runs its main-checkout steps (`in: main`) with the flai built from main (`scripts/flai.sh`) or the installed release. Neither has them until S-0423 is accepted, so step 2 would fail on an unknown flag.
- S-0423 can't move to review while T-1783 is open, and the instance can't run correctly until S-0423 is accepted.
- Accepting the instance and pushing the release are yours anyway (ADR-0067).

## Question: where does the live run go?

1. **Recommended:** I cancel T-1783 here and remove criterion 4 from S-0423. I write a draft story in E-0024 with `after: S-0423` that holds T-1783's Work and criterion 4's words. The epic still gets its one run before it is accepted. You finalize the draft.
2. Alternative: keep T-1783. S-0423 then stays in progress, blocked, until you accept it some other way. Under the definition of done that is not possible, so I don't recommend it.

## What I'm doing meanwhile

- I'm working the other seven tasks to TH-0583's answers.
- I'm not bumping `template/template.yaml`'s version, though T-1780's Work asks for it. `flai release` bumps it at publish (D-0456).

### 2026-11-02T04:13:18Z orchestrator (recommendation)
**Recommended:** take option 1.

- S-0423's agent cancels T-1783 and removes criterion 4 from S-0423.
- It writes a draft story in E-0024 with `after: S-0423`. The draft holds T-1783's Work and criterion 4's words, and its Notes say why it waits on S-0423.
- S-0423's Notes record the move.
- The epic still gets its one live run before it is accepted, as S-0423's goal asks.

**Why:** you have answered this before (D-0422, TH-0717, confirmed 2026-10-29). Some criteria need a change that exists only on the story branch to run on the host or the main checkout. Those criteria cannot be verified before review. Such a run moves to a follow-up story after the story. S-0536's day-long run moved this way, and so did S-0537's criterion 5 (D-0423). Here, `flai release --evaluate --require-met` and the `push` gate in `flai release --pending` exist only on `story/S-0423`. The instance's `in: main` steps cannot run them until S-0423 is accepted. If those steps run the installed flai rather than `scripts/flai.sh`, they also have to wait until S-0423 is published.

**Why I recommend rather than answer:** this takes a criterion out of a finalized story and adds a story. That changes scope, so it is your call. Finalizing the draft is yours too.

**Alternative (option 2):** keep T-1783. I don't recommend it. S-0423 could never meet the definition of done, so it would stay in progress.

Source: design/digest/D-0422-plan-for-s-0536-eight-tasks-in-five-layers-and-two-questions-on-criterion-6.md

### 2026-11-02T04:13:28Z alex
Confirmed the recommendation of 2026-11-02T04:13:18Z orchestrator.

Source: design/digest/D-0422-plan-for-s-0536-eight-tasks-in-five-layers-and-two-questions-on-criterion-6.md

### 2026-11-02T04:14:24Z orchestrator
You confirmed option 1, so that is the answer (D-0461). It is S-0423's agent's to carry out:

- cancel T-1783, and remove criterion 4 from S-0423;
- write a draft story in E-0024, with `after: S-0423`, that holds T-1783's Work and criterion 4's words, and Notes saying why it waits;
- record the move in S-0423's Notes.

The orchestrator has nothing more to do here. Finalizing the draft is yours.

Source: design/digest/D-0461-s-0423-the-live-publish-run-t-1783-criterion-4-cannot-finish-before-s-0423-is.md

### 2026-11-02T04:15:09Z agent-S-0423
Done. T-1783 is cancelled and S-0423's criterion 4 is removed; the docs criterion is now number 4. The live run is the draft S-0562 in E-0024, with `after: S-0423`, waiting for you to finalize it. S-0423's Notes and narrative record the move.

### 2026-11-02T04:15:15Z agent-S-0423
Resolved: Option 1 applied: T-1783 cancelled, criterion 4 removed, live run moved to the draft S-0562
