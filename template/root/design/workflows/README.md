# workflows

What flai runs as data rather than as its own code, one folder or file per kind:

| Path | Holds |
|--------|-------|
| [roles/](roles/README.md) | One file per strategic role: the agent `flai serve` starts behind a host action, what starts it, the calls it may make, and its harness definition and prompt. |
| [board.md](board.md) | The board: its states, each with a kind, the moves between them with who may make each, the rules each passes and the effects flai gives it, and the timers measured over every item's history. The template ships the baseline flai carries; a project without the file runs that baseline. flai reads the file of the checkout it runs in, a story's worktree its branch's and the main checkout main's, and refuses to run on one it cannot, save for `flai check`, which says what to fix. `flai upgrade` adds it when it is missing and never changes it after. |
| [formulas/](formulas/README.md) | One file per formula: a repeatable process whose steps each run a flai command, a script, or a role, from which `flai formula run` makes an ordinary story with one task per step. |

A file here is code the operator runs: a change to one is reviewed as a change to a script is, and reaches main only through a story the operator accepts. No role may edit this folder. [ADR-0207](../adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) records the role file, [ADR-0238](../adrs/0238-a-formula-is-a-markdown-file-under-design-workflows-formulas-whose-front-matter.md) the formula file. A story that changes `board.md` is accepted by the operator alone, never by the orchestrator or the merge queue. Convention: [strategic-agents.md](../conventions/strategic-agents.md).
