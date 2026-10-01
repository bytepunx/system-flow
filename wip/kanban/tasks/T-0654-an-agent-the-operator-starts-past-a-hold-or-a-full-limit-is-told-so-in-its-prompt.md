---
id: T-0654
type: task
nature: remediation
title: An agent the operator starts past a hold or a full limit is told so in its prompt
status: done
parent: S-0182
owner: arobson
created: 2026-10-01T09:16:25Z
updated: 2026-10-01T09:21:42Z
transitions:
  - to: ready
    at: 2026-10-01T09:16:34Z
    by: agent-S-0182
  - to: in-progress
    at: 2026-10-01T09:19:38Z
    by: agent-S-0182
  - to: done
    at: 2026-10-01T09:21:42Z
    by: agent-S-0182
stream: S-0182
tags: []
touches: [flai/internal/serve/start.go, flai/internal/serve/restart.go, flai/internal/serve/start_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
usage:
  source: log
  seconds: 124
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 42
      output: 12063
      cache_read: 3462422
      cache_write: 42689
      cost: 1.2596
---
# T-0654 An agent the operator starts past a hold or a full limit is told so in its prompt

## Work

`flai serve agent start` (the story page's Start agent) overrides a hold and the in-progress limit by design, but the agent it starts gets the launcher's prompt, "started by flai serve ... because it entered ready", and the warning goes only to the command's log. An agent started that way cannot tell it was started on the operator's word past a hold, which is what I-0050 records. Pass the hold's reason, or that the limit was full, to the harness and say in the prompt that the operator started the story past it and the agent works it as asked, without narrowing or changing the claim to clear the hold.

## Done when

Tests show the prompt for a story started past a hold names the operator and the hold's reason, past a full limit says so, and a start that overrides nothing keeps today's prompt; `go test ./internal/serve ./internal/harness` passes.

## Notes
