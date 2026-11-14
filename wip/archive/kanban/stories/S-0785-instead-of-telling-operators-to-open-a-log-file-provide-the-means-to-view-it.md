---
id: S-0785
type: story
nature: feature
title: Instead of telling operators to open a log file, provide the means to view it
status: done
owner: alex
created: 2026-11-14T10:23:44Z
updated: 2026-11-14T17:31:39Z
transitions:
  - to: ready
    at: 2026-11-14T10:23:45Z
    by: alex
  - to: in-progress
    at: 2026-11-14T10:24:02Z
    by: agent-S-0785
  - to: review
    at: 2026-11-14T11:36:07Z
    by: agent-S-0785
  - to: done
    at: 2026-11-14T17:31:39Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/internal/hostapi/logs.go, flai/internal/hostapi/logs_test.go, flai/internal/hostapi/writes_methods.go, flai/internal/hostapi/guardreads_test.go, design/adrs/README.md, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/log/+server.ts, flaiover/src/routes/api/log/log.test.ts, flaiover/src/lib/server/logs.ts, flaiover/src/lib/logs.ts, flaiover/src/lib/logs.test.ts, flaiover/src/lib/components/LogViewer.svelte, flaiover/src/lib/components/LogViewer.svelte.test.ts, flaiover/src/lib/components/LogLink.svelte, flaiover/src/lib/components/LogText.svelte, flaiover/src/lib/components/LogText.svelte.test.ts, flaiover/src/routes/ops/log/+page.svelte, flaiover/src/lib/sitemenu.ts, design/system/flaiover-dashboard.md, design/components/flaiover/src/lib/components/LogViewer.examples.json, design/components/flaiover/src/lib/components/LogLink.examples.json, design/components/flaiover/src/lib/components/LogText.examples.json, flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/routes/ops/klusters/+page.svelte, flaiover/src/lib/components/DeploymentRow.svelte, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/PlannerPanel.svelte, flaiover/src/lib/components/PlanAction.svelte, flaiover/src/lib/components/AnalyzerPanel.svelte, flaiover/src/lib/components/RolePage.svelte, flaiover/src/routes/board/+page.svelte, flaiover/src/lib/components/AgentStream.svelte, docs/users/flaiover.md, flaiover/src/lib/markdown.ts, flaiover/src/lib/markdown.test.ts, flaiover/src/lib/filelinks.test.ts, design/adrs/0300-the-dashboard-reads-a-log-flai-wrote-through-log-read-confined-to-the-project-s.md, flai/internal/hostapi/images.go, flai/internal/hostapi/writes_actions.go, flai/internal/serve/serve.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4337
  priced_by: harness
  turns:
    - day: 2026-11-14
      ceremony: 2
      work: 51
  models:
    - model: claude-haiku-5-5
      input: 30
      output: 18987
      cache_read: 907784
      cache_write: 116450
      cost: 0.0893
    - model: claude-opus-5-5
      input: 646
      output: 218789
      cache_read: 42861524
      cache_write: 1055499
      cost: 19.1932
  strategic:
    - kind: planner
      seconds: 466
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 84
          output: 42376
          cache_read: 4104887
          cache_write: 159758
          cost: 0.303
        - model: claude-opus-5-5
          input: 60
          output: 18537
          cache_read: 2779686
          cache_write: 123396
          cost: 1.9141
    - kind: orchestrator
      seconds: 91
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 34
          output: 219
          cache_read: 1229459
          cache_write: 204220
          cost: 0.4245
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: alex
    at: 2026-11-14T10:23:44Z
  value: 12.5
  by: planner-S-0785
  at: 2026-11-14T10:31:19Z
forecast:
  duration: 1h23m
  delivery: 2026-11-14T11:58:00Z
  basis: flai forecast's median of 110 s per unit over 131 large feature stories on claude-opus-5-5, times size 45, in line with S-0757's 65 minutes for a like host read, route, and Ops page.
  by: planner-S-0785
  at: 2026-11-14T10:31:19Z
---
# S-0785 Instead of telling operators to open a log file, provide the means to view it

## Goal

Right now errors in several areas (agents status, deploy status, issues with review, etc.) will instruct the operator to view the contents of a log file. The issue is that this requires the operator to have shell access to the host machine, which is not always the case.

These log paths should work as hyperlinks that can be vied within flaiover instead.

## Acceptance criteria

- [x] No status/error messages should present a path to a log file as a plain string
- [x] There should be a log viewer component that flai displays log contents in
- [x] The path to the log file in error and status messages is displayed as a hyperlink that takes the operator to page with the log viewer set to the log file's contents

## Tasks

- T-3354 The host API method log.read answers a part of a log flai wrote, confined to the project's and flai serve's log folders
- T-3355 The dashboard's log page and its LogViewer, LogLink, and LogText components show a log a message names
- T-3356 Every log path the dashboard shows in a status or error message is a link to the log page
- T-3357 A log path in a thread entry or a message that flai wrote renders as a link to the log page

## Notes

### Planning

Touches, by where each came from:

- Declared: `flaiover/src` and `flai/cmd`, the operator's. Both stay. `flaiover/src` is a folder touch the story keeps because the work may reach more rendering sites below it than its tasks name today: T-3356 is told to find the rest of the refusals, notices, and whys that can name a log. `flai/cmd` is a folder touch kept because it is declared. No task names a file there yet: the Go message builders that print a log path (`flai/cmd/kluster.go`, `flai/cmd/image.go`, `flai/cmd/serve.go`) are terminal output, which the plan leaves as it is.
- Design and layout, from the code search: T-3354's files under `flai/internal/hostapi/`, the registration in `writes_methods.go`, the read listed in `guardreads_test.go`, and `flaiover/src/lib/server/agent.ts`'s `REQUIRED_METHODS`, which D-0612 says every new host API method the dashboard calls needs. `design/adrs/README.md` is there for the ADR T-3354 writes; the ADR's own file is named only once `flai adr new` makes it, and `task_done` widens to it.
- Layout: T-3355's route, server helper, `$lib/logs.ts`, components, page, `sitemenu.ts`, and the examples under `design/components/flaiover/src/lib/components/`. T-3356's rendering sites and the examples of the components among them that have examples. T-3357's `markdown.ts` and `filelinks.ts`, through which thread entries and messages render.
- Design: `design/system/flaiover-dashboard.md`, whose Views and API tables `design-routes.test.ts` reads, and `docs/users/flaiover.md`, the user guide for the dashboard's pages.
- Co-change: `flai touches suggest` ranks `docs/users/flai-reference.md`, `docs/users/flai.md`, and `design/system/flai-cli.md` first, from the commit history of the two declared folders. None is added: the plan adds no command or flag.

Forecast: with the two declared touches alone, `flai forecast` gave 10m for size 5, too low for work in Go and SvelteKit. With the 42 touches above it gives 1h23m for size 45 (110 s a unit, the median of 131 large feature stories on claude-opus-5-5), delivered at 2026-11-14T11:58:00Z, counted from the story's start at 10:25. That figure stands: it fits S-0757, which took 65 minutes from in-progress to review in six tasks for a similar host API read, API route, component, and Ops page. This story has four tasks but more rendering sites and a confinement rule to test.

Cost of delay: `flai cod` gives 12.50 USD a week, from the operator's input of 5m time lost per cycle at the project's 150 USD an hour, one cycle a week. It stands as worked out.

No `after`: the story reads no other open story's output.
