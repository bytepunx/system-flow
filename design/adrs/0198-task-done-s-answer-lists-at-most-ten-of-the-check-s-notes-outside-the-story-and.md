---
id: ADR-0198
title: "task_done's answer lists at most ten of the check's notes outside the story and counts the rest in notes_omitted"
status: proposed
date: 2026-10-23
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0179]
---

# ADR-0198 task_done's answer lists at most ten of the check's notes outside the story and counts the rest in notes_omitted

## Context

`flai task done` runs `flai check --strict` scoped to the story and answers its findings inside the story as `check.findings` and those outside it as `check.notes` ([ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md)). A note never fails the close, and the close records no issue from it: the close-out does that. [ADR-0179](0179-task-done-and-story-start-carry-an-inbox-that-lists-the-story-s-threads-and.md) fits the whole answer to 40,000 bytes of JSON, but its fit cuts only the inbox, so the notes are carried whole however many there are.

[I-0202](../issues/I-0202-task-done-s-answer-carries-every-check-note-outside-the-story-so-it-exceeds-the-40-000-bytes-it-promises-and-overflows-the-tool-result.md) records S-0491, whose four `task_done` calls each answered about 90,500 characters, over Claude Code's tool output limit, so each answer went to a file the agent read back with `jq`. The bulk was 399 notes outside the story: 382 `doc.updated` and 17 `components.examples`.

## Decision

The check step of `flai task done` keeps at most ten of the check's notes outside the story, the first ten in the order the check gives them, as `check.notes`, and counts the rest in a new field `check.notes_omitted`, omitted when zero.

- **Findings stay whole.** The findings inside the story are listed in full: they fail the close and the agent fixes them.
- **The CLI.** `flai task done --json` prints the same result. The text output counts every note in its check line, lists the ten it keeps, and says how many it left out and that `flai check --story S-nnnn` lists them all.
- **The fit.** ADR-0179's fit is unchanged and still cuts only the inbox. With the notes bounded, the inbox fit keeps the whole answer within 40,000 bytes.
- **The close-out.** `flai verify` and the close-out keep every note, since `--record-issues` records them; this decision is the task close's answer only.

## Consequences

- An answer with hundreds of notes outside the story stays within the tool's output limit, and the agent reads only what it can act on.
- An agent that wants every note runs `flai check --story S-nnnn`.

## Alternatives considered

- **Cut the notes in the fit, the fewest that make the answer fit, before the inbox.** Keeps every note when the answer is small, but still answers up to about 35,000 bytes of notes the agent cannot act on, in place of the inbox's changes it can.
- **Count the notes by rule and list none.** Smallest, but drops the paths and messages a reader uses to see whether a note is the story's after all.
