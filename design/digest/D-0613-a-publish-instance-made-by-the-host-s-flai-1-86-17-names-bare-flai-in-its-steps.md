---
id: D-0613
title: A publish instance made by the host's flai 1.86.17 names bare flai in its steps, not scripts/flai.sh, and S-0774 and S-0776 ran them as written with the installed 1.86.17, which has every flag they use
kind: decision
sources: [S-0776, S-0774, S-0770]
topics: [formulas, release, flai]
paths: [design/workflows/formulas/publish.md, scripts/flai.sh]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0613 A publish instance made by the host's flai 1.86.17 names bare flai in its steps, not scripts/flai.sh, and S-0774 and S-0776 ran them as written with the installed 1.86.17, which has every flag they use

`design/system/strategic-agents.md` § Working an instance says that where the manifest names `flai.command`, a flai step's task names that command in place of `flai` (S-0681, I-0332). This repository's manifest names `scripts/flai.sh`. The publish instances made on 2026-11-13 do not: their tasks say `flai release --evaluate --require-met --judged ""` and `flai release --pending`, bare, as the instances S-0770 to S-0776 show. The host's installed flai, which made them, is 1.86.17, behind 1.94 in the project's history.

The agents split on it. S-0770 ran each step as `scripts/flai.sh`, as D-0559 did. S-0774 and S-0776 ran them as written, with the installed flai in the main checkout. From S-0776's Decisions (`wip/archive/agents/S-0776.md`):

```text
Layers are the formula's: T-3316, then T-3317, then T-3318, each waiting on the one before; worked by the story's agent directly, no sub-agents, with the installed flai in the main checkout as S-0774 did.
```

Both passed. T-3317's note: "exited 0 in the main checkout (flai 1.86.17): threshold met". T-3318's note: "exited 0 in the main checkout (flai 1.86.17), `pushed: true`", tagging `flai/v1.94.1` and `flaiover/v0.60.0`. S-0774 tagged `flai/v1.93.1` and `flaiover/v0.59.0` the same way. So 1.86.17 has every flag these steps use, `--judged` included, which 1.86.6 lacked (I-0332).

Running a step as written works while the installed flai has its flags. An `unknown flag` exit is the wrong binary's, not the step's: run the step again as `scripts/flai.sh`, as D-0559 says, rather than block the task. The bare `flai` in the tasks lasts until the host's flai is upgraded past S-0681 (`flai host upgrade`, the operator's).

Sources: S-0776 (`wip/archive/agents/S-0776.md`, Decisions; T-3317, T-3318 Notes); S-0774 (`wip/archive/agents/S-0774.md`); S-0770 (`wip/archive/agents/S-0770.md`, Decisions).
