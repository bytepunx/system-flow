---
id: I-0210
title: flai guard lets a session flai serve started run the flai serve sub-commands that write the host configuration
class: defect
status: closed
count: 1
first_reported: 2026-10-26T04:49:22Z
last_reported: 2026-10-26T04:49:22Z
updated: 2026-10-30T00:25:27Z
---

# I-0210 flai guard lets a session flai serve started run the flai serve sub-commands that write the host configuration

## Description
flai guard lets a session flai serve started run the flai serve sub-commands that write the host configuration

## Instances

### 2026-10-26T04:49:22Z
Story: S-0498.
Found while building flai serve notify (T-2130). flai/internal/guard/guard.go has no rule for flai serve's own sub-commands: a story's agent, a strategic agent, or a sub-agent that flai serve started can run flai serve notify set or clear, flai serve forge, flai serve agent provider, or flai serve enable. agentEnv drops FLAI_CONFIG from a served agent, so such a call writes the operator's ~/.flai/config.json. That breaks safety.md (no agent changes the operator's global tool configuration), and ADR-0205's statement that the reminder channel is the operator's to set in a shell. The ADR-0039 settings host action is about the dashboard, not agents. A fix: guard refuses, to any session flai serve started, every flai serve sub-command that writes the host configuration or turns a host action on or off.

## Remediation

Story S-0511 remediates this issue, created from it at 2026-10-26T16:01:48Z.
Closed 2026-10-30T00:25:27Z: S-0511: ADR-0220 and the guard rule hostConfig in flai/internal/guard/hostconfig.go, run in Decide for every caller, refuse every sub-agent and every session flai serve started the flai serve sub-commands that write the host configuration or turn a host action on or off (enable, disable, agent set and clear, harness and provider with a setting or --reset, checks set, timeout, clear, components idle and max-servers, forge with a setting or --reset, notify set and clear, import add and remove), and flai config set and flai template use; the reads stay open and the operator's own session is not held. hostconfig_test.go reproduces the instance with the story's agent's flai serve notify set.
