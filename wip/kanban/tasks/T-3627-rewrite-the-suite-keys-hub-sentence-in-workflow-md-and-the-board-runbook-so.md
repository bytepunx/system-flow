---
id: T-3627
type: task
nature: remediation
title: Rewrite the suite keys' hub sentence in workflow.md and the board runbook so each step has its own sentence, with no running list of stories
status: backlog
parent: S-0856
owner: alex
created: 2026-11-27T23:58:15Z
updated: 2026-11-28T00:01:15Z
transitions: []
stream: S-0856
tags: [docs]
touches: [design/system/workflow.md, docs/operators/runbooks/change-the-board.md]
---
# T-3627 Rewrite the suite keys' hub sentence in workflow.md and the board runbook so each step has its own sentence, with no running list of stories

## Work

- In `design/system/workflow.md` § Suites a step needs and runs, replace the paragraph that opens "Since S-0807 a story's move and its close-out read `verify`" with one sentence or short paragraph per step that reads the keys: the move and the close-out, the acceptance's gates, the batch, a state's `runs`, the publish, and `flai test --at`. Each says in the present tense what reads which key and links its own paragraph or section; none lists which stories have landed or are still to come. S-0808 was accepted on 2026-11-27 and left the paragraph naming every step "since S-nnnn" and then listing them again by story: read it on main as it is now.
- In the same section, rewrite "A tier's `queue` and `when` keep their meaning until S-0813 retires them" to say what they mean now, with no forward reference to a story.
- In `docs/operators/runbooks/change-the-board.md`, after "Since S-0807 a story's move and its close-out read `verify`", replace the sentence "Nothing else reads the keys yet but ..." in the same way: one bullet or sentence per step, beside the ones for S-0810 and S-0811 already there, and no sentence that sums up which steps do and do not read the keys. That sentence still says "the acceptance's gates run as they did", which S-0808 made untrue: say that they read the accept transition's `verify`.
- In `design/system/workflow.md` § Branches and collisions, under the hub files headings, add a short paragraph on hub lines: a sentence every sibling story rewrites whole conflicts at each acceptance even with the section and word merges ([ADR-0287](../../../design/adrs/0287-a-section-merge-merges-a-paragraph-both-sides-change-by-its-sentences-keeping.md), [ADR-0314](../../../design/adrs/0314-a-section-merge-merges-a-line-both-sides-change-by-its-words-keeping-both.md)), as I-0387 and I-0414 found, so each step is written on its own line.
- Waits for no task: it runs beside the check task, and neither changes the other's files.

## Done when

- Neither file holds a sentence that lists which stories have landed or are still to come, nor names a story not yet started.
- Each step that reads the suite keys has its own sentence or paragraph in both files, true to the code on main.
- `flai test design/system/workflow.md docs/operators/runbooks/change-the-board.md` passes.

## Notes
