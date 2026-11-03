---
id: D-0471
title: The orchestrator's acceptance step still says flai leaves the main checkout's uncommitted files out of the acceptance commit, which ADR-0248 made untrue for the shared set, and S-0566 left the sentence for a follow-up
kind: decision
sources: [S-0566]
topics: [orchestration, review, conventions]
paths: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, flai/internal/roles/baseline/orchestrator.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, design/adrs/0248-every-acceptance-commits-the-main-checkout-s-uncommitted-files-under-wip-design.md]
date: 2026-11-03
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0471 The orchestrator's acceptance step still says flai leaves the main checkout's uncommitted files out of the acceptance commit, which ADR-0248 made untrue for the shared set, and S-0566 left the sentence for a follow-up

S-0566 made every acceptance commit the main checkout's uncommitted files under `wip/`, `design/`, and `claims.shared`, the shared set ([ADR-0248](../adrs/0248-every-acceptance-commits-the-main-checkout-s-uncommitted-files-under-wip-design.md)). It updated the orchestrator's harness prompt in Go but not the words of its role file. From S-0566's Decisions (`wip/archive/agents/S-0566.md`):

```text
The baseline role file `flai/internal/roles/baseline/orchestrator.md` (and its copies under `design/workflows/roles/` and `template/root/`) still says, in the orchestrator's acceptance step, that flai "leaves them out of the acceptance" of the main checkout's uncommitted files; true for the rest, not the shared set. Left alone: S-0554 and S-0425 touched those files, the sentence instructs no operator, and the harness prompt in Go was updated. Open for a follow-up.
```

The convention says the same today: `design/conventions/strategic-agents.md` § As the orchestrator › Accepting a story, step 1, reads "flai leaves them out of the lint fix and the acceptance commit". Where those words and ADR-0248 differ, ADR-0248 holds: the shared set is committed with the acceptance, or in a commit of its own before the merge when the story's change would overwrite one of its files, and only the files outside it are left as they are. A story that changes the orchestrator's role file or the convention carries the fix.

Source: S-0566.
