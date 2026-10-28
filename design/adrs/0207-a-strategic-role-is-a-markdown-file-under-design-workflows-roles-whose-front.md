---
id: ADR-0207
title: "A strategic role is a markdown file under design/workflows/roles whose front matter names its action, trigger, item, topics, permissions, calls, and harness definition, whose body is its definition's prompt, and whose rules beyond an allow-list are named conditions flai implements"
status: accepted
date: 2026-10-28
supersedes: []
superseded_by: []
refines: [ADR-0060, ADR-0082, ADR-0087, ADR-0099, ADR-0130, ADR-0196]
topics: [planning, orchestration, analysis, conventions, cli, template]
---

# ADR-0207 A strategic role is a markdown file under design/workflows/roles whose front matter names its action, trigger, item, topics, permissions, calls, and harness definition, whose body is its definition's prompt, and whose rules beyond an allow-list are named conditions flai implements

## Context

The planner, the orchestrator, the analyzer, and the curator are written in Go. Their host actions, triggers, item kinds, and topics, and the calls `flai guard` lets them make, live in `flai/internal/serve`, `flai/internal/context/strategic.go`, and `flai/internal/guard/guard.go`. Their start prompts live in `flai/internal/harness/harness.go`. Their harness definitions are hand-written `.claude/agents/*.md` files, which `subAgents` in `flai/internal/harness/adapters.go` passes to Claude Code as `--agents` JSON. A new role, or a change to one, needs a flai release. E-0024 makes a role data that a project holds.

The operator's answers on TH-0471 settle three things this decision rests on. Role files and formula files live under one `design/workflows/` folder, in `roles/` and `formulas/` (answer 2). A role's model is `roles.<name>.agent` in the manifest, merged over `agent` as `planning.agent` is (answer 4). flai generates a role's harness definition from its role file rather than reading a hand-written one, and writes nothing of its own directly in `.claude/` (answer 6). TH-0567 places the generation in S-0419: a session `flai serve` starts takes the definition through `--agents`, which writes nothing.

Most of a role's rules are allow-lists: the tools, commands, and folders it may use. Some are not. The planner creates a story only as a draft and moves an item only to backlog. The analyzer edits only under one folder. The orchestrator replies on a thread only as `answer_threads` allows, and accepts a story only as itself. Those rules need code. [repository-layout.md](../system/repository-layout.md) adds no `design/` subfolder without an ADR.

## Decision

**A strategic role is one markdown file, `design/workflows/roles/<name>.md`. Its front matter names the role, its harness definition, its host action, its triggers, the items it takes, its topics, how its usage is charged, and the calls the guard allows it, each plain or under a named condition that flai implements, with each call only the operator may allow inside a `permissions` entry naming the operator's setting. Its body is its harness definition's prompt, with the baseline marker below which a project adds its own.** This refines [ADR-0060](0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md), [ADR-0082](0082-flai-serve-starts-the-planner-for-an-epic-or-a-story-behind-the-plan-host.md), [ADR-0087](0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md), [ADR-0099](0099-the-analyzer-runs-behind-the-analyze-host-action-and-writes-one-report-under.md), [ADR-0130](0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md), and [ADR-0196](0196-the-curator-runs-behind-the-curate-host-action-and-takes-in-what-the-archive.md); the rest of each stands.

### 1. Where it lives

