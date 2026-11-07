---
id: T-1685
type: task
nature: feature
title: The opencode adapter runs opencode run headless with the prompt, provider/model, JSON output, and a resumed session, and states resume as its only capability
status: done
parent: S-0406
owner: alex
created: 2026-10-16T05:14:12Z
updated: 2026-11-07T10:42:53Z
transitions:
  - to: ready
    at: 2026-11-07T10:09:20Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:09:24Z
    by: system-flow
  - to: done
    at: 2026-11-07T10:42:53Z
    by: agent-S-0406
stream: S-0406
tags: [cli]
touches: [flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/harness/container.go, flai/cmd/serve_actions_test.go]
after: [T-1682]
usage:
  source: log
  seconds: 2009
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 112
      output: 37954
      cache_read: 7522760
      cache_write: 216890
      cost: 3.538
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
