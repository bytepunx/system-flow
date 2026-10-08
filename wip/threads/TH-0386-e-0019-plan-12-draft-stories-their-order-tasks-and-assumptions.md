---
id: TH-0386
title: "E-0019 plan: 12 draft stories, their order, tasks, and assumptions"
anchor:
  path: wip/kanban/epics/E-0019-agent-adapters-for-litellm-and-openrouter.md
  item: E-0019
status: open
participants: [planner-E-0019]
created: 2026-10-08T09:01:45Z
updated: 2026-10-08T09:01:45Z
---

# TH-0386 E-0019 plan: 12 draft stories, their order, tasks, and assumptions

On wip/kanban/epics/E-0019-agent-adapters-for-litellm-and-openrouter.md.

## Entries

### 2026-10-08T09:01:45Z planner-E-0019
E-0019 now has 12 draft stories in the backlog, each with tasks, touches, a forecast, and a cost of delay value. The values come from your 500 USD a week penalty (TH-0375), shared by forecast. Each story's `### Planning` note says where every figure came from.

## Stories and order

| Story | After | Tasks by layer | Forecast | Value/week |
|-------|-------|----------------|----------|------------|
| S-0349 provider on the agent, providers on the host, manifest override | none | T-1371; T-1372; T-1373; T-1374 | 55m | 60.84 |
| S-0350 claude-code starts through the provider; key never logged | S-0349 | T-1376; T-1378; T-1380 | 29m | 32.08 |
| S-0351 OpenRouter check with guard and permission_prompt, key-gated tests, findings recorded | S-0350 | T-1384; T-1385 + T-1386 | 45m | 49.78 |
| S-0352 Capabilities and usage Reader on the adapter | S-0351 | T-1387 + T-1388; T-1389; T-1390 | 45m | 49.78 |
| S-0353 guard Decide on a neutral Call, Claude Code reader | S-0351 | T-1391; T-1392; T-1393 | 50m | 55.31 |
| S-0354 hold-and-ask package, ask verdict, protected list per harness | S-0353 | T-1394 + T-1395; T-1396; T-1397 | 35m | 38.72 |
| S-0355 harness without a guard: no roles unless guard none | S-0352, S-0359 | T-1398; T-1399 + T-1400; T-1401 | 34m | 37.61 |
| S-0356 LiteLLM check, both routes, key-gated tests, findings recorded | S-0351 | T-1402; T-1403 + T-1404 | 30m | 33.19 |
| S-0357 priced from the gateway's spend log, `priced_by` | S-0352, S-0356 | T-1405 + T-1406 + T-1408; T-1407; T-1409; T-1410 | 42m | 46.46 |
| S-0358 `priced_by` in metrics.md, flai stats, flai show, dashboard | S-0357 | T-1411; T-1412; T-1413 + T-1414 | 30m | 33.19 |
| S-0359 dashboard: Providers settings and the agent's provider field | S-0349 | T-1415; T-1416 + T-1417; T-1418 | 37m | 40.93 |
| S-0360 claude version check through the project's provider (addition) | S-0351 | T-1426; T-1427 | 20m | 22.12 |

`;` separates layers and `+` joins tasks that run together. The epic's order holds: the provider split and OpenRouter first, then the contracts, then cost. Docs ship with each story.

## Assumptions

1. S-0351 and S-0356 need your keys in the story agent's environment: `OPENROUTER_API_KEY`, and `LITELLM_BASE_URL` with `LITELLM_API_KEY`. Without them the agent asks on a thread. The tests skip, naming the variable, wherever a key is missing.
2. S-0356 waits only for S-0351, so LiteLLM may run beside the contracts rather than after them. The pull order is yours.
3. S-0355 changes the `command` harness. Its roles are refused unless `guard: none` is set. Its starts are refused unless auto-approve or the new `deny_protected` is on. `deny_protected` stands for ADR-0131's "the host's arguments deny the protected paths", since flai cannot read another harness's arguments.
4. S-0357 adds `spend_key_env` for a key that can read the spend log, as ADR-0132's consequences allow.
5. S-0354 keeps only Claude Code's protected paths. `.codex/`, `AGENTS.md`, and the rest come with their harness.
6. The planner, the orchestrator, and the analyzer start through a provider (S-0350), but ADR-0131's checks apply to story agents only.
7. S-0355 waits for S-0359 because both change the Settings page, `settings.go`, and `serve_actions.go`.

## Proposals

- **Add S-0360 (drafted).** Today's check of each new `claude` version calls Anthropic directly. On a capped subscription, the case behind your penalty, it would fail and open a thread on every project.
- **No split, merge, or drop.** S-0349 is the largest (size 33, 27 touches). If it grows, T-1373's CLI and MCP surface splits off cleanly.

## What you need to do

Review and finalize the drafts. Configure the OpenRouter key before S-0351 is pulled, and the LiteLLM proxy before S-0356.
