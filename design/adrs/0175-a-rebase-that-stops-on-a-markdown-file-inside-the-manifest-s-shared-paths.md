---
id: ADR-0175
title: "A rebase that stops on a markdown file inside the manifest's shared paths merges it by sections, and the merge queue keeps a story whose stop it cannot settle enqueued, messages a conflict with a story ahead in the batch, and sends back one with main"
status: proposed
date: 2026-10-19
supersedes: []
superseded_by: []
refines: [ADR-0098, ADR-0121, ADR-0126, ADR-0157, ADR-0173]
---

# ADR-0175 A rebase that stops on a markdown file inside the manifest's shared paths merges it by sections, and the merge queue keeps a story whose stop it cannot settle enqueued, messages a conflict with a story ahead in the batch, and sends back one with main

## Context

The merge queue (ADR-0173) rebases a copy of each enqueued story onto the batch in turn. Today a rebase that stops there is aborted and the story is left out, whatever the paths. `flai stream sync` and `flai accept` settle some stops themselves: `design/issues/summary.md` written again (ADR-0098), issue files merged by their instances (ADR-0126), and `design/adrs/README.md` merged row by row (ADR-0157). The batch ran none of them.

Every story changes the files in `claims.shared`, such as `design/system/flai-cli.md`, `docs/users/flai.md`, and `template/CHANGELOG.md`. Two stories that both edit one of them near the same lines stop git's rebase although their changes do not contradict each other: both add a section at the end, both add an item to the same list, or both bump the front matter's `updated` date. I-0150's third story stopped this way on `design/system/flai-cli.md`. A batch that leaves such a story out, or sends it back, stops the flow for a merge no person needs to judge.

The landing of a green batch runs the acceptance of today, whose rebase onto main must settle the same stops the batch's did, or the landing stops where the batch did not.

## Decision

A rebase that stops on a markdown file inside `claims.shared` merges it by sections, in `flai stream sync`, `flai accept`, and the merge queue's batch alike; the batch also runs the summary, issue, and ADR index steps, and treats a stop it still cannot settle by where its paths conflict.

### One set of merge steps

The batch's rebase of each story runs the steps `flai stream sync` and `flai accept` run on a stop, in the batch's worktree, from the same code: it merges the issue files and the ADR index, writes `summary.md` again from them, stages them, and continues, for as long as each stop is on those files alone. The fold after a clean rebase (ADR-0126) is left to the landing's acceptance, which runs it on the story's own branch.

### A markdown file in the shared paths merges by sections

A stop whose paths include a markdown file inside `claims.shared` that none of those steps covers merges it from the stop's three versions by sections:

- The file is cut at its headings into the front matter and sections, each keyed by its heading path.
- A section only one side changed takes that side's text; one both changed alike takes it; one both added appears once per side, ours first.
- A section both changed differently is merged line by line against its base. Where the two sides' changes do not touch the same base lines, both are kept. Where both insert lines at the same place and remove none, the batch's or main's lines come first, then the story's. Where both change or remove the same base lines, the file is not settled.
- In the front matter, `updated` takes the later of the two dates; any other key both changed differently leaves the file not settled.
- A section one side removed and the other changed leaves the file not settled.

A file not settled leaves the stop as it is: `flai stream sync` hands it to the agent, as it does today, and the batch treats it as below. A shared path that is not markdown is never merged.

The trial merge still reports a conflict on these files between two stories: unlike the issue files and the ADR index, a section merge can fail.

### What the batch does with a stop it cannot settle

The batch aborts the stop, leaves the story out, and records its outcome with the paths, by where they conflict:

| The stop's unsettled paths | Outcome | What happens |
|---|---|---|
| All inside `claims.shared` | `waiting` | The story stays enqueued and is shown waiting with the paths. It is never sent back for them. The queue tries it again once the main branch or its head moves. |
| Any outside, changed by a story ahead in the batch | `left_out` | The story stays enqueued. flai messages it from each story ahead whose changes conflict, in the two stories' conversation, `about` those paths, as a trial-merge conflict is told (ADR-0121). |
| Any outside, changed by no story ahead | `sent_back` | The conflict is with the main branch itself, as when the story ahead has landed. The story goes back to `in-progress`, by `flai`, with the paths as the reason, and leaves the queue. |

### Cases it cannot settle

- Both sides change or remove the same lines of one section.
- One side removes or renames a section the other changes.
- Both change a front matter key other than `updated` differently.
- A generated file outside the issues folder, such as `docs/users/flai-reference.md`, merges as markdown, never by being generated again.
- A shared path that is not markdown.

## Consequences

- A batch no longer stops, and never sends a story back, on stops in the shared paths that a section merge settles; `flai stream sync` and `flai accept` stop on fewer of them too.
- A section merge of independent edits can produce a document that reads oddly, such as two added sections in an order neither story chose. The story's own review and the batch's tiers, the markdown lint among them, are the check.
- A story left waiting on a shared path needs its agent to sync it, or the story ahead to land and the queue to try again.
- S-0400's landing and bisect use the same outcomes and send-back.

## Alternatives considered

- Union every conflicting hunk of a shared file: it would duplicate or interleave lines that both stories rewrote. Lost on correctness.
- Merge the shared files in the batch only: the landing's acceptance would stop on the same files. Lost on the flow.
- Send back a story whose shared-path stop is not settled: the story's own change is sound, and the criterion forbids it. Lost.
- Write `docs/users/flai-reference.md` again from the built flai: the batch would have to build flai mid-rebase. Deferred.
