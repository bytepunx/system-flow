---
id: T-2187
type: task
nature: remediation
title: The user, operator, and design documents say the host configuration is the operator's and which flai serve sub-commands the guard refuses an agent
status: done
parent: S-0511
owner: alex
created: 2026-10-28T08:20:07Z
updated: 2026-10-30T00:24:31Z
transitions:
  - to: ready
    at: 2026-10-30T00:13:50Z
    by: agent-S-0511
  - to: in-progress
    at: 2026-10-30T00:13:51Z
    by: agent-S-0511
  - to: done
    at: 2026-10-30T00:24:31Z
    by: agent-S-0511
stream: S-0511
tags: [cli, guard, docs]
touches: [docs/users/flai.md, docs/operators/index.md, docs/operators/settings.md, design/system/flai-cli.md, design/system/agent-adapters.md]
after: [T-2184]
usage:
  source: log
  seconds: 640
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 51
      output: 14421
      cache_read: 2472419
      cache_write: 120313
      cost: 1.5543
---
# T-2187 The user, operator, and design documents say the host configuration is the operator's and which flai serve sub-commands the guard refuses an agent

## Work

Write the rule where the documents already say what the guard refuses, one or two sentences each, linking the ADR T-2184 recorded by relative path, and bump each `updated`. In `docs/users/flai.md`: under `## Sub-agents`, beside the sentence that `flai guard` refuses the shared paths and `flai forge push` to every sub-agent and every served session; under `## Configuration`, where it says each `flai serve` setting has its own command run in a shell on the host, that an agent `flai serve` starts is refused those commands; and under `### flai serve: flai on the host, for the dashboards`, where the host actions are turned on, that no agent turns one on or off. In `docs/operators/index.md`: under `## Reminding the operator of what awaits them`, that `notify set` and `clear` are refused to every agent; under `### Starting an agent when a story becomes ready`, where it says the agent's flai reads its own configuration, that the guard refuses the agent the commands that would write it; and under `### The forge host action: pull requests on the move to review (S-0393)`, beside the `flai forge push` sentence, that `flai serve forge` with a setting is refused too. In `docs/operators/settings.md`, in the `FLAI_CONFIG` row under `### Read by flai`, that the guard refuses a served agent the writes to the file. In `design/system/flai-cli.md`, in the `flai guard` row of `## Commands`, the rule in one clause with the ADR. In `design/system/agent-adapters.md` § The guard, add the rule to the list of what the hook refuses every served session. Name the reading forms that stay open in the user document only. Waits for T-2184 for the ADR's number; runs beside T-2185 and T-2186, with which it shares no path.

## Done when

Each named section carries the rule with a link to the ADR, `flai check --strict` and the markdown lint are clean, and nothing in the documents still says a served agent may run those commands.

## Notes

`docs/users/flaiover.md` and the dashboard are untouched: the `settings` host action and the dashboard's writes through `flai manifest set` are not this rule's.
