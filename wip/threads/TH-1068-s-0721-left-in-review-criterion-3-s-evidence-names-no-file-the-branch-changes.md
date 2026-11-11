---
id: TH-1068
title: "S-0721 left in review: criterion 3's evidence names no file the branch changes"
anchor:
  path: wip/kanban/stories/S-0721-every-mcp-server-flai-builds-resolves-all-its-tools-json-schemas-again-which-is.md
  item: S-0721
status: resolved
participants: [orchestrator, alex]
created: 2026-11-11T05:47:15Z
updated: 2026-11-11T05:47:58Z
---

# TH-1068 S-0721 left in review: criterion 3's evidence names no file the branch changes

On wip/kanban/stories/S-0721-every-mcp-server-flai-builds-resolves-all-its-tools-json-schemas-again-which-is.md.

## Entries

### 2026-11-11T05:47:15Z orchestrator
**Recommended:** accept S-0721 yourself with `flai accept S-0721 --by alex`. Every criterion is met at 5d9c4a81a8f2.

**Why I did not accept it:** `accept` refused my evidence:

```text
the evidence for acceptance criterion 3 (...) names no file the branch changes (it names wip/agents/S-0721.md), so it cannot be checked against the diff
```

Criterion 3 asks for the CPU figures in the story's narrative. The narrative is `wip/agents/S-0721.md` in the main checkout, and no story branch changes it. So no changed file can meet this criterion. I named none in its place.

**What the verifier found:**

| Criterion | Met by |
|---|---|
| 1 | `flai/internal/mcpserver/server.go`, `flai/internal/mcpserver/folder.go` |
| 2 | `flai/internal/mcpserver/schema_cache_test.go` |
| 3 | `wip/agents/S-0721.md` (main checkout): `-race` user CPU went from 544–566 s to 105–106 s, 19% of before. The plain run went from 47 s to 24 s. |
| 4 | `design/tech/go-libraries.md`; `go.mod` and `go.sum` unchanged |

- The close-out passed every step at the branch head, 5d9c4a81a8f2.
- The dry run listed no blocker.
- The narrative measured its own before figure, 544–566 s, not the Notes' 608 s. The result is under half of either.

**Alternative:** a later story could let flai take, as evidence, a narrative path for a criterion about the narrative. Until then, a criterion of that kind stops my acceptance.

### 2026-11-11T05:47:58Z alex
Resolved: S-0721 was accepted
