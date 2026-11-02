---
id: ADR-0249
title: "Strategic usage is per role: a role's activity is charged by its role file's charge rule under its name in usage.strategic, flai stats reports every defined role, and the dashboard draws a series per role"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0083, ADR-0095, ADR-0100]
topics: [cli, dashboard, orchestration, planning, analysis]
---

# ADR-0249 Strategic usage is per role: a role's activity is charged by its role file's charge rule under its name in usage.strategic, flai stats reports every defined role, and the dashboard draws a series per role

## Context

Strategic usage is charged and reported per kind, and the kinds are four built-in agents: the planner's activity is charged to the item it planned ([ADR-0083](0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md)), the orchestrator's evenly to the work items it names ([ADR-0095](0095-an-orchestrator-activity-s-usage-is-charged-evenly-to-the-work-items-it-named.md)), the analyzer's evenly to the issues it names ([ADR-0100](0100-an-analyzer-activity-s-usage-is-charged-evenly-to-the-issues-it-names-under.md)), and the curator's to nothing ([ADR-0196](0196-the-curator-runs-behind-the-curate-host-action-and-takes-in-what-the-archive.md)).

Since [ADR-0207](0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) a project defines a strategic role in a file under `design/workflows/roles`, and the file names how its usage is charged in `charge`: `item`, `named`, `issues`, or `none`. A defined role keeps an activity document and logs its activities as the built-in agents do, but nothing reads its `charge`: `flai serve` charges its activity to nothing, `workitem` and `issues` refuse its name as a kind of `usage.strategic`, and the dashboard's strategic charts draw only the planner, the orchestrator, and the analyzer. A formula's role step ([ADR-0239](0239-a-formula-s-instance-is-worked-by-its-story-s-agent-step-by-step-a-role-step.md)) runs a defined role on the formula's instance, and its usage reaches the instance only through this charge. `design/system/metrics.md` is the contract between `flai stats` and the dashboard, and changes only with an ADR.

## Decision

**Strategic usage is per role. A kind of strategic usage is a role the project defines, or one of the four built-in agents, and a role's activity is charged by the rule its role file names in `charge`, under the role's name in `usage.strategic`. `flai stats` reports every role the project defines, and the dashboard draws one series per role.**

1. **The rules.** `charge` takes one of four rules, each charging the activity's share of its run's usage, as the built-in agents' charges do:
   - `item`: to the item the run was started for, as `flai serve` recorded the run in `serve/agents.json`, and to every item above it. A run with no item recorded, such as one a person ran by hand, charges nothing.
   - `named`: evenly to the epics, stories, and tasks the activity's entry names, each share to its item and every item above it.
   - `issues`: evenly to the issues the entry names, under each issue's own `usage.strategic`, carried to the story made from an issue as ADR-0100 says.
   - `none`: to nothing.
2. **The built-in agents.** The baseline role files carry the rules the four have today: `item` on the planner, `named` on the orchestrator, `issues` on the analyzer, and `none` on the curator. A built-in agent with no role file keeps that rule by its name. ADR-0083, ADR-0095, ADR-0100, and ADR-0196 stand as the rule of each; what they say of their own kind is now what their rule says of any role.
3. **A formula's role step.** A role run started for a formula's instance has the instance as its item, so a role whose rule is `item` is charged to the instance and the items above it. The analyzer started for a role step ([ADR-0245](0245-a-formula-s-role-step-that-names-the-built-in-analyzer-starts-it-through-the.md)) keeps its own rule, `issues`.
4. **What stays in the project.** What no item or issue carries stays in the role's project strategic total: activities under `none`, those that named or reached nothing, and any charge that failed.
5. **The kind.** An entry of `usage.strategic` on an item or an issue takes the role's name as its `kind`. Entries are ordered planner, orchestrator, analyzer, curator, then every other kind by name, wherever they are listed: on an item, in `usage.strategic.kinds`, in `strategic`, and in `strategic_days[].agents`. A charge is refused for a name that is neither a built-in agent nor a role the project defines. `flai check` accepts on an item or an issue any kind of a role's name form, defined now or not, since a charge outlives the role file it was made under.
6. **`flai stats`.** `strategic` holds one entry for every role the project defines and every activity document that exists: a role with no document yet is listed with zeros, an empty `log`, and `last_run` `""`. Each entry splits its totals into `items`, `issues`, and `project` as for the built-in agents. `usage.strategic.kinds` and `strategic_days[].agents` carry each role by name. The existing fields keep their shape.
7. **The dashboard.** The Strategic Cost and Strategic Use charts draw one series per kind `strategic_days[].agents` carries, the built-in agents first, each other role in the next palette slot. A role's page draws its own cost and use per day from `flai stats`.
8. **`flai check`.** It reports a role file whose `charge` is none of the four rules, and one whose `charge` is `item` while its `item` is `none`, since such a run is started for no item.

## Consequences

- A project's own role is charged, reported, and charted as the built-in agents are, with no flai release for a new role.
- The planner, the orchestrator, the analyzer, and the curator charge as before; a project can change one's rule in its role file.
- A role whose run named no item and whose rule is `item` charges nothing, as a planner run started by hand does today.
- A role removed from `design/workflows/roles` leaves its charges on the items and issues that carry them; `flai stats` keeps reporting it while its activity document exists.
- `metrics.md` § Strategic usage, § Strategic agents, § Strategic use per day, and § Charts say the same. `usage.strategic` carries no new field, so no `flai.minimum` is raised.

## Alternatives considered

- **Charging by the role's `item` alone, with no rule in the file.** A role that names the items it acts on, as the orchestrator does, could not be charged to them, and the analyzer's charge to issues could not be expressed.
- **A rule in the manifest rather than the role file.** The role file is where a role is defined (ADR-0207); a second place to keep in step adds a check and nothing else.
- **Keeping the four kinds and charging a defined role's usage under a single `role` kind.** Two defined roles would be one series, and their spend could not be told apart on an item.
- **Refusing on an item a kind no role file now defines.** Removing a role file would turn every item it was charged on into a `flai check` finding, archived ones included.
