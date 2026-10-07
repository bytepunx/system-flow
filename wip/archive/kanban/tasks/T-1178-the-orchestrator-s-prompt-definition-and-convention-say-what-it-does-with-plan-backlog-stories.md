---
id: T-1178
type: task
nature: feature
title: The orchestrator's prompt, definition, and convention say what it does with plan_backlog_stories
status: done
parent: S-0328
owner: alex
created: 2026-10-07T19:40:36Z
updated: 2026-10-07T20:17:40Z
transitions:
  - to: ready
    at: 2026-10-07T20:02:02Z
    by: agent-S-0328
  - to: in-progress
    at: 2026-10-07T20:02:02Z
    by: agent-S-0328
  - to: done
    at: 2026-10-07T20:17:40Z
    by: agent-S-0328
stream: S-0328
tags: []
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md]
after: [T-1176]
usage:
  source: log
  seconds: 938
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 54
      output: 23704
      cache_read: 4021064
      cache_write: 101080
      cost: 1.9105
---
# T-1178 The orchestrator's prompt, definition, and convention say what it does with plan_backlog_stories

## Work

`orchestratePrompt`, the convention's "As the orchestrator" in both copies, and both `orchestrator.md` copies say, in each one's style:

- With `plan_backlog_stories`, after the epics, start the planner with `plan` for each story `flai plan --candidates` lists, one at a time.
- On a thread a story's planner opened, while the permission is on: approve a plan whose tasks, touches, and figures fit the story by replying and resolving the thread; answer its questions. For a cost of delay input it asks for, take its recommended figure unless the thread or the story gives a reason for one of its alternatives. Set it with `item_edit` `cost_of_delay` on the story, giving the inputs only, never a value: the planner works the value out on its next run. Then reply naming what it set, and resolve. Log each with `activity_log`. Never confirm a recommendation.
- The `answer_threads` rule that a cost of delay input is the operator's judgement does not apply to a story planner's thread while `plan_backlog_stories` is on.

`harness_test.go` pins the new sentences in the prompt and in both `orchestrator.md` copies.

The two `orchestrator.md` files are under `.claude/`: the sub-agent does not write them. It returns each whole new file in its final message, and the story's agent writes them.

It waits on T-1176, whose permission it names. It runs in layer 2 beside T-1177 and T-1179, sharing no path with them.

## Done when

- [ ] `flai test` passes on the paths changed, with both `orchestrator.md` files written.
- [ ] The prompt, both definitions, and both convention copies say the same.

## Notes
