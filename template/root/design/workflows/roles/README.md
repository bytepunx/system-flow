# roles

One file per strategic role, `<name>.md`: an agent `flai serve` starts as a session of its own behind a host action. flai reads every `*.md` here but this README. [ADR-0207](../../adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) defines the format.

A role file's front matter names the role (`name`, `role`), the harness definition flai generates for it (`definition`: `description`, `tools`, `model`), the host action that gates it (`action`), what starts it (`trigger`), the items a run takes (`item`), its `topics`, how its usage is charged (`charge`), and the calls the guard allows it: flai's MCP tools (`tools`), flai's commands (`commands`), the folders it may edit (`files`), and the calls only an operator's setting allows (`permissions`). A call is named plainly or under a condition flai implements. The body is the definition's prompt, with the baseline marker; a project adds its own below it, under `## Project additions`, and keeps it through `flai upgrade`, as it keeps its `topics`.

The template ships four baseline roles. A role file the template does not ship is the project's own. `flai role list` prints every role; `flai check` reports `roles.file`, `roles.grant`, and `roles.definition`. Convention: [strategic-agents.md](../../conventions/strategic-agents.md).

| File | Role | Action | Trigger |
|------|------|--------|---------|
| [planner.md](planner.md) | `plan` | `plan` | Asked; an edit that makes a story's figures stale; an acceptance, cancellation, or reorder while `planning.replan` is `agent`; `planning.schedule`, for every ready story |
| [orchestrator.md](orchestrator.md) | `orchestrate` | `orchestrate` | Asked; a thread that awaits it or the operator; a move to review or to backlog; an edit of a draft; a planner's run that wrote items; `release-due`; `ready-empty`; a reorder; `orchestration.schedule`. Each run acts on its reasons and ends |
| [analyzer.md](analyzer.md) | `analyze` | `analyze` | Asked; `analysis.schedule` |
| [curator.md](curator.md) | `curate` | `curate` | Asked; an acceptance; `digest.schedule` |

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
