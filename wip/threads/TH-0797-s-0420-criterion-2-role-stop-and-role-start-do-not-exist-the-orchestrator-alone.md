---
id: TH-0797
title: "S-0420 criterion 2: role.stop and role.start do not exist; the orchestrator alone has a hold"
anchor:
  path: wip/kanban/stories/S-0420-the-dashboard-s-planner-orchestrator-and-analyzer-pages-become-one-role-page-per-defined-role-with-its-settings-panel-its-activity-document-and-its-runs-and-agent-status-carries-every-role-by-name.md
  item: S-0420
status: resolved
participants: [agent-S-0420, orchestrator, alex]
created: 2026-11-02T00:16:23Z
updated: 2026-11-02T00:47:58Z
---

# TH-0797 S-0420 criterion 2: role.stop and role.start do not exist; the orchestrator alone has a hold

On wip/kanban/stories/S-0420-the-dashboard-s-planner-orchestrator-and-analyzer-pages-become-one-role-page-per-defined-role-with-its-settings-panel-its-activity-document-and-its-runs-and-agent-status-carries-every-role-by-name.md.

## Entries

### 2026-11-02T00:16:23Z agent-S-0420
Criterion 2 says the role page offers **Run**, **Stop**, and **Start** "each through `role.run`, `role.stop`, and `role.start`". Only `role.run` exists. S-0418 dropped `role.stop`, `role.start`, and the `continuous` trigger (ADR-0209 §1). The only hold today is the orchestrator's, through `orchestrate.stop` and `orchestrate.start` (S-0228).

**Recommendation:** the role page offers **Run** for any role through `role.run`, and **Stop** and **Start** only for a role that has a hold, which today is the orchestrator, through `orchestrate.stop` and `orchestrate.start`. I would reword criterion 2 to say so and tick it on that reading. A hold for every role would be a story of its own: it needs `flai serve role stop|start`, a hold every event and schedule start obeys, and host methods. That is outside S-0420's claim (`flai/internal/hostapi/writes.go`, `flai/internal/serve/role*.go`, `flai/cmd/serve*.go`).

**Alternative:** build `role.stop` and `role.start` in S-0420 as a new task: a hold per role kept in flai serve's state, which the event and schedule starts obey, and the two host methods. This widens the touches and adds about an hour.

Until you answer, I go on with everything else, and build the page so that Stop and Start show where `agent.status` says a role can be held.

### 2026-11-02T00:17:47Z orchestrator (recommendation)
**Recommended:** take the agent's recommendation. Reword S-0420 criterion 2 as follows. The role page offers **Run** for any role through `role.run`. It offers **Stop** and **Start** only for a role that has a hold, shown where `agent.status` says the role can be held. Today that is the orchestrator alone, through `orchestrate.stop` and `orchestrate.start`. S-0420's agent then ticks the criterion on that reading.

## Why

- **The ADR keeps the only hold.** ADR-0209 § Decision point 1 makes every role's trigger a list of events and rejects `continuous`. Point 7 keeps `flai serve orchestrate stop|start` as the one hold. D-0425 (TH-0749) dropped `role stop|start` from S-0418 for this reason. Criterion 2 was written before that and names methods that no longer exist.
- **A per-role hold is a decision not yet recorded.** It would need a hold that every event and schedule start obeys, `flai serve role stop|start`, and two host methods. That belongs in its own story, with an ADR if wanted. It does not fit S-0420's claim.
- **The page is ready for a hold later.** The agent is building Stop/Start to show wherever `agent.status` reports `held`. A later per-role hold needs no page change.

Rewording a finalized story's criterion changes its scope, so this is a recommendation for you rather than an answer.

Source: design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md § decision

### 2026-11-02T00:40:41Z alex
Confirmed the recommendation of 2026-11-02T00:17:47Z orchestrator.

Source: design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md § decision

### 2026-11-02T00:47:58Z agent-S-0420
Resolved: Criterion 2 reworded as confirmed: Run for any role through role.run; Stop and Start only where a role has a hold, today the orchestrator through orchestrate.stop and orchestrate.start.
