---
id: ADR-0250
title: "A role file's start prompt is a text/template in its front matter's start field, rendered over the run's fields with a project's additions below the marker, and the built-in roles start on their files' triggers through their own launchers"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0207]
topics: [planning, orchestration, analysis, conventions, cli, template]
---

# ADR-0250 A role file's start prompt is a text/template in its front matter's start field, rendered over the run's fields with a project's additions below the marker, and the built-in roles start on their files' triggers through their own launchers

## Context

[ADR-0207](0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) made a strategic role a file under `design/workflows/roles`, whose body is its harness definition's prompt. Its § 4 kept the start prompt in flai: `planPrompt`, `orchestratePrompt`, `analyzePrompt`, and `curatePrompt` in `flai/internal/harness/harness.go`, and, for a role a project defines, one flai builds from the run's item and inputs (`definedPrompt`). The start prompt says what this run is for; the body says what the role is. Its § 3 named the triggers each baseline role file carries, but only the orchestrator's ([ADR-0209](0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md)) and a defined role's are read. The planner still starts on the rules `flai/internal/serve/replan.go` holds ([ADR-0084](0084-flai-serve-plans-again-on-its-own-behind-the-plan-host-action-on-an-edit-when.md)), the analyzer on `analysis.schedule` as `analysis_schedule.go` reads it ([ADR-0099](0099-the-analyzer-runs-behind-the-analyze-host-action-and-writes-one-report-under.md)), and the curator on acceptances and `digest.schedule` as `curate_schedule.go` reads them ([ADR-0199](0199-the-curator-runs-after-acceptances-coalesced-and-on-digest-schedule-marks.md)). The trigger conditions `stale-figures`, `replan-agent`, and `each-ready-story` are named in `flai/internal/roles/conditions.go` but not implemented.

On TH-0760 the operator asked that the built-in roles be defined by their files alone, so that a project changes what a built-in role is told, or when it runs, with no flai release. The orchestrator's confirmed recommendation there settled the shape: a start prompt templated over the run's fields, under an ADR of its own, with the body left as the definition's prompt, and the replan rules and the schedules read as `event` and `schedule` triggers.

## Decision

**A role file's start prompt is data: a Go `text/template` in its front matter's `start` field, rendered over the run's fields, with a project's additions under a `### Start prompt` heading below the marker. The planner, the analyzer, and the curator start on their files' `trigger` entries, judged by the functions the defined roles' launchers use, while their own launchers keep their queues, records, and journal entries.** This refines ADR-0207 § 3 and § 4; the rest of ADR-0207 stands, and so do ADR-0084, ADR-0099, ADR-0199, and ADR-0209, but for where each reads its triggers.

### 1. Where the start prompt lives

- `start` is a front matter field: a string, written as a YAML block scalar (`start: |`), holding a Go `text/template`. The body stays the definition's prompt, whole, as ADR-0207 § 4 has it, and S-0553's generation reads it as before.
- It is a front matter field rather than a fenced section of the body because the body is the definition's prompt, which flai would otherwise have to cut a section out of before generating it, and because a field is parsed, validated, and merged by `flai upgrade` as every other field is.
- A project's additions go under a `### Start prompt` heading in the body's `## Project additions`, below the marker. The text from that heading to the next heading of level three or above, or the end of the body, is rendered with the same fields and appended to the rendered `start` after a blank line. It is left out of the definition's prompt, so a project adds to what a run is told without adding to what the role is.
- A built-in role whose file has no `start`, as a project's file from before this decision has, renders the baseline's `start` that flai carries for it. A defined role whose file has no `start` renders flai's default start prompt, which says what `definedPrompt` says today: what its file's body says, on its item, with its inputs.

### 2. The template and the fields

The template is Go's `text/template`, with no functions beyond its built-ins. The rendered text is trimmed of the whitespace around it. Its fields:

| Field | Type | What it holds |
|-------|------|---------------|
| `.Name` | string | The run's `FLAI_AGENT`: the role's name, or `<name>-<item>` for a run on an item. |
| `.Root` | string | The project's directory, where the run works. |
| `.Role` | string | The role's verb, its file's `role`. |
| `.RoleName` | string | The role's name, its file's `name`. |
| `.File` | string | The role file's path from the project root. |
| `.Body` | string | The definition's prompt, the body without the marker, for a start prompt that repeats it, as the default does. |
| `.Item` | string | The epic or story the run is on, empty for none. |
| `.Kind` | string | The item's kind, `epic` or `story`, empty for none. |
| `.AskedOn` | string | The thread the operator wrote on since the item was planned, when the orchestrator asked for it again ([ADR-0222](0222-with-plan-backlog-epics-or-plan-backlog-stories-on-the-orchestrator-may-ask-the.md)); empty otherwise. |
| `.Reasons` | list | Why the run was started, each with `.Line`, its line, and `.Prime`, the form `flai prime --reason` takes it in, empty when it has none ([ADR-0217](0217-a-pack-primed-for-the-orchestrator-s-reasons-loads-the-role-s-conventions-with.md)). |
| `.Primed` | bool | Whether every reason has a form, so that the run primes with them. |
| `.PrimeForms` | list of strings | Each reason's form once, in the order first given. |
| `.Queue` | bool | Whether the merge queue runs for the project ([ADR-0173](0173-a-merge-queue-in-flai-serve-with-no-model-runs-behind-the-merge-queue-host.md)). |
| `.Focus` | string | What the analyzer looks for, one of `bottlenecks`, `intent`, and `risk`; empty for all three. |
| `.Since` | string | The window's start since the role's last run, as the curator's run takes it; empty on a first run. |
| `.Inputs` | list | A defined role's run's inputs, in the order given, each with `.Name` and `.Value`, its whitespace folded to single spaces. |

