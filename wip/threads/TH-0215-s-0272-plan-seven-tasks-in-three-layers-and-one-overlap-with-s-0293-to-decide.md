---
id: TH-0215
title: "S-0272 plan: seven tasks in three layers, and one overlap with S-0293 to decide"
anchor:
  path: wip/kanban/stories/S-0272-an-agent-with-an-open-question-ends-instead-of-waiting-flai-serve-restarts-it-on-the-answer-and-wait-for-events-keeps-a-timeout-only-for-an-agent-with-work-in-hand.md
  item: S-0272
status: answered
participants: [planner-S-0272, alex]
created: 2026-10-06T22:54:30Z
updated: 2026-10-06T23:34:52Z
---

# TH-0215 S-0272 plan: seven tasks in three layers, and one overlap with S-0293 to decide

On wip/kanban/stories/S-0272-an-agent-with-an-open-question-ends-instead-of-waiting-flai-serve-restarts-it-on-the-answer-and-wait-for-events-keeps-a-timeout-only-for-an-agent-with-work-in-hand.md.

## Entries

### 2026-10-06T22:54:30Z planner-S-0272
I've drafted seven tasks for S-0272 in three layers. I've also revisited its touches, forecast and cost of delay. The story's words are unchanged; only its `### Planning` notes were rewritten.

**The plan:**

| Layer | Task | What | Waits for |
|------:|------|------|-----------|
| 1 | T-1087 | `wait_for_events` answers `end: true` and `why` (MCP server, tool description, instructions) | none |
| 1 | T-1090 | The story agent's prompt: write the narrative's state and end on a question; end on `end` | none |
| 1 | T-1094 | `work-management.md`, `delegation.md`, their template copies, and `CHANGELOG` | none |
| 1 | T-1097 | ADR and `metrics.md`: the definition of an empty wake | none |
| 2 | T-1101 | `flai stats` counts empty wakes per story and for the window | T-1097 |
| 2 | T-1104 | End-to-end test: end on a question, restart on the answer, the answer in the first `inbox` | T-1087 |
| 3 | T-1112 | The design (`workflow.md`, `agent-narrative.md`, `flai-cli.md`) and the user guide | all the others |

**Figures:**

- **Forecast:** 50m, delivery 2026-10-07T08:45Z. flai gave 39m. I kept the earlier 50m for the seven tasks, the ADR and a test that spans serve and MCP.
- **Cost of delay:** 96 USD a week, unchanged. flai gives 115.98. I kept the epic planner's share by the turns each story removes, so the siblings still add up to 900.
- **Touches:** I added `workflow.md`, `agent-narrative.md` and `usage/log_test.go`. `design/adrs` stays a folder touch because the ADR's name is not known yet.

**Assumptions:**

1. The caller's story comes from `FLAI_STORY`, which flai serve sets. An agent you run by hand, which serve would not restart, keeps its timeout. So do the planner, the orchestrator and the analyzer, whose prompts still say to hold for an answer.
2. `end` comes only after any events already behind the cursor have been reported, so an answer that has just arrived is never swallowed.
3. An empty wake is a main-agent `wait_for_events` that timed out with no events and no changed paths.
4. Serve's half of the restart already exists and is tested (`TestAnAgentThatEndedAskingIsStartedAgainWhenAnswered`, in `serve/agents.go`). T-1104 adds the MCP half on the same lab. `serve/restart.go`, declared earlier, is the operator's restart and no task changes it. I kept it only because it was declared.

**Proposals, for your word:**

1. **Criterion 4 overlaps S-0293**, which classifies every turn, empty wakes among them. TH-0176 proposed moving it there, and the thread was resolved with the criterion still here. So I planned it here: T-1097 and T-1101, about 10 of the 50m. My recommendation is to keep it here, so the saving is measured as soon as this story lands, with S-0293 building on T-1101's parsing. If you would rather move it, say so: I'd propose dropping criterion 4, cancelling T-1097 and T-1101, and taking T-1101 off T-1112's `after`.
2. **S-0294** (a draft) proposes that serve restart an agent that ended without finishing. That is next to this story's restart on an answer, but it is a separate decision under ADR-0043. I'd keep them apart. No change to either.

### 2026-10-06T23:34:52Z alex
accept proposals
