---
id: T-2898
type: task
nature: improvement
title: scripts/agent-container.sh --probe runs the image as probe-S-nnnn with a probe start request
status: done
parent: S-0682
owner: alex
created: 2026-11-10T08:46:14Z
updated: 2026-11-10T09:08:14Z
transitions:
  - to: ready
    at: 2026-11-10T08:59:43Z
    by: agent-S-0682
  - to: in-progress
    at: 2026-11-10T08:59:44Z
    by: agent-S-0682
  - to: done
    at: 2026-11-10T09:08:14Z
    by: agent-S-0682
stream: S-0682
tags: [scripts, containers]
touches: [scripts/agent-container.sh]
after: [T-2896]
usage:
  source: log
  seconds: 485
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 3
      output: 27
      cache_read: 439656
      cache_write: 5721
      cost: 0.1639
---
# T-2898 scripts/agent-container.sh --probe runs the image as probe-S-nnnn with a probe start request

## Work

Waits for T-2896: pilot decodes the start request with unknown fields refused, so the `probe` field this script writes exists only once `harness.Request` has it.

- Take `--probe` before the story: `scripts/agent-container.sh --probe S-nnnn` and `scripts/agent-container.sh --probe S-nnnn -- cmd [arg]...`. The usage line and the header comment say so. Without `--probe` the script runs as it does now.
- With `--probe`, name the agent `probe-S-nnnn` in place of `agent-S-nnnn`, write `"probe":true` in the start request, always start a new session (ignore `FLAI_AGENT_SESSION`, which would resume the story's), and in the `--` mode pass `FLAI_AGENT=probe-S-nnnn` and no `FLAI_STORY`. The mounts stay the story's, so the probe sees what the story's agent would.
- Say on standard error, before the run, that it is a probe that works no item.
- `scripts/agent-container.sh` is on the container harness's protected list (S-0508): the write asks the story's owner through `permission_prompt` unless auto-approve is on, and the story is accepted by the operator.
- Run `flai test scripts/agent-container.sh` (ShellCheck).

## Done when

- `--probe` makes the agent `probe-S-nnnn`, the request carry `"probe":true`, the session new, and the `--` mode run without `FLAI_STORY`; a run without `--probe` is unchanged.
- ShellCheck passes on the script.

## Notes
