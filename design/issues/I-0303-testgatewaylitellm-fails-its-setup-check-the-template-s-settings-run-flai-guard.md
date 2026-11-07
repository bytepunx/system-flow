---
id: I-0303
title: "TestGatewayLiteLLM fails its setup check: the template's settings run flai guard 6 times, and it wants 4"
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-11-07T10:19:48Z
last_reported: 2026-11-07T10:19:48Z
updated: 2026-11-07T23:50:03Z
---

# I-0303 TestGatewayLiteLLM fails its setup check: the template's settings run flai guard 6 times, and it wants 4

## Description
TestGatewayLiteLLM fails its setup check: the template's settings run flai guard 6 times, and it wants 4

## Instances

### 2026-11-07T10:19:48Z
Story: S-0543.
2026-11-07: S-0543 ran `GOMAXPROCS=2 go test -race -count=3 ./internal/serve/` without `-skip '^TestGateway'`. TestGatewayLiteLLM/unified failed all 3 times at gateway_test.go:528: "the template's settings run flai guard 6 times, want 4 (PreToolUse twice, SubagentStart, SubagentStop)". The template's .claude/settings.json now has more guard hooks than the test expects, so the manual gateway tier (`flai test gateway`) fails before it reaches a model call. The haiku alias resolved to a free OpenRouter model (gateway_test.go:480). Running the gateway test by mistake was S-0543's error; the failing check is a separate, pre-existing defect.

## Remediation

Story S-0591 remediates this issue, created from it at 2026-11-07T15:15:23Z.
Closed 2026-11-07T23:50:03Z: Fixed by S-0591 (T-2592, commit 2c04f2b79): guardSettings in flai/internal/serve/gateway_test.go now counts and rewrites only the hooks' runs of flai guard, matched as a command after "$(" or "; ", so the two refusal messages S-0518 added ("flai guard did not run") are neither counted nor rewritten. TestGuardSettingsRunsTheBuiltFlai checks it on the template's own settings under -short, so the go-test tier catches the next drift rather than the manual gateway tier. TestGatewayOpenRouter, which shares guardSettings, is fixed with it.
