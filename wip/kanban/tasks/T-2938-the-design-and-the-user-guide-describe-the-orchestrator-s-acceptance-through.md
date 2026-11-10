---
id: T-2938
type: task
nature: remediation
title: The design and the user guide describe the orchestrator's acceptance through the host in a container, and I-0256 is closed
status: in-progress
parent: S-0687
owner: alex
created: 2026-11-10T08:49:59Z
updated: 2026-11-10T09:25:52Z
transitions:
  - to: ready
    at: 2026-11-10T09:25:51Z
    by: agent-S-0687
  - to: in-progress
    at: 2026-11-10T09:25:52Z
    by: agent-S-0687
stream: S-0687
tags: [docs, design, issues]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0256-a-strategic-agent-s-container-runs-flai-with-no-host-configuration-so-the.md, design/issues/summary.md]
after: [T-2934, T-2935]
---
# T-2938 The design and the user guide describe the orchestrator's acceptance through the host in a container, and I-0256 is closed

## Work

Waits for T-2934 and T-2935, whose tool and refusal this describes. In `design/system/strategic-agents.md`, rewrite the paragraph of § In a container (S-0509) that says the orchestrator's own `flai verify` and `flai accept` run inside its container: they now run on the host through the MCP tools `accept` and `verify`, and the shell's refuse there; change the I-0256 line under § The live run (S-0567) to say it is closed, and § Accepting a story (S-0221) to name the tool. In `design/system/flai-cli.md`, change the strategic container section's statement of the `FLAI_CONFIG` gap to the refusal and its reason. In `docs/users/flai.md`, add `accept` to the MCP tools table and say in § When the orchestrator accepts, and in the container harness section, which route it takes and what the shell's refusal says. Link the ADR T-2932 wrote. Then close the issue with `flai issue close I-0256 --reason` naming the MCP tool, the refusal, and the test that reproduces it, which also writes `design/issues/summary.md`.

## Done when

- [ ] The two design documents and the user guide describe the MCP route and the refusal, linking the ADR
- [ ] I-0256 is closed with a reason naming what fixed it, and `design/issues/summary.md` no longer lists it
- [ ] The markdown lint and `flai check --strict` pass

## Notes
