# workflows

What flai runs as data rather than as its own code, one folder per kind:

| Folder | Holds |
|--------|-------|
| [roles/](roles/README.md) | One file per strategic role: the agent `flai serve` starts behind a host action, what starts it, the calls it may make, and its harness definition and prompt. |
| `formulas/` | Not yet: flai reads no formula file today. |

A file here is code the operator runs: a change to one is reviewed as a change to a script is, and reaches main only through a story the operator accepts. No role may edit this folder. [ADR-0207](../adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) records the role file. Convention: [strategic-agents.md](../conventions/strategic-agents.md).
