---
id: S-0584
type: story
nature: improvement
title: Analyze the project and list the issues the report names
status: review
owner: alex
created: 2026-11-07T08:03:23Z
updated: 2026-11-07T08:34:34Z
transitions:
  - to: ready
    at: 2026-11-07T08:03:40Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T08:04:17Z
    by: agent-S-0584
  - to: review
    at: 2026-11-07T08:34:34Z
    by: agent-S-0584
blocked:
  - from: 2026-11-07T08:22:38Z
    until: 2026-11-07T08:32:15Z
    reason: awaiting the operator on TH-0859
tags: []
touches: [wip/kanban/stories/S-0584-analyze-the-project-and-list-the-issues-the-report-names.md, design/issues/I-0299-a-role-run-the-stdio-mcp-server-starts-through-role-run-is-never-reaped-so-flai.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1238
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 3
      work: 124
  models:
    - model: claude-haiku-5-5
      input: 24
      output: 9046
      cache_read: 505237
      cache_write: 64563
      cost: 0.0176
    - model: claude-opus-5-5
      input: 268
      output: 23411
      cache_read: 21852420
      cache_write: 212818
      cost: 6.5423
  strategic:
    - kind: orchestrator
      seconds: 65
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 21
          output: 127
          cache_read: 998952
          cache_write: 100380
          cost: 0.272
formula: analysis-cadence
---
# S-0584 Analyze the project and list the issues the report names

## Goal

Run the analyzer on its focus, or on all three, and list the open issues its report names, for the operator to choose stories from at acceptance.

## Acceptance criteria

- [x] Run the analyzer and write its report
- [x] List the open issues the analyzer's report names

## Tasks
- T-2545 Run the analyzer and write its report
- T-2547 List the open issues the analyzer's report names

## Notes

Made by `flai formula run analysis-cadence` from `design/workflows/formulas/analysis-cadence.md`, with the inputs `focus=`.

### Report

The analyzer's run on S-0584 ended with its work done at 2026-11-07T08:22:44Z (T-2545). Its report is `design/analysis/2026-11-07-all.md`, focus all, in the main checkout. Step 2 (T-2547), `flai issue list --report latest --json`, exited 0 there and named these open issues:

| Issue | Class | Title |
|-------|-------|-------|
| I-0230 | efficiency | flai guard refuses the orchestrator a read or a permitted flai command when the command carries a stderr redirection such as 2>&1 or 2>/dev/null |
| I-0246 | defect | TestAHeldWaitFailsOnATruncatedStory and TestServeStartsADefinedRoleOnAnEventAndOnASchedule time out under flai test's Go run when the host's load is near 80 |
| I-0257 | impression | A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host |
| I-0268 | impression | ADR-0254 §5a's send-back of a story whose kluster deploy or checks failed, and its kluster acceptance blocker, are not built |
| I-0284 | efficiency | orchestration.permissions.set_cost_of_delay is off, so a cost of delay input the operator confirmed waits a day for the operator to write it |
| I-0285 | efficiency | after chains on planned stories hold ready stories for days while in-progress stays under two of five |
| I-0286 | efficiency | Declared touches match what a story changes in 16 percent of completed stories, so claims mislead the pull hold, the trial merge, and the conversations |
| I-0287 | efficiency | The merge queue builds one-story batches and runs three verifications per story landed |
| I-0288 | impression | The board stood empty for three days while the backlog held only drafts waiting on operator decisions nobody asked for again |
| I-0289 | efficiency | flai guard refuses the orchestrator git tag --list and git worktree list, which only read |
| I-0290 | defect | One line holding a NUL byte in flai serve's journal drops the whole close-outs section of flai stats |
| I-0291 | efficiency | Under the judgement release policy the publish permission publishes nothing: every release since 2026-11-02 was the operator's by hand |
| I-0292 | impression | No analysis.schedule or analysis.formula is set, so the analyzer ran once in thirty days |
| I-0293 | impression | Design documents describe a /conventions page, a /streams page, a mounted-repository readiness check, a flai OTEL span, a dashboard scan of wip/kanban, and a single lockfile that are not there |
| I-0294 | impression | The dashboard renders repository markdown with markdown-it html: true and no sanitizer, so a document any agent writes can run script as the operator |
| I-0295 | impression | The close-out job and the manifest's test commands run the story worktree's scripts on the host, so a container harness does not keep an agent's code off the host |
| I-0296 | defect | flai serve records a Claude Code usage-limit refusal as a failed permission_prompt check and never checks the version again |
| I-0297 | impression | CI actions and base images are pinned by tag, three workflows have no permissions block, and flaiover tracks two lockfiles |
| I-0298 | impression | Five source files under flai are over 2,000 lines |

The operator chooses at acceptance which of them become stories.

### Recorded while working

- I-0299 (defect): a role run that the stdio MCP server starts through `role_run` is never reaped, so flai serve counts its zombie as running and never settles it. The analyzer's run stayed unsettled from about 08:18Z until this agent's session ended (TH-0859).
