---
id: T-2577
type: task
nature: remediation
title: An ADR decides the built-in default agent, the template's default, per-project host agent settings, and the refusal's links to the settings pages
status: done
parent: S-0587
owner: alex
created: 2026-11-07T10:11:44Z
updated: 2026-11-07T10:31:11Z
transitions:
  - to: ready
    at: 2026-11-07T10:28:06Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:28:11Z
    by: system-flow
  - to: done
    at: 2026-11-07T10:31:11Z
    by: agent-S-0587
stream: S-0587
tags: [cli, adr]
touches: [design/adrs/README.md]
usage:
  source: log
  seconds: 180
  estimated: true
  priced_by: estimate
  models:
    - model: claude-haiku-5-5
      input: 40
      output: 184
      cache_read: 1244236
      cache_write: 99357
      cost: 0.0966
    - model: claude-opus-5-5
      input: 16
      output: 118
      cache_read: 2387820
      cache_write: 19317
      cost: 0.8634
---
# T-2577 An ADR decides the built-in default agent, the template's default, per-project host agent settings, and the refusal's links to the settings pages

## Work

Write one ADR with `flai adr new` that refines ADR-0037 and ADR-0251. It is the first layer because every other task builds on what it decides. It decides four things:

1. **The built-in default agent.** When neither the story nor `system-flow.yaml` names a harness, the story resolves to a default built into flai, recommended `harness: claude-code` with the model the template ships. Say which model, and that a project's `agent:` and a story's own agent still win, field by field, as `Agent.With` does. ADR-0251 kept the refusal "with neither"; this ADR replaces that case.
2. **The template's default.** `template/root/system-flow.yaml.tmpl` ships the `agent:` block uncommented, so a new or imported project states its default in its manifest.
3. **Per-project host agent settings.** The host command, the agent's name, and each harness's program and arguments are kept per project root in the host configuration, with a host-wide entry (`*`) as the fallback, as `host_actions` are. Setting them from a project's folder or a project's settings page changes only that project. Changing every project takes an explicit `--all-projects`, or the page's "everywhere" choice. Providers stay host-wide, since their keys are the host's, and the ADR says so against criterion 4.
4. **The refusal's links.** Each "names no harness" refusal, and the waiting reason the launcher records, names the missing setting and the dashboard page where it is entered: `/settings/agents` for the default agent and `/settings/harnesses` for a harness's program or the host command. It also gives the CLI command. The dashboard renders those paths as links.

Run `flai adr accept` on it once it is written. `design/adrs/README.md` is generated from the ADR files.

## Done when

- The ADR is accepted, refines ADR-0037 and ADR-0251, and states the four decisions with the alternatives it rejected.
- `flai check --strict` passes, and `design/adrs/README.md` lists the ADR.

## Notes
