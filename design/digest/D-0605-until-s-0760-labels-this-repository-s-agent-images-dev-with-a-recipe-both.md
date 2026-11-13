---
id: D-0605
title: Until S-0760 labels this repository's agent images dev with a recipe, both report stale on this host, so flai check's agent.image stale warning and flai serve's container agent image stale log line are expected, and make agent-images does not clear them
kind: decision
sources: [S-0735, S-0760]
topics: [docker, cli, code]
paths: [scripts/agent-image.sh, flai/internal/check/image.go, flai/internal/serve/container.go, flai/internal/agentimage/status.go]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0605 Until S-0760 labels this repository's agent images dev with a recipe, both report stale on this host, so flai check's agent.image stale warning and flai serve's container agent image stale log line are expected, and make agent-images does not clear them

S-0735 built `flai image status` and the label judgement behind it (ADR-0295). It noted what that judgement finds on this host, from its Decisions (`wip/archive/agents/S-0735.md`):

```text
On this host `flai image status` now reports both agent images stale. They carry flai 1.86.12, from `scripts/agent-image.sh`'s `git describe`, against the tree's `dev`, and no recipe or base label. `make agent-images` does not clear that until S-0760 makes the script label `dev` and the recipe. Until then, `flai check` on the main checkout warns `agent.image` stale (advisory), and each container start logs `container agent image stale` and goes on.
```

So a story worked before S-0760 is accepted should not treat that warning or that log line as its own defect, and should not try to clear it by rebuilding the images. A start still refuses a missing image, naming `flai image status`. S-0760 changes `scripts/agent-image.sh` so that `make agent-images` builds images `flai image status` reports current.
