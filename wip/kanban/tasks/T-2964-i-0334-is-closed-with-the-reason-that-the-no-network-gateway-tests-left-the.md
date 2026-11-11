---
id: T-2964
type: task
nature: remediation
title: I-0334 is closed with the reason that the no-network gateway tests left the TestGateway prefix and a test guards it
status: done
parent: S-0708
owner: alex
created: 2026-11-10T10:27:51Z
updated: 2026-11-10T11:03:46Z
transitions:
  - to: ready
    at: 2026-11-10T11:03:36Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T11:03:36Z
    by: system-flow
  - to: done
    at: 2026-11-10T11:03:46Z
    by: agent-S-0708
stream: S-0708
tags: [issues]
touches: [design/issues/I-0334-testgatewaystreamisjudged-and-testgatewaymissingnameswhatismissing-run-in-no.md, design/issues/summary.md]
after: [T-2961, T-2962]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 491
      cache_read: 558256
      cache_write: 9164
      cost: 0.1932
---
# T-2964 I-0334 is closed with the reason that the no-network gateway tests left the TestGateway prefix and a test guards it

## Work

Run `flai issue close I-0334 --reason "<reason>"`. The reason names the two tests' new names, which the go-test tier now runs, and the guard test in `flai/internal/serve/gateway_test.go` that fails on any `TestGateway*` function `scripts/gateway-smoke.sh` does not run, as the reproduction. It says the skip in the go-test tier and `scripts/integration.sh` stays `^TestGateway`, so a new paid check is never run by a close-out, and names the rule in `design/conventions/code-quality.md`. This rewrites the issue and its row in `design/issues/summary.md`. It waits for T-2961 and T-2962, so that the reason can name the tests and the rule as they are.

## Done when

- I-0334 is closed, with that reason.
- `design/issues/summary.md` shows it closed.
- `flai check --strict` is clean.

## Notes
