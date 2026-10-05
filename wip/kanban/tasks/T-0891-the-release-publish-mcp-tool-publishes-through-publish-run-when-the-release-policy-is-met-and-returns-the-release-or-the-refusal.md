---
id: T-0891
type: task
nature: feature
title: The release_publish MCP tool publishes through publish.run when the release policy is met, and returns the release or the refusal
status: backlog
parent: S-0222
owner: alex
created: 2026-10-05T04:46:40Z
updated: 2026-10-05T04:46:40Z
transitions: []
stream: S-0222
tags: [flai]
touches: [flai/internal/mcpserver/orchestrate.go, flai/internal/mcpserver/orchestrate_test.go, flai/internal/mcpserver/folder.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go]
after: [T-0885]
---
# T-0891 The release_publish MCP tool publishes through publish.run when the release policy is met, and returns the release or the refusal

## Work

Add the MCP tool `release_publish` to `flai/internal/mcpserver/orchestrate.go`, beside S-0217's `release_evaluate`. Register it in `addProjectTools` in `flai/internal/mcpserver/folder.go`. It takes `reason`, one sentence that the orchestrator also logs: the policy figure that was met, or its judgement.

Before it changes anything, it refuses, saying why and what would allow it:

- when the project's `push` host action is off, since `publish.run` needs it
- when the evaluation holds the batch back under `whole_epics`, naming the stories and their epics
- when the policy is `threshold` or `theme` and the evaluation is not met, with its figures
- when the policy is `judgement` and `reason` is empty
- when nothing is accepted and unreleased

Otherwise it publishes as `publish.run` in `flai/internal/hostapi/writes.go` does: `flai release --pending` on the host, as the host action, never forced. Record the orchestrator as the caller in the run's log, so a release it cut is told apart from the operator's. It returns the policy and its figures, the versions and tags released, and the items bundled.

When `flai release --pending` refuses with exit 3, the tool returns its message unchanged as a conflict. That covers S-0174's check for newer release tags on the remote, and a remote that moved under the push. The orchestrator logs the message and opens a thread with it.

This task waits for T-0885, whose `held_by_epic` the tool refuses on. It runs with the prompt task, whose paths it does not share.

## Done when

- MCP tests on fixtures publish under `threshold` met, `theme` met, and `judgement` with a reason, and return the versions, tags, and items
- MCP tests refuse under `threshold` and `theme` not met, a batch held under `whole_epics`, `judgement` with no reason, and the `push` action off, and change no file, tag, or remote
- an MCP test on a fixture whose remote has newer release tags returns S-0174's refusal as a conflict, and one whose remote moved returns the moved-remote refusal
- `go test ./internal/mcpserver/ ./internal/hostapi/` passes

## Notes
