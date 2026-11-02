---
id: D-0459
title: flai serve starts a formula's scheduled instance by running flai formula run --ready --autocommit through os.Executable, since internal/serve cannot import the cmd package that makes the instance
kind: decision
sources: [S-0422, T-1778, T-2012, S-0424]
topics: [cli, go, code, flai]
paths: [flai/internal/serve/formula_schedule.go, flai/cmd/formula.go, flai/internal/serve/mergequeue.go, flai/internal/serve/analysis_schedule.go]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: [S-0424, T-1786]
---

# D-0459 flai serve starts a formula's scheduled instance by running flai formula run --ready --autocommit through os.Executable, since internal/serve cannot import the cmd package that makes the instance

S-0422 decided how `flai serve` makes a formula's instance when its `schedule` comes round (ADR-0239 § 5). The code that makes an instance lives on `app` in `flai/cmd/formula.go`, with the `--ready` flag T-2012 added, and `flai/internal/serve` cannot import `cmd`.

From S-0422's narrative, `## Decisions`:

```text
T-1778 now waits for T-2012 and runs `flai formula run <name> --ready --autocommit` through `os.Executable`, as the merge queue runs `flai accept`. The creation lives on `app` in `cmd/formula.go`, which `internal/serve` cannot import, and T-2012 owns that file and the `--ready` flag.
```

So `formulaScheduler` in `flai/internal/serve/formula_schedule.go` takes `exe: os.Executable` and runs the command, as `flaiAccepter` in `mergequeue.go` and the verify job do. A later change that needs a `cmd` command from `serve` follows the same pattern rather than moving the code.

S-0424 followed it for `analysis.formula`: "`analysis.formula` is started behind the `analyze` host action, the analysis schedule's own, not `formula`; it runs `flai formula run` through `os.Executable` as the formula scheduler does (D-0459)." The analysis scheduler in `flai/internal/serve/analysis_schedule.go` is the second caller.

Sources: S-0422 (wip/archive/agents/S-0422.md), T-1778, T-2012; S-0424 (wip/archive/agents/S-0424.md).
