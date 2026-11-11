---
id: T-2961
type: task
nature: remediation
title: The two no-network gateway tests leave the TestGateway prefix, and a test fails on a TestGateway function gateway-smoke.sh does not run
status: done
parent: S-0708
owner: alex
created: 2026-11-10T10:27:38Z
updated: 2026-11-10T11:03:11Z
transitions:
  - to: ready
    at: 2026-11-10T10:54:01Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T10:54:01Z
    by: system-flow
  - to: done
    at: 2026-11-10T11:03:11Z
    by: agent-S-0708
stream: S-0708
tags: [flai, go, testing]
touches: [flai/internal/serve/gateway_test.go]
usage:
  source: log
  seconds: 550
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 895
      cache_read: 1017379
      cache_write: 16700
      cost: 0.3521
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

On the old names, `go test -count=1 -run '^TestOnlyThePaidChecksAreNamedTestGateway$' ./internal/serve/` failed naming exactly the two tests, which reproduces I-0334:

```text
--- FAIL: TestOnlyThePaidChecksAreNamedTestGateway (0.06s)
    gateway_test.go:392: TestGatewayStreamIsJudged (gateway_test.go:401:1) runs in no tier: ...
    gateway_test.go:392: TestGatewayMissingNamesWhatIsMissing (gateway_test.go:496:1) runs in no tier: ...
```

After the renames to `TestStreamOfAGatewayRunIsJudged` and `TestMissingGatewayInputsAreNamed`, `flai test flai/internal/serve/gateway_test.go` passed (gofmt, vet, golangci-lint, go-test), and the three tests pass under the go-test tier's `-short -skip '^TestGateway'`.
