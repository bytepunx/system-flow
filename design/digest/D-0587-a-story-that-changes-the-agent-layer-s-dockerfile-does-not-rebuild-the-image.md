---
id: D-0587
title: "A story that changes the agent layer's Dockerfile does not rebuild the image under its :local tag, which running agents use; it checks the change on the built image with the environment overridden"
kind: decision
sources: [S-0686, S-0740, S-0760]
topics: [docker, harness, release]
paths: [flai/images/system-flow/Dockerfile, scripts/agent-image.sh, flai/internal/agentimage/build.go, flai/cmd/image.go]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: [S-0760]
---

# D-0587 A story that changes the agent layer's Dockerfile does not rebuild the image under its :local tag, which running agents use; it checks the change on the built image with the environment overridden

S-0686 dropped `GOFLAGS=-buildvcs=false` from the project's layer, `flai/images/system-flow/Dockerfile`. It did not rebuild the image to check the change. From S-0686's Decisions (`wip/archive/agents/S-0686.md`):

```text
The layer image is not rebuilt under its :local tag, which running agents use; the effect of dropping GOFLAGS is checked with GOFLAGS emptied on the built image, which is the same environment.
```

The reason: `scripts/agent-image.sh` tags the layer `flai-agent-sf-claude-code:local`. That is the image `flai serve`'s container harness starts every story's agent from, so a rebuild in one story changes what the other stories in progress run on. When a Dockerfile change only sets or drops an environment variable, run the built image with that variable overridden (`docker run -e GOFLAGS= ...`). That gives the same environment without a rebuild. The new image reaches the agents when the operator rebuilds the images after acceptance (`make agent-images`).

S-0686 also put its Docker-needing test in a script of its own, `scripts/agent-vcs-test.sh`, run by `scripts/smoke.sh`. It did not extend `scripts/agent-container-test.sh`, which runs with no Docker. The new script skips when Docker or the layer image is missing.

S-0740 held to the same rule for `flai image ensure` and `flai image build`, which tag the harness image `flai-agent-claude-code:local`, the name the start runs. From S-0740's Decisions (`wip/archive/agents/S-0740.md`):

```text
Criteria ticked on test evidence over a fake docker and an httptest release; a real pull and build were not run here: no published release holds the agent image lists yet (S-0734 is accepted, not released), and a real build replaces the image the running agents use, which is the operator's to choose.
```

So a story that changes how the images are pulled or built proves it over a fake docker and a test release, and leaves the real run on the host to the operator.

S-0760, which changed `scripts/agent-image.sh` to build and label the images, held to it too. From its Decisions (`wip/archive/agents/S-0760.md`):

```text
No real image build in this story: rebuilding the :local images changes what running agents use, which is the operator's to choose (D-0587). The script is tested over a fake docker.
```

Sources: S-0686, S-0740, S-0760.
