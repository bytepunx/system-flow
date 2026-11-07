---
id: T-1716
type: task
nature: feature
title: The scripted OpenCode session and its OpenRouter generation fixture are written from the pinned release's event types and OpenRouter's documented answer
status: cancelled
parent: S-0412
owner: alex
created: 2026-10-16T05:20:26Z
updated: 2026-11-07T21:53:43Z
transitions:
  - to: cancelled
    at: 2026-11-07T21:53:43Z
    by: alex
stream: S-0412
tags: [cli]
touches: [flai/tests/integration/testdata/opencode/session.jsonl, flai/tests/integration/testdata/opencode/generation.json]
---
# T-1716 The scripted OpenCode session and its OpenRouter generation fixture are written from the pinned release's event types and OpenRouter's documented answer

## Work

- Write `flai/tests/integration/testdata/opencode/session.jsonl`: the scripted session the stub replays, one line per step, in the pinned OpenCode release's `--format json` event types that S-0406 recorded in `design/system/agent-adapters.md` § OpenCode. It starts the story, closes one task, makes one sub-agent `flai` write, one protected write, and moves the story to review, with a `step_finish` per model call naming a non-Claude OpenRouter model as `openrouter/<vendor>/<model>` and the generation ID S-0410 joins on.
- Write `flai/tests/integration/testdata/opencode/generation.json`: OpenRouter's documented `/api/v1/generation` answer for those generation IDs, in the shape `flai/internal/usage/testdata/openrouter/generation-a.json` has.
- Hand-write both from the documentation and the source; nothing is recorded from a paid call, and neither holds a key (TH-0533).
- Waits for nothing in this story: the event types and the join come from S-0406 and S-0410 on main.

## Done when

Both fixtures exist, each says in the test that reads it which release or page it follows, and neither holds a key.

## Notes
- 2026-11-07T21:53:43Z: moved to cancelled: E-0023 cancelled: OpenCode isn't going to receive first-class support given its aversion to supporting before hooks for tool calls (amongst other things)
