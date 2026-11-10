---
id: I-0256
title: A strategic agent's container runs flai with no host configuration, so the orchestrator's flai accept and flai verify inside it see no host actions and lack the test tiers' tools
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-02T19:06:36Z
last_reported: 2026-11-02T19:06:36Z
updated: 2026-11-10T08:37:14Z
---

# I-0256 A strategic agent's container runs flai with no host configuration, so the orchestrator's flai accept and flai verify inside it see no host actions and lack the test tiers' tools

## Description
A strategic agent's container runs flai with no host configuration, so the orchestrator's flai accept and flai verify inside it see no host actions and lack the test tiers' tools

## Instances

### 2026-11-02T19:06:36Z
Story: S-0509.
S-0509 (ADR-0247) starts the planner, the orchestrator, the analyzer, and the curator in a container with FLAI_CONFIG and FLAI_CACHE_DIR in the role's cache folder `.flai-cache/agents/<role>`. The container's flai cannot read the host's `~/.flai/config.json`, which is never mounted, so it sees no host actions. The orchestrator runs `flai accept --by orchestrator` and `flai verify` from Bash in the main checkout. Inside the container, `flai accept` cannot tell that the merge queue batches, and `flai verify` cannot run the Go and Node tiers, since the image has neither (S-0510). Seen in the design, not in a live run: S-0567 makes the live run. Remedy to weigh: give the orchestrator's acceptance an MCP route that runs on the host, as `verify` already does through `flai serve`'s job, or mount a read-only copy of the host actions.

## Remediation

Story S-0687 remediates this issue, created from it at 2026-11-10T08:37:14Z.
