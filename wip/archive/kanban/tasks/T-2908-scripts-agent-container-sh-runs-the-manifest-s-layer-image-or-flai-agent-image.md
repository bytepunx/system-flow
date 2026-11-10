---
id: T-2908
type: task
nature: remediation
title: scripts/agent-container.sh runs the manifest's layer image, or FLAI_AGENT_IMAGE, and refuses an image not built
status: done
parent: S-0670
owner: alex
created: 2026-11-10T08:47:36Z
updated: 2026-11-10T09:29:36Z
transitions:
  - to: ready
    at: 2026-11-10T09:28:31Z
    by: agent-S-0670
  - to: in-progress
    at: 2026-11-10T09:28:32Z
    by: agent-S-0670
  - to: done
    at: 2026-11-10T09:29:36Z
    by: agent-S-0670
stream: S-0670
tags: [scripts, docker]
touches: [scripts/agent-image.sh, scripts/agent-container.sh]
usage:
  source: log
  seconds: 64
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 2142
      cache_read: 793247
      cache_write: 19049
      cost: 0.3492
---
# T-2908 scripts/agent-container.sh runs the manifest's layer image, or FLAI_AGENT_IMAGE, and refuses an image not built

## Work

The fix for I-0279. `scripts/agent-container.sh` sets `IMAGE="flai-agent-claude-code:local"` and never reads `agent.image.layer`, so a close-out run by hand in it lacks Go, golangci-lint, and ShellCheck.

- `scripts/agent-image.sh` gains `--tag`, which prints the image a story's agent runs in and exits: the layer's tag, `flai-agent-<key>-claude-code:local` with the key in lower case, as `build_layer` and `manifest.LayerImage` name it, when the manifest names `agent.image.layer`; else `flai-agent-claude-code:local`. It reuses `manifest_layer` and the key's reading, so the two scripts cannot name the tag apart. It builds nothing.
- `scripts/agent-container.sh` takes `IMAGE` from `FLAI_AGENT_IMAGE` when it is set, else from `scripts/agent-image.sh --tag`. It gains `--image`, which prints that image and exits before it needs a story, a worktree, or a server, so a test can read the choice without Docker.
- Before `docker run`, it refuses an image `docker image inspect` does not find, naming the command that builds it: `scripts/agent-image.sh --layer` for the layer, `make agent-images` otherwise.
- The header comments of both scripts name the new flags and `FLAI_AGENT_IMAGE`.

Waits for nothing: it is the first layer.

## Done when

- [ ] In this repository, `scripts/agent-container.sh --image` prints `flai-agent-sf-claude-code:local`, and with `FLAI_AGENT_IMAGE=x` prints `x`
- [ ] `scripts/agent-image.sh --tag` prints the same tag and builds nothing
- [ ] A run whose image is not built stops before `docker run`, naming the build command
- [ ] `flai test scripts/agent-image.sh scripts/agent-container.sh` passes, ShellCheck included

## Notes
