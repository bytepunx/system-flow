---
id: TH-0646
title: Plan for S-0497
anchor:
  path: wip/kanban/stories/S-0497-custom-dockerized-agents.md
  item: S-0497
status: resolved
participants: [planner-S-0497, orchestrator]
created: 2026-10-25T21:28:30Z
updated: 2026-10-25T21:28:41Z
---

# TH-0646 Plan for S-0497

On wip/kanban/stories/S-0497-custom-dockerized-agents.md.

## Entries

### 2026-10-25T21:28:30Z planner-S-0497
The story's agent started S-0497 and wrote T-2118 and T-2119 before this plan; both match what the outcome needs, so I revisited them and created nothing.

Tasks, order, and layers:

- Layer 1: T-2118, the agent adapters design gains a section on a container harness of flai's own (a baseline image flai builds and starts) measured against flai's concerns and weighed against flai today and Docker Agent. Touches `design/system/agent-adapters.md`, beside S-0496's Docker Agent section, so that the ADR cites the comparison rather than carrying it.
- Layer 2: T-2119, after T-2118, the draft ADR, written with `flai adr new` and left `proposed`, with the trade-offs, limitations, security concerns, and the Kubernetes and signet directions the story's note names. Touches `design/adrs/README.md`; the ADR file takes its name from the title, so the story keeps the folder touch `design/adrs` for it, noted under `### Planning`.

Both tasks have empty `## Work` and `## Done when`. They are the story's agent's, so I did not fill them. Proposed, for the agent to take or change: T-2118 Work: a `### A flai baseline image` section after `### Docker Agent` in `agent-adapters.md`, a concerns table in the shape of the Docker Agent one (headless and resume, MCP, sub-agents and their models, tool-refusal hook, permission handler, usage stream, instruction files), and a three-column comparison, flai today (claude-code on the host, ADR-0129 to ADR-0132), Docker Agent (ADR-0202), and the baseline image, covering what flai controls (mounts, tools, flai MCP, flai guard, cost capture per harness, model, and agent), secrets and configuration (signet), and Kubernetes as a later runtime. Done when: the section is there, every claim about Docker Agent cites the S-0496 survey or ADR-0202, and `flai check --strict` and the markdown lint pass. T-2119 Work: `flai adr new "<title>" --refines 130 --refines 202 --body-stdin --commit`, status `proposed`, with Context, Decision, Consequences (trade-offs, limitations, security concerns), and Alternatives considered naming flai today, Docker Agent, and an agent loop of flai's own. Done when: the ADR is proposed, its index row is in `design/adrs/README.md`, and both criteria can be ticked. T-2119 carries nature `feature`; `research` fits a research story's ADR, as the agent may set.

Figures:

- Touches: `design/system/agent-adapters.md` and `design/adrs/README.md` (layout, from S-0496), `design/adrs` (folder kept for the ADR file). `flai touches suggest` refused for want of seeds.
- Forecast: flai gave 41m and delivery 2026-10-26T01:01Z; set to 25m, delivery 2026-10-25T22:10Z, because S-0496 (same shape, agent, and model) was in progress 10m15s and done 17m35s after it started.
- Cost of delay: no inputs and no epic; asked on TH-0645 (recommended `penalty_per_week: 50` as S-0496 had). The value is written when an input is set.

Assumptions: the analysis lives in `agent-adapters.md` as S-0496's did; no code, tech document, or operator documentation changes, since nothing is built; the ADR is a draft (`proposed`) that the operator accepts or not; signet is not described anywhere in the design, so the ADR names it as the secrets source the operator intends and settles nothing about it.

### 2026-10-25T21:28:37Z orchestrator
Approved. The plan fits S-0497.

- **The tasks fit.** The story's agent wrote T-2118 and then T-2119, in two layers, and you created none. T-2118's section in `agent-adapters.md` carries criterion 2's three-way analysis. T-2119's proposed ADR carries criterion 1. Your proposed Work and Done when are the agent's to take or change.
- **The touches fit.** There are three, and the one folder touch, `design/adrs`, is kept only for the ADR file and is named under `### Planning`. `agent-adapters.md` is also changed by S-0496, which waits in the merge queue. A story in review holds nothing, so S-0497 syncs onto it once it lands.
- **The figures fit.** The forecast is 25m, lowered from flai's 41m on S-0496's observed time, with the reason stated.
- **The value waits on TH-0645,** where I recommended your `penalty_per_week: 50` to you, since S-0497 is in progress.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-25T21:28:41Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the tasks, touches, and forecast fit S-0497. The cost of delay input stays on TH-0645.
