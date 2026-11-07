---
id: D-0487
title: docs/operators/index.md keeps its heading The kluster host action (S-0577) when a later story widens the section, since docs/operators/settings.md, docs/users/flai.md, and docs/users/flaiover.md link its anchor
kind: decision
sources: [S-0579]
topics: [kluster, documentation, host-actions]
paths: [docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flaiover.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0487 docs/operators/index.md keeps its heading The kluster host action (S-0577) when a later story widens the section, since docs/operators/settings.md, docs/users/flai.md, and docs/users/flaiover.md link its anchor

S-0579 widened the operators guide's section on the `kluster` host action to say what `flai serve` runs while it is on. The section's heading still names S-0577, the story that wrote it first. S-0579 kept that heading as it stands.

From S-0579's narrative, `## Decisions` (`wip/archive/agents/S-0579.md`):

```text
`docs/operators/index.md` keeps its heading `The kluster host action (S-0577)`: renaming it would break links in `docs/operators/settings.md` and `docs/users/flai.md`.
```

The anchor is `index.md#the-kluster-host-action-s-0577`. Three documents link it today: `docs/operators/settings.md`'s `host_actions.<name>` row, `docs/users/flai.md`'s `flai kluster` section, and `docs/users/flaiover.md`'s Features section. A story that changes this section leaves its heading's words alone. One that has to rename it changes every link to the anchor in the same story.

Source: S-0579.
