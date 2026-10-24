---
id: T-1823
type: task
nature: feature
title: "An ADR records the curator: its host action, how it starts and runs, its pack, what the guard lets it do, where it logs, and how its run is settled and charged"
status: done
parent: S-0431
owner: alex
created: 2026-10-16T06:05:20Z
updated: 2026-10-23T23:42:38Z
transitions:
  - to: ready
    at: 2026-10-23T23:41:25Z
    by: agent-S-0431
  - to: in-progress
    at: 2026-10-23T23:41:25Z
    by: agent-S-0431
  - to: done
    at: 2026-10-23T23:42:38Z
    by: agent-S-0431
stream: S-0431
tags: [cli]
touches: [design/adrs, design/adrs/README.md]
usage:
  source: log
  seconds: 73
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 12
      output: 111
      cache_read: 973671
      cache_write: 11815
      cost: 0.3351
---
# T-1823 An ADR records the curator: its host action, how it starts and runs, its pack, what the guard lets it do, where it logs, and how its run is settled and charged

## Work

Write the ADR with `flai adr new`, refining ADR-0099 and ADR-0075, on TH-0476's answer to question 3. It records, on the analyzer's pattern:

- The host action `curate`, off by default; `flai curate [--since]`, `curate.run`, and the MCP tool `curate`, which an agent `flai serve` started is refused; one run per project.
- Its agent, `digest.agent` merged over the project's `agent`; its environment (`FLAI_AGENT=curator`, `FLAI_ROLE=curate`, no item or story); its definition `.claude/agents/curator.md`.
- Its pack: role `curate`, the topic it primes with, and the conventions whose `roles` list it.
- What it reads (the harvest since its last activity entry) and writes (entries through `flai digest new` and `edit` only, `asked_again` among them).
- Its guard: the reads, the digest commands, `inbox`, `activity_log`, `thread_open`, `thread_reply`, `wait_for_events`; no file edits, no work item writes.
- Its run in `serve/agents.json` under `curator`, its log, its outcome (`worked` or `failed`), its activity entry in `wip/agents/curator.md` naming the entries it wrote, and its cost: the kind's project strategic total, charged to no item.

First layer: every code task builds what this decides.

## Done when

The ADR is accepted and listed in `design/adrs/README.md`, and `flai check --strict` reports nothing on it.

## Notes

The ADR's file name is not known until `flai adr new` makes it, so the task touches the shared `design/adrs` folder.
