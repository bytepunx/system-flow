---
id: D-0451
title: "A go test run by hand on flai/internal/serve without -skip '^TestGateway' runs the paid gateway checks, which flai test's go-test tier always skips"
kind: decision
sources: [S-0549, S-0685, I-0334, S-0708]
topics: [testing, go, flai]
paths: [flai/internal/serve/gateway_test.go, system-flow.yaml, design/conventions/delegation.md, design/conventions/code-quality.md, design/issues/I-0334-testgatewaystreamisjudged-and-testgatewaymissingnameswhatismissing-run-in-no.md]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: [S-0708, TH-1015]
---

# D-0451 A go test run by hand on flai/internal/serve without -skip '^TestGateway' runs the paid gateway checks, which flai test's go-test tier always skips

S-0549 diagnosed 42 failing real-git tests with three `go test` runs by hand, against `design/conventions/delegation.md`, which says not to run `go test` by hand. Its reason was that `flai test` runs only the whole integration tier. Its narrative records the cost: "one of them omitted `-skip '^TestGateway'` and ran the paid gateway test, which failed on OpenRouter's 402 and spent $0.00 by the key's usage" (wip/archive/agents/S-0549.md, Decisions).

The gateway checks are `TestGateway*` in `flai/internal/serve/gateway_test.go`. They call OpenRouter or the LiteLLM proxy with the operator's key, and run only by name, as `flai test gateway` (S-0548). The go-test tier in `system-flow.yaml` passes `-skip ^TestGateway` for that reason. So a hand `go test` of `flai/internal/serve`, or of `./...`, must pass `-skip '^TestGateway'` too, or it spends on the operator's key. When the key is out of credit, the run fails on a 402 that has nothing to do with the story (I-0237, D-0448).

Source: S-0549.

The other side of the pattern, that the `TestGateway` prefix names only a paid check `scripts/gateway-smoke.sh` runs by name and a no-network gateway test takes another name, is now stated in `design/conventions/code-quality.md` § Project additions. S-0708 (T-2961, T-2962) renamed the two no-network tests I-0334 recorded to `TestStreamOfAGatewayRunIsJudged` and `TestMissingGatewayInputsAreNamed`, and added `TestOnlyThePaidChecksAreNamedTestGateway`, which fails on a `TestGateway*` function the script does not run; I-0334 is closed. S-0685's `TestCreditRefusalIsFoundInAGatewayRun` was the first test named that way (`wip/archive/agents/S-0685.md`, § Decisions; `wip/archive/agents/S-0708.md`, § Decisions).

Sources: S-0549; S-0685; I-0334; S-0708.