A field outside the table is `roles.file`: `flai check` parses the template and executes it over a run with every field empty and over one with every field set, and reports the template's error with the field it names. So is a template that does not parse, and so is a `### Start prompt` addition that does either.

### 3. Upgrade

`flai upgrade` merges a baseline role file as ADR-0207 § 8 says: the front matter, `start` among it, comes from the template, and the project's text from the marker on stays, its `### Start prompt` with it. A project that wants a built-in role told more adds it there, and keeps it across upgrades. A project's edit of `start` itself, or of `trigger`, holds until the next `flai upgrade` replaces it, as an edit of any other field of a baseline file but `topics` does. A project that wants a built-in role told otherwise, or started otherwise, for good turns the role's action off and defines a role of its own, as ADR-0207 § 8 says.

### 4. The built-in roles' triggers

Each built-in role's launcher in `flai serve` reads its role file, the project's or the baseline flai carries, at each look, and starts the role only for the trigger entries the file names, under their conditions. The orchestrator's are read so already (ADR-0209), and its reasons-and-end run stays. The functions that judge an event or a schedule against a trigger entry are the ones the defined roles' launchers use (`flai/internal/serve/role_events.go` and `role_schedule.go`), so a built-in role and a defined one are judged alike. Each launcher keeps how its role's runs start: the planner's queue, its coalescing and `planCheck`, and its records under `plans`; the analyzer's `analyzeCheck`, its formula under `analysis.formula`, and its record; and the curator's coalesced acceptances, its usage-limit retry, and its record. Their journal entries stay as they are.

| Role | Entries its launcher serves | Started |
|------|-----------------------------|---------|
| Planner | `{event: edited}`, `{event: accepted}`, `{event: cancelled}`, `{event: reordered}`, with or without a condition; `{schedule: <setting>, when: each-ready-story}` | Once per item, through its queue. An event without a condition queues the event's item when it is of a kind the role takes; `reordered` has none. |
| Analyzer | `{schedule: <setting>}` | Once each time the setting comes round, for all three focuses, or the formula `analysis.formula` names. |
| Curator | `{event: accepted}`; `{schedule: <setting>}` | Once for the acceptances seen meanwhile, coalesced, and once each time the setting comes round. |

`asked` stays every role's on-demand start. An entry a built-in role's launcher does not serve is `roles.file` in `flai check`, naming the entries it serves, so that nothing a file says is silently passed over. A schedule entry on a role that takes an item, but for `each-ready-story`, is `roles.file` too: such a run has no item to be on. Taking an entry out of the file stops the role's starts on it: a project's planner file without `{schedule: planning.schedule, when: each-ready-story}` queues nothing when `planning.schedule` comes round.

### 5. The trigger conditions

The three conditions ADR-0207 § 5 names are implemented in `flai/internal/roles/conditions.go`, as the rules the replanner holds today, taking what each needs as arguments so that they stay free of `flai serve`:

- `stale-figures`, on `{event: edited}`: the edit is of a backlog or ready story's goal, criteria, or touches, or of its cost of delay's inputs, by someone other than a planner (`planner-<item>`), and the story's forecast or cost of delay value is older than the edit (ADR-0084).
- `replan-agent`, on `{event: accepted}`, `{event: cancelled}`, and `{event: reordered}`: `planning.replan` is `agent`, and it queues each story whose forecast delivery the replay of the forecasts moved, rather than the event's item.
- `each-ready-story`, on a schedule: the ready stories, in the board's pull order, each queued once.

Replaying the forecasts after an acceptance, a cancellation, or a reorder stays flai's own work while `planning.replan` is `flai` or `agent`: it is no run of the planner, and no trigger entry gates it.

## Consequences

- The four built-in roles are defined by their files alone: what a run is told, and when it starts, can change in a project with no flai release, through a story the operator accepts, as any change to a role file does (ADR-0207 § 7).
- `harness.go` keeps the story agent's start prompt and the fields it renders a role's start prompt over, and none of the strategic roles' text.
- Three copies of each baseline role file, here, in the template, and embedded in flai, carry the start prompt, and the baseline test keeps them equal, as it keeps their bodies.
- A start prompt is a template, so a project's mistake in one is found by `flai check`, not by a run that starts with a broken prompt.
- The records stay keyed as they are, `plans`, `orchestrator`, `analyzer`, and `curator`; keying them by the role's name is left to a later story.

## Alternatives considered

- **A fenced section of the body.** It keeps the start prompt beside the definition's prompt, but the body is that prompt, and every reader of it, the definition's generation, the packs, and the dashboard, would have to cut the section out. Rejected.
- **A file of its own beside the role file**, such as `planner.start.md`. It keeps the front matter short, but splits one role over two files that `flai upgrade`, `flai check`, and the copies in the template and flai must keep together. Rejected.
- **Moving the starts of the built-in roles into the defined roles' launchers whole**, so that `role_events.go` and `role_schedule.go` start the planner, the analyzer, and the curator too. It would leave one launcher, but the planner's per-item queue and coalescing, the curator's coalesced acceptances and retry, and the analyzer's formula would have to move with them, and their records and journal entries would change with no gain to a project. Rejected: the launchers share how a trigger is judged, not how a run starts.
- **Keeping a project's `trigger` across `flai upgrade`, as its `topics` are.** It would make a project's change to a built-in role's triggers last, but it changes `flai upgrade` for one field and leaves a project's file silently apart from the template's. Left to a later decision if a project needs it.
