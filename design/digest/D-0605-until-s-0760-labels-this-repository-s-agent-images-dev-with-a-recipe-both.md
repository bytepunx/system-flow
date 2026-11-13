---
id: D-0605
title: Until S-0760 labels this repository's agent images dev with a recipe, both report stale on this host, so flai check's agent.image stale warning and flai serve's container agent image stale log line are expected, and make agent-images does not clear them
kind: decision
sources: [S-0735, S-0760, S-0751]
topics: [docker, cli, code]
paths: [scripts/agent-image.sh, flai/internal/check/image.go, flai/internal/serve/container.go, flai/internal/agentimage/status.go, flai/internal/serve/imagejob.go]
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

So a story worked before S-0760 is accepted should not treat that warning or that log line as its own defect, and should not try to clear it by rebuilding the images. S-0760 changes `scripts/agent-image.sh` so that `make agent-images` builds images `flai image status` reports current.

S-0751 changed what a start does with a stale or missing image, once a flai with it runs `flai serve` on this host. A start in `flai serve` no longer logs the stale image and goes on: it asks for `flai serve`'s image job, which builds the images as `flai image ensure` does, and is held until the job ends. Outside `flai serve`, as for a start by hand through the CLI or MCP, a stale image is still warned of as `container agent image stale` and the start goes on, and a missing one refuses, naming `flai image ensure` and `flai image status` (`wip/archive/agents/S-0751.md`, `design/system/flai-cli.md` § The container harness (S-0508)). `flai check`'s `agent.image` warning now says that `flai serve` makes or renews the image at the next start that needs it.

Sources: S-0735, S-0760, S-0751.
