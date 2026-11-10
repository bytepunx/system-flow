---
id: I-0334
title: "TestGatewayStreamIsJudged and TestGatewayMissingNamesWhatIsMissing run in no tier, because go-test and integration skip ^TestGateway and the gateway tier runs only the two paid checks by name"
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-10T09:01:59Z
last_reported: 2026-11-10T09:01:59Z
updated: 2026-11-10T09:55:53Z
---

# I-0334 TestGatewayStreamIsJudged and TestGatewayMissingNamesWhatIsMissing run in no tier, because go-test and integration skip ^TestGateway and the gateway tier runs only the two paid checks by name

## Description
TestGatewayStreamIsJudged and TestGatewayMissingNamesWhatIsMissing run in no tier, because go-test and integration skip ^TestGateway and the gateway tier runs only the two paid checks by name

## Instances

### 2026-11-10T09:01:59Z
Story: S-0685.
S-0685 (T-2893) added a no-network test of the gateway check's 402 handling to flai/internal/serve/gateway_test.go. Named TestGateway..., it would never run: the go-test tier in system-flow.yaml and scripts/integration.sh pass -skip '^TestGateway', and scripts/gateway-smoke.sh runs only -run '^TestGatewayOpenRouter$' and '^TestGatewayLiteLLM$'. S-0685 named its test TestCreditRefusalIsFoundInAGatewayRun so go-test runs it. The two older no-network tests in the same file, TestGatewayStreamIsJudged and TestGatewayMissingNamesWhatIsMissing, match the skip and so run in no tier: a change that breaks the stream parsing or the skip messages passes every close-out. A fix: narrow the skip to '^TestGateway(OpenRouter|LiteLLM)$' in the go-test tier and integration.sh, or rename the two tests.

## Remediation

Story S-0708 remediates this issue, created from it at 2026-11-10T09:55:53Z.
