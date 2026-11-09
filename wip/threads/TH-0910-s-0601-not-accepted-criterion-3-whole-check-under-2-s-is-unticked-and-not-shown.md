---
id: TH-0910
title: "S-0601 not accepted: criterion 3 (whole check under 2 s) is unticked and not shown"
anchor:
  path: wip/kanban/stories/S-0601-flai-check-runs-its-rules-at-once-and-keeps-the-component-props-node-read-so-a.md
  item: S-0601
status: resolved
participants: [orchestrator, alex, agent-S-0601]
created: 2026-11-08T02:15:36Z
updated: 2026-11-09T08:23:36Z
---

# TH-0910 S-0601 not accepted: criterion 3 (whole check under 2 s) is unticked and not shown

On wip/kanban/stories/S-0601-flai-check-runs-its-rules-at-once-and-keeps-the-component-props-node-read-so-a.md.

## Entries

### 2026-11-08T02:15:36Z orchestrator
Recommended: reword criterion 3 to what the story showed, or accept it as it stands with criterion 3 left to I-0312. If you reword it, I accept on my next run.

Why it is not accepted:

- `flai accept S-0601 --by orchestrator --verified f44b468e26aa --dry-run` is blocked: acceptance criterion 3 of S-0601 is unticked.
- Criterion 3 asks for a warm whole `flai check --strict` under 2 s. design/system/server-performance.md § Cause 11 records the before and after, but neither is under 2 s.

| Warm `flai check --strict`, five runs | Wall | CPU |
|---|---|---|
| main 3087d9cf4 | 7.9 to 11.7 s | 4.9 to 5.3 s |
| S-0601 1582a0dd3 | 3.3 to 4.9 s | 3.8 to 3.9 s |

The host load average was 84 to 103 on 24 cores. The remaining critical path is the item list, `workItems`, and the digest rule, recorded as I-0312.

What is met, verified at f44b468, where the close-out passed every step:

- Criterion 1: `check.Run` runs the rules at once and merges them in the old order. The test is `TestRunAtOnceMatchesTheRulesOneAfterAnother` in flai/internal/check/concurrent_test.go. It uses no lock, as the story's notes allow.
- Criterion 2: the kept props reads are in flai/internal/components/propscache.go, keyed by path, size, mtime, and the svelte version. The test is `TestReadPropsStartsNodeOnlyForAComponentWhoseKeyChanged` in propscache_test.go.

Options:

1. Reword criterion 3, for example: "halves a warm whole `flai check --strict` against main on the same host, and server-performance.md records the before and after". The agent or you then tick it, and I accept.
2. Keep criterion 3 as written. The story stays in review until someone measures it on a quieter host.
3. Accept it by hand with criterion 3 unticked, and leave I-0312 to a follow-up story.

### 2026-11-09T08:19:31Z alex
I reworded criterion 3 and ticked it

### 2026-11-09T08:19:39Z agent-S-0601
Noted. Criterion 3 now asks S-0601 to halve a warm whole check against main on the same host. The measurements meet it: 3.3 to 4.9 s against main's 7.9 to 11.7 s, the two binaries taking turns. S-0601's Notes and narrative now say so. The branch is rebased onto main for TH-0919, and its close-out is running. S-0601 goes back to review when the close-out passes.

### 2026-11-09T08:23:36Z orchestrator
Resolved: The operator reworded criterion 3 to halving a warm whole check against main on the same host, and ticked it. server-performance.md records 3.3 to 4.9 s against main's 7.9 to 11.7 s, and the verifier matched it at aaebd6f.
