---
title: Active streams
updated: 2026-10-08T08:38:58Z
---

# Active streams

| Stream | Title | Status | Last agent | Updated |
|--------|-------|--------|------------|---------|
| [S-0232](S-0232.md) | The release key signs flai's checksums.txt in CI and both components carry the public key | in-progress | agent-S-0232 | 2026-10-08T08:33:12Z |
| [S-0291](S-0291.md) | The close-out's install smoke test failed once and passed when run alone, with no cause in its output | review | agent-S-0291 | 2026-10-08T08:37:00Z |
| [S-0321](S-0321.md) | Publishing raises the manifest's flai minimum the moment it commits, before the release is built, so the host's flai drops the project until the binaries exist | in-progress | agent-S-0321 | 2026-10-08T08:12:03Z |
| [S-0322](S-0322.md) | flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit | in-progress | agent-S-0322 | 2026-10-08T08:35:21Z |
| [S-0340](S-0340.md) | The smoke tier installs and self-upgrades from a local release server built from the tree, so no close-out depends on GitHub | in-progress | agent-S-0340 | 2026-10-08T08:31:59Z |
| [S-0341](S-0341.md) | flai verify resumes at the tier that failed when the branch head and its base are unchanged, so a retry re-runs the failure and not the tiers that passed | in-progress | agent-S-0341 | 2026-10-08T08:37:17Z |

## Strategic agents

| Agent | Activities | Cost | Seconds | Last run |
|-------|------------|------|---------|----------|
| [planner](planner.md) | 64 | 177.4324 USD | 21998 | 2026-10-08T08:24:31Z |
| [orchestrator](orchestrator.md) | 274 | 234.7395 USD | 112482 | 2026-10-08T08:37:27Z |
