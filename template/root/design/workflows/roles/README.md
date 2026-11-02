# roles

One file per strategic role, `<name>.md`: an agent `flai serve` starts as a session of its own behind a host action. flai reads every `*.md` here but this README. [ADR-0207](../../adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) defines the format.

A role file's front matter names the role (`name`, `role`), the harness definition flai generates for it (`definition`: `description`, `tools`, `model`), the host action that gates it (`action`), what starts it (`trigger`), the items a run takes (`item`), its `topics`, how its usage is charged (`charge`), and the calls the guard allows it: flai's MCP tools (`tools`), flai's commands (`commands`), the folders it may edit (`files`), the calls only an operator's setting allows (`permissions`), and what one run is told (`start`). A call is named plainly or under a condition flai implements. The body is the definition's prompt, with the baseline marker; a project adds its own below it, under `## Project additions`, and keeps it through `flai upgrade`, as it keeps its `topics`.

The template ships four baseline roles. A role file the template does not ship is the project's own. `flai role list` prints every role; `flai check` reports `roles.file`, `roles.grant`, and `roles.definition`, and advises `roles.baseline` for each of the four the project lacks, whose baseline flai runs in its place. Convention: [strategic-agents.md](../../conventions/strategic-agents.md).

The trigger column says what each baseline file's `trigger` names. `flai serve` starts each role on its file's entries alone ([Triggers](#triggers)).

| File | Role | Action | Trigger |
|------|------|--------|---------|
| [planner.md](planner.md) | `plan` | `plan` | Asked; an edit that makes a story's figures stale; an acceptance, cancellation, or reorder while `planning.replan` is `agent`; `planning.schedule`, for every ready story |
| [orchestrator.md](orchestrator.md) | `orchestrate` | `orchestrate` | Asked; a thread that awaits it or the operator; a move to review or to backlog; an edit of a draft; a planner's run that wrote items; `release-due`; `ready-empty`; a reorder; `orchestration.schedule`. Each run acts on its reasons and ends |
| [analyzer.md](analyzer.md) | `analyze` | `analyze` | Asked; `analysis.schedule` |
| [curator.md](curator.md) | `curate` | `curate` | Asked; an acceptance; `digest.schedule` |

## The harness definition

flai generates a role's whole harness definition from its file: `name`, `definition.description`, `definition.tools`, `definition.model`, and the body without the marker line as the prompt. `flai serve` puts the role's agent's model over `definition.model` and passes the definition in `--agents` at every start, reading nothing under `.claude/`. So a role's prompt and definition are changed in its file alone ([ADR-0236](../../adrs/0236-a-built-in-role-s-claude-code-definition-is-generated-into-flai-cache-agents-by.md)).

- `flai role definition <name>` prints the generated definition as a Claude Code agent file.
- `flai role definition --write [<name>...]` writes each named role's definition, or every built-in role's when none is named, to `.flai-cache/agents/<name>.md`, which git ignores. For each built-in role whose `.claude/agents/<name>.md` is not a link to it, it prints the command that makes the link, `ln -sf ../../.flai-cache/agents/<name>.md .claude/agents/<name>.md`. flai never writes under `.claude/`, so the operator runs each command once and commits the links.
- A session the operator starts by hand reads `.claude/agents/<name>.md`, the link. Never edit the generated file: after a role file changes, run `flai role definition --write` again. Until it runs in a checkout, such as a fresh clone or a story's worktree, the link dangles and a session started there has no such agent.
- `flai check` reports `roles.definition` for a built-in role's `.claude/agents/<name>.md` that is neither the link nor equal to the generated definition, as a copy made before the template stopped shipping one may be. A missing file or a dangling link is no finding.

## The start prompt

The body says what the role is. The start prompt says what one run is for: the item it is on, why it was started, and what to do first. It is the front matter's `start` field, a Go `text/template` written as a block scalar, `start: |` ([ADR-0250](../../adrs/0250-a-role-file-s-start-prompt-is-a-text-template-in-its-front-matter-s-start-field.md)). flai renders it over the run's fields at every start and trims the whitespace around it. The template has no functions beyond `text/template`'s own.

| Field | What it holds |
|-------|---------------|
| `.Name` | The run's `FLAI_AGENT`: the role's name, or `<name>-<item>` for a run on an item |
| `.Root` | The project's directory, where the run works |
| `.Role` | The role's verb, the file's `role` |
| `.RoleName` | The role's name, the file's `name` |
| `.File` | The role file's path from the project root |
| `.Body` | The definition's prompt: the body without the marker and without the `### Start prompt` section |
| `.Item` | The epic or story the run is on; empty for none |
| `.Kind` | The item's kind, `epic` or `story`; empty for none |
| `.AskedOn` | The thread the operator wrote on since the item was planned, when the orchestrator asked for it again; empty otherwise |
| `.Reasons` | Why the run was started, each with `.Line`, its line, and `.Prime`, the form `flai prime --reason` takes; empty for none |
| `.Primed` | Whether every reason has a form, so that the run primes with them |
| `.PrimeForms` | Each reason's form once, in the order first given |
| `.Queue` | Whether the merge queue runs for the project |
| `.Focus` | What the analyzer looks for: `bottlenecks`, `intent`, or `risk`; empty for all three |
| `.Since` | The start of the curator's window, its last run; empty on a first run |
| `.Inputs` | A defined role's run's inputs, in the order given, each with `.Name` and `.Value` |

A project adds to a role's start prompt under a `### Start prompt` heading in its `## Project additions`, below the marker. The text runs to the next heading of level three or above, or the end of the body. flai renders it over the same fields and appends it to the rendered `start` after a blank line. It stays out of the definition's prompt, so it changes what a run is told and not what the role is. `flai upgrade` keeps it, as it keeps all the text from the marker on.

A built-in role's file with no `start`, such as one from before ADR-0250, renders the baseline's `start` that flai carries. A role the project defines with no `start` renders flai's default: who the role is, its file's body, its item, and its inputs. `flai check` reports `roles.file` for a `start` or a `### Start prompt` section that does not parse, or that names a field not in the table.

## Triggers

A `trigger` entry is `asked`, `{event: <event>, when?: <condition>}`, or `{schedule: <setting>, when?: <condition>}`. Each event is a change flai's `inbox` and `wait_for_events` report, or one kind of it, such as a move to done. A schedule names the manifest setting that holds it: `planning.schedule`, `analysis.schedule`, `digest.schedule`, `orchestration.schedule`, or the role's own `roles.<name>.schedule`.

`continuous` is no trigger since [ADR-0209](../../adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md): every role runs for a reason and ends. `flai check` reports it on any role as `roles.file`.

`flai serve` reads each built-in role's triggers from its file, the project's or the baseline flai carries, at every look ([ADR-0250](../../adrs/0250-a-role-file-s-start-prompt-is-a-text-template-in-its-front-matter-s-start-field.md) §4). It starts the role only on the entries the file names, judged as a defined role's are. Each built-in role keeps a launcher of its own, which serves only some entries. `asked` is every role's.

| Role | Entries its launcher serves | Started |
|------|-----------------------------|---------|
| Planner | `{event: edited}`, `{event: accepted}`, `{event: cancelled}`, and `{event: reordered}`, each with or without a condition; `{schedule: <setting>, when: each-ready-story}` | Once per item, through its queue. An event with no condition queues the event's item when the planner takes its kind; `reordered` has none |
| Orchestrator | The events its file names, under their conditions, and `{schedule: orchestration.schedule}` | One run for the reasons queued meanwhile, as [ADR-0209](../../adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md) says |
| Analyzer | `{schedule: <setting>}` | Once each time the setting comes round, for all three focuses, or the formula `analysis.formula` names |
| Curator | `{event: accepted}`; `{schedule: <setting>}` | Once for the acceptances seen meanwhile, and once each time the setting comes round |

`flai check` reports `roles.file` for an entry a built-in role's launcher does not serve, naming those it serves. It also reports a schedule entry without `when: each-ready-story` on a role that takes an item, since such a run has no item to be on. Taking an entry out of a file stops the role's starts on it: a planner's file without `{schedule: planning.schedule, when: each-ready-story}` queues nothing when `planning.schedule` comes round. Replaying the forecasts after an acceptance, a cancellation, or a reorder is flai's own work under `planning.replan`, and no entry gates it.

A project changes a built-in role's start prompt for good by adding a `### Start prompt` section below the marker. An edit of the file's `start` or `trigger` holds only until the next `flai upgrade`, which replaces the front matter with the template's, all but `topics`. A project that wants a built-in role told otherwise, or started otherwise, for good turns the role's action off and defines a role of its own. A role file changes only through a story the operator accepts.

| Event | Happens when |
|-------|--------------|
| `accepted` | A story moved to done |
| `cancelled` | An item was cancelled |
| `ready` | A story moved to ready |
| `review` | A story moved to review |
| `edited` | An item's own words were edited |
| `reordered` | The pull order changed |
| `thread-opened` | A thread was opened |
| `moved` | An item moved to a state |
| `thread` | A thread's newest entry makes it await the reader, or, in the orchestrator's role, the operator |
| `planned` | A planner's run wrote items beyond the item it planned |
| `ready-empty` | The ready count fell below the in-progress room, once per draining |
| `release-due` | An acceptance changed what `flai release --evaluate` answers |

[ADR-0207](../../adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) defines the first seven. [ADR-0211](../../adrs/0211-the-orchestrator-s-trigger-vocabulary-gains-the-events-moved-thread-planned.md) adds `moved`, `thread`, `planned`, `ready-empty`, and `release-due`.

A condition narrows a trigger to the events it holds for. flai implements each one.

| Condition | Qualifies | Holds when |
|-----------|-----------|------------|
| `stale-figures` | `edited` | Someone other than a planner or the story's own agent edited a backlog or ready story's goal, criteria, or touches when its forecast or cost of delay value is older than the edit, or its cost of delay's inputs after its value |
| `replan-agent` | `accepted`, `cancelled`, `reordered` | `planning.replan` is `agent`; for each story whose delivery moved, in the order the replay moved them |
| `each-ready-story` | a schedule | Once for every ready story, in pull order |
| `to-review` | `moved` | The item moved to review |
| `to-backlog` | `moved` | The item moved to backlog; on a call, `item_move` and `flai move`, it allows only a move to backlog |
| `draft-changed` | `edited` | The item edited is a story that is a draft after the edit |
| `awaits-orchestrator` | `thread` | The thread awaits the orchestrator or the operator |

ADR-0211 adds the last four.
