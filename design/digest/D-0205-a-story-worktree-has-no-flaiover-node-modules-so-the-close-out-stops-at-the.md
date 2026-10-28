---
id: D-0205
title: A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there
kind: cause
sources: [I-0080]
topics: [all]
paths: [scripts/flaiover-install.sh, scripts/flaiover-test.sh, scripts/flaiover-unit.sh, scripts/README.md, design/system/devex.md, docs/operators/index.md, design/issues/I-0080-a-story-worktree-has-no-flaiover-node-modules-so-the-close-out-stops-at-the-flaiover-step-until-flaiover-install-runs-there.md, design/issues/summary.md]
date: 2026-10-07
status: absorbed
absorbed_by: "design/system/devex.md#used-by-this-monorepo"
asked_again: []
cited: []
---

# D-0205 A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there

I-0080 was closed on 2026-10-07.

The cause:

```text
A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there
```

The fix:

```text
S-0281: scripts/flaiover-install.sh --if-needed installs flaiover's dependencies when node_modules is missing or older than pnpm-lock.yaml; flaiover-test.sh, close-out's flaiover step, always calls it and flaiover-unit.sh calls it in a story worktree, so a worktree made by git worktree add no longer stops at prettier: not found or skips vitest
```

Source: I-0080.
