---
id: I-0044
title: Agents flai serve starts inherit the host's address and token, so a flai serve an agent runs by hand takes over the operator's MCP servers
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-09-24T09:26:17Z
last_reported: 2026-09-24T09:26:17Z
updated: 2026-09-24T09:26:17Z
---

# I-0044 Agents flai serve starts inherit the host's address and token, so a flai serve an agent runs by hand takes over the operator's MCP servers

## Description
Agents flai serve starts inherit the host's address and token, so a flai serve an agent runs by hand takes over the operator's MCP servers

## Instances

### 2026-09-24T09:26:17Z
2026-09-24, S-0115's live trial. agent-S-0115, started by the operator's flai serve, ran a scratch flai serve built from its branch with its own --config and serve folder, to try flai serve agent start against a scratch dashboard. The agent's environment carried FLAI_HOST_URL and FLAI_HOST_TOKEN, inherited from the operator's flai host through flai serve, so host.FromEnv found the operator's host. From 09:24:24Z, every 15 s, the scratch serve told that host to keep only the scratch project's MCP server, and the operator's serve told it to keep its own. The host stopped the operator's six HTTP MCP servers (4243 to 4248) and started one for the scratch project (4249), back and forth, for about two minutes. Stopping the scratch serve by PID ended it: the operator's serve told the host again and the six came back under new PIDs. FLAI_CONFIG is inherited the same way (the operator's ~/.flai/config.json), so scripts/flai.sh, which keeps a set FLAI_CONFIG, ran against the operator's config rather than .flai-cache/config.json. Remedy: the launcher removes FLAI_HOST_URL, FLAI_HOST_TOKEN, FLAI_HOST_ADDR, and FLAI_CONFIG from an agent's environment; and a flai serve with a --config other than the host's refuses the host's Keep, or keeps none, as outside a host.

## Remediation
