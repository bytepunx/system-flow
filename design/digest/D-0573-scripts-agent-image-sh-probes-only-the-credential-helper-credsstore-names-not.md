---
id: D-0573
title: scripts/agent-image.sh probes only the credential helper credsStore names, not credHelpers, since the base images come from Docker Hub, which credsStore serves
kind: decision
sources: [S-0713, I-0339]
topics: [docker, harness]
paths: [scripts/agent-image.sh, scripts/agent-image-test.sh, docs/operators/runbooks/agent-container.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0573 scripts/agent-image.sh probes only the credential helper credsStore names, not credHelpers, since the base images come from Docker Hub, which credsStore serves

S-0713 fixed I-0339: `make agent-images` stopped at the first base image when Docker Desktop's credential helper, `desktop.exe`, could not start from WSL. Before its first build, `scripts/agent-image.sh` now asks the helper `credsStore` names to list its logins, ten seconds at most, and when it is missing, fails, or hangs, builds with no login from a scratch `DOCKER_CONFIG`. `scripts/README.md` and `docs/operators/runbooks/agent-container.md` describe that fallback; neither says why `credHelpers` is left out.

From S-0713's Decisions (`wip/archive/agents/S-0713.md`):

```text
Only `credsStore` is probed, not `credHelpers`. I-0339 names only `credsStore`, and the base images come from Docker Hub, which `credsStore` serves.
```

A later story that meets a per-registry helper under `credHelpers` failing the same way, for a base image from another registry, would widen the probe; until then it is left out on purpose. The same story also chose that only the build paths fall back, the default build and `--layer`, since `--tag` and `--check` pull nothing, and that `scripts/agent-image-test.sh` does not cover a helper that hangs: it would add ten seconds to every smoke run, and a hang reaches the same path as an exit 1, since `timeout` exits 124.

Sources: S-0713, I-0339.
