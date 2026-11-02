---
id: TH-0815
title: "S-0509 criterion 4: how to run the planner-in-a-container proof before this code is released"
anchor:
  path: wip/kanban/stories/S-0509-the-planner-the-orchestrator-the-analyzer-and-the-curator-run-in-a-container-when-their-harness-is-container-with-the-main-checkout-mounted-as-an-adr-settles-pilot-passing-their-role-and-their-activity-logged-as-on-the-host.md
  item: S-0509
status: resolved
participants: [agent-S-0509, orchestrator, alex]
created: 2026-11-02T18:18:27Z
updated: 2026-11-02T19:02:24Z
---

# TH-0815 S-0509 criterion 4: how to run the planner-in-a-container proof before this code is released

On wip/kanban/stories/S-0509-the-planner-the-orchestrator-the-analyzer-and-the-curator-run-in-a-container-when-their-harness-is-container-with-the-main-checkout-mounted-as-an-adr-settles-pilot-passing-their-role-and-their-activity-logged-as-on-the-host.md.

## Entries

### 2026-11-02T18:18:27Z agent-S-0509
Criterion 4 asks for a planner run in a container that plans a story and logs its activity in `wip/agents/planner.md`. flai serve cannot do that run during this story:

- The host runs the installed flai 1.80.2, which lacks this story's code. Its container adapter refuses a start with no story.
- `~/.flai/config.json` has no `container` harness entry (no login, no image). Setting it is a host-config write, which the guard refuses me (ADR-0220).
- Your rule is never to start a second flai serve.

**Recommended: option 1.**

1. **Hand run now, serve run after release.** I add a role form to `scripts/agent-container.sh`. It runs the planner on one backlog story you name, in the tree's `flai-agent-claude-code:local`, with the mounts the adapter builds and your login read-only. That proves the mounts, the guard, and the plan. A serve test feeds a recorded container log to the planner's run end and checks the activity entry's duration and cost. The flai-serve run itself becomes an issue for after publish, and criterion 4 is ticked only for what the hand run shows.
2. **You run it after acceptance.** I build and test everything with fakes. Criterion 4 stays unticked, with the reason in the story's notes. After you publish, set `flai serve agent harness container --login <path>` and the planner's harness, then ask for a plan.
3. **A second flai serve from the tree,** for one planner run. This breaks your rule. Choose it only if you lift the rule for this run.

For option 1, please name the backlog story to plan, or let me pick the first backlog story with no plan. A planner session costs about what a host planner run costs.

Layers 1 to 3 (ADR, adapter, start paths) do not wait on this. I am working them now.

### 2026-11-02T18:19:28Z orchestrator (recommendation)
**Recommended:** option 2, written as you answered for S-0508 on TH-0812 (D-0466): move the live planner run to a follow-up draft rather than leave criterion 4 unticked.

- S-0509 keeps what it can prove: the four start paths on a fake adapter, the argv and mounts per role, and the serve test that feeds a recorded container log to the planner's run end and checks the activity entry's duration and cost.
- S-0509's agent narrows criterion 4 to that, and writes a draft story in E-0029 with `after: S-0509`, holding the live part: `flai serve` starts a planner on `harness: container`, it plans one backlog story, and `wip/agents/planner.md` gains its entry. Its Notes say it waits for the release and for your host settings (`flai serve agent harness container --login <path>`).
- S-0509's Notes record the move.

**Why:** the same case as D-0466 and D-0461. The role start paths exist only on `story/S-0509`, and this host's `flai serve` is the installed 1.80.2. An unticked criterion would keep S-0509 from the definition of done.

**Why not option 1 now:** it mounts your Claude login into a container and spends a planner session, for a proof the follow-up gives on the real path. If you want it anyway, say so and name the story to plan. **Option 3** breaks your rule against a second `flai serve`; I don't recommend it.

**Why I recommend rather than answer:** it changes a finalized story's criteria and adds a story. That is your call, and so is finalizing the draft.

Source: design/digest/D-0466-s-0508-the-live-container-run-criterion-6-t-2169-s-first-half-cannot-run-before.md

### 2026-11-02T18:56:22Z alex
Confirmed the recommendation of 2026-11-02T18:19:28Z orchestrator.

Source: design/digest/D-0466-s-0508-the-live-container-run-criterion-6-t-2169-s-first-half-cannot-run-before.md

### 2026-11-02T19:02:24Z agent-S-0509
Resolved: Done as confirmed: criterion 4 now covers a planner's container run with docker faked, logged with duration and cost, plus the other three start paths on fakes. The live flai-serve planner run moved to draft S-0567 in E-0029, after S-0509. S-0509's Notes record the move.
