---
id: TH-0810
title: "S-0424: the two live runs (criterion 4, T-1787, T-1788) cannot run before S-0424 is accepted; move them to a follow-up story?"
anchor:
  path: wip/kanban/stories/S-0424-the-template-ships-baseline-formulas-for-importing-a-repository-from-the-board-and-for-the-analyzer-s-cadence-the-board-s-import-starts-its-formula-and-each-is-run-once-here.md
  item: S-0424
status: open
participants: [agent-S-0424]
created: 2026-11-02T05:29:51Z
updated: 2026-11-02T05:29:51Z
---

# TH-0810 S-0424: the two live runs (criterion 4, T-1787, T-1788) cannot run before S-0424 is accepted; move them to a follow-up story?

On wip/kanban/stories/S-0424-the-template-ships-baseline-formulas-for-importing-a-repository-from-the-board-and-for-the-analyzer-s-cadence-the-board-s-import-starts-its-formula-and-each-is-run-once-here.md.

## Entries

### 2026-11-02T05:29:51Z agent-S-0424
Criterion 4 asks for each formula to be run once here, its instance accepted, before S-0424 goes to review. Neither run can start while S-0424 is open:

- `design/workflows/formulas/import.md` and `analysis-cadence.md` exist only on `story/S-0424`. `flai formula run` reads the main checkout's folder.
- `flai import --formula` (T-2017) and `flai issue list --report` (T-1786) exist only on the story's branch. The board's Import runs the installed flai.

**Recommended:** option 1, as on TH-0803 (D-0461) and D-0422.

- I cancel T-1787 and T-1788 and remove criterion 4 from S-0424.
- I write a draft story in E-0024 with `after: S-0424`. It holds both tasks' Work and criterion 4's words, and its Notes say why it waits.
- S-0424's Notes record the move. The epic still gets both live runs before it is accepted.

**Alternative (option 2):** keep criterion 4. S-0424 then cannot meet the definition of done and stays in progress until it is accepted some other way.

I go on with T-2016, T-2017, T-1785, T-1786, and T-1789 meanwhile.
