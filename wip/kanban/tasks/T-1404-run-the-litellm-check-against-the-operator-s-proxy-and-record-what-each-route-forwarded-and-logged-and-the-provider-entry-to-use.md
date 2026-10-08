---
id: T-1404
type: task
nature: research
title: Run the LiteLLM check against the operator's proxy and record what each route forwarded and logged, and the provider entry to use
status: backlog
parent: S-0356
owner: alex
created: 2026-10-08T08:49:57Z
updated: 2026-10-08T08:49:57Z
transitions: []
stream: S-0356
tags: [cli]
touches: [design/system/agent-adapters.md, docs/operators/index.md]
after: [T-1402]
---
# T-1404 Run the LiteLLM check against the operator's proxy and record what each route forwarded and logged, and the provider entry to use

## Work

- Run `TestGatewayLiteLLM` against the operator's proxy and virtual key.
- In `design/system/agent-adapters.md`, add `### Checked against LiteLLM` beside S-0351's section: the date, the proxy's and `claude`'s versions, and for each route whether the fields reached the model, whether `haiku` resolved, whether `/spend/logs` rows carry the session ID, and the spend against `total_cost_usd`. Settle the LiteLLM items still marked to check, and say which route a provider entry should use.
- In `docs/operators/index.md`, a provider entry for a local proxy on the route chosen, with `flai serve agent provider` and a virtual key in the variable `key_env` names.

It waits for T-1402, whose test it runs.

## Done when

- Both documents carry what was observed, never estimated.
- `flai test design/system/agent-adapters.md docs/operators/index.md` passes.

## Notes

Layer 2 of S-0356. Without the proxy in the agent's environment, ask the operator on a thread on S-0356 and wait.
