---
id: TH-0372
title: "S-0340 criterion 6: keep S-0291 for its retries, or cancel it?"
anchor:
  path: wip/kanban/stories/S-0340-the-smoke-tier-installs-and-self-upgrades-from-a-local-release-server-built-from-the-tree-so-no-close-out-depends-on-github.md
  item: S-0340
status: resolved
participants: [agent-S-0340, orchestrator, alex]
created: 2026-10-08T08:16:13Z
updated: 2026-10-08T08:22:10Z
---

# TH-0372 S-0340 criterion 6: keep S-0291 for its retries, or cancel it?

On wip/kanban/stories/S-0340-the-smoke-tier-installs-and-self-upgrades-from-a-local-release-server-built-from-the-tree-so-no-close-out-depends-on-github.md.

## Entries

### 2026-10-08T08:16:13Z agent-S-0340
Criterion 6 of S-0340 asks you whether S-0291 is cancelled or closed as taken over.

**Recommended: keep S-0291 in progress for its retries, minus T-1340.** S-0340 takes the smoke tier off GitHub, but users' `install.sh` and `flai self-upgrade` still reach GitHub and still drop the 1 MB listing. S-0291's retries (T-1338, T-1339) fix that, and they also steady S-0340's new `scripts/install-published-test.sh`. S-0291's T-1340 (print `install.sh`'s output in `install-test.sh`) is done by S-0340's rewrite. S-0340 closes I-0086. I proposed this split to S-0291's agent on MS-0013.

**Alternative: cancel S-0291** as taken over by S-0340. Its retries are then lost unless a new story picks them up.

Until you answer, I work everything else in S-0340. I will record your answer in the narrative and act on it in T-1357.

### 2026-10-08T08:16:19Z orchestrator (recommendation)
Recommendation for the operator: keep S-0291 for its retries (T-1338, T-1339, T-1341) and drop T-1340, as the agent recommends. I recommended the same on TH-0370. Users' `install.sh` and `flai self-upgrade` still reach GitHub, and S-0340's new published-install check would flake without the retries. S-0340 closes I-0086. Criterion 6 leaves this to you. One answer here also settles TH-0370's question 2.

Source: wip/kanban/stories/S-0340-the-smoke-tier-installs-and-self-upgrades-from-a-local-release-server-built-from-the-tree-so-no-close-out-depends-on-github.md

### 2026-10-08T08:18:49Z alex
Confirmed the recommendation of 2026-10-08T08:16:19Z orchestrator.

Source: wip/kanban/stories/S-0340-the-smoke-tier-installs-and-self-upgrades-from-a-local-release-server-built-from-the-tree-so-no-close-out-depends-on-github.md

### 2026-10-08T08:22:10Z alex
Resolved.
