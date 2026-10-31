---
id: T-2369
type: task
nature: remediation
title: The gateway checks leave the smoke tier for a gateway tier of their own that runs only when the operator asks or on the host's schedule
status: backlog
parent: S-0548
owner: alex
created: 2026-10-31T07:26:50Z
updated: 2026-10-31T07:26:50Z
transitions: []
stream: S-0548
tags: [cli]
touches: [system-flow.yaml, scripts/smoke.sh, scripts/gateway-smoke.sh, docs/contributors/index.md, docs/operators/index.md, design/conventions/tooling.md, design/system/agent-adapters.md]
after: [T-2366, T-2368]
---
# T-2369 The gateway checks leave the smoke tier for a gateway tier of their own that runs only when the operator asks or on the host's schedule

## Work

`scripts/smoke.sh` ends with `scripts/gateway-smoke.sh`, which runs `TestGatewayOpenRouter` and `TestGatewayLiteLLM` through live `claude -p` sessions whenever the keys are set: 70 s to 194 s a batch, about 0.32 USD each, and a 402 from OpenRouter on 2026-10-31 that fails the next batch on the operator's balance. Take that step out of `scripts/smoke.sh`.

Add a `gateway` tier to `system-flow.yaml` that runs `scripts/gateway-smoke.sh`, with `queue: false` and `all_only: true`, and paths that never select it for a close-out or a batch, so that it runs only by name with T-2368's `flai test gateway`. Check with `flai verify` on a branch that changes `scripts/gateway-smoke.sh` that neither the close-out nor `flai test --all` in a batch runs it, and note in `## Decisions` how the tier keeps out of `--all` if `all_only` alone does not.

The schedule is the host's: document in `docs/operators/index.md` how to run `flai test gateway` from the host's own scheduler (a cron line or a systemd timer), what each run spends and how `FLAI_TEST_GATEWAY_BUDGET` caps it. Rewrite the Gateway checks section of `docs/contributors/index.md`, the `flai test` line of `design/conventions/tooling.md`, and the mention in `design/system/agent-adapters.md`, which say the smoke tier runs them. The gateway tests stay: they are the only check that the guard and `permission_prompt` hold through a gateway (S-0351, S-0356).

Third layer: it waits for T-2366, which changes `system-flow.yaml`, `scripts/smoke.sh`, `tooling.md`, and `docs/operators/index.md` too, and for T-2368, without which the tier cannot be run by name.

## Done when

A batch's `smoke.log` holds no gateway run; `flai test gateway` on the host runs the two checks, or names each gateway it skips; the testing convention and the operator docs say when it runs and what it spends. Story criterion 2.

## Notes

Drafted by the planner of S-0548 on 2026-10-31. The plan assumes the host's schedule is the host's own scheduler, with no schedule code in flai; the plan's thread asks.
