---
id: T-1147
type: task
nature: improvement
title: The design, the user guide, the reference, and an ADR say that archiving an item resolves the threads still open on it
status: done
parent: S-0277
owner: alex
created: 2026-10-07T01:20:49Z
updated: 2026-10-07T03:06:09Z
transitions:
  - to: ready
    at: 2026-10-07T03:00:37Z
    by: agent-S-0277
  - to: in-progress
    at: 2026-10-07T03:00:38Z
    by: agent-S-0277
  - to: done
    at: 2026-10-07T03:06:09Z
    by: agent-S-0277
stream: S-0277
tags: [flai, docs]
touches: [design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md, design/adrs]
after: [T-1144, T-1146]
usage:
  source: log
  seconds: 331
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 85
      output: 20905
      cache_read: 3192960
      cache_write: 116716
      cost: 1.8336
---
# T-1147 The design, the user guide, the reference, and an ADR say that archiving an item resolves the threads still open on it

## Work

Write down the behaviour T-1144 and T-1146 built, in the same story, as `documentation.md` requires.

- Add an ADR with `flai adr new`, run from the story's worktree. It records the decision and its reason. Archiving an item, by acceptance or by `flai archive`, resolves every thread still `open` or `answered` on it, with an entry saying why. An archived item cannot take an answer, and a thread left open on it is a `threads.archived` finding at every story's close-out (I-0073). It also records the alternatives:
  - refusing an operator's acceptance while a thread is open, as ADR-0093 refuses the orchestrator's
  - leaving the thread and only reporting it
  - a thread archive folder
- In `design/system/flai-cli.md`, update the `flai accept` and `flai archive` rows, and the `flai check` rule for `threads.archived` with its new message.
- In `design/system/workflow.md`, under "Transitions and who makes them", say that acceptance resolves the story's open threads.
- In `docs/users/flai.md`, update "Accept and release", "Archiving", and the sentence under "Threads" that says open threads on archived items are flagged.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference` for the changed `flai archive` help, and the `flai accept` help if T-1144 changed it.

This task waits for T-1144 and T-1146, because it describes what they built and the reference is generated from their help text.

## Done when

- The ADR is accepted and listed in `design/adrs/README.md`. Each document above says what flai now does, and links the ADR where the document links ADRs.
- `docs/users/flai-reference.md` matches `make flai-reference`.
- The markdown lint and `flai check --strict` pass on the changed files.

## Notes
