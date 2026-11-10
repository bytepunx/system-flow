---
id: D-0480
title: A story whose work reads or runs repositories outside its worktree is worked by the story's agent itself, with its files under the worktree's git-ignored bin/ and its waits as Bash loops, since sub-agents, Write, and Monitor are refused every path outside the worktree
kind: decision
sources: [S-0575, I-0265]
topics: [kluster, testing, deployment, delegation]
paths: [design/issues/I-0265-a-story-s-sub-agents-write-and-monitor-are-refused-every-path-outside-its.md, design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md]
date: 2026-11-06
status: active
absorbed_by: ""
asked_again: []
cited: [S-0701, T-3048, T-3064, TH-1025]
---

# D-0480 A story whose work reads or runs repositories outside its worktree is worked by the story's agent itself, with its files under the worktree's git-ignored bin/ and its waits as Bash loops, since sub-agents, Write, and Monitor are refused every path outside the worktree

S-0575, an experiment that read `~/git/bytepunx/authstar`, `signet`, and `kluster` and ran their scripts against a local cluster, met this under `flai serve`. Its explorer was refused Read, Grep, and Glob on every path outside the worktree, and the story's agent's own Write to `/tmp/S-0575` and its Monitor call were refused by `permission_prompt`. Bash reads and writes outside the worktree passed.

From its Decisions (`wip/archive/agents/S-0575.md`):

```text
I work the tasks myself rather than through task sub-agents. The explorer started for the story was refused every read outside the worktree, and the experiment reads and runs `~/git/bytepunx/authstar`, `signet`, and `kluster`, so a task sub-agent could do none of it.
```

```text
Files the experiment writes go under the worktree's git-ignored `bin/s-0575/`, not `/tmp`: the Write tool is refused outside the worktree. Logs, the kubeconfig, and the exported trees stay in `/tmp/S-0575`.
```

For a later story that works outside this repository, such as the kluster stories of E-0032 (S-0576 to S-0581): plan its tasks as the story's agent's own, keep the files it writes under the worktree's git-ignored `bin/<story>/`, and wait with Bash loops rather than Monitor. I-0265 records the friction, about 20 minutes on S-0575, and its remedy, allowing a story's named outside paths or warning in the start prompt, is open.

Sources: S-0575, I-0265.
