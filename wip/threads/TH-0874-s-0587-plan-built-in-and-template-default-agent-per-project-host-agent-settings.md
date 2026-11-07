---
id: TH-0874
title: "S-0587 plan: built-in and template default agent, per-project host agent settings, linked refusals; eight tasks in five layers, 3h"
anchor:
  path: wip/kanban/stories/S-0587-default-harness-and-agent-should-work-out-of-the-box.md
  item: S-0587
status: resolved
participants: [planner-S-0587, orchestrator]
created: 2026-11-07T10:13:50Z
updated: 2026-11-07T10:15:22Z
---

# TH-0874 S-0587 plan: built-in and template default agent, per-project host agent settings, linked refusals; eight tasks in five layers, 3h

On wip/kanban/stories/S-0587-default-harness-and-agent-should-work-out-of-the-box.md.

## Entries

### 2026-11-07T10:13:50Z planner-S-0587
S-0587 is planned: eight tasks in five layers, 48 touches, forecast 3h. The cost of delay is still waiting on TH-0872.

## Cause

- A new project has no default agent. `template/root/system-flow.yaml.tmpl` ships `agent:` commented out, and neither `flai new` nor `flai import` writes one.
- flai has no built-in fallback, so `Repo.AgentOf` (ADR-0251) has nothing to fall back to.
- The `claude-code` the dashboard shows is a form placeholder, not a stored default.
- The host command and harness programs are host-wide in `~/.flai/config.json`, so setting them in one project changes every project.

## Tasks and layers

| Layer | Task | What |
|-------|------|------|
| 1 | T-2577 | ADR: built-in default, template default, per-project host agent settings, refusal links |
| 2 | T-2578 | `AgentOf` falls back to flai's built-in default agent |
| 2 | T-2579 | The template ships `agent:`; `flai new` and `flai import` tests |
| 2 | T-2580 | Host command, name, and harness programs kept per project root, `*` fallback, `--all-projects` |
| 3 | T-2581 | Harnesses settings page saves per project, "every project" on request (after T-2580) |
| 3 | T-2582 | Each refusal and waiting reason names `/settings/agents` or `/settings/harnesses` and the command (after T-2578, T-2580) |
| 4 | T-2583 | The dashboard renders those paths as links (after T-2582) |
| 5 | T-2584 | Design and docs (after all code tasks) |

## Assumptions to confirm or change

1. **The built-in default is `claude-code`, with the model the template ships.** It is the only harness that runs without an image or a host command. A project's `agent:` and a story's own agent still win. This replaces ADR-0251's "refused with neither" case.
2. **Criterion 4 covers the host's command, agent name, and harness programs, not providers.** These become per project root with a host-wide fallback, as `host_actions` are. Providers stay host-wide because their keys are the host's. The alternative is to keep them host-wide and only require `--all-projects` to change them. That is cheaper, but a change still reaches every project, so it misses criterion 4 as written.
3. **The refusal links to dashboard paths, and the dashboard turns them into links.** The CLI keeps the command alongside. With a built-in default, the refusal mostly fires when the harness program is missing on the host, so T-2582 adds that case.

The story keeps the folder touch `design/adrs`, since the ADR's number is allocated only when it is written. It has no `after`.

Reply here to change an assumption. Otherwise the story's agent works the plan as written.

### 2026-11-07T10:15:03Z orchestrator
Approved as planned. The story's agent works the eight tasks in five layers as written.

**Why the plan fits S-0587:**

- **Assumption 1, a built-in `claude-code` default.** Criteria 1 and 2 say a story is always given a default agent and harness. That rules out ADR-0251's "with neither the start is refused" case. T-2577 records the change as a new ADR that refines ADR-0251. A story's own agent and the project's `agent:` still win, so what ADR-0251 decided for those cases stands.
- **Assumption 2, per-project command, agent name, and harness programs, with providers left host-wide.** This meets criterion 4 as written. The cheaper alternative would still let one project's change reach every project. Providers carry the host's keys, so they are credentials, not a project's harness or agent setting.
- **Assumption 3, dashboard paths in the refusal, the command beside them.** This meets criterion 3. T-2582 also covers the case still left once a default exists: the harness program is missing on the host.
- **Touches and figures.** Each task's touches sit in its own layer without overlap. Keeping the folder touch `design/adrs` is right, because `flai adr new` allocates the number. The 3h forecast has its stated basis: flai's 1h31m, raised for three test tiers and a host-config migration.

**Still open:** the cost of delay input on TH-0872, which waits on your answer to my recommendation there. It does not hold this plan. The policy is `throughput`, so the value does not move S-0587.

Source: design/adrs/0251-a-story-whose-agent-names-no-harness-is-started-with-the-project-s-default.md § Decision

Source: design/adrs/0251-a-story-whose-agent-names-no-harness-is-started-with-the-project-s-default.md § Decision

### 2026-11-07T10:15:22Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the eight tasks, touches, and 3h forecast fit S-0587's four criteria; the cost of delay input stays on TH-0872.
