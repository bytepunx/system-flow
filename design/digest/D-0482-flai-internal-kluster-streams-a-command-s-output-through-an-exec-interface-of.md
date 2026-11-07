---
id: D-0482
title: flai/internal/kluster streams a command's output through an Exec interface of its own, because execx.Runner returns a command's combined output only when it ends
kind: decision
sources: [S-0577, T-2474]
topics: [kluster, flai, go, code]
paths: [flai/internal/kluster/kluster.go, flai/internal/execx/execx.go]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0482 flai/internal/kluster streams a command's output through an Exec interface of its own, because execx.Runner returns a command's combined output only when it ends

S-0577 built `flai kluster up`, `down`, and `status` on a package of its own, `flai/internal/kluster`. A `kluster up` takes minutes, and `flai kluster up` streams its output as it comes (`design/system/kluster.md` § Commands (S-0577, S-0578)).

`execx.Runner` in `flai/internal/execx/execx.go` cannot do that. Its `Run` and `RunInput` return a command's combined stdout and stderr as one string once the command has ended. So the kluster package runs commands through an `Exec` interface it defines, which streams.

From S-0577's narrative, `## Decisions`:

```text
`flai/internal/kluster` streams through an `Exec` interface of its own, because `execx.Runner` returns output only at the end and a bring-up takes minutes (T-2474).
```

A later command that must show a long run's output as it comes, such as S-0578's `flai kluster deploy`, uses that `Exec` or a streaming runner like it, not `execx.Runner`. D-0389 records the other limit of `execx.Runner`: it trims the output it returns.

Sources: S-0577 (wip/archive/agents/S-0577.md), T-2474.
