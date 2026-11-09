---
id: D-0538
title: flai guard's HostAPIReads leaves out the host API methods that mark, pull, or start something or were not checked, so a task sub-agent is refused them until a story checks one's code and moves it to the list
kind: decision
sources: [S-0612, I-0307]
topics: [flai, cli, go, code]
paths: [flai/internal/guard/guard.go, flai/internal/hostapi/guardreads_test.go, design/system/flai-cli.md]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0538 flai guard's HostAPIReads leaves out the host API methods that mark, pull, or start something or were not checked, so a task sub-agent is refused them until a story checks one's code and moves it to the list

S-0612 (I-0307) let every agent run `flai hostapi` alone and with a method on `guard.HostAPIReads`, `--timing` included; every other method is refused by name as a write. The rule is in `design/system/flai-cli.md`; which methods stayed off the list, and why, is in the story's narrative alone:

```text
T-2702: `guard.HostAPIReads` lists only methods that write nothing beyond the caches a guard read such as `flai board` already keeps. Borderline methods stay off it and so are refused: `components.status`, which marks dev servers used; `dashboard.check`, which pulls an image; the host-level and `checks.*` methods; `file.get`; the `*.template` methods; and `kluster.deployments`, `kluster.workloads`, and `kluster.log`. A later story can move one to the list once its code is checked.
```

```text
T-2703: `settings.shared_check` reads, as `flai shared check` does, but it stays on `hostAPIWrites`, with a comment saying why, so that the guard's list stays T-2702's reviewed one.
```

So a task sub-agent that needs one of these methods, for a measurement such as `flai hostapi --timing`, is refused; the story's agent runs it, or a story checks that the method writes nothing and adds it to `HostAPIReads`, which `flai/internal/hostapi/guardreads_test.go` keeps in step with the method table.

Source: S-0612's narrative, `wip/archive/agents/S-0612.md` § Decisions; I-0307.
