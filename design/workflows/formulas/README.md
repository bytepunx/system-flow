# formulas

One file per formula, `<name>.md`: a repeatable process as data, which `flai formula run` makes an ordinary story from. flai reads every `*.md` here but this README. [ADR-0238](../../adrs/0238-a-formula-is-a-markdown-file-under-design-workflows-formulas-whose-front-matter.md) defines the format.

A formula's front matter names it (`name`, the file's stem), the nature of the story it makes (`nature`, `improvement` by default), that story's goal (`description`), the values a run takes (`inputs`), its steps (`steps`), and, optionally, when `flai serve` makes its instance (`schedule`). Each step has an `id`, a `title`, what it `run`s, the steps it waits for (`after`), and, optionally, where it runs (`in`): `worktree`, the default, the story's worktree, or `main`, the project's main checkout, for a step such as `flai release --pending` that tags and pushes the branch checked out there; `flai formula run` names that checkout in the step's task. An input is used by its name between double curly braces, as ADR-0238 §4 shows, in the title, a step's title, and a step's run. A `schedule` is a five-field cron expression in UTC, or `daily`; while the `formula` host action is on, `flai serve` makes the instance in ready each time it comes round, with no inputs, so a formula that declares inputs cannot take one. The body is for a reader, with the baseline marker; a project adds its own below it, under `## Project additions`.

A step runs one of three kinds:

| Kind | `run` |
|------|-------|
| A flai command | `{flai: <the command line after flai>}` |
| A script under `scripts/` | `{script: <command line>}` |
| A role a file under [roles/](../roles/README.md) defines | `{role: <name>, prompt: <text>}` |

```bash
flai formula list                                    # every formula: its nature, steps, and file
flai formula show <name>                             # its inputs, its steps, and their layers
flai formula run <name> --input key=value [--epic E-nnnn]  # one backlog story, one task per step
flai formula run <name> --ready                      # the story in ready, so flai serve starts its agent
flai check                                           # formulas.file on a formula that could not run
```

The template ships one baseline formula, [publish.md](publish.md), which the orchestrator starts under its `publish` permission: it fetches, evaluates the release policy, and runs `flai release --pending`, each step in the main checkout ([ADR-0240](../../adrs/0240-under-publish-the-orchestrator-starts-the-baseline-publish-formula-whose-main.md)).

The story `flai formula run` makes carries `formula: <name>` and is ordered, pulled, worked by a story's agent, and accepted as any story. Its agent runs each step itself, in the steps' layers, and closes each task as its step passes; a step that fails blocks its task and asks on a thread. A role file may let its role run a formula: its `commands` list `flai formula run <name>`. A formula file is code the operator runs: it changes only through a story the operator accepts, and no role may edit this folder. Convention: [strategic-agents.md](../../conventions/strategic-agents.md).