- `design/workflows/` is a new subfolder of the design folder (the manifest's `layout.design`, `design/` by default). It holds what flai runs as data: `roles/` now, `formulas/` when E-0024's formula stories land. Its `README.md` says so.
- `design/workflows/roles/` holds one file per role, `<name>.md`, and a `README.md` that says what a role file holds and lists the roles. flai reads every `*.md` in it but `README.md`.
- The template ships `planner.md`, `orchestrator.md`, `analyzer.md`, and `curator.md` as baseline role files. A role file the template does not ship is the project's own.

### 2. The front matter

Every field is plain YAML. `flai check` reports `roles.file` for a file that does not parse, misses a required field, or gives a value outside the allowed ones.

| Field | Type | Meaning | Allowed values |
|-------|------|---------|----------------|
| `title` | string | The document's title, as the documentation standard asks. | Any; double-quoted when it holds a colon followed by a space. Required. |
| `updated` | date | The day of the last meaningful edit. | `YYYY-MM-DD`. Required. |
| `status` | string | The documentation standard's status. It gates nothing: the host action does. | `draft`, `active`, `deprecated`. Required. |
| `name` | string | The role's name. It names the activity document `wip/agents/<name>.md`, the harness definition (`--agent <name>`), the dashboard page, the run's `FLAI_AGENT` (`<name>`, or `<name>-<item>` for a run on an item), the `usage.strategic` kind, and the manifest's `roles.<name>` block. | Kebab-case; equal to the file stem; unique; not `explorer` or `verifier`, the sub-agents' definitions. Required. |
| `role` | string | The verb that `flai prime --role`, `FLAI_ROLE`, and a convention's `roles` ([ADR-0068](0068-a-convention-s-roles-list-every-agent-that-reads-it-the-story-s-agent-included.md)) take. | Kebab-case; unique; not `story`, `explore`, or `verify`. Required. |
| `definition` | mapping | The fields flai generates the harness definition from (§4). | `description`, `tools`, `model`. Required. |
| `definition.description` | string | When the harness should use the agent, as a definition's `description`. | Not empty. Required. |
| `definition.tools` | list of strings | The harness tools the session has, in Claude Code's names, flai's MCP tools as `mcp__flai__<tool>`. Another harness's adapter translates them through its tool-name table (ADR-0130). | Not empty. Required. |
| `definition.model` | string | The definition's model. The role's agent's model (§6) wins over it. | `inherit` or a model ID. Default `inherit`. |
| `action` | string | The host action that gates every start of the role. | A built-in role's action (`plan`, `orchestrate`, `analyze`, `curate`) for that role, or a new kebab-case name. Never a host action that is no role's: `push`, `forge`, `auto-publish`, `auto-approve`, `auto-upgrade`, `agent`, `dashboard`, `checks`, `host`, `settings`, `components`, `merge-queue`. Required. |
| `trigger` | list | What starts the role while its action is on (§3). | One entry or more. Required. |
| `item` | string or list | The item kinds a run takes. | `none`, or a list of `epic` and `story`. Required. |
| `topics` | list of strings | The topics the role's pack selects, first in the pack, as `flai prime --role` gives them. | Topic names `flai check` accepts. Default `[]`. |
| `works_story` | boolean | Whether the session may work a story as a story's agent does: `story_start`, a worktree, commits on its branch, `task_done`, a move to review. | Default `false`. |
| `charge` | string | How the usage of the role's activity is charged. S-0425 generalizes [ADR-0083](0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md), [ADR-0095](0095-an-orchestrator-activity-s-usage-is-charged-evenly-to-the-work-items-it-named.md), and [ADR-0100](0100-an-analyzer-activity-s-usage-is-charged-evenly-to-the-issues-it-names-under.md) by it. | `item` (the item its run was for), `named` (evenly to the work items its activity names), `issues` (evenly to the issues it names), `none` (its project strategic total). Default `none`. |
| `tools` | list of call entries | flai's MCP tools the session may call besides the reads every agent has. | Names of flai's MCP tools (§5). Default `[]`. |
| `commands` | list of call entries | flai commands the session may run besides the reads every agent has. | `flai <command>`, `flai <command> <subcommand>`, or `flai <command> *` (§5). Default `[]`. |
| `files` | list of strings | Folders under which the session may edit a file with `Edit`, `MultiEdit`, `Write`, or `NotebookEdit`. Each entry is held by `under-folder`. | Folder paths from the project root, ending in `/`; `{design}/` stands for the design folder. Not the root, and not `.git/`, `.claude/`, or `{design}/workflows/`, nor a folder holding one. Default `[]`. |
| `permissions` | list of mappings | The calls only an operator's setting allows, each listed under the setting that allows it (§5). | Each entry has `setting` and one or both of `tools` and `commands`. Default `[]`. |

A call entry, in `tools`, `commands`, or a permission's `tools` and `commands`, is either the call's name, which allows it as named, or a mapping `{call: <name>, when: <condition>}`, which allows it only as the condition says. Entries for one call are alternatives: the call passes when any of them passes.

### 3. Triggers

A `trigger` entry is one of these.

| Entry | Starts the role | Logged as |
|-------|-----------------|-----------|
| `asked` | On demand: by the operator, through `flai role run`, the host API, or the MCP tool; or by another role whose permission lets it ask, as the orchestrator asks the planner with `plan`. | `asked`, or the asker |
| `{event: <event>, when?: <condition>}` | When the event happens. A role that takes an item is queued once per item; a role with `item: none` runs once for the events seen meanwhile ([ADR-0084](0084-flai-serve-plans-again-on-its-own-behind-the-plan-host-action-on-an-edit-when.md), [ADR-0199](0199-the-curator-runs-after-acceptances-coalesced-and-on-digest-schedule-marks.md)). | `<event> <ID>` |
| `{schedule: <setting>, when?: <condition>}` | Each time the schedule comes round. The setting is the manifest key holding a five-field cron expression in UTC or `daily`: `planning.schedule`, `analysis.schedule`, `digest.schedule`, or `roles.<name>.schedule`. | `schedule <spec>` |
| `continuous` | While its action is on: started again when it ends, stopped when the action goes off (ADR-0087). It is then the role's only trigger. | `continuous` |

The events are `accepted` (a story moved to done), `cancelled` (an item cancelled), `ready` (a story moved to ready), `review` (a story moved to review), `edited` (an item's own words edited), `reordered` (the pull order changed), and `thread-opened` (a thread opened). An unknown kind, event, or setting is `roles.file`. S-0418 reads the triggers.

### 4. The body and the harness definition

- The body is the definition's prompt: the system prompt that a `.claude/agents/<name>.md` body is today. It carries the baseline marker `<!-- system-flow:end-of-baseline -->`. Above it is the baseline text; below it, under `## Project additions`, the project's own. The prompt is both, without the marker line.
- Every role file has the marker, the project's own included. `roles.file` then checks one rule, and a project's file keeps its additions if the template ever ships a file of the same name.
- The start prompt stays flai's: `planPrompt`, `orchestratePrompt`, `analyzePrompt`, and `curatePrompt` today, and for a defined role one flai builds from the run's item, trigger, and inputs. The start prompt says what this run is for; the body says what the role is.
- flai generates the definition from `name`, `definition.description`, `definition.tools`, `definition.model` with the role's agent's model over it, and the body as the prompt. For a session `flai serve` starts on `claude-code`, flai passes it in the `--agents` JSON that `subAgents` already builds, keyed by `name`, and selects it with `--agent <name>`. It writes nothing under `.claude/`.
- A session the operator starts by hand reads `.claude/agents/`. If Claude Code follows a symbolic link there, `.claude/agents/<name>.md` is a link to the definition flai generates outside `.claude/`, and the operator creates it once, since `.claude/` is protected. If Claude Code does not follow one, a hand-started session takes the definition through `--agents` too. S-0419 generates the definitions and finds out which.
- Until S-0419, `.claude/agents/*.md` stay the definitions. `flai check` reports `roles.definition` for a role file whose `definition.description` or `definition.tools` is missing, or whose `definition.tools` lists a tool its rules refuse: an `mcp__flai__<tool>` that is neither a read every agent has nor in `tools` or a permission's `tools`, or `Edit`, `MultiEdit`, `Write`, or `NotebookEdit` while `files` is empty.

### 5. Calls, conditions, and permissions

**What every role has, and what no file changes.** Every role may make the reads every agent has: flai's read MCP tools (`MCPReads` in the guard), flai's read commands and read forms, and git's reads. It writes only as itself, as [ADR-0143](0143-flai-guard-refuses-a-story-s-agent-the-planner-and-the-analyzer-a-call-that.md) holds the planner: no `--by` or `--placed-by` naming another, and no change to `FLAI_AGENT` or `FLAI_SESSION`. Its sub-agents are held as every sub-agent is (ADR-0060). Anything beyond that is what its file lists.

**Commands.** `flai <command>` allows the command with no subcommand, or in any form when the command has no subcommands (`flai edit`, `flai move`). `flai <command> <subcommand>` allows that subcommand. `flai <command> *` allows every subcommand. flai matches entries against its own command tree. A name that is no flai MCP tool or command is `roles.file`.

**Conditions.** A condition is a rule flai implements in Go, in `flai/internal/roles/conditions.go`, named in kebab-case. A file names one; it cannot define one. An unknown name, or a condition on a call it does not qualify, is `roles.file`. A gated condition is allowed only inside a `permissions` entry whose setting is the permission the table names.

| Condition | Qualifies | Allows | Gated by |
|-----------|-----------|--------|----------|
| `draft-only` | `item_new`, `flai story new` | A story only as a draft (`draft` true, `--draft`); an epic or a task as asked. | No |
| `to-backlog` | `item_move`, `flai move` | A move to backlog and nowhere else. | No |
| `ungated-edit` | `item_edit`, `flai edit` | Any edit that neither finalizes a draft (`draft` false, `--no-draft`) nor sets a cost of delay input. | No |
| `own-thread` | `thread_reply`, `thread_resolve`, `flai thread reply`, `flai thread resolve` | On a thread the role opened: a follow-up, never a recommendation or a source; and resolving it. | No |
| `under-folder` | each `files` entry | An edit of a file under the folder, after `..` and symbolic links are resolved, never of the folder itself. | No |
| `plan-epic` | `plan`, `flai plan` | Starting the planner on an epic. | `plan_backlog_epics` |
| `plan-story` | `plan`, `flai plan` | Starting the planner on a story. | `plan_backlog_stories` |
| `cost-inputs` | `item_edit`, `flai edit` | Giving a story cost of delay inputs and changing nothing else: no value, no clearing ([ADR-0119](0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md)). | `plan_backlog_stories` |
| `planner-thread` | `thread_reply`, `thread_resolve`, `flai thread reply`, `flai thread resolve` | On a thread a story's planner opened (`planner-S-nnnn`): an answer, with or without a source, or a recommendation; and resolving it (ADR-0119). | `plan_backlog_stories` |
| `finalize-draft` | `item_edit`, `flai edit` | Finalizing a draft and changing nothing else. | `finalize_drafts` |
| `to-ready` | `item_move`, `flai move` | A move to ready and nowhere else. | `promote_to_ready` |
| `by-policy` | `flai order` | `flai order --by <policy> --apply`, placing no story by hand. | `order_ready` |
| `answer-threads` | `thread_reply`, `flai thread reply` | On a thread the role did not open: a recommendation while the setting is `recommend`; a recommendation, or an answer citing a source, while it is `autonomous`; never resolving it ([ADR-0090](0090-a-thread-entry-is-marked-a-recommendation-in-its-heading-and-cites-its-source.md)). | `answer_threads` |
| `accept-as-self` | `flai accept`, `flai move` | Accepting a story in review, or moving it to done, as the role itself (`--by <name>`), never with `--batch` ([ADR-0093](0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md)). | `accept_reviews` |
| `enqueue-as-self` | `flai queue add` | Enqueuing a story for the merge queue as the role itself (`--by <name>`). | `accept_reviews` |

**Trigger conditions.** Three conditions qualify a trigger entry rather than a call. They express the planner's replanning (ADR-0084).

| Condition | Qualifies | Starts the role |
|-----------|-----------|-----------------|
| `stale-figures` | `{event: edited}` | Only for an edit of a backlog or ready story's goal, criteria, or touches by someone other than a planner, when the story's forecast or cost of delay value is older than the edit. |
| `replan-agent` | `{event: accepted}`, `{event: cancelled}`, `{event: reordered}` | Only while `planning.replan` is `agent`, for each story whose delivery moved. |
| `each-ready-story` | `{schedule: <setting>}` | Once for every ready story. |

**Gated calls.** Some calls only the operator's settings may allow. A role file names each only inside a `permissions` entry whose `setting` is the permission that allows it and, where the table names one, under that condition. Anywhere else, `flai check` reports `roles.grant`.

| Call | Permission | How a file grants it |
|------|------------|----------------------|
| Starting the planner on an epic: `plan`, `flai plan E-nnnn` | `plan_backlog_epics` | `plan-epic` |
| Starting the planner on a story: `plan`, `flai plan S-nnnn` | `plan_backlog_stories` | `plan-story` |
| Setting a story's cost of delay inputs: `item_edit`, `flai edit` | `plan_backlog_stories` | `cost-inputs` |
| Finalizing a draft: `item_edit`, `flai edit`; creating a story that is not a draft: `item_new`, `flai story new` | `finalize_drafts` | `finalize-draft`; creating, plain |
| A move past backlog: `item_move`, `flai move` | `promote_to_ready` | `to-ready` |
| Applying an order: `flai order` | `order_ready` | `by-policy` |
| A reply on another's thread | `answer_threads` | `answer-threads` |
| An acceptance: `flai accept`, `flai move S-nnnn done`, `queue_add`, `flai queue add` | `accept_reviews` | `accept-as-self`; `queue_add` plain; `enqueue-as-self` |
| A publish: `release_publish` | `publish` | plain |
| A change to the sub-projects: `subproject_edit`, `flai subproject add`, `set`, `remove` | `edit_subprojects` | plain |

A plain `item_edit`, `flai edit`, `item_move`, `flai move`, `item_new`, `flai story new`, `plan`, `flai plan`, `flai order`, or `flai accept` outside a permission is therefore `roles.grant`: each would allow a gated call. A role replies on threads plainly, as the planner does; `answer-threads` is how the orchestrator's replies are held to its setting.

**Never granted.** No setting allows these, and a file that names them anywhere is `roles.grant`:

- `flai release` but with `--evaluate`, `flai push`, `git push`, `git tag`, and git's other writes;
- `flai accept --batch`, and `flai queue remove`, `hold`, and `resume`;
- placing a story by hand with `flai order`, and `flai thread confirm`;
- `item_move` to done;
- `shared_paths_edit`, and `flai shared add` and `remove` ([ADR-0096](0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md));
- `flai forge push` (S-0393);
- while `works_story` is `false`, the calls of a story's agent: `story_start`, `task_done`, `criteria_tick`, `stream_state`, the `message_` tools, and a move to in-progress or review.

**Permissions.** A `permissions` entry is `{setting, tools?, commands?}`. `setting` names the operator's switch in the manifest. The orchestrator's file names `orchestration.permissions.<p>`. Any other role's file names only `roles.<its name>.permissions.<p>`, so that turning on one of the orchestrator's permissions grants no other role anything. `<p>` is one of `plan_backlog_epics`, `plan_backlog_stories`, `finalize_drafts`, `promote_to_ready`, `order_ready`, `answer_threads`, `accept_reviews`, `publish`, and `edit_subprojects`. `answer_threads` takes `off`, `recommend`, or `autonomous`; the others are booleans. The guard reads the setting at each call, as it reads `orchestration.permissions` today. The `roles.<name>.permissions` block reaches the manifest with S-0416, which is the first to read it.

### 6. The agent

A role's agent is `roles.<name>.agent` in the manifest, merged over the project's `agent` as `planning.agent` is. `planning.agent`, `orchestration.agent`, `analysis.agent`, and `digest.agent` stay valid and win for their roles. S-0417 implements it.

### 7. It runs as the operator

A defined role's session runs on the host, as the built-in roles' sessions do: as the operator's user, with the operator's git and forge credentials, the provider's key, and every file that user can read. `flai guard` holds it to its file for the calls the guard sees: flai's MCP tools, flai and git command lines, and file edits. The guard is not a shell. It does not judge a program other than flai or git that the session runs through `Bash`, so a role whose `definition.tools` includes `Bash` has the operator's shell. A role file is code the operator runs. A change to one is reviewed as a change to a script is, and reaches main only through a story the operator accepts. `files` never includes the role files' own folder, so a role cannot widen itself.

### 8. Upgrade

`flai upgrade` never touches a role file the template does not ship: `Compute` classifies only the files the template renders. A baseline role file carries the marker, so `flai upgrade` merges it as it merges a convention (`hasMarker(existing) && hasMarker(c.content)` in `Compute`, then `merge`):

- the front matter and the body above the marker come from the template;
- the project's text from the marker on stays;
- the project's `topics` stay when they differ from the topics the lock records the template last gave the file;
- a project's change to any other field is replaced;
- a baseline role file the project lacks is added.

A project that wants a baseline role to act otherwise turns its action or its permissions off and writes a role of its own under another name.

### 9. No behaviour changes in S-0414

In S-0414 the role files are data. flai keeps deciding from its Go until S-0415 (packs), S-0416 (the guard), S-0417 and S-0418 (starts and triggers), and S-0419 (the built-in roles and their generated definitions) read the files. `flai role list` and `flai role show` print the roles, and `flai check` reports `roles.file`, `roles.grant`, and `roles.definition`.

### The baseline roles

| Field | `planner.md` | `orchestrator.md` | `analyzer.md` | `curator.md` |
|-------|--------------|-------------------|---------------|--------------|
| `role` | `plan` | `orchestrate` | `analyze` | `curate` |
| `action` | `plan` | `orchestrate` | `analyze` | `curate` |
| `trigger` | `asked`; `edited` with `stale-figures`; `accepted`, `cancelled`, `reordered` with `replan-agent`; `planning.schedule` with `each-ready-story` | `continuous` | `asked`; `analysis.schedule` | `asked`; `accepted`; `digest.schedule` |
| `item` | `[epic, story]` | `none` | `none` | `none` |
| `topics` | `[planning]` | `[orchestration]` | `[analysis]` | `[digest, planning]` |
| `charge` | `item` | `named` | `issues` | `none` |
| `tools` | `activity_log`, `inbox`, `item_new` (`draft-only`), `item_edit` (`ungated-edit`), `item_move` (`to-backlog`), `thread_open`, `thread_reply`, `wait_for_events` | `activity_log`, `inbox`, `thread_open`, `wait_for_events`, `thread_reply` and `thread_resolve` (`own-thread`) | `activity_log`, `inbox`, `issue_new`, `issue_bump`, `thread_open`, `thread_reply`, `wait_for_events` | `activity_log`, `inbox`, `thread_open`, `thread_reply`, `wait_for_events` |
| `commands` | `flai story new` (`draft-only`), `flai epic new`, `flai task new`, `flai edit` (`ungated-edit`), `flai touches *`, `flai thread new`, `flai thread reply`, `flai issue new`, `flai issue bump`, `flai move` (`to-backlog`) | `flai issue new`, `flai issue bump`, `flai queue`, `flai thread new`, `flai thread reply` and `flai thread resolve` (`own-thread`) | `flai issue new`, `flai issue bump` | `flai digest harvest`, `new`, `edit`, `absorb`, `cited`, `compact` |
| `files` | none | none | `{design}/analysis/` | none |
| `permissions` | none | all nine, as below | none | none |

### Example: `orchestrator.md`

The body is cut short here; the baseline file carries today's `.claude/agents/orchestrator.md` body.

```markdown
---
title: Orchestrator
updated: 2026-10-28
status: active
name: orchestrator
role: orchestrate
definition:
  description: Keeps the work of this project moving, for flai serve, which runs it for as long as the orchestrate host action is on. It plans, finalizes, promotes, orders, answers threads, accepts, publishes, and changes sub-projects, each only while its permission in orchestration.permissions is on. It writes work items and threads through flai only, never edits a file, and never works a story.
  tools: [Read, Grep, Glob, Bash, Agent, mcp__flai__prime, mcp__flai__inbox, mcp__flai__board, mcp__flai__item_get, mcp__flai__item_edit, mcp__flai__item_move, mcp__flai__doc_get, mcp__flai__doc_search, mcp__flai__thread_get, mcp__flai__thread_open, mcp__flai__thread_reply, mcp__flai__thread_resolve, mcp__flai__who_touches, mcp__flai__activity_log, mcp__flai__wait_for_events, mcp__flai__plan, mcp__flai__order_by_policy, mcp__flai__promote_candidates, mcp__flai__release_evaluate, mcp__flai__release_publish, mcp__flai__verify, mcp__flai__queue_add, mcp__flai__subprojects]
  model: inherit
action: orchestrate
trigger: [continuous]
item: none
topics: [orchestration]
works_story: false
charge: named
tools:
  - activity_log
  - inbox
  - thread_open
  - wait_for_events
  - {call: thread_reply, when: own-thread}
  - {call: thread_resolve, when: own-thread}
commands:
  - flai issue new
  - flai issue bump
  - flai queue
  - flai thread new
  - {call: flai thread reply, when: own-thread}
  - {call: flai thread resolve, when: own-thread}
files: []
permissions:
  - setting: orchestration.permissions.plan_backlog_epics
    tools: [{call: plan, when: plan-epic}]
    commands: [{call: flai plan, when: plan-epic}]
  - setting: orchestration.permissions.plan_backlog_stories
    tools:
      - {call: plan, when: plan-story}
      - {call: item_edit, when: cost-inputs}
      - {call: thread_reply, when: planner-thread}
      - {call: thread_resolve, when: planner-thread}
    commands:
      - {call: flai plan, when: plan-story}
      - {call: flai edit, when: cost-inputs}
      - {call: flai thread reply, when: planner-thread}
      - {call: flai thread resolve, when: planner-thread}
  - setting: orchestration.permissions.finalize_drafts
    tools: [{call: item_edit, when: finalize-draft}]
    commands: [{call: flai edit, when: finalize-draft}]
  - setting: orchestration.permissions.promote_to_ready
    tools: [{call: item_move, when: to-ready}]
    commands: [{call: flai move, when: to-ready}]
  - setting: orchestration.permissions.order_ready
    commands: [{call: flai order, when: by-policy}]
  - setting: orchestration.permissions.answer_threads
    tools: [{call: thread_reply, when: answer-threads}]
    commands: [{call: flai thread reply, when: answer-threads}]
  - setting: orchestration.permissions.accept_reviews
    tools: [queue_add]
    commands:
      - {call: flai accept, when: accept-as-self}
      - {call: flai move, when: accept-as-self}
      - {call: flai queue add, when: enqueue-as-self}
  - setting: orchestration.permissions.publish
    tools: [release_publish]
  - setting: orchestration.permissions.edit_subprojects
    commands: [flai subproject add, flai subproject set, flai subproject remove]
---

You are the orchestrator: the agent flai serve runs to keep the work of this project moving while the operator has the orchestrate host action on. You act through flai. You never change a file through the file system, and you never work a story.

1. Call the flai MCP tool `prime` with role `orchestrate` before anything else. ...
2. Call `inbox`, and read the board with `board`.
3. Act only within `orchestration.permissions` in `system-flow.yaml`, each off by default, and by `orchestration.policy`. ...

<!-- system-flow:end-of-baseline -->

## Project additions
```

## Consequences

- A project defines a role, or changes what a built-in role is told, without a flai release. Its rules are reviewed as a file in a story, not as Go.
- The four built-in roles become four files. Once S-0415 to S-0419 land, the guard, the packs, and `flai serve` read one shape for every role.
- Every rule a file can state is one flai implements and tests. A project that needs a rule no condition states asks flai for a new condition. That is slower than writing the rule, and it keeps the guard's rules in tested code.
- `flai check` refuses a file that grants its role an acceptance, a publish, a move past backlog, finalizing a draft, or cost of delay inputs outside the permission the operator holds for it.
- The planner's file states rules the guard enforces only in part today. The guard refuses the planner `--no-draft` on the command line, but lets its `item_edit` with `draft` false, and any cost of delay inputs it sets, through. `ungated-edit` holds both. Once S-0416 reads the file, the guard enforces them, as the strategic-agents convention already asks.
- Every role comes under ADR-0143's own-name rule once S-0416 reads its file, the orchestrator included. Today the orchestrator is held by its own `--by orchestrator` checks.
- A baseline role file's front matter is the template's. A project's change to it, other than its topics, is lost on upgrade. A project changes a baseline role through its additions below the marker, its topics, its permissions, and its action, or writes a role of its own.
- Until S-0419 generates the definitions, `.claude/agents/planner.md`, `orchestrator.md`, `analyzer.md`, and `curator.md` duplicate the role files' definition fields, and `roles.definition` keeps the two consistent with the rules.
- `repository-layout.md`, `design/README.md`, and `CLAUDE.md`, in the template and here, list `design/workflows/`.

## Alternatives considered

- **`design/roles/` and `design/formulas/`, two folders.** The planner's recommendation on TH-0471; the operator chose one `design/workflows/` folder.
- **A role file that names a hand-written `.claude/agents/<name>.md`.** The operator chose generation (TH-0471, answer 6). One file then holds a role, and nothing under `.claude/` needs the operator's write.
- **Rules as expressions in the file**, such as `item_move: {to: [backlog]}`. The guard would need a small language to evaluate, and a file could write a rule no test covers. Named conditions keep the logic in tested Go and the file readable.
- **Several conditions per entry, all of which must hold.** No rule of today's needs two. One condition per entry, with entries as alternatives, is simpler to read and to check.
- **A permission as a flag on each call entry**, such as `{call: release_publish, permission: publish}`. A `permissions` block shows in one place everything a setting unlocks, as the Settings page lists it, and gives `roles.grant` one place to check.
- **A defined role naming `orchestration.permissions`.** Turning on the orchestrator's `accept_reviews` would then grant every role that names it. Each role gets its own switches.
- **`status` gating whether a role runs.** The host action already does; two switches could disagree.
- **The marker in baseline files only.** One rule for every file is simpler to check, and a project's file keeps its additions if the template later ships one of its name.
