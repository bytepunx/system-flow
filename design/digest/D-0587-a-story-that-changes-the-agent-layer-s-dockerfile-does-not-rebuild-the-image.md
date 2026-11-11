---
id: D-0587
title: "A story that changes the agent layer's Dockerfile does not rebuild the image under its :local tag, which running agents use; it checks the change on the built image with the environment overridden"
kind: decision
sources: [S-0686]
topics: [docker, harness]
paths: [flai/images/system-flow/Dockerfile, scripts/agent-image.sh]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0587 A story that changes the agent layer's Dockerfile does not rebuild the image under its :local tag, which running agents use; it checks the change on the built image with the environment overridden

S-0686 dropped `GOFLAGS=-buildvcs=false` from the project's layer, `flai/images/system-flow/Dockerfile`. It did not rebuild the image to check the change. From S-0686's Decisions (`wip/archive/agents/S-0686.md`):

```text
The layer image is not rebuilt under its :local tag, which running agents use; the effect of dropping GOFLAGS is checked with GOFLAGS emptied on the built image, which is the same environment.
```

The reason: `scripts/agent-image.sh` tags the layer `flai-agent-sf-claude-code:local`. That is the image `flai serve`'s container harness starts every story's agent from, so a rebuild in one story changes what the other stories in progress run on. When a Dockerfile change only sets or drops an environment variable, run the built image with that variable overridden (`docker run -e GOFLAGS= ...`). That gives the same environment without a rebuild. The new image reaches the agents when the operator rebuilds the images after acceptance (`make agent-images`).

S-0686 also put its Docker-needing test in a script of its own, `scripts/agent-vcs-test.sh`, run by `scripts/smoke.sh`. It did not extend `scripts/agent-container-test.sh`, which runs with no Docker. The new script skips when Docker or the layer image is missing.

Sources: S-0686.
