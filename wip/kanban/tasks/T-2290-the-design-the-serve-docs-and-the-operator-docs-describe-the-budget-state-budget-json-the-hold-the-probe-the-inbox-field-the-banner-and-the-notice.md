---
id: T-2290
type: task
nature: improvement
title: The design, the serve docs, and the operator docs describe the budget state, budget.json, the hold, the probe, the inbox field, the banner, and the notice
status: backlog
parent: S-0534
owner: alex
created: 2026-10-28T21:27:35Z
updated: 2026-10-28T21:27:35Z
transitions: []
stream: S-0534
tags: [docs, design, budget]
touches: [design/system/flai-cli.md, design/system/project-manifest.md, design/system/flaiover-dashboard.md, design/system/strategic-agents.md, design/system/agent-narrative.md, docs/operators/index.md, docs/operators/runbooks/backup.md, docs/users/flai.md, docs/users/flaiover.md]
after: [T-2284, T-2285, T-2286, T-2288, T-2289]
---
# T-2290 The design, the serve docs, and the operator docs describe the budget state, budget.json, the hold, the probe, the inbox field, the banner, and the notice

## Work

Describe what the story built, as built, in the documents criterion 6 names and those that say a failed start is retried.

The design:

- `design/system/flai-cli.md` § Commands: the `flai serve` rows on the kept files, the journal, `orchestrateRetry`, and `auto_restarts`. Cover the budget state per provider in `budget.json`, what sets, holds, probes, clears, and extends it, its journal entries, and `flai serve status`'s line.
- `design/system/project-manifest.md`: that an agent's `provider`, or its harness when it names none, is the key the host's budget state holds starts by.
- `design/system/flaiover-dashboard.md` § Views or § Runtime: the budget banner, where its data comes from, and when it shows.
- `design/system/strategic-agents.md`: each role's run outcome. A run that ended on the limit is not retried; its trigger, queue, or schedule waits for the state to clear.
- `design/system/agent-narrative.md` § What an agent is told through MCP: the `budget` field and the instruction to hold.

The docs:

- `docs/operators/index.md`: § Where things are gets a row for `budget.json`. The host action sections for agent, plan, orchestrate, analyze, and curate say a start is held while the provider's budget is exhausted. § Reminding the operator gets the budget notice and its JSON.
- `docs/operators/runbooks/backup.md`: `budget.json` among the serve files. Restoring it is optional, since the next refusal sets it again.
- `docs/users/flai.md`: the `inbox` field, the refusal a hand start gives, and `flai serve status`'s line.
- `docs/users/flaiover.md`: the banner.

Cite ADR-0209 §5. No ADR is added, since §5 decides this. No new command or flag is added, so `docs/users/flai-reference.md` and the settings index stay as they are. If an earlier task did add one, run `make flai-reference` and add both files here.

This task waits for every other task, so it describes what was built.

## Done when

- [ ] Each document listed says what the story built, linking ADR-0209 where it decides it.
- [ ] `flai check --strict` and the markdown lint are clean on the changed files.

## Notes

Planned by planner-S-0534.
