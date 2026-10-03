---
id: ADR-0079
title: "The planner, the orchestrator, and the analyzer each log their activities in one document per project under wip/agents, written by flai with its totals in front matter that flai stats reports"
status: accepted
date: 2026-10-03
supersedes: []
superseded_by: []
topics: [cli, orchestration, planning, analysis]
---

# ADR-0079 The planner, the orchestrator, and the analyzer each log their activities in one document per project under wip/agents, written by flai with its totals in front matter that flai stats reports

## Context

A story's agent keeps a narrative and its usage is written on its story (ADR-0051). The planner, the orchestrator, and the analyzer (ADR-0075, E-0016) work above a story and have neither: nothing records what they did, how long it took, or what it cost, and the Strategic Cost and Strategic Use charts the designer asked for (S-0216) need those numbers per day. The designer decided on 2026-10-02 (S-0206) that each gets a log document with totals in its front matter. The planner runs once per item and ends; the orchestrator runs for a long time and does many things in one run, so a run is not always one activity.

## Decision

Each strategic agent has one activity document per project: `wip/agents/planner.md`, `wip/agents/orchestrator.md`, and `wip/agents/analyzer.md`. flai writes them; no agent edits them by hand.

1. **Front matter.** `kind` (the file's name), `accrued_cost` in US dollars, `accrued_seconds` of wall-clock activity, `tasks_completed`, the number of activities logged, and `last_run`, when the newest activity ended. They are parsed strictly, like a work item's.
2. **Log.** Under `## Log`, one entry per activity, newest last: when it ended, a one-line summary the agent gives, the items it touched, its wall-clock seconds, and its cost, marked estimated when the cost was apportioned or priced rather than reported. Appending an entry adds it to the totals.
3. **Measured as a task is.** flai serve keeps a strategic agent's run logs as `<key>-<kind>-*.log` beside the story runs' logs. An activity's span runs from the later of the newest run's start and the last entry to its end, and its usage is the run's apportioned to that span as a task's is to its intervals (ADR-0051), so a session resumed in a later run, whose result reports cumulative totals, is charged only its share.
4. **Two ways an activity ends.** An agent whose run spans activities reports each with the MCP tool `activity_log`, giving its kind, summary, and items; flai measures and writes the entry. When a strategic run ends, flai serve logs the time since the last entry as one activity, with the run's last result as its summary, unless nothing was spent in it, so an agent that runs once and ends needs no call.
5. **Listed, checked, and reported.** `wip/agents/index.md` lists the documents that exist under a heading of their own. `flai check` validates their front matter and passes over them as narratives. `flai stats --json` reports each kind's totals and its log entries in the window as `strategic`, defined in `design/system/metrics.md`.

## Consequences

- The strategic agents' cost and time are recorded without anyone entering them, and S-0205 can chart them per day from the log entries.
- The documents live in the main checkout's `wip/`, as narratives do, and are committed with the work around them.
- An activity's cost is an apportionment of the run's and says so, as a task's does.
- Front matter is parsed strictly, so a flai older than this reports the three documents as orphan narratives until the host's flai is upgraded.
- `activity_log` names the kind rather than inferring it from the caller's name, since how strategic agents are named is S-0208's to settle.

## Alternatives considered

- A narrative per strategic agent, in the story narrative's shape: its `## Current state` and `## Next steps` describe a stream of work toward review, which these agents do not have, and its front matter holds no totals.
- Usage on the items a strategic agent touched: S-0225 to S-0227 attribute cost to items that way, but an activity can touch none, or many, and the totals per agent are what the strategic charts compare against delivery.
- One activity per run: the orchestrator's run spans many decisions, and one entry per run would hide when each was made and what it cost.
- A log file under `.flai-cache`: the totals are metrics inputs and belong in the repository, where `flai stats` and the dashboard read everything else.
