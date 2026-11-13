---
id: T-3207
type: task
nature: improvement
title: system-flow.yaml names this repository's layer build arguments the general way, and flai image build --layer builds the layer as scripts/agent-image.sh does
status: done
parent: S-0760
owner: alex
created: 2026-11-12T16:06:38Z
updated: 2026-11-13T19:48:32Z
transitions:
  - to: ready
    at: 2026-11-13T19:48:12Z
    by: agent-S-0760
  - to: in-progress
    at: 2026-11-13T19:48:12Z
    by: agent-S-0760
  - to: done
    at: 2026-11-13T19:48:32Z
    by: agent-S-0760
stream: S-0760
tags: [cli]
touches: [system-flow.yaml, flai/images/system-flow/Dockerfile]
usage:
  source: log
  seconds: 20
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 1301
      cache_read: 256282
      cache_write: 6502
      cost: 0.1198
---
# T-3207 system-flow.yaml names this repository's layer build arguments the general way, and flai image build --layer builds the layer as scripts/agent-image.sh does

## Work

Write the layer's build arguments into `system-flow.yaml`'s `agent.image` with the flai command S-0745 gives, never by hand: the manifest's other keys are flai's. Use the values `scripts/agent-image.sh` passes today:

- `GO_VERSION` from `flai/go.mod`;
- `GOLANGCI_LINT_VERSION`, `SHELLCHECK_VERSION`, `PNPM_VERSION`, and `GORELEASER_VERSION` from `scripts/install-tools.sh`.

Give `flai/images/system-flow/Dockerfile` the same values as `ARG` defaults, so a build with none still works. Say in the task's Notes how the pins are kept in step with `scripts/install-tools.sh` and `flai/go.mod`; `dependabot` moves the Dockerfile's `FROM` alone.

`flai/images` is a protected folder for a container's agent, so the write there asks the owner. It is the first layer.

## Done when

- `flai check --strict` passes with the arguments in the manifest.
- The docker arguments `flai image build --layer --dry-run`, or S-0745's test, records match `scripts/agent-image.sh --layer`'s.

## Notes
