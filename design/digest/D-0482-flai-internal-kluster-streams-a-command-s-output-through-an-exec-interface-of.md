---
id: D-0482
title: flai/internal/kluster streams a command's output through an Exec interface of its own, because execx.Runner returns a command's combined output only when it ends
kind: decision
sources: [S-0577, T-2474, S-0578]
topics: [kluster, flai, go, code, testing]
paths: [flai/internal/kluster/kluster.go, flai/internal/kluster/deploy.go, flai/internal/execx/execx.go]
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

S-0578 built `deploy`, `teardown`, `promote`, and `deployments` in the same package and kept that `Exec`, widened to pass the environment and standard input the steps need. The code-quality addition names `execx.Runner` for git and docker, but a second runner in one package would split its tests, which fake `kluster`, `kubectl`, and the rest on `PATH`. From S-0578's narrative, `## Decisions`:

```text
The package keeps S-0577's `kluster.Exec` and its fakes-on-PATH tests rather than `execx.Runner`, widened to pass the environment and standard input the steps need; the code-quality addition names `execx.Runner` for git and docker, but S-0577's package already chose its own streaming `Exec`, and a second runner in one package would split the fakes.
```

A later change in `flai/internal/kluster`, such as S-0579's jobs in flai serve, runs its commands through that `Exec`, not `execx.Runner`. D-0389 records the other limit of `execx.Runner`: it trims the output it returns.

Sources: S-0577 (wip/archive/agents/S-0577.md), T-2474, S-0578 (wip/archive/agents/S-0578.md).
