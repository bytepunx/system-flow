# roles

One file per strategic role, `<name>.md`: an agent `flai serve` starts as a session of its own behind a host action. flai reads every `*.md` here but this README. [ADR-0207](../../adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) defines the format.

A role file's front matter names the role (`name`, `role`), the harness definition flai generates for it (`definition`: `description`, `tools`, `model`), the host action that gates it (`action`), what starts it (`trigger`), the items a run takes (`item`), its `topics`, how its usage is charged (`charge`), and the calls the guard allows it: flai's MCP tools (`tools`), flai's commands (`commands`), the folders it may edit (`files`), and the calls only an operator's setting allows (`permissions`). A call is named plainly or under a condition flai implements. The body is the definition's prompt, with the baseline marker; a project adds its own below it, under `## Project additions`, and keeps it through `flai upgrade`, as it keeps its `topics`.

The template ships four baseline roles. A role file the template does not ship is the project's own. `flai role list` prints every role; `flai check` reports `roles.file`, `roles.grant`, and `roles.definition`. Convention: [strategic-agents.md](../../conventions/strategic-agents.md).

| File | Role | Action | Trigger |
|------|------|--------|---------|
| [planner.md](planner.md) | `plan` | `plan` | Asked; an edit that makes a story's figures stale; an acceptance, cancellation, or reorder while `planning.replan` is `agent`; `planning.schedule`, for every ready story |
| [orchestrator.md](orchestrator.md) | `orchestrate` | `orchestrate` | Continuous while the action is on |
| [analyzer.md](analyzer.md) | `analyze` | `analyze` | Asked; `analysis.schedule` |
| [curator.md](curator.md) | `curate` | `curate` | Asked; an acceptance; `digest.schedule` |
