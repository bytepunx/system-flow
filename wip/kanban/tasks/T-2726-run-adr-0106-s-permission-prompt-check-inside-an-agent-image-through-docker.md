---
id: T-2726
type: task
nature: remediation
title: Run ADR-0106's permission_prompt check inside an agent image through docker
status: in-progress
parent: S-0618
owner: alex
created: 2026-11-09T10:48:07Z
updated: 2026-11-09T10:53:17Z
transitions:
  - to: ready
    at: 2026-11-09T10:53:17Z
    by: agent-S-0618
  - to: in-progress
    at: 2026-11-09T10:53:17Z
    by: agent-S-0618
stream: S-0618
tags: [flai, serve, container, claude-code]
touches: [flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go]
---
# T-2726 Run ADR-0106's permission_prompt check inside an agent image through docker

## Work

Add to `flai/internal/serve/claudecheck.go` a check of the Claude Code an agent image holds, made as the host's check is (ADR-0106 § Decision part 2), but inside the image.

- Build the scratch project with `newClaudeScratch` as the host's check does.
- Run the image's `claude` in a `docker run --rm` of the image, with `--entrypoint claude`.
  - Mount the scratch folder at its host path and work in its project.
  - Give the container the network, login mount, and provider variables a story's container gets, so the model call goes where that agent's would.
  - Build that argv in `flai/internal/harness/container.go`, beside `dockerRun` and `containerStart`, as an exported helper, rather than copying its mounts and provider-key refusal into `serve`.
- The headless argv is `harness.ClaudeCodeCheck` with the image's own `claude` and, as the MCP server on stdio, the image's own `flai` (`/usr/local/bin/flai`, which the bedrock image ships) with `--config` the scratch configuration and `mcp --agent flai-check`. Stdio is used because over HTTP permission_prompt auto-approves only the agent flai serve runs for the story (ADR-0243), which a scratch project has none of.
- Record the outcome in `claude-checks.json` under the version and the provider it ran through, as the host's check does, with the image named on the record (a new `image` field on `ClaudeCheck`, omitted for the host's).
- Tell a failure with `tell`, naming the image in the thread's text.
- Add a method on `claudeChecker`, `considerImage` or the like, that starts this check in the background, one at a time with the host's check, and never waits for it. It starts none when the version has a record through the provider, whatever its outcome.

Test it in `claudecheck_test.go` with a fake docker runner: the argv mounts the scratch folder, names the image and `claude`, and carries the permission arguments; a written file is recorded passed with the image; a missing file is recorded failed and told on a thread; a recorded version is not checked again. Test the new harness helper in `flai/internal/harness/container_test.go`.

This task waits for none: it is the first layer.

## Done when

- An agent image's Claude Code version with no record through the project's provider can be checked inside the image, its outcome recorded in `claude-checks.json` with the image, and a failure told on a thread.
- The tests above pass under `flai test` on the four files.

## Notes
