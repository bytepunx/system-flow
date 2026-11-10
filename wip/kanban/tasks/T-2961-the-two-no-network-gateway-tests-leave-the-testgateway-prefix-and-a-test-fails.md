---
id: T-2961
type: task
nature: remediation
title: The two no-network gateway tests leave the TestGateway prefix, and a test fails on a TestGateway function gateway-smoke.sh does not run
status: in-progress
parent: S-0708
owner: alex
created: 2026-11-10T10:27:38Z
updated: 2026-11-10T10:54:01Z
transitions:
  - to: ready
    at: 2026-11-10T10:54:01Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T10:54:01Z
    by: system-flow
stream: S-0708
tags: [flai, go, testing]
touches: [flai/internal/serve/gateway_test.go]
---
# T-2961 The two no-network gateway tests leave the TestGateway prefix, and a test fails on a TestGateway function gateway-smoke.sh does not run

## Work

In `flai/internal/serve/gateway_test.go`, write the guard test first, named outside the prefix, such as `TestOnlyThePaidChecksAreNamedTestGateway`. It parses the package's `_test.go` files with `go/parser`, collects every function named `TestGateway*`, reads `scripts/gateway-smoke.sh` from the repository root, and fails naming each such function the script does not run by name. It needs no network and no `-short` gate, so the go-test tier runs it. Run it once on the current names and record in this task's Notes that it fails naming `TestGatewayStreamIsJudged` and `TestGatewayMissingNamesWhatIsMissing`: that is I-0334 reproduced.

Then rename the two tests out of the prefix, as S-0685 named `TestCreditRefusalIsFoundInAGatewayRun`: for example `TestStreamOfAGatewayRunIsJudged` and `TestMissingGatewayInputsAreNamed`. Update their doc comments and any reference to the old names in the file. Leave `TestGatewayOpenRouter`, `TestGatewayLiteLLM`, the go-test tier in `system-flow.yaml`, and `scripts/integration.sh` as they are. It waits for no task.

## Done when

- The guard test failed on the old names, as Notes records, and passes now.
- `flai test flai/internal/serve/gateway_test.go` passes, and its go-test tier runs the guard test and the two renamed tests.
- No function named `TestGateway*` in the package is missing from `scripts/gateway-smoke.sh`.

## Notes
