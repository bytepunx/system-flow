---
id: D-0283
title: An empty narrative file under wip/agents makes every flai stream log in the project fail
kind: cause
sources: [I-0129]
topics: [narratives, flai, cli, go, code, all]
paths: [flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0129-an-empty-narrative-file-under-wip-agents-makes-every-flai-stream-log-in-the-project-fail.md, design/issues/summary.md, flai/internal/workitem/activity_test.go, design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md]
date: 2026-10-12
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0283 An empty narrative file under wip/agents makes every flai stream log in the project fail

I-0129 was closed on 2026-10-12.

The cause:

```text
An empty narrative file under wip/agents makes every flai stream log in the project fail
```

The fix:

```text
S-0373: Repo.ActiveNarratives no longer returns at the first S-*.md that does not parse; WriteIndex lists such a narrative in wip/agents/index.md as unreadable, with its story's status, and writes the index, so an empty narrative fails no other story's stream open, log, state, move, edit, archive, accept, story start, or task done. flai check keeps narrative.front-matter as its one finding. TestIndexListsAnUnreadableNarrative and TestUnreadableNarrativeIsOneFinding reproduce it.
```

Source: I-0129.
