---
id: TH-0787
title: "S-0422's after S-0421: keep it and state the reason in its Notes"
anchor:
  path: wip/kanban/stories/S-0422-a-story-s-agent-works-a-formula-s-instance-on-its-branch-a-role-step-runs-its-role-charged-to-the-instance-a-failed-step-asks-on-a-thread-and-a-formula-starts-from-its-schedule-or-a-role.md
  item: S-0422
status: resolved
participants: [agent-S-0551, orchestrator, alex, agent-S-0422]
created: 2026-10-31T09:54:36Z
updated: 2026-11-02T03:47:21Z
---

# TH-0787 S-0422's after S-0421: keep it and state the reason in its Notes

On wip/kanban/stories/S-0422-a-story-s-agent-works-a-formula-s-instance-on-its-branch-a-role-step-runs-its-role-charged-to-the-instance-a-failed-step-asks-on-a-thread-and-a-formula-starts-from-its-schedule-or-a-role.md.

## Entries

### 2026-10-31T09:54:36Z agent-S-0551
**Recommendation:** keep `after: [S-0421]` and add one line under `### Planning` in S-0422's Notes:

> After S-0421: S-0422 reads the formula file format, `flai formula run`, and the instance story's `formula` field, which S-0421 writes, to work a formula's instance step by step.

## Why

- ADR-0230 (S-0551) says `after` names a story whose output this one reads, with the reason in Notes. `flai check` advises `wip.after-reason` when the Notes give none.
- S-0422 works an instance that `flai formula run` makes, and its Notes already say "the `formula` field is S-0421's". This is a read dependency, so the `after` is right.
- The check passes today only because the Notes mention S-0421 by ID in passing. The line above makes the reason explicit.

The story's words are yours, so I have not edited them.

### 2026-10-31T09:56:32Z orchestrator (recommendation)
**Recommended:** keep `after: [S-0421]` on S-0422, and add agent-S-0551's line under `### Planning` in S-0422's Notes, as written above.

**Why:**

- ADR-0188 §6 says `after` names only a story whose accepted change the story needs on main.
- S-0422 works an instance that `flai formula run` makes. It reads the formula format and the `formula` field, and S-0421 adds both. Its Notes already say "the `formula` field is S-0421's". This is a real dependency.
- The added line states the reason outright, as S-0551's criterion 2 asks, instead of leaving it as a passing mention.

**Why you:** S-0422 is finalized, so its words are yours. My role cannot edit it.

Source: design/adrs/0188-a-task-takes-no-folder-touch-the-planner-s-folder-touch-on-a-story-needs-the.md § Decision

### 2026-11-02T02:21:48Z alex
Confirmed the recommendation of 2026-10-31T09:56:32Z orchestrator.

Source: design/adrs/0188-a-task-takes-no-folder-touch-the-planner-s-folder-touch-on-a-story-needs-the.md § Decision

### 2026-11-02T02:22:58Z orchestrator
Your confirmation makes the recommendation the answer: S-0422 keeps `after: [S-0421]` (D-0454).

The line is not in S-0422's Notes yet. I checked at hash 9ba64612. S-0422 is finalized, so its words are yours, and my role cannot edit them.

Two ways it lands:

- You add it now, under `### Planning` in S-0422's Notes:

  > After S-0421: S-0422 reads the formula file format, `flai formula run`, and the instance story's `formula` field, which S-0421 writes, to work a formula's instance step by step.

- Or S-0422's agent adds it when it pulls the story. This thread is on S-0422, so it reaches that agent's inbox.

Nothing is held up meanwhile. `flai check` passes `wip.after-reason` on S-0422 today. S-0422 waits for S-0421, which is in progress.

Source: design/conventions/strategic-agents.md § Rules

Source: design/conventions/strategic-agents.md § Rules

### 2026-11-02T03:03:21Z agent-S-0422
Added under `### Planning` in S-0422's Notes, as confirmed: S-0422 keeps `after: [S-0421]`, with the reason that it reads the formula file format, `flai formula run`, and the `formula` field S-0421 writes.

### 2026-11-02T03:47:21Z alex
Resolved.
