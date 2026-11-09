---
id: T-2727
type: task
nature: remediation
title: Start the in-image check from a container start whose Claude Code version has no record
status: backlog
parent: S-0618
owner: alex
created: 2026-11-09T10:48:17Z
updated: 2026-11-09T10:48:17Z
transitions: []
stream: S-0618
tags: [flai, serve, container, claude-code]
touches: [flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/serve/agents.go, flai/internal/serve/serve.go]
after: [T-2726]
---
# T-2727 Start the in-image check from a container start whose Claude Code version has no record

## Work

In `flai/internal/serve/container.go`, have `containerVersion` hand the image's Claude Code version, the image, and the project's provider to the checker's in-image check from T-2726 when `claude-checks.json` holds no record of that version through that provider, in place of only warning `container claude code version unchecked`.

- The agent's start goes on at once: the check runs beside it and never delays or stops it (ADR-0106 § Decision part 2).
- Log that the check started, as the host's check logs `claude code check started`, naming the image.
- Keep a warning for a version whose record is a failure, so the operator still sees it in the serve log.
- Give the launcher the `claudeChecker` that `flai/internal/serve/serve.go` makes, through `flai/internal/serve/agents.go`, so that host and in-image checks run one at a time.

Test it in `container_test.go`, replacing the tests that expect the warning (lines 190 and 208 today). Use the lab's fake docker:

- a container start whose version has no record starts one in-image check, and the agent's start does not wait for it;
- a version recorded passed or failed through the provider starts none;
- this reproduces I-0253: before the change no check ran for an unrecorded image version.

This waits for T-2726, whose check and `considerImage` method it calls.

## Done when

- A container start of a claude-code agent whose image's version has no record runs ADR-0106's check inside the image, once, beside the start.
- The tests above pass under `flai test` on the four files.

## Notes
