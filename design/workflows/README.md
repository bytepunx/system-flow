# workflows

What flai runs as data rather than as its own code, one folder per kind:

| Folder | Holds |
|--------|-------|
| [roles/](roles/README.md) | One file per strategic role: the agent `flai serve` starts behind a host action, what starts it, the calls it may make, and its harness definition and prompt. |
| [formulas/](formulas/README.md) | One file per formula: a repeatable process whose steps each run a flai command, a script, or a role, from which `flai formula run` makes an ordinary story with one task per step. |

A file here is code the operator runs: a change to one is reviewed as a change to a script is, and reaches main only through a story the operator accepts. No role may edit this folder. [ADR-0207](../adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) records the role file, [ADR-0238](../adrs/0238-a-formula-is-a-markdown-file-under-design-workflows-formulas-whose-front-matter.md) the formula file. Convention: [strategic-agents.md](../conventions/strategic-agents.md).
