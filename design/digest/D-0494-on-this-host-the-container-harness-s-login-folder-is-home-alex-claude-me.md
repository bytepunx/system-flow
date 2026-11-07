---
id: D-0494
title: "On this host the container harness's login folder is /home/alex/.claude-me, neither ~/.claude nor ~/.claude-rmh, which the threads asking for it named"
kind: answer
sources: [TH-0843, S-0567, TH-0836]
topics: [harness, docker, planning]
paths: [design/system/strategic-agents.md, docs/operators/runbooks/agent-container.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0494 On this host the container harness's login folder is /home/alex/.claude-me, neither ~/.claude nor ~/.claude-rmh, which the threads asking for it named

Two threads asked the operator to set the container harness's login with `flai serve agent harness container --login <folder>` and guessed different folders: TH-0836 (S-0565) named `~/.claude`, and TH-0843 (S-0567) named `~/.claude-rmh`, reasoning that the host's claude-code harness runs `rmhclod`. The orchestrator's recommendation on TH-0843, which alex confirmed at 2026-11-07T05:38:10Z, said to use the folder that holds the `.credentials.json` the operator's login uses: ADR-0247 mounts it read-only at `/run/claude-login`, and only that file is read.

What the operator set was neither guess. S-0567's Notes, `### Live run`, record: "The container harness's login is `/home/alex/.claude-me`, mounted read-only. It is neither `~/.claude` nor `~/.claude-rmh`, the two folders the threads named." The planner's container run on S-0559 then worked with it (wip/archive/kanban/stories/S-0567-flai-serve-starts-the-planner-in-a-container-on-harness-container-and-it-plans.md).

This is one host setting, shared by every container run on the host. A later story that needs it should check `flai serve agent harness` rather than ask for a folder by name. If it reads `login: not set`, ask the operator for the folder that holds their login's credentials.

The operator, alex, answered TH-0836 again at 2026-11-07T06:17:59Z.

The answer:

```text
move S-0556
```

The operator, alex, answered TH-0836 again at 2026-11-07T06:25:56Z.

The answer:

```text
**Recommended, alex:** press Stop on S-0556's activity page now, as agent-S-0565 asks. Wait about 30 seconds, press Retry, and say here what the card showed after each.

- You moved S-0556 to ready yourself at 06:17:39Z, so I moved nothing. Its agent took it to in-progress at 06:18:19Z, in its container.
- The Stop and the Retry are yours alone: I have no tool that stops or restarts an agent (ADR-0058).
- The Stop asks to confirm first. Once stopped, S-0556 gets no new agent until you press Retry.

TH-0828, my reminder to keep S-0556 out of ready, has done its job, so I resolve it.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0058-the-agent-host-action-lets-the-dashboard-stop-a-story-s-agent-and-a-stop.md § Decision
```

The operator, alex, answered TH-0836 again at 2026-11-07T06:27:39Z.

The recommendation it overruled:

```text
**Recommended, alex:** reply with one line saying what S-0556's card showed after Stop and after Retry. Only you saw the card, so I cannot answer for you.

**What to compare against:** ADR-0058 §4 says that after a stop, the dashboard shows the story's agent as failed, "stopped by the operator", with Retry. After Retry it should show the agent as working again. If what you saw differs, say how. agent-S-0565 then records the difference as an issue.

**If you don't remember:** agent-S-0565 records the card's two states in S-0565's Notes as "not observed by the operator". Beside them it records the `stopped` and `started` events it already has, and what ADR-0058 §4 says the card should show. It does not write down a state nobody saw.

I have moved nothing. S-0556's agent runs in its new container. Its restart after an answered thread is the second half of criterion 2.

Source: design/adrs/0058-the-agent-host-action-lets-the-dashboard-stop-a-story-s-agent-and-a-stop.md § Decision
```

The operator's answer:

```text
It showed as stopped by the operator. I did not see a retry button on the story's agent page, but when I right-clicked the card, it let me retry.
```

The operator, alex, answered TH-0836 again at 2026-11-07T06:54:03Z.

The answer:

```text
the agent asked the question only once it had completed it's last task (somewhat of a race condition) it's agent is continuing with close out steps and will move the story to review. we can mark our criteria here as completed
```

Sources: TH-0843, TH-0836, S-0567.
