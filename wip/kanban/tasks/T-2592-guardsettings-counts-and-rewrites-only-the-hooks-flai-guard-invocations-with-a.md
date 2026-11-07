---
id: T-2592
type: task
nature: remediation
title: guardSettings counts and rewrites only the hooks' flai guard invocations, with a test on the template's settings that -short runs
status: backlog
parent: S-0591
owner: alex
created: 2026-11-07T21:48:10Z
updated: 2026-11-07T21:48:10Z
transitions: []
stream: S-0591
tags: [flai, test]
touches: [flai/internal/serve/gateway_test.go]
---
# T-2592 guardSettings counts and rewrites only the hooks' flai guard invocations, with a test on the template's settings that -short runs

## Work

`guardSettings` in `flai/internal/serve/gateway_test.go` counts every `flai guard` substring in `template/root/.claude/settings.json` and wants 4, then rewrites each to the built flai's path. Since S-0518 the two refusing hooks also print `flai guard did not run (exit %s)`, so it counts 6 and fails `TestGatewayLiteLLM` and `TestGatewayOpenRouter` before any model call (I-0303).

- Split the count and the rewrite out of `guardSettings` into a function that takes the settings' text and the flai path and needs no `testing.T`, and make it match only an invocation of the command: `flai guard` where the hook runs it, after `$(` or after a semicolon and a space, never the words inside the message. Keep the want of 4 (PreToolUse twice, SubagentStart, SubagentStop) and its error text.
- Leave the message text unrewritten, so the scratch project's hooks still say `flai guard did not run`.
- Add a test, not skipped under `-short`, that reads the template's real settings and checks 4 invocations rewritten to a given path, the two messages left as they are, and no `flai guard` invocation left. Run it against the current settings before the fix to see it fail with 6.
- Run `flai test flai/internal/serve/gateway_test.go`. The paid `gateway` tier is not run.

Waits for nothing: the first layer.

## Done when

- The new test passes under `go test -short ./internal/serve/ -run <its name>` and failed on the old counting.
- `guardSettings` uses the new function, and `go vet` and golangci-lint are clean on the file.

## Notes
