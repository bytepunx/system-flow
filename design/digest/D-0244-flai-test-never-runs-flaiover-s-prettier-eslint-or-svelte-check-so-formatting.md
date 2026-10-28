---
id: D-0244
title: flai test never runs flaiover's prettier, eslint, or svelte-check, so formatting faults surface only at the close-out
kind: cause
sources: [I-0105]
topics: [testing, flai, cli, go, code, all]
paths: [scripts/flaiover-lint.sh, system-flow.yaml, flai/internal/verify/select_test.go, design/system/devex.md, design/conventions/tooling.md, design/issues/I-0105-flai-test-never-runs-flaiover-s-prettier-eslint-or-svelte-check-so-formatting-faults-surface-only-at-the-close-out.md, design/issues/summary.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md]
date: 2026-10-08
status: active
absorbed_by: ""
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

Source: I-0105.
