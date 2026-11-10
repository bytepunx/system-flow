---
id: T-2917
type: task
nature: remediation
title: The system-flow layer carries curl, openssl, and GoReleaser at the version scripts/install-tools.sh pins
status: done
parent: S-0672
owner: alex
created: 2026-11-10T08:48:03Z
updated: 2026-11-10T09:59:27Z
transitions:
  - to: ready
    at: 2026-11-10T09:55:10Z
    by: agent-S-0672
  - to: in-progress
    at: 2026-11-10T09:55:10Z
    by: agent-S-0672
  - to: done
    at: 2026-11-10T09:59:27Z
    by: agent-S-0672
stream: S-0672
tags: [docker, scripts]
touches: [flai/images/system-flow/Dockerfile, scripts/agent-image.sh]
usage:
  source: log
  seconds: 257
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 10
      output: 45
      cache_read: 1412762
      cache_write: 6580
      cost: 0.5224
---
# T-2917 The system-flow layer carries curl, openssl, and GoReleaser at the version scripts/install-tools.sh pins

## Work

In `flai/images/system-flow/Dockerfile`, add `curl` and `openssl` to the final stage's `apk add`: `install.sh` refuses without curl, and without openssl, which it and `scripts/install-test.sh` use to sign and verify the test releases. Add GoReleaser at `GORELEASER_VERSION`, v2.18.1 by default, with a `dev.system-flow.goreleaser` label beside the others. GoReleaser v2.18.1 builds only with Go 1.27, and the layer holds Go 1.26.0 with `GOTOOLCHAIN=local`. So take the release archive for `TARGETARCH`, checked against the release's `checksums.txt`, or build it in a stage of its own on a Go 1.27 image. Write the choice and its reason in the narrative's `## Decisions`.

In `scripts/agent-image.sh`, pass `GORELEASER_VERSION` from `scripts/install-tools.sh`'s pin as a build argument, as `build_layer` already passes the other pins.

Both paths are on the protected list (T-2157): write them yourself, not through a sub-agent. Waits for no task.

## Done when

- `scripts/agent-image.sh --layer` builds `flai-agent-sf-claude-code:local`.
- `docker run --rm --entrypoint sh flai-agent-sf-claude-code:local -c 'curl --version && openssl version && goreleaser --version'` exits 0, and GoReleaser reports the version `scripts/install-tools.sh` pins.
- ShellCheck passes on `scripts/agent-image.sh`.

## Notes
