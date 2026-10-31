---
id: T-2412
type: task
nature: improvement
title: Review the open after chain S-0420 to S-0422 against the rule, keeping each after that reads another story's output with its reason and proposing the rest
status: backlog
parent: S-0551
owner: alex
created: 2026-10-31T07:44:42Z
updated: 2026-10-31T07:44:42Z
transitions: []
stream: S-0551
tags: [flai, planning]
touches: [wip/kanban/stories/S-0420-the-dashboard-s-planner-orchestrator-and-analyzer-pages-become-one-role-page-per-defined-role-with-its-settings-panel-its-activity-document-and-its-runs-and-agent-status-carries-every-role-by-name.md, wip/kanban/stories/S-0421-a-formula-is-a-markdown-file-under-design-workflows-formulas-that-flai-check-validates-and-flai-formula-run-makes-a-story-from-it-with-one-task-per-step-and-the-steps-after-between-them.md, wip/kanban/stories/S-0422-a-story-s-agent-works-a-formula-s-instance-on-its-branch-a-role-step-runs-its-role-charged-to-the-instance-a-failed-step-asks-on-a-thread-and-a-formula-starts-from-its-schedule-or-a-role.md]
after: [T-2393, T-2408]
---
# T-2412 Review the open after chain S-0420 to S-0422 against the rule, keeping each after that reads another story's output with its reason and proposing the rest

## Work

S-0417 to S-0419 are archived, so the open chain is S-0420 (after S-0419), S-0421 (after S-0420), and S-0422 (after S-0421). For each, read what the story reads of the one it waits for. Where it reads that story's output, such as a loader or a file format, ask the operator on a thread on that story to add the reason to its Notes, naming the story, since a finalized story's words are the operator's. Where it does not, ask there to drop the `after`. Run `flai check` to see the new advisory and record what it names.

It waits for the check advisory and the convention task, which define the rule it applies.

## Done when

- Each of S-0420, S-0421, and S-0422 has an `after` with a reason in its Notes, or a thread asking the operator to add the reason or drop the `after`.
- The story's narrative records the review, story by story.

## Notes

S-0425's `after` on S-0423 and S-0424 is outside the chain the criterion names; record it in the narrative if the advisory names it.
