---
id: I-0303
title: "TestGatewayLiteLLM fails its setup check: the template's settings run flai guard 6 times, and it wants 4"
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-11-07T10:19:48Z
last_reported: 2026-11-07T10:19:48Z
updated: 2026-11-07T10:19:48Z
---

# I-0303 TestGatewayLiteLLM fails its setup check: the template's settings run flai guard 6 times, and it wants 4

## Description
TestGatewayLiteLLM fails its setup check: the template's settings run flai guard 6 times, and it wants 4

## Instances

### 2026-11-07T10:19:48Z
Story: S-0543.
2026-11-07: S-0543 ran `GOMAXPROCS=2 go test -race -count=3 ./internal/serve/` without `-skip '^TestGateway'`. TestGatewayLiteLLM/unified failed all 3 times at gateway_test.go:528: "the template's settings run flai guard 6 times, want 4 (PreToolUse twice, SubagentStart, SubagentStop)". The template's .claude/settings.json now has more guard hooks than the test expects, so the manual gateway tier (`flai test gateway`) fails before it reaches a model call. The haiku alias resolved to a free OpenRouter model (gateway_test.go:480). Running the gateway test by mistake was S-0543's error; the failing check is a separate, pre-existing defect.

## Remediation
