---
id: D-0506
title: "S-0406 is a complete draft, but I leave finalizing it to you: it needs OpenCode installed on the host first (TH-0467 question 2): fixtures only for now - I have no intention of configuring OpenCode…"
kind: answer
sources: [TH-0532]
topics: [agents, cli, template, flai, go, conventions, code, all]
paths: [flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/opencodeconfig.go, flai/internal/harness/opencodeconfig_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/protected/protected.go, flai/internal/protected/protected_test.go, flai/internal/check/instructions.go, flai/internal/check/instructions_test.go, flai/internal/check/check.go, flai/internal/serve/agents_test.go, flai/tests/integration/opencode_test.go, scripts/opencode-smoke.sh, scripts/smoke.sh, scripts/README.md, design/adrs, design/system/agent-adapters.md, design/system/flai-cli.md, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, template/root/opencode.json, template/CHANGELOG.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [S-0406, S-0412, T-1682, TH-0532, TH-0533, TH-0866]
---

# D-0506 S-0406 is a complete draft, but I leave finalizing it to you: it needs OpenCode installed on the host first (TH-0467 question 2): fixtures only for now - I have no intention of configuring OpenCode…

OpenCode is not installed on this host, and the operator does not intend to configure it for now: S-0406 and the rest of E-0023 test the opencode adapter on recorded fixtures only, never against a real `opencode` binary.

The operator, alex, answered TH-0532 at 2026-11-07T09:27:45Z: "fixtures only for now - I have no intention of configuring OpenCode at this time".

What was asked: the orchestrator left S-0406 a draft because its first criterion checked the four open OpenCode findings against a real `opencode` binary, and flai never installs one. It recommended installing OpenCode, or answering TH-0467 question 2 with "fixtures only". planner-E-0023 then re-planned S-0406 to settle the findings from OpenCode's source at a pinned release (TH-0866), and the orchestrator finalized it.

See D-0507 for the same answer on OpenRouter and S-0412.

Sources: TH-0532, TH-0467.
