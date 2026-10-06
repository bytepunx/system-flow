---
id: ADR-0101
title: "With the settings host action on, the dashboard edits the strategic agents' settings in the manifest through flai manifest set"
status: proposed
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0039]
topics: [cli, dashboard, orchestration]
---

# ADR-0101 With the settings host action on, the dashboard edits the strategic agents' settings in the manifest through flai manifest set

## Context

The strategic agents' settings live in `system-flow.yaml`: what the orchestrator may do (`orchestration.permissions`), how it orders the ready column and when a release is due (`orchestration.policy`, `orchestration.release`), and the planner's and the analyzer's agents and schedules (`planning`, `analysis`). Until S-0229 they were set by hand. [ADR-0039](0039-a-settings-host-action-turned-on-only-in-a-shell-lets-the-dashboard-change-the.md) lists what the `settings` host action lets a project's dashboard change: its other host actions, its default agent, and its MCP token. Every dashboard write runs a flai command hostapi builds from values it checks ([ADR-0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md)). Letting the dashboard edit the strategic settings widens what a dashboard token can do: with `settings` on, it can let the orchestrator accept stories and publish releases.

## Decision

**With the `settings` host action on for a project, its dashboard may also change the manifest's strategic settings, through `flai manifest set`, from the orchestrator's, planner's, and analyzer's pages.**

- `flai manifest set <key>=<value>... [--unset <key>]...` writes the keys of its catalog: `orchestration.permissions.*`, `orchestration.policy`, `orchestration.release.*`, `planning.agent`, `replan`, `schedule`, `hour_rate`, `cycle`, `default_duration`, and `analysis.agent` and `schedule`. It keeps every other key, the order, and the comments. It checks the manifest whole, as it would be after the change, with the validation `flai check` uses, and refuses with each field and its reason, writing nothing.
- `planning.currency` is not among them: changing it re-denominates every amount on the items, so it stays a hand edit.
- The host API's `settings.manifest` runs it with `--autocommit` and the dashboard's trailer, as `settings.default_agent` runs `flai agent set`. flai's refusal reaches the dashboard as Refused with each `{field, reason}`; with `settings` off it answers Disabled. `settings.get` gives each key's value, default, kind, values, meaning, and for a permission its risk, from the same catalog.
- A save does not restart the orchestrator. `flai guard` reads the permissions on each call, so a change holds from the orchestrator's next call, and its prompt tells it to read its permissions and policies again before each decision.

## Consequences

- With `settings` on for a project, the dashboard token can grant the orchestrator `accept_reviews` and `publish`: whoever holds it can have work merged and released without the operator's review. Each permission's risk is shown beside it on the page and in the operator guide, and the shell remains the only way to turn `settings` on or off.
- The panels and `flai check` read one catalog and one validation, so the dashboard cannot offer a value flai would refuse, and a value it refuses is said beside its field.
- A problem elsewhere in the manifest blocks every save until it is fixed by hand, because the manifest is checked whole.

## Alternatives considered

- **A host action of its own for the strategic settings.** Finer consent, but one more switch for the same person at the same shell; `settings` already means the operator trusts the dashboard with the project's configuration.
- **Writing the manifest from hostapi directly.** Faster, but it would break ADR-0029's rule that every write is a flai command, and the validation would live in two places.
- **Restarting the orchestrator on a save.** It would lose its context mid-decision; reading the permissions on each call already holds it to the change.
