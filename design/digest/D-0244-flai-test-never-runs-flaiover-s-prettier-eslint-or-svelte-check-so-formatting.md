---
id: D-0244
title: flai test never runs flaiover's prettier, eslint, or svelte-check, so formatting faults surface only at the close-out
kind: cause
sources: [I-0105, I-0206]
topics: [testing, flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [scripts/flaiover-lint.sh, scripts/flaiover-check.sh, scripts/README.md, system-flow.yaml, flai/internal/verify/select_test.go, design/system/devex.md, design/conventions/tooling.md, design/issues/I-0105-flai-test-never-runs-flaiover-s-prettier-eslint-or-svelte-check-so-formatting-faults-surface-only-at-the-close-out.md, design/issues/I-0206-flai-test-on-flaiover-paths-runs-no-svelte-check-so-a-type-error-a-task-s-tests-passed-over-first-shows-in-the-close-out.md, design/issues/summary.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md]
date: 2026-10-08
status: absorbed
absorbed_by: "design/conventions/tooling.md#project-additions"
asked_again: []
cited: [S-0530, TH-0702]
---

# D-0244 flai test never runs flaiover's prettier, eslint, or svelte-check, so formatting faults surface only at the close-out

I-0105 was closed on 2026-10-08.

The cause:

```text
flai test never runs flaiover's prettier, eslint, or svelte-check, so formatting faults surface only at the close-out
```

The fix:

```text
S-0319: system-flow.yaml's tests gain a `flaiover-lint` tier, outside `--all` and before vitest, that runs `scripts/flaiover-lint.sh` (prettier `--check --ignore-unknown` and eslint `--no-warn-ignored`) on the changed files under `flaiover/`, so `flai test` finds a formatting or lint fault in a flaiover file before the close-out. svelte-check takes no file list and stays in the `--all` `flaiover` tier. `TestRepositoryManifestLintsChangedFlaioverFiles` in `flai/internal/verify/select_test.go` fails without the tier.
```

The svelte-check part stayed open as I-0206, closed on 2026-10-28:

```text
flai test on flaiover paths runs no svelte-check, so a type error a task's tests passed over first shows in the close-out
```

Its fix:

```text
S-0530: system-flow.yaml's tests gain a svelte-check tier, outside --all and after vitest, that runs scripts/flaiover-check.sh (pnpm check --incremental, svelte-kit sync then svelte-check on the whole flaiover project) when a changed flaiover/src .ts, .js, or .svelte file, a root flaiover .ts or .js file, flaiover/tsconfig.json, flaiover/package.json, or the script itself selects it, so flai test finds a type error such as TS2783 at the task, in 8 to 10 s warm, before the close-out. The flaiover tier covers it, so --all and the close-out run svelte-check once. TestRepositoryManifestTypeChecksChangedFlaioverFiles in flai/internal/verify/select_test.go fails without the tier.
```

S-0530's narrative (`wip/archive/agents/S-0530.md`) measured `--incremental` at about 28 s without it, 31 s cold and 8 to 10 s warm with it, and found that a warm cache still reports a prop type change in its callers. `scripts/flaiover-test.sh`, the `--all` `flaiover` tier, keeps plain svelte-check.

Sources: I-0105, I-0206.
