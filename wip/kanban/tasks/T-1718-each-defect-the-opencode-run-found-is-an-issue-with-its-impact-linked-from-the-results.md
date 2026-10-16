---
id: T-1718
type: task
nature: experiment
title: Each defect the OpenCode run found is an issue with its impact, linked from the results
status: backlog
parent: S-0412
owner: alex
created: 2026-10-16T05:20:39Z
updated: 2026-10-16T05:20:39Z
transitions: []
stream: S-0412
tags: [cli]
touches: [design/issues, design/issues/summary.md, design/experiments/S-0412-a-non-claude-model-works-one-story-of-this-repository-end-to-end-on-opencode-through-openrouter-with-the-guard-and-the-permission-prompt-on-and-the-results-record-what-broke-and-what-it-cost.md]
after: [T-1717]
---
# T-1718 Each defect the OpenCode run found is an issue with its impact, linked from the results

## Work

- For each defect T-1717 recorded, run `flai issue new` from the story's worktree with its impact, or `flai issue bump` the open issue that already records it.
- Link each issue from its finding in the results document, and write the document's conclusion: whether the contracts held, and what to fix first.
- Waits for T-1717: it files what that run found.

## Done when

Every defect in the results has an issue linked from it, and the conclusion is written.

## Notes
