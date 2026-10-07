---
id: ADR-0109
title: Archiving an item resolves the threads still open on it
status: accepted
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0093]
---

# ADR-0109 Archiving an item resolves the threads still open on it

## Context

A thread anchored on a work item follows the item into `wip/archive` and stays in `wip/threads` ([ADR-0020](0020-files-plus-mcp.md)): threads have no archive folder. `flai accept` and `flai archive` moved an item and left its threads `open` or `answered`, and `flai check` warned `threads.archived` on each until someone resolved it by hand. An archived item takes no more work, so a thread left open on it awaits an answer nothing will act on.

The warning came back at every story's close-out. [ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md) made such a finding outside the story a note that the close-out records, and [I-0073](../issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md) recorded it 30 times between 2026-10-05 and 2026-10-07, each on a thread left by an acceptance that no story of the close-out's had touched. The finding said "resolve it or move it", and there was nowhere to move it.

[ADR-0093](0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md) already refuses the orchestrator's acceptance while a thread on the story or one of its tasks is not resolved. Nothing stopped the operator's acceptance, or `flai archive`, leaving one open.

## Decision

Archiving an item resolves every thread still `open` or `answered` that is anchored on it, with a dated entry saying why.

1. **By acceptance.** After it archives (step 2), `flai accept` resolves each thread still open or answered on what it archived: the story and its tasks, and an epic that followed the story to done with its cancelled stories and their tasks. The entry is `Resolved: <id> was accepted`, by whoever accepts (`--by`, else the author). The thread files go into the acceptance commit. A rerun that finishes an acceptance whose commit failed resolves the threads left open first, those an older flai left among them. `--dry-run` says `would resolve TH-nnnn, open or answered on what it archives`, and `--json` lists them in `resolved_threads`. A thread the preview cannot read is a blocker; one that cannot be resolved after the archive stops the acceptance before the commit, naming it and saying to run `flai accept` again.
2. **By `flai archive`.** It resolves each thread still open or answered on each item it archives, as `<id> was archived`, by the author. It commits nothing, as before. `--dry-run` names the threads it would resolve, and `--json` lists them in `resolved_threads`.
3. **The narrative.** A story's narrative that is still live, when a task is archived without its story, is mirrored again, so its `## Open questions` drops the threads resolved. An archived narrative is not.
4. **The orchestrator's refusal stands.** ADR-0093's `thread_open` blocker still refuses the orchestrator's acceptance while a thread on the story or one of its tasks is not resolved: the orchestrator does not settle a question by accepting over it. Its acceptance resolves only what that blocker does not cover, the threads on an epic that follows the story and on that epic's cancelled stories. The operator, who may accept over an open thread, settles it by accepting.
5. **The check.** `threads.archived` stays, for a thread left open by a hand edit or an older flai, and now says `<TH-nnnn> is <status> but <id> is archived; resolve it with flai thread resolve <TH-nnnn>`.

## Consequences

- No acceptance or `flai archive` leaves a thread open on an archived item, and the close-outs stop recording `threads.archived` for one.
- A question still awaiting an answer when its item is accepted closes without one. Its file keeps every entry and the reason, and `flai thread list --all` shows it; a question that still matters is opened again as a new thread on a live item.
- The acceptance commit holds thread files as well as the work items and the archive.
- An acceptance can stop after archiving on a thread it cannot write, as it can on a commit; rerunning `flai accept` finishes it.

## Alternatives considered

- **Refuse the operator's acceptance while a thread is open, as ADR-0093 refuses the orchestrator's.** The operator reads the threads and accepts over them on purpose; a refusal would make every such acceptance two steps, resolve and then accept, for the same outcome.
- **Leave the thread open and only report it.** That was the status quo, and it is what I-0073 recorded at every close-out.
- **A thread archive folder.** Moving the thread would hide it from `flai check` and leave it unresolved, and every reader of `wip/threads`, the dashboard and the MCP server among them, would need a second place to look.
