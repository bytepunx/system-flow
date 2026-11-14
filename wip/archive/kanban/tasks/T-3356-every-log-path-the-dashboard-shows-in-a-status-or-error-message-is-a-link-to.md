---
id: T-3356
type: task
nature: feature
title: Every log path the dashboard shows in a status or error message is a link to the log page
status: done
parent: S-0785
owner: alex
created: 2026-11-14T10:29:09Z
updated: 2026-11-14T11:16:35Z
transitions:
  - to: ready
    at: 2026-11-14T10:56:02Z
    by: agent-S-0785
  - to: in-progress
    at: 2026-11-14T10:56:02Z
    by: agent-S-0785
  - to: done
    at: 2026-11-14T11:16:35Z
    by: agent-S-0785
stream: S-0785
tags: []
touches: [flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/routes/ops/klusters/+page.svelte, flaiover/src/lib/components/DeploymentRow.svelte, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/PlannerPanel.svelte, flaiover/src/lib/components/PlanAction.svelte, flaiover/src/lib/components/AnalyzerPanel.svelte, flaiover/src/lib/components/RolePage.svelte, flaiover/src/routes/board/+page.svelte, flaiover/src/lib/components/AgentStream.svelte, docs/users/flaiover.md, flaiover/src/lib/activity.test.ts, flaiover/src/lib/activity.ts, flaiover/src/lib/components/ActivityView.svelte, flaiover/src/lib/components/ActivityView.svelte.test.ts, flaiover/src/lib/components/AgentStream.svelte.test.ts, flaiover/src/lib/components/AnalyzerPanel.svelte.test.ts, flaiover/src/lib/components/DeploymentRow.svelte.test.ts, flaiover/src/lib/components/DismissibleNotice.svelte, flaiover/src/lib/components/DismissibleNotice.svelte.test.ts, flaiover/src/lib/components/HostAgentNotice.svelte, flaiover/src/lib/components/ImageJob.svelte, flaiover/src/lib/components/ImageJob.svelte.test.ts, flaiover/src/lib/components/Inbox.svelte.test.ts, flaiover/src/lib/components/InboxView.svelte, flaiover/src/lib/components/OrchestratorPanel.svelte, flaiover/src/lib/components/PlanAction.svelte.test.ts, flaiover/src/lib/components/PlannerPanel.svelte.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/RolePage.svelte.test.ts, flaiover/src/lib/components/StoryAgent.svelte.test.ts, flaiover/src/lib/components/StrategicAgentPanel.svelte, flaiover/src/routes/activity/+page.svelte, flaiover/src/routes/activity/activity.svelte.test.ts, flaiover/src/routes/board/cardmenu.svelte.test.ts, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts, flaiover/src/routes/ops/images/+page.svelte, flaiover/src/routes/ops/klusters/klusters.svelte.test.ts]
after: [T-3355]
usage:
  source: log
  seconds: 1233
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 288
      output: 97439
      cache_read: 19088614
      cache_write: 470072
      cost: 8.5478
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 20791
          cache_write: 579
          cost: 0.0063
---
# T-3356 Every log path the dashboard shows in a status or error message is a link to the log page

## Work

Replace every place the dashboard shows a log path as plain text with `LogLink`, and render the free text from flai that can name a log (statuses, whys, errors, an agent run's or a job's detail) with `LogText`, so that no status or error message the dashboard shows names a log only as a string.

- Known sites: `routes/ops/deployments/+page.svelte` (step and job logs), `routes/ops/klusters/+page.svelte` ("Its output"), `lib/components/DeploymentRow.svelte` (the deploy log tooltip and `checks_log`), `lib/components/StoryAgent.svelte`, `PlannerPanel.svelte`, `PlanAction.svelte`, `AnalyzerPanel.svelte`, `RolePage.svelte`, and `routes/board/+page.svelte` ("its output is in ... on the host").
- Find the rest: the places the dashboard renders a refusal or error message from flai (notices, toasts, form errors), an agent's why on the Activity and inbox views, and an image job's or kluster job's detail; render each with `LogText`.
- Change the wording "its output is in X on the host" to name the log as a link the operator can open.
- Keep each component's examples under `design/components/` in step with its props, and update `docs/` where the dashboard's guide says to open a log on the host.
- Waits for T-3355, whose `LogLink` and `LogText` it uses.

## Done when

- A search of `flaiover/src` finds no log path rendered as text outside `LogLink`, and each changed component's tests pass with `flai test`.

## Notes
