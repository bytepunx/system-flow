---
id: E-0024
type: epic
nature: feature
title: "Workflows as data: a project defines a strategic role or a repeatable formula of steps in files flai runs, and the built-in roles are expressed the same way, so a new role or a release process needs no flai release"
status: in-progress
owner: alex
created: 2026-10-16T04:28:02Z
updated: 2026-10-28T13:05:49Z
transitions:
  - to: ready
    at: 2026-10-20T19:01:17Z
    by: alex
  - to: backlog
    at: 2026-10-21T05:57:42Z
    by: alex
  - to: ready
    at: 2026-10-28T13:01:57Z
    by: alex
  - to: in-progress
    at: 2026-10-28T13:05:49Z
    by: agent-S-0414
tags: [cli, dashboard, template]
topics: [cli, orchestration, planning, template, conventions]
touches: [flai/internal/serve, flai/internal/guard, flai/internal/context, flai/internal/hostapi, flai/cmd, flaiover/src, design/system/strategic-agents.md, design/system/metrics.md, design/conventions/strategic-agents.md, ".claude/agents", docs/operators/index.md, docs/users/flai.md, template]
usage:
  source: sum
  seconds: 14144
  priced_by: harness
  turns:
    - day: 2026-10-28
      ceremony: 9
      hand_edits: 7
      work: 188
  models:
    - model: claude-opus-5-5
      input: 1440
      output: 594108
      cache_read: 99426368
      cache_write: 4464831
      cost: 57.3015
  strategic:
    - kind: planner
      seconds: 1866
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 1090
          output: 50181
          cache_read: 10121204
          cache_write: 346401
          cost: 11.9783
        - model: claude-haiku-5-5
          input: 84
          output: 28177
          cache_read: 2955958
          cache_write: 109701
          cost: 0.123
        - model: claude-opus-5-5
          input: 51
          output: 5737
          cache_read: 3726197
          cache_write: 294016
          cost: 1.282
    - kind: orchestrator
      seconds: 12627
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1830
          output: 28911
          cache_read: 426414047
          cache_write: 767694
          cost: 96.1237
        - model: claude-sonnet-5-5
          input: 44
          output: 199
          cache_read: 791923
          cache_write: 190371
          cost: 0.6752
---
# E-0024 Workflows as data: a project defines a strategic role or a repeatable formula of steps in files flai runs, and the built-in roles are expressed the same way, so a new role or a release process needs no flai release

## Outcome

The planner, the orchestrator, and the analyzer are fixed in Go: their prompts, their host actions, the packs `flai prime --role` gives them, the permissions they hold, the tools the guard lets them call, their activity documents, and their pages in the dashboard ([strategic-agents.md](../../../design/system/strategic-agents.md), [ADR-0075](../../../design/adrs/0075-the-planner-the-orchestrator-and-the-analyzer-prime-by-role-plan-orchestrate-or.md), [ADR-0079](../../../design/adrs/0079-the-planner-the-orchestrator-and-the-analyzer-each-log-their-activities-in-one.md), [ADR-0143](../../../design/adrs/0143-flai-guard-refuses-a-story-s-agent-the-planner-and-the-analyzer-a-call-that.md)). A fourth role, a reviewer, a release manager, a documentation gardener, is a flai release away, and a repeatable multi-step process such as publishing, importing a repository, or the analyzer's cadence is spelled out in prose and shell scripts. Gas Town defines formulas and molecules in TOML and instantiates them with tracked steps; Spec Kit has extensions, presets, and a catalog; Paperclip ships installable teams of roles and routines. This epic makes a role and a formula data that flai runs, expresses the three built-in roles in that form so one path serves both, and ships the baseline through the template.

