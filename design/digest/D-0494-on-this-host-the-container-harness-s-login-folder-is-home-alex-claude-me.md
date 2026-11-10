---
id: D-0494
title: "On this host the container harness's login folder is /home/alex/.claude-me, neither ~/.claude nor ~/.claude-rmh, which the threads asking for it named"
kind: answer
sources: [TH-0843, S-0567, TH-0836, S-0583]
topics: [harness, docker, planning]
paths: [design/system/strategic-agents.md, docs/operators/runbooks/agent-container.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: [TH-1000]
cited: [S-0583]
---

# D-0494 On this host the container harness's login folder is /home/alex/.claude-me, neither ~/.claude nor ~/.claude-rmh, which the threads asking for it named

Two threads asked the operator to set the container harness's login with `flai serve agent harness container --login <folder>` and guessed different folders: TH-0836 (S-0565) named `~/.claude`, and TH-0843 (S-0567) named `~/.claude-rmh`, reasoning that the host's claude-code harness runs `rmhclod`. The orchestrator's recommendation on TH-0843, which alex confirmed at 2026-11-07T05:38:10Z, said to use the folder that holds the `.credentials.json` the operator's login uses: ADR-0247 mounts it read-only at `/run/claude-login`, and only that file is read.

What the operator set was neither guess. S-0567's Notes, `### Live run`, record: "The container harness's login is `/home/alex/.claude-me`, mounted read-only. It is neither `~/.claude` nor `~/.claude-rmh`, the two folders the threads named." The planner's container run on S-0559 then worked with it (wip/archive/kanban/stories/S-0567-flai-serve-starts-the-planner-in-a-container-on-harness-container-and-it-plans.md).

This is one host setting, shared by every container run on the host. A later story that needs it should check `flai serve agent harness` rather than ask for a folder by name. If it reads `login: not set`, ask the operator for the folder that holds their login's credentials.

S-0583 relied on it without asking: its narrative (wip/archive/agents/S-0583.md, Decisions) records that its probe runs in `flai-agent-claude-code:local` "use the login folder `/home/alex/.claude-me`, the operator's container login setting (D-0494), mounted read-only at `/run/claude-login`".

TH-0836's later answers (06:17Z to 06:54Z) are about S-0565's stop-and-retry proof on S-0556, not the login: S-0565's Notes record them, and I-0278 records that the story's agent page offered no Retry after a stop while the card's right-click menu did.

Sources: TH-0843, TH-0836, S-0567, S-0583.
