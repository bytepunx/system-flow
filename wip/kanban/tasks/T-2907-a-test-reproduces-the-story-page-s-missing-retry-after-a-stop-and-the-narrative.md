---
id: T-2907
type: task
nature: improvement
title: A test reproduces the story page's missing Retry after a stop, and the narrative names its cause and the fix proposed
status: backlog
parent: S-0669
owner: alex
created: 2026-11-10T08:47:33Z
updated: 2026-11-10T08:47:33Z
transitions: []
stream: S-0669
tags: [flaiover, agent, test]
touches: [flaiover/src/lib/components/StoryAgent.svelte.test.ts, "flaiover/src/routes/items/[id]/item.svelte.test.ts"]
---
# T-2907 A test reproduces the story page's missing Retry after a stop, and the narrative names its cause and the fix proposed

## Work

This task waits for none: it finds the cause the fix builds on. I-0278 was not reproduced: the page's `activity` and `status` were not captured when the operator saw no Retry. The page and the card's menu both decide Retry with `retryable()` in `flaiover/src/lib/activity.ts`. So the two places must differ in what they pass it, or in what the page does after the rule.

Compare the two inputs for a stopped run of a story in progress: `state: 'failed'`, `outcome: 'stopped'`, `stopped` set, and exit 143.

- The card's menu: `cardMenu` in `flaiover/src/lib/cardmenu.ts`, given `card.status`, the board's `writable`, and `hostAgent.enabled` from `flaiover/src/routes/board/+page.svelte`.
- The story page: `StoryAgent.svelte`'s `canRetry`, given `item.status`, `writable && !item.archived` from `flaiover/src/routes/items/[id]/+page.svelte`, and `status.enabled` from its own `/api/host-agent` answer. It also checks `activity.run.started !== retried`, which hides Retry once pressed for that run.

Check what `agent.status` gives the page for such a run against what it gives the board. If they differ, the cause is in flai (`flai/internal/serve/agents.go`), and the test belongs there.

Write a test where the cause is: in `StoryAgent.svelte.test.ts` or `item.svelte.test.ts`, or in flai's `agents_test.go` when the cause is there. It must fail on main for that cause. Record in the narrative's `## Decisions` the cause, the evidence, and the fix proposed. If the fix would change what ADR-0058 §4 says a stopped run offers, ask on a thread on S-0669 before building it.

## Done when

- [ ] A test fails on main because the story page offers no Retry for a stopped run of a story in progress where the card's menu offers one.
- [ ] The narrative's `## Decisions` names the cause, the evidence for it, and the fix proposed.

## Notes
