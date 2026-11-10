---
id: T-2914
type: task
nature: remediation
title: The agent-container runbook and scripts/README.md say the script runs the manifest's layer image
status: backlog
parent: S-0670
owner: alex
created: 2026-11-10T08:47:57Z
updated: 2026-11-10T08:47:57Z
transitions: []
stream: S-0670
tags: [docs, scripts]
touches: [docs/operators/runbooks/agent-container.md, scripts/README.md]
after: [T-2908]
---
# T-2914 The agent-container runbook and scripts/README.md say the script runs the manifest's layer image

## Work

Documentation changes with the behaviour, in the same story (`documentation.md`).

- `docs/operators/runbooks/agent-container.md`: its closing section says the close-out's tiers need the layer image, "which this script does not run". Say instead that the script runs `flai-agent-<key>-claude-code:local` when the manifest names `agent.image.layer`, `flai-agent-sf-claude-code:local` here, else `flai-agent-claude-code:local`, that `FLAI_AGENT_IMAGE` names another, and that `--image` prints the choice. Add the layer image to its images table and to the removal command, and name the refusal and its build command.
- `scripts/README.md`: the `agent-container.sh` row says it runs `flai-agent-claude-code:local`; give it the image choice, `FLAI_AGENT_IMAGE`, and `--image`. The `agent-image.sh` row gains `--tag`. A row for `agent-container-test.sh` says what it tests and that `scripts/smoke.sh` runs it.

Waits for T-2908, whose flags and refusal it documents. It runs beside T-2911: they share no path.

## Done when

- [ ] Neither document says the script runs `flai-agent-claude-code:local` alone, nor that it does not run the layer image
- [ ] Both name `FLAI_AGENT_IMAGE` and `--image`, and `scripts/README.md` has a row for `agent-container-test.sh`
- [ ] `flai test docs/operators/runbooks/agent-container.md scripts/README.md` passes the markdown lint

## Notes
