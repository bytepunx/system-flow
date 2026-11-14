# workflows

What flai runs as data rather than as its own code, one folder per kind:

| Folder | Holds |
|--------|-------|
| [roles/](roles/README.md) | One file per strategic role: the agent `flai serve` starts behind a host action, what starts it, the calls it may make, and its harness definition and prompt. |
| `board.md` | The board: its states, each with a kind, the transitions between them with who may make each, the rules each passes and the effects flai gives it, and the timers measured over every item's transitions. A project without one runs the baseline flai carries, which is today's board. No project has one yet: flai reads it once S-0755 lands, and until S-0737 flai's own code is the baseline. |
| [formulas/](formulas/README.md) | One file per formula: a repeatable process whose steps each run a flai command, a script, or a role, from which `flai formula run` makes an ordinary story with one task per step. |

A file here is code the operator runs: a change to one is reviewed as a change to a script is, and reaches main only through a story the operator accepts. No role may edit this folder. [ADR-0207](../adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) records the role file, [ADR-0238](../adrs/0238-a-formula-is-a-markdown-file-under-design-workflows-formulas-whose-front-matter.md) the formula file, [ADR-0301](../adrs/0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md) the workflow file, `board.md`, which only the operator accepts a change to. Convention: [strategic-agents.md](../conventions/strategic-agents.md).
