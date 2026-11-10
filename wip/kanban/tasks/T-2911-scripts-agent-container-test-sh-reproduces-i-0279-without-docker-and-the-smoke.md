---
id: T-2911
type: task
nature: remediation
title: scripts/agent-container-test.sh reproduces I-0279 without Docker, and the smoke tier runs it
status: backlog
parent: S-0670
owner: alex
created: 2026-11-10T08:47:47Z
updated: 2026-11-10T08:47:47Z
transitions: []
stream: S-0670
tags: [scripts, test]
touches: [scripts/agent-container-test.sh, scripts/smoke.sh]
after: [T-2908]
---
# T-2911 scripts/agent-container-test.sh reproduces I-0279 without Docker, and the smoke tier runs it

## Work

The test criterion 1 asks for. A POSIX `sh` script under `set -eu`, `scripts/agent-container-test.sh`, that needs no Docker, no MCP server, and no worktree:

- In this repository, `scripts/agent-container.sh --image` prints the layer's tag, `flai-agent-sf-claude-code:local`, the case I-0279 records: before the fix it printed `flai-agent-claude-code:local`.
- With `FLAI_AGENT_IMAGE` set, it prints that value.
- In a scratch tree under `mktemp -d`, holding copies of `scripts/env.sh`, `scripts/agent-image.sh`, and `scripts/agent-container.sh` beside a `system-flow.yaml` that names no `agent.image.layer`, it prints `flai-agent-claude-code:local`; with a layer and an upper-case key, it prints the tag with the key in lower case.
- Each case prints what it expected and what it got when it fails, and the script exits non-zero.

`scripts/smoke.sh` runs it beside `scripts/ci-pins.sh`, with an `echo "smoke: ..."` line as its other steps have, so a change to either script, which the smoke tier's `scripts/**/*.sh` selects, runs it.

Waits for T-2908, whose `--image` flag and image choice it tests.

## Done when

- [ ] `scripts/agent-container-test.sh` passes on the fixed scripts
- [ ] It fails, naming the case, when `scripts/agent-container.sh` is given back its fixed `IMAGE="flai-agent-claude-code:local"`, checked by hand and recorded in this task's Notes
- [ ] `scripts/smoke.sh` runs it, and `flai test scripts/agent-container-test.sh scripts/smoke.sh` passes

## Notes
