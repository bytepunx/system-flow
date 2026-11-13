---
id: T-3223
type: task
nature: feature
title: The plan, board, stream, touches, role, analyze, curate, forge, kluster, and agent-run commands read a state's kind
status: backlog
parent: S-0744
owner: alex
created: 2026-11-12T16:07:38Z
updated: 2026-11-12T16:07:38Z
transitions: []
stream: S-0744
tags: [cli]
touches: [flai/cmd/plan.go, flai/cmd/board.go, flai/cmd/stream.go, flai/cmd/stream_sync.go, flai/cmd/stream_restore.go, flai/cmd/touches.go, flai/cmd/role.go, flai/cmd/analyze.go, flai/cmd/curate.go, flai/cmd/forge.go, flai/cmd/kluster.go, flai/cmd/serve_agent_run.go]
---
# T-3223 The plan, board, stream, touches, role, analyze, curate, forge, kluster, and agent-run commands read a state's kind

## Work

Replace the one to three references to a state's name in each of these commands with a read of the kind or the effect: `flai board` prints the workflow's columns in order, `flai plan` takes its candidates from the queue states, and the stream, touches, and agent-run commands tell an open story from a closed one by kind.

## Done when

These files name no state, and the existing `cmd` tests pass unchanged.

## Notes
