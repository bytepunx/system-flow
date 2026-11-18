---
id: I-0405
title: flai kluster deploy --job, teardown --job, and flai checks run --job from a story's container write their request under the main checkout's .flai-cache/kluster/requests, which the container does not mount
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-18T08:00:16Z
last_reported: 2026-11-18T08:00:16Z
updated: 2026-11-18T08:00:16Z
---

# I-0405 flai kluster deploy --job, teardown --job, and flai checks run --job from a story's container write their request under the main checkout's .flai-cache/kluster/requests, which the container does not mount

## Description
flai kluster deploy --job, teardown --job, and flai checks run --job from a story's container write their request under the main checkout's .flai-cache/kluster/requests, which the container does not mount

## Instances

### 2026-11-18T08:00:16Z
Story: S-0834.
Found by reading the code in S-0834 (T-3549); not run. `askKlusterJob` (flai/cmd/kluster.go, `root := mainRootOf(repo)` then `kluster.Ask(ctx, root, …)`) is called by deploy --job, teardown --job, and `flai checks run --job` (flai/cmd/checks.go). `kluster.Ask` writes the request under `RequestsDir(root)`, `<main>/.flai-cache/kluster/requests` (flai/internal/kluster/request.go, `RequestsDir`, `MkdirAll`, the request write). `HasAskedJob` and `Await`, which `--wait` uses, read the same folder. A story's container mounts no part of `.flai-cache/kluster` (ADR-0253, ADR-0319). So from a container's shell, the request either fails at MkdirAll or lands in the container's own filesystem: flai serve never sees it, and after ten seconds `Ask` says flai serve is not running. The MCP tools `kluster_deploy`, `kluster_checks`, and `kluster_teardown` run on the host and work, as ADR-0319's Consequences say. A fix might write the request through the agent cache folder, or have the command, in a container, say to use the MCP tool.

## Remediation
