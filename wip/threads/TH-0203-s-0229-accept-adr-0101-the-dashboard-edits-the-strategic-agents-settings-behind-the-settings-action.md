---
id: TH-0203
title: "S-0229: accept ADR-0101, the dashboard edits the strategic agents' settings behind the settings action?"
anchor:
  path: wip/kanban/stories/S-0229-the-strategic-agents-settings-are-edited-in-the-dashboard-permissions-policy-release-policy-schedules-and-agents.md
  item: S-0229
status: resolved
participants: [agent-S-0229, alex, agent-S-0277]
created: 2026-10-06T22:20:23Z
updated: 2026-10-07T03:06:18Z
---

# TH-0203 S-0229: accept ADR-0101, the dashboard edits the strategic agents' settings behind the settings action?

On wip/kanban/stories/S-0229-the-strategic-agents-settings-are-edited-in-the-dashboard-permissions-policy-release-policy-schedules-and-agents.md.

## Entries

### 2026-10-06T22:20:23Z agent-S-0229
ADR-0101 (proposed, on story/S-0229: `design/adrs/0101-with-the-settings-host-action-on-the-dashboard-edits-the-strategic-agents.md`) refines ADR-0039: with the `settings` host action on for a project, its dashboard may also change the manifest's `orchestration`, `planning` (not `currency`), and `analysis` keys through `flai manifest set`, committed with the dashboard's trailer. The consequence to weigh: a dashboard token can then grant the orchestrator `accept_reviews` and `publish`. Each permission's risk is shown beside it, and only the shell turns `settings` on or off. This is assumption 6 on the plan thread TH-0149, which you resolved.

Recommended: accept it as written. Reply `accept` and I will set it accepted. Otherwise it stays proposed, and the story goes to review with it proposed.

### 2026-10-06T22:51:48Z alex
accept

### 2026-10-07T03:06:18Z agent-S-0277
Resolved: S-0229 was accepted; left open by a flai older than ADR-0109 (S-0277)
