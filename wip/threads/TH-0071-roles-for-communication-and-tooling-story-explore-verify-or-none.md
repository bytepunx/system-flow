---
id: TH-0071
title: "Roles for communication and tooling: [story, explore, verify] or none?"
anchor:
  path: wip/kanban/stories/S-0196-a-convention-s-roles-name-every-agent-that-reads-it-the-story-s-agent-included.md
  item: S-0196
status: resolved
participants: [agent-S-0196, alex]
created: 2026-10-02T23:25:37Z
updated: 2026-10-02T23:41:35Z
---

# TH-0071 Roles for communication and tooling: [story, explore, verify] or none?

On wip/kanban/stories/S-0196-a-convention-s-roles-name-every-agent-that-reads-it-the-story-s-agent-included.md.

## Entries

### 2026-10-02T23:25:37Z agent-S-0196
S-0196's third criterion leaves two choices open. Under ADR-0068 the two options select the same agents today (story, explore, verify all read the file). They differ only for agents added later, such as the orchestrator, planner, or analyzer, which check will accept:

- `communication.md`: `[story, explore, verify]`, or no roles (every agent, the strategic agents included)
- `tooling.md`: `[story, explore, verify]`, or no roles

My recommendation is **no roles for both**. Communication (how to write and report) and tooling (how to run flai and the scripts) apply to any agent, the planner, orchestrator, and analyzer included, the same way safety and documentation do. An explicit list would have to be extended each time a new agent arrives.

Everything else goes in as the criterion lists it: code-quality [story, verify]; continuous-improvement, decisions, git, session-start, and work-management [story]; delegation [story, explore, verify]; logging [story, verify]; documentation, safety, telemetry, and README with no roles. One consequence to check: git.md loses `verify`, so the verifier no longer reads git.md in its pack. Reply "none", "list", or a value for each.

### 2026-10-02T23:41:24Z alex
I agree - no roles is correct so that all agents follow the tooling and comms guidnace.

### 2026-10-02T23:41:35Z agent-S-0196
Resolved: The designer chose no roles for communication and tooling, so every agent follows them.
