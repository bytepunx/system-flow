---
id: T-2909
type: task
nature: improvement
title: The story page offers Retry for a stopped agent whenever the card's menu does
status: backlog
parent: S-0669
owner: alex
created: 2026-11-10T08:47:44Z
updated: 2026-11-10T08:47:44Z
transitions: []
stream: S-0669
tags: [flaiover, agent]
touches: [flaiover/src/lib/activity.ts, flaiover/src/lib/activity.test.ts, flaiover/src/lib/components/StoryAgent.svelte, "flaiover/src/routes/items/[id]/+page.svelte", design/system/flaiover-dashboard.md]
after: [T-2907]
---
# T-2909 The story page offers Retry for a stopped agent whenever the card's menu does

## Work

This task waits for T-2907, because it builds the fix proposed there, against the test written there. Make the change the narrative's `## Decisions` names, where its cause is. Keep `retryable()` in `flaiover/src/lib/activity.ts` the one rule for the story page and the card's menu (S-0202).

Add a case to `activity.test.ts` for a stopped run of a story in progress, with `outcome: 'stopped'`, `stopped` set, and exit 143. It checks that `retryable` and `agentAction` both answer Retry for it.

If the cause turned out to be in flai's `agent.status`, fix it in `flai/internal/serve/agents.go` instead. `flai task done` widens the touches with what the task commits.

If what the page offers after a stop changes, update the "Stopping an agent (S-0170)" and "Each story's agent (S-0104)" paragraphs of `design/system/flaiover-dashboard.md` in the same commit.

## Done when

- [ ] The test T-2907 wrote passes.
- [ ] `activity.test.ts` checks that a stopped run of a story in progress gets Retry from both `retryable` and `agentAction`.
- [ ] `flai test` passes on the paths the task changed.

## Notes
