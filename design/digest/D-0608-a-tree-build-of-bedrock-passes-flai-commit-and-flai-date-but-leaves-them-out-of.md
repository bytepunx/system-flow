---
id: D-0608
title: A tree build of bedrock passes FLAI_COMMIT and FLAI_DATE but leaves them out of the image's recipe, since ADR-0295 does not compare a development image's commit and the commit has a label of its own
kind: decision
sources: [S-0740]
topics: [docker, cli, code]
paths: [flai/internal/agentimage/build.go, flai/internal/agentimage/agentimage.go, scripts/agent-image.sh]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: [S-0760]
---

# D-0608 A tree build of bedrock passes FLAI_COMMIT and FLAI_DATE but leaves them out of the image's recipe, since ADR-0295 does not compare a development image's commit and the commit has a label of its own

S-0740 built the tree build of the agent images, bedrock's `tree` target from the tree's own Dockerfiles and source. From S-0740's Decisions (`wip/archive/agents/S-0740.md`):

```text
The build labels each image with dev.system-flow.recipe, and claude-code with dev.system-flow.base, bedrock's image ID. FLAI_COMMIT and FLAI_DATE are passed in a tree build but left out of the recipe, since ADR-0295 §3 does not compare the commit and the commit has its own label.
```

`Recipe` (`flai/internal/agentimage/agentimage.go`) itself leaves out only `BASE` and `FLAI_VERSION`, as `design/system/flai-cli.md` says. The tree build (`flai/internal/agentimage/build.go`) keeps `FLAI_COMMIT` and `FLAI_DATE` out of the argument list it hands to the recipe, while still passing them to `docker build`. So a new commit changes no recipe and makes no development image stale; `dev.system-flow.commit` records it instead.

A build of its own that must produce images `flai image status` reports current, such as `scripts/agent-image.sh` under S-0760 or an operator's build command under S-0759, should compute the recipe the same way: without `BASE`, `FLAI_VERSION`, `FLAI_COMMIT`, or `FLAI_DATE`.

Sources: S-0740.
