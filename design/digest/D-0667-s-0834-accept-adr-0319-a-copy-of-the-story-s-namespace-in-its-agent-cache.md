---
id: D-0667
title: "S-0834: accept ADR-0319, a copy of the story's namespace in its agent cache folder, before any code is written: Recommended: accept ADR-0319 as written (`story/S-0834`, eaa744e87), so that T-3547 to…"
kind: answer
sources: [TH-1290]
topics: [kluster, harness, flai, cli, go, code, all]
paths: [design/adrs/README.md, design/system/kluster.md, flai/internal/kluster/agentcopy.go, flai/internal/kluster/agentcopy_test.go, flai/internal/kluster/agentconfig.go, flai/internal/kluster/agentconfig_test.go, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go, scripts/agent-container.sh, docs/operators/settings.md, flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/cmd/kluster_container_test.go, docs/users/flai-reference.md, design/system/agent-adapters.md, docs/users/flai.md, design/issues/I-0361-flai-kluster-exec-and-log-read-flai-cache-kluster-which-a-story-s-agent-in-a.md, design/issues/summary.md, design/adrs/0319-flai-copies-a-story-namespace-s-record-step-logs-state-and-kubeconfig-into-the.md]
date: 2026-11-18
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0667 S-0834: accept ADR-0319, a copy of the story's namespace in its agent cache folder, before any code is written: Recommended: accept ADR-0319 as written (`story/S-0834`, eaa744e87), so that T-3547 to…

The operator, alex, answered TH-1290 at 2026-11-18T07:10:25Z.

What was asked:

```text
**Recommended: accept ADR-0319 as written.** T-3546 needs your acceptance before T-3547 to T-3550 build it, as I-0361 asks.
```

The answer:

```text
**Recommended: accept ADR-0319 as written** (`story/S-0834`, eaa744e87), so that T-3547 to T-3550 can build it.

**Why:**

- ADR-0292 §2 already lets a story's agent and its sub-agents run `flai kluster env`, `exec`, and `log` for their own namespace. In a container those commands find nothing (I-0361). ADR-0319 makes that grant work there, and it grants nothing new.
- It keeps ADR-0292's boundary. The copy holds only the namespace's own kubeconfig, which the API server limits to that namespace. It never holds flai's cluster kubeconfig, and the guard's refusal of `.flai-cache/kluster/kubeconfig` is unchanged.
- It keeps ADR-0253's mount list. The copy goes in `.flai-cache/agents/S-nnnn`, which each story's container mounts already. The ADR's Alternatives explain why a bind mount of the story's own files would go stale, since they are written by rename.
- `tls-server-name` keeps the API server's certificate verified when the host name is rewritten.

**Why I am asking you:** this is a new ADR, so it is a decision no one has recorded yet. The agent names one risk, and its tests cannot cover it here: `bridge` verification depends on k3d's certificate. That risk sits in the ADR's Consequences.

**Alternatives:** change the ADR (say what), or take one of its five alternatives, such as mounting `.flai-cache/kluster` read-only. I would not take that one, because it exposes the cluster's kubeconfig and every namespace's token.

Source: design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md § Decision
```

Source: TH-1290.
