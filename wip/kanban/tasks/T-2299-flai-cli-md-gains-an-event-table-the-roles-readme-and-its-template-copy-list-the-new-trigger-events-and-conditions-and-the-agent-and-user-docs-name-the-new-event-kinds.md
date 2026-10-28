---
id: T-2299
type: task
nature: feature
title: flai-cli.md gains an event table, the roles README and its template copy list the new trigger events and conditions, and the agent and user docs name the new event kinds
status: backlog
parent: S-0535
owner: alex
created: 2026-10-28T21:39:08Z
updated: 2026-10-28T21:39:08Z
transitions: []
stream: S-0535
tags: [flai, events, docs, template]
touches: [design/system/flai-cli.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, template/CHANGELOG.md, design/system/agent-narrative.md, docs/users/flai.md]
after: [T-2293, T-2298]
---
# T-2299 flai-cli.md gains an event table, the roles README and its template copy list the new trigger events and conditions, and the agent and user docs name the new event kinds

## Work

S-0535's criterion 6 asks for `flai-cli.md`'s event table and `design/workflows/roles/README.md`'s trigger table. `flai-cli.md` has no event table yet. The MCP events are described in `agent-narrative.md` § What an agent is told through MCP. The roles README's table lists each role's trigger, with no list of events or conditions.

- `design/system/flai-cli.md`: add a table of the events `wait_for_events` and `inbox` report. Give each kind, who emits it, where it is kept (`workitem.Changes`, `edits.jsonl`, `overlaps.jsonl`, the new `events.jsonl`, the thread files, messages), and its fields. Link T-2291's ADR.
- `design/workflows/roles/README.md` and `template/root/design/workflows/roles/README.md`: list the trigger events, the new ones among them, and the trigger conditions with the event each qualifies. Leave the orchestrator's row as `continuous`: S-0536 changes it.
- `template/CHANGELOG.md`: an entry for the template's README change.
- `design/system/agent-narrative.md` § What an agent is told through MCP: a bullet for each new kind, beside `edited`, `overlapped`, and `message`.
- `docs/users/flai.md` § Serving agents over MCP: name the new kinds.

It waits for T-2293, whose conditions it documents, and for T-2298, whose descriptions it must agree with.

## Done when

- Every event and condition S-0535 adds is named, with its fields, in each document above.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes

The template copy of the roles README is a baseline file, so it changes with the repository's copy in the same task.
