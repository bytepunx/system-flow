# roles

One file per strategic role, `<name>.md`: an agent `flai serve` starts as a session of its own behind a host action. flai reads every `*.md` here but this README. [ADR-0207](../../adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) defines the format.

A role file's front matter names the role (`name`, `role`), the harness definition flai generates for it (`definition`: `description`, `tools`, `model`), the host action that gates it (`action`), what starts it (`trigger`), the items a run takes (`item`), its `topics`, how its usage is charged (`charge`), and the calls the guard allows it: flai's MCP tools (`tools`), flai's commands (`commands`), the folders it may edit (`files`), and the calls only an operator's setting allows (`permissions`). A call is named plainly or under a condition flai implements. The body is the definition's prompt, with the baseline marker; a project adds its own below it, under `## Project additions`, and keeps it through `flai upgrade`, as it keeps its `topics`.

The template ships four baseline roles. A role file the template does not ship is the project's own. `flai role list` prints every role; `flai check` reports `roles.file`, `roles.grant`, and `roles.definition`, and advises `roles.baseline` for each of the four the project lacks, whose baseline flai runs in its place. Convention: [strategic-agents.md](../../conventions/strategic-agents.md).

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

## Triggers

A `trigger` entry is `asked`, `{event: <event>, when?: <condition>}`, or `{schedule: <setting>, when?: <condition>}`. Each event is a change flai's `inbox` and `wait_for_events` report, or one kind of it, such as a move to done. A schedule names the manifest setting that holds it: `planning.schedule`, `analysis.schedule`, `digest.schedule`, `orchestration.schedule`, or the role's own `roles.<name>.schedule`.

`continuous` is no trigger since [ADR-0209](../../adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md): every role runs for a reason and ends. `flai check` reports it on any role as `roles.file`.

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
| `stale-figures` | `edited` | Someone other than a planner edited a backlog or ready story's goal, criteria, or touches, and its forecast or cost of delay value is older than the edit |
| `replan-agent` | `accepted`, `cancelled`, `reordered` | `planning.replan` is `agent`; for each story whose delivery moved |
| `each-ready-story` | a schedule | Once for every ready story |
| `to-review` | `moved` | The item moved to review |
| `to-backlog` | `moved` | The item moved to backlog; on a call, `item_move` and `flai move`, it allows only a move to backlog |
| `draft-changed` | `edited` | The item edited is a story that is a draft after the edit |
| `awaits-orchestrator` | `thread` | The thread awaits the orchestrator or the operator |

ADR-0211 adds the last four.
