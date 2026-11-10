---
id: T-2937
type: task
nature: remediation
title: The orchestrator's role text and the strategic-agents convention, in every copy, accept through the MCP tool accept and verify through the MCP tool verify
status: in-progress
parent: S-0687
owner: alex
created: 2026-11-10T08:49:53Z
updated: 2026-11-10T09:25:50Z
transitions:
  - to: ready
    at: 2026-11-10T09:25:48Z
    by: agent-S-0687
  - to: in-progress
    at: 2026-11-10T09:25:50Z
    by: agent-S-0687
stream: S-0687
tags: [roles, conventions, template]
touches: [design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md]
after: [T-2934]
---
# T-2937 The orchestrator's role text and the strategic-agents convention, in every copy, accept through the MCP tool accept and verify through the MCP tool verify

## Work

Waits for T-2934, whose tool's name and inputs this text names. In the orchestrator's role file, in its three copies (`design/workflows/roles/orchestrator.md`, `template/root/design/workflows/roles/orchestrator.md`, `flai/internal/roles/baseline/orchestrator.md`), add `mcp__flai__accept` to its tools and rewrite the `accept_reviews` permission, the accept step, and step 12: verify with the MCP tool `verify` and read the result with `flai verify <S-nnnn> --last --json`, in place of a shell `flai verify` that runs tiers; read the blockers with the MCP tool `accept` and `dry_run`, and accept with it, giving the commit and the evidence, in place of `flai accept --dry-run` and `--evidence -` from the shell. Correct the text that says the MCP `verify` answers the commit verified, which it does not. Keep the `{{if .Queue}}` paragraph that turns to `queue_add`. Make the same change in § Accepting a story of the strategic-agents convention, both copies, above the marker in the template's as the baseline. Add a line for the template's change to `template/CHANGELOG.md` under its top `## Unreleased`. Keep the three role copies identical where the baseline test compares them.

## Done when

- [ ] Every copy of the orchestrator's role text and of the convention names the MCP tools `accept` and `verify` for verifying and accepting, and no shell `flai accept` or tier-running `flai verify`
- [ ] `template/CHANGELOG.md` says so under `## Unreleased`
- [ ] The roles tests and the markdown lint pass with `flai test` on the changed paths

## Notes
