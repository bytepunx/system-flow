---
id: T-1182
type: task
nature: improvement
title: Strategic agent panels, the story agent, plan actions, host agent notices, and settings show their times in local time
status: done
parent: S-0329
owner: alex
created: 2026-10-07T19:49:33Z
updated: 2026-10-07T20:37:50Z
transitions:
  - to: ready
    at: 2026-10-07T20:26:26Z
    by: agent-S-0329
  - to: in-progress
    at: 2026-10-07T20:26:26Z
    by: agent-S-0329
  - to: done
    at: 2026-10-07T20:37:50Z
    by: agent-S-0329
stream: S-0329
tags: [dashboard]
touches: [flaiover/src/lib/components/StrategicAgentPanel.svelte, flaiover/src/lib/components/OrchestratorPanel.svelte, flaiover/src/lib/components/OrchestratorPanel.svelte.test.ts, flaiover/src/lib/components/PlannerPanel.svelte.test.ts, flaiover/src/lib/components/AnalyzerPanel.svelte.test.ts, flaiover/src/lib/components/PlanAction.svelte, flaiover/src/lib/components/PlanAction.svelte.test.ts, flaiover/src/lib/components/HostAgentNotice.svelte, flaiover/src/lib/components/HostAgentNotice.svelte.test.ts, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts]
after: [T-1180]
usage:
  source: log
  seconds: 684
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 57
      output: 20759
      cache_read: 3532759
      cache_write: 90835
      cost: 1.6787
---
# T-1182 Strategic agent panels, the story agent, plan actions, host agent notices, and settings show their times in local time

## Work

These components each define their own `at = (s) => s.replace('T', ' ').replace(/:\d\dZ$/, ' UTC')`, or print a time raw. Replace each with `localTime` from `flaiover/src/lib/localtime.ts`, and update the tests that assert `... UTC` to the zone the tests run in.

| File | What it shows today |
|------|---------------------|
| `flaiover/src/lib/components/StrategicAgentPanel.svelte` | exports `at`; the planner's and analyzer's runs and activity entries (`PlannerPanel`, `AnalyzerPanel` render through it) |
| `flaiover/src/lib/components/OrchestratorPanel.svelte` | imports `at` from the panel for its decisions and runs |
| `flaiover/src/lib/components/PlanAction.svelte` | its own `at`: `planner running since ... UTC` |
| `flaiover/src/lib/components/HostAgentNotice.svelte` | its own `at`: agent started and ended |
| `flaiover/src/lib/components/StoryAgent.svelte` | its own `at`: the story agent started and ended |
| `flaiover/src/lib/components/SettingsPanel.svelte` | the planning schedule's next run, `p.next`, raw |

Remove the local `at` helpers rather than keep a second formatter; `OrchestratorPanel.svelte` imports `localTime` instead of the panel's `at`. A schedule is a cron expression flai reads in UTC: its text, `in UTC`, stays, since it names how the expression is read, not a time shown; its next run is a time and shows local.

It waits for T-1180 because it calls its formatter. It shares no file with the other two view tasks, so the three run together.

## Done when

- None of the files above shows a time in UTC, and none defines its own time formatter.
- `PlannerPanel.svelte.test.ts`, `AnalyzerPanel.svelte.test.ts`, and the other tests above assert the local text in the tests' pinned zone, and `flai test` on the files above passes.

## Notes

- The planner, orchestrator, and analyzer pages under `flaiover/src/routes/workflow/` render these panels; their route tests may assert times too. Add them to this task's touches if they do.
