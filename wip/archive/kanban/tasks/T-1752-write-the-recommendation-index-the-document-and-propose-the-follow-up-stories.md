---
id: T-1752
type: task
nature: research
title: Write the recommendation, index the document, and propose the follow-up stories
status: cancelled
parent: S-0413
owner: alex
created: 2026-10-16T05:46:00Z
updated: 2026-10-16T06:01:29Z
transitions:
  - to: cancelled
    at: 2026-10-16T06:01:29Z
    by: agent-S-0413
stream: S-0413
tags: [prose, research]
touches: [design/system/agent-prose.md, design/system/README.md]
after: [T-1749]
---
# T-1752 Write the recommendation, index the document, and propose the follow-up stories

## Work

Close the research with one recommendation the operator can act on.

- Open `design/system/agent-prose.md` with a short summary and a recommendation section: the order to try the linter, the model, and the alternatives in, and what each step should show before the next starts.
- List the follow-up stories the recommendation implies, each with a title, a nature, and a sentence of scope, such as a linter in `flai check`, a training experiment, or a length budget. Raise them with the operator on a thread on S-0413 rather than creating them: which to make is the operator's choice.
- Bring the document's row in `design/system/README.md`, which T-1750 adds, in line with the recommendation.
- Run the markdown lint on both files.

Waits for T-1749: the recommendation weighs every section before it.

## Done when

- `design/system/agent-prose.md` opens with the summary and the recommendation, and its row in `design/system/README.md` describes it.
- The thread on S-0413 names the proposed follow-up stories.
- `flai test` on the two files passes.

## Notes
- 2026-10-16T06:01:29Z: moved to cancelled: Duplicates agent-S-0413's plan (T-1750, T-1755, T-1758, T-1760), as the planner recommends on TH-0472
