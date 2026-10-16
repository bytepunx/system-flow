---
id: T-1699
type: task
nature: feature
title: The OpenCode integration test shows a protected write held and allowed by the owner's answer, or refused by the deny rule
status: backlog
parent: S-0408
owner: alex
created: 2026-10-16T05:16:41Z
updated: 2026-10-16T05:16:41Z
transitions: []
stream: S-0408
tags: [cli]
touches: [flai/tests/integration/opencode_test.go]
after: [T-1697, T-1698]
---
# T-1699 The OpenCode integration test shows a protected write held and allowed by the owner's answer, or refused by the deny rule

## Work

- Extend `flai/tests/integration/opencode_test.go`: the agent writes a file under `.opencode/`; under hold-and-ask, answer the thread it opens as the owner and assert the write lands; under deny rules, assert it is refused with the thread to open.
- Keep the skip on a missing or unknown `opencode`.
- Waits for T-1697 (the rules and the capability) and T-1698 (the plugin's reply).

## Done when

The test passes on the known `opencode` version and is skipped, saying so, without it.

## Notes
