---
id: D-0475
title: A criterion that needs a run this host cannot make, on macOS or another git version, is reworded to what the close-out here verifies, and the flai workflow's macOS and ubuntu jobs on the acceptance commit make the run itself
kind: decision
sources: [S-0572]
topics: [testing, ci, planning, review]
paths: [".github/workflows/flai.yml"]
date: 2026-11-03
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0475 A criterion that needs a run this host cannot make, on macOS or another git version, is reworded to what the close-out here verifies, and the flai workflow's macOS and ubuntu jobs on the acceptance commit make the run itself

S-0572's criteria 1 and 2 first promised a pass on macOS and on git 2.56. No host here can run those: this host has git 2.47.3 and no macOS. The operator decided, in conversation, logged in S-0572's Decisions (`wip/archive/agents/S-0572.md`) at 2026-11-03T00:16:27Z:

```text
criteria 1 and 2 are reworded to what the close-out on this host verifies (the shape of the fix command and of the assertion, and a pass here), and ticked; the macOS and git 2.56 runs stay the `flai` workflow's on the acceptance commit.
```

S-0572's Notes say the same. A story whose fix is for another platform or tool version words its criteria as the shape of the fix and a pass on this host. It leaves the run on that platform to CI after acceptance. The story does not wait in review for that run. This is the CI case of the operator's answers on live runs that cannot happen before acceptance, which moved those runs to a follow-up story (D-0461, D-0466).

Source: S-0572.
