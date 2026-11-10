---
id: TH-1009
title: "S-0687 verified and ready; its acceptance stops on an MD051 finding in MS-0457, uncommitted in main"
anchor:
  path: wip/kanban/stories/S-0687-a-strategic-agent-s-container-runs-flai-with-no-host-configuration-so-the.md
  item: S-0687
status: resolved
participants: [orchestrator, alex]
created: 2026-11-10T09:53:28Z
updated: 2026-11-10T09:55:47Z
---

# TH-1009 S-0687 verified and ready; its acceptance stops on an MD051 finding in MS-0457, uncommitted in main

On wip/kanban/stories/S-0687-a-strategic-agent-s-container-runs-flai-with-no-host-configuration-so-the.md.

## Entries

### 2026-11-10T09:53:28Z orchestrator
**Recommended:** fix the link fragment in `wip/messages/MS-0457-s-0682-changed-scripts-agent-container-sh-and-its-runbook-it-is.md` line 53 in the main checkout, then tell me here, or run `flai accept S-0687 --by orchestrator --verified 8bb4d7f5f676cee6cde957019a05b0ccf24fdb8d` yourself. S-0687 itself has no blocker.

**What stopped it:** `flai accept` lints the main checkout's uncommitted files before it commits them. One finding is one its fix cannot correct:

```text
wip/messages/MS-0457-s-0682-changed-scripts-agent-container-sh-and-its-runbook-it-is.md:53:168 error MD051/link-fragments Link fragments should be valid [Context: "[Settings](#settings)"]
```

The line is agent-S-0670's entry of 2026-11-10T09:34:21Z. It quotes `[Settings](#settings)` from the runbook, and the message file has no `Settings` heading. Rewriting it as plain text, such as `the agent image (Settings)`, clears it. Nothing was merged: S-0687 stays in review at 8bb4d7f.

**What I checked:**

- `flai verify S-0687` passed every step at 8bb4d7f5f676cee6cde957019a05b0ccf24fdb8d, the head of `story/S-0687`. It reused eight tiers from the passing run at 524e080; main brought only `wip/` paths since.
- The verifier matched each criterion to the diff:
  - 1: the MCP tool `accept` (`flai/internal/mcpserver/accept.go`, `folder.go`, `flai/internal/hostapi/writes_outcome.go`, `flai/cmd/mcp.go`, `mcp_http.go`), the guard (`flai/internal/guard/*`), and the container refusals in `flai/cmd/accept.go` and `flai/cmd/verify.go`, reproduced by `TestTheOrchestratorInARolesContainerRefusesToAcceptByHand` and `TestVerifyRefusesTheOrchestratorsTiersInARolesContainer`; ADR-0278.
  - 2: `design/issues/I-0256-…md` closed with its reason, `design/issues/summary.md`.
- The dry run listed no blocker.

**Also worth a story:** flai's message write let a line through that its own lint now refuses. `tooling.md` says flai refuses a `wip` write that would bring a markdown finding.

### 2026-11-10T09:55:47Z alex
Resolved: S-0687 was accepted
