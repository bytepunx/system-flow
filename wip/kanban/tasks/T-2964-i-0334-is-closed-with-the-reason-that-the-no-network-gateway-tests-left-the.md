---
id: T-2964
type: task
nature: remediation
title: I-0334 is closed with the reason that the no-network gateway tests left the TestGateway prefix and a test guards it
status: backlog
parent: S-0708
owner: alex
created: 2026-11-10T10:27:51Z
updated: 2026-11-10T10:27:51Z
transitions: []
stream: S-0708
tags: [issues]
touches: [design/issues/I-0334-testgatewaystreamisjudged-and-testgatewaymissingnameswhatismissing-run-in-no.md, design/issues/summary.md]
after: [T-2961, T-2962]
---
# T-2964 I-0334 is closed with the reason that the no-network gateway tests left the TestGateway prefix and a test guards it

## Work

Run `flai issue close I-0334 --reason "<reason>"`. The reason names the two tests' new names, which the go-test tier now runs, and the guard test in `flai/internal/serve/gateway_test.go` that fails on any `TestGateway*` function `scripts/gateway-smoke.sh` does not run, as the reproduction. It says the skip in the go-test tier and `scripts/integration.sh` stays `^TestGateway`, so a new paid check is never run by a close-out, and names the rule in `design/conventions/code-quality.md`. This rewrites the issue and its row in `design/issues/summary.md`. It waits for T-2961 and T-2962, so that the reason can name the tests and the rule as they are.

## Done when

- I-0334 is closed, with that reason.
- `design/issues/summary.md` shows it closed.
- `flai check --strict` is clean.

## Notes
