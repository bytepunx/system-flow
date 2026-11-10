---
id: T-1719
type: task
nature: feature
title: agent-adapters.md and the operator docs say what the fixture run proves and what only a live run would
status: cancelled
parent: S-0412
owner: alex
created: 2026-10-16T05:20:42Z
updated: 2026-11-07T21:53:43Z
transitions:
  - to: cancelled
    at: 2026-11-07T21:53:43Z
    by: alex
stream: S-0412
tags: [cli]
touches: [docs/operators/index.md, design/system/agent-adapters.md]
after: [T-1718]
---
# T-1719 agent-adapters.md and the operator docs say what the fixture run proves and what only a live run would

## Work

- Add the fixture run to `design/system/agent-adapters.md` § OpenCode: what it proves composes (the adapter's start, the guard, the permission prompt, the usage and its pricing, the status), and what stays unproven until a live run against a real `opencode` and OpenRouter, which is not part of this repository for now (TH-0533). Say there which later harnesses the contracts look ready for.
- Update the table in `docs/operators/index.md` § Harnesses: what a story on `opencode` has, and that each piece is tested on fixtures only.
- Waits for T-1718: it describes what the test's assertions showed.

## Done when

The docs say what the fixture run proves and what it does not, and the markdown lint passes.

## Notes
- 2026-11-07T21:53:43Z: moved to cancelled: E-0023 cancelled: OpenCode isn't going to receive first-class support given its aversion to supporting before hooks for tool calls (amongst other things)