- [ ] A role is a file under a folder the manifest names: its name, its harness prompt, the host action that starts it and when (on the operator's word, on an event, on a schedule), the topics its pack selects, the permissions it may hold, the MCP tools and commands the guard lets it call, and whether it may work a story. `flai serve` starts a role so defined as it starts the planner, the orchestrator, and the analyzer, with an activity document under `wip/agents/<role>.md` and the same journal entries.
- [ ] The planner, the orchestrator, and the analyzer are expressed in that format, shipped by the template as baseline files with a marker below which a project adds its own, and behave exactly as today: their tests pass unchanged, and `flai upgrade` brings a new baseline role to a project without touching what it added.
- [ ] `flai guard` reads a role's permitted calls from its file and refuses the rest as it refuses the built-in roles' today (ADR-0143), and `flai prime --role <name>` gives a defined role its pack within the budget ([ADR-0049](../../../design/adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md)).
- [ ] A formula is a file: a named sequence of steps with `after` between them, each a flai command, a script under `scripts/`, or a role run with a prompt and inputs; `flai formula run <name>` instantiates it as a story with its steps as tasks, so that it is worked, measured, and archived as any story is, and a formula may be started by a host action, a schedule, or a role.
- [ ] The template ships baseline formulas for what is prose today: publishing a release behind `orchestration.release`, importing a repository from the board, and the analyzer's cadence; each is run once in this repository before the epic is accepted.
- [ ] The dashboard's Planner, Orchestrator, and Analyzer pages become one page per defined role, with the settings panel, the activity document, and the runs, and `agent.status` carries every role by name.
- [ ] `metrics.md` generalizes, by an ADR, the strategic usage it reports per kind to usage per role, so a project's own role is charged and charted as the built-in ones are ([ADR-0083](../../../design/adrs/0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md), [ADR-0095](../../../design/adrs/0095-an-orchestrator-activity-s-usage-is-charged-evenly-to-the-work-items-it-named.md), [ADR-0100](../../../design/adrs/0100-an-analyzer-activity-s-usage-is-charged-evenly-to-the-issues-it-names-under.md)).
- [ ] `flai check` validates role and formula files, with a refusal for a role that grants itself a call the operator's permissions do not allow, and `docs/operators/index.md`, `docs/users/flai.md`, the convention `strategic-agents.md`, and `design/system/strategic-agents.md` describe how to define one.

## Stories

To be drafted by the planner. Order: the role file format and its validation first, read by `flai prime --role` and the guard, with the three built-in roles expressed in it and their tests green; then `flai serve` starting a defined role; then the dashboard's generic role page and status; then formulas and `flai formula run`; then the baseline formulas, each run here once; then metrics per role with its ADR; documentation and the template with each story.
- S-0414 A strategic role is a markdown file under design/workflows/roles that flai reads and flai check validates, and the template ships the planner, the orchestrator, and the analyzer as baseline role files
- S-0415 flai prime --role gives a role a project defines its pack within the budget, from the topics and the item kind its role file names, and the three built-in roles prime from their files unchanged
- S-0416 flai guard reads a role's permitted calls from its role file and refuses the rest, and the planner's, the orchestrator's, and the analyzer's rules move into their files with their guard tests unchanged
- S-0417 flai serve starts a role a project defines on the operator's word, behind the host action its file names, with an activity document under wip/agents and the same journal entries as the built-in roles
- S-0418 A role a project defines starts on an event, on a schedule, or for as long as its host action is on, as its role file's trigger says, coalesced and logged with its trigger as the planner's runs are
- S-0419 The planner, the orchestrator, and the analyzer start through the path a defined role takes, with their prompts read from their role files, and behave exactly as today
- S-0420 The dashboard's Planner, Orchestrator, and Analyzer pages become one role page per defined role, with its settings panel, its activity document, and its runs, and agent.status carries every role by name
- S-0421 A formula is a markdown file under design/workflows/formulas that flai check validates, and flai formula run makes a story from it with one task per step and the steps' after between them
- S-0422 A story's agent works a formula's instance on its branch, a role step runs its role charged to the instance, a failed step asks on a thread, and a formula starts from its schedule or a role
- S-0423 The template ships a baseline publish formula that fetches, evaluates orchestration.release, and publishes what is pending, the orchestrator starts it in place of calling release_publish, and it is run once here
- S-0424 The template ships baseline formulas for importing a repository from the board and for the analyzer's cadence, the board's import starts its formula, and each is run once here
- S-0425 metrics.md generalizes strategic usage per kind to usage per role by an ADR, so a project's own role is charged on its items, reported by flai stats, and charted as the built-in roles are

## Notes

- From the comparison of system-flow with Gas Town, Spec Kit, and Paperclip on 2026-10-16. Split from "workflow formulas and a memory digest of the archive" because the two share no code; the digest is its own epic.
- The planner reads first: `design/system/strategic-agents.md`; `design/conventions/strategic-agents.md`; `.claude/agents/*.md`, which stay the definitions' source for the harness (ADR-0130); ADR-0075, ADR-0079, ADR-0082, ADR-0084, ADR-0087, ADR-0099, ADR-0143; `flai/internal/serve` where each role is started; `flai/internal/guard/guard.go`.
- Open to the designer: where role and formula files live, since the layout allows no new top-level folder and `design/conventions` is for norms; whether a formula's instance is a story of nature `improvement` or a nature of its own; whether a role may be given a model and provider apart from the project's default agent, as a story's roles are.
- A role defined by a project runs as the operator, as every agent does: the safety note under `docs/operators/index.md` § Starting an agent when a story becomes ready applies to it, and the convention says so.
- 2026-10-28T13:01:57Z: moved to ready: follows S-0414, which moved to ready
- 2026-10-28T13:05:49Z: moved to in-progress: follows S-0414, which moved to in-progress
