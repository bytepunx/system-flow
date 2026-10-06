---
id: ADR-0090
title: "A thread entry is marked a recommendation in its heading and cites its source in its last line, and the operator confirms a pending recommendation to make it the answer"
status: proposed
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0020]
topics: [cli, orchestration]
---

# ADR-0090 A thread entry is marked a recommendation in its heading and cites its source in its last line, and the operator confirms a pending recommendation to make it the answer

## Context

S-0220 lets the orchestrator reply on threads that await the operator, as `orchestration.permissions.answer_threads` allows: with `recommend` its reply is a recommendation the operator confirms, and with `autonomous` it answers when it can cite a source and otherwise escalates as a recommendation. A thread is a markdown file under `wip/threads/` (ADR-0020) whose entries are headed `### <time> <author>`, and its status in the front matter says who it awaits: any reply by someone other than the opener sets `answered`. A recommendation must not do that, since the operator has not answered, and both a recommendation and an answer must say what they rest on. The marks have to live in the file, pass the project's markdown lint, which flai applies to every thread write (S-0179), and leave a file an older flai still reads, since the operator's installed flai and an agent's may be older than the tree.

## Decision

A thread entry carries two optional marks, and a pending recommendation becomes the answer by one confirming entry.

1. **The recommendation mark is a suffix on the entry's heading**: `### 2026-10-06T10:00:00Z orchestrator (recommendation)`. flai reads the author without the suffix and sets the entry's `recommendation`. The heading differs from the author's plain one, so a recommendation never joins a plain entry of the same second (I-0043) and never repeats a heading (MD024).
2. **The source is the entry's last paragraph, one line**: `Source: <path>` or `Source: <path> § <heading>`, a repository path such as an ADR, a design document, or a convention, validated as a thread's anchor is: the path exists, and the heading is in it. Any entry may cite one; flai reads it into the entry's `source` as `{ path, heading }` and keeps it in the entry's text.
3. **Status.** A recommendation leaves the status as it was: the thread still awaits the operator. It is refused on a resolved thread and from the thread's opener. A plain reply that cites a source sets the status as any reply does: `answered` from anyone but the opener.
4. **Pending.** A recommendation is pending while the thread is not resolved, it is the newest entry by someone other than the opener, and no later entry confirms it. Any later reply by someone other than the opener takes its place, so the operator answering differently withdraws it.
5. **Confirm.** `threads.Confirm` appends an entry by the confirming author, `Confirmed the recommendation of <time> <author>.`, citing the recommendation's source, and sets `answered`. It is refused when no recommendation is pending and to the recommendation's own author.
6. **What readers get.** Each entry in the JSON of `threads.View`, `flai thread show --json`, the host API, and the MCP tools carries `recommendation` (a bool) and `source` (absent when none); a thread's view carries `pending_recommendation`, the pending entry with its `at`, `author`, `text`, `recommendation`, and `source`, or null.

## Consequences

- The file says what the orchestrator did, in a form the operator reads without flai: the heading says recommendation, and the last line names the source.
- An older flai reads a recommendation as a plain entry whose author is `orchestrator (recommendation)`, and the source line as part of its text. The status is in the front matter, so its inbox and lists are unchanged; code that compares the last entry's author to an agent's name, such as whether the orchestrator wrote last, sees another author until it is upgraded.
- A plain entry whose last paragraph is one line starting with `Source:` and a space now reads as citing that source, whoever wrote it.
- Two recommendations by one author in one second join under one heading, as plain entries do, and the joined entry cites the later source.
- The CLI, MCP, host API, guard, metrics, and dashboard build on these marks in later tasks of S-0220.

## Alternatives considered

- **A marker in the entry's first line, such as `Recommendation:`.** An older flai would keep the author intact, but a same-second entry by the same author would join the previous one under its heading, burying the marker mid-text, and a reply that starts with the word would read as a recommendation.
- **An HTML comment, such as `<!-- recommendation -->`.** Hidden when rendered, so the operator reading the file or an older dashboard could not tell a recommendation from an answer.
- **A field in the front matter.** Marks belong to entries, not threads; a thread holds several replies, and the front matter would have to point at one.
- **A status `recommended`.** Every reader of the closed status list, older flai included, would refuse the file, and the thread awaits the operator exactly as `open` says.
- **Taking the source line out of the entry's text.** Readers that print only the text, such as `flai thread show`, would lose the source until each learns the field.
