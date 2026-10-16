---
id: T-1685
type: task
nature: feature
title: The opencode adapter runs opencode run headless with the prompt, provider/model, JSON output, and a resumed session, and states resume as its only capability
status: backlog
parent: S-0406
owner: alex
created: 2026-10-16T05:14:12Z
updated: 2026-10-16T05:14:12Z
transitions: []
stream: S-0406
tags: [cli]
touches: [flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
---
# T-1685 The opencode adapter runs opencode run headless with the prompt, provider/model, JSON output, and a resumed session, and states resume as its only capability

## Work

- Add the `opencode` adapter in `flai/internal/harness/opencode.go` and register it in `Adapters` in `harness.go`.
- `DefaultHost` is `opencode`. `Start` builds `run <prompt> --model <provider/model> --format json`, with `--session <id>` on a resume, then the host's arguments, in the story's worktree, with the environment `roleEnv` gives.
- `Capabilities` say resume only. `Reader` returns the unread reader the `command` adapter uses, until the usage story.
- Add `opencode` to `TestWhichAdapterAStoryGets` and `TestEachAdapterStatesItsCapabilities`, and test the arguments for a first start and a resume.
- Waits for nothing. The configuration the agent reads is the next layer's.

## Done when

`harness.For` gives a story on `opencode` this adapter, its start arguments are tested for a first start and a resume, and the harness tests pass.

## Notes
