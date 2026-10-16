---
id: T-1848
type: task
nature: experiment
title: The results document weighs what else changed, cites threads an entry answered, says whether the hypothesis holds, and recommends keeping, changing, or dropping the digest
status: backlog
parent: S-0435
owner: alex
created: 2026-10-16T06:09:04Z
updated: 2026-10-16T06:09:04Z
transitions: []
stream: S-0435
tags: [cli]
touches: [design/experiments/S-0435-digest-priming-and-threads-to-the-operator.md, design/experiments/README.md, design/issues/summary.md]
after: [T-1847]
---
# T-1848 The results document weighs what else changed, cites threads an entry answered, says whether the hypothesis holds, and recommends keeping, changing, or dropping the digest

## Work

In the results document:

- Name what else changed between the windows: the epics and stories accepted, the stories' mix by nature, and the orchestrator's `answer_threads` and other permissions, read from `system-flow.yaml`'s history.
- List three threads of the after window that an entry answered, each with the entry cited, from the entries' `asked_again` and the threads' sources, or say that none did; and the entries the curator absorbed or compacted in the window (`flai digest list --json`).
- Say whether the result supports the hypothesis, refutes it, or cannot tell, and why.
- Recommend keeping the digest, changing it, or dropping it. For a change, file it with `flai issue new` and link it.
- Add the document to `design/experiments/README.md`.

Waits for T-1847, whose figures it weighs; the two touch the same document.

## Done when

Every claim cites a figure, a thread, or an entry; the README lists the document; `flai check --strict` and the markdown lint report nothing on it.

## Notes

`flai issue new` writes `design/issues/summary.md` and a new issue file; the close-out commits what it writes.
