---
id: S-0536
type: story
nature: feature
title: "The orchestrator runs for a reason and ends: its trigger names events, flai check rejects continuous, flai serve coalesces the reasons into one run naming them, wait_for_events ends the run when they are done, and ADR-0209 is accepted"
status: backlog
parent: E-0030
owner: alex
created: 2026-10-28T21:02:32Z
updated: 2026-10-28T21:02:32Z
transitions: []
tags: [flai, serve, orchestrator]
topics: [orchestration, cli, template, conventions]
touches: [flai/internal/serve, flai/internal/mcpserver, flai/internal/harness, flai/internal/manifest, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/system/project-manifest.md, design/adrs, docs]
after: [S-0532, S-0535]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0536 The orchestrator runs for a reason and ends: its trigger names events, flai check rejects continuous, flai serve coalesces the reasons into one run naming them, wait_for_events ends the run when they are done, and ADR-0209 is accepted

## Goal

ADR-0209 §1, §2, §4, and §7. The orchestrator is the one role with `trigger: [continuous]`: it primes, decides, holds `wait_for_events`, and decides again without ending, and every wake re-reads a transcript of about 450k tokens, 75 to 80 percent of the time to wait again. This story makes it run for a reason and end, like the planner and the curator. `orchestrator.md`, here and in the template, names the events of the previous story as its trigger; `flai check` rejects `continuous` on any role. flai serve keeps a reason queue per project, starts one run no sooner than `orchestration.coalesce` (a Go duration, `30s` by default) after the first queued reason with every reason seen meanwhile, queues reasons that arrive during a run for the next, and stops a run that holds past `orchestration.run_timeout` (`20m` by default). The start prompt lists each reason with its item or thread; the role's prompt step 15 becomes: act on each reason, log each action, then end. `wait_for_events` answers the orchestrator `end: true`, `why: reasons done` when nothing queued remains for it, as it answers a story's agent that asked (ADR-0154). The guard's hold, the permissions, the activity log and its charge, the orchestrator's own thread, `flai serve orchestrate stop|start`, and the `orchestrate` host action stay as ADR-0087 and ADR-0189 have them. The ADR, drafted as proposed on 2026-10-28, is accepted in this story, and the living design, the convention and its template copy, and the docs change with it.

## Acceptance criteria
- [ ] `orchestrator.md` and its template copy name the trigger events; `flai check` reports `continuous` on any role as a finding, with a test
- [ ] flai serve queues the reasons a project's events give, starts one run per batch after `orchestration.coalesce`, with the reasons in the prompt and the journal entry, and starts the next run for reasons that arrived meanwhile when the run ends; a test drives three events inside one window into one run
- [ ] `wait_for_events` answers the orchestrator `end: true` with `why: reasons done` when nothing queued remains, and a run past `orchestration.run_timeout` is stopped and logged
- [ ] `orchestration.coalesce`, `orchestration.run_timeout`, and `orchestration.schedule` are manifest keys `flai check` validates and `flai manifest set` sets, with the defaults the ADR gives
- [ ] The orchestrator's prompt, the strategic-agents convention and its template copy, `flai-cli.md`, `strategic-agents.md`, `project-manifest.md`, and the user docs say the orchestrator runs for a reason and ends
- [ ] ADR-0209's status is accepted, and the orchestrator in this repository runs under the new trigger for a day with its runs and reasons in the journal, recorded in the narrative with the context tokens per action from `flai stats`

## Tasks

## Notes

Evidence: ADR-0209's context section. The budget hold is the second story of E-0030; reasons queued while it is set wait. The reason-specific pack is the next story: until it lands, a run primes the general `orchestrate` pack as today, which the per-action figure will show.
