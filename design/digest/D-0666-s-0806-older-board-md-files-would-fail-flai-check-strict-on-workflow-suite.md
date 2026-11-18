---
id: D-0666
title: "S-0806: older board.md files would fail flai check --strict on workflow.suite-unnamed: Recommended: agent-S-0806's first option, with one condition"
kind: answer
sources: [TH-1289]
topics: [verify, workflow, cli, template, flai, go, conventions, code, all]
paths: [flai/internal/workflow/workflow.go, flai/internal/workflow/workflow_test.go, flai/internal/workflow/baseline.go, flai/internal/workflow/suites.go, flai/internal/workflow/suites_test.go, flai/internal/workflow/baseline/board.md, flai/internal/workflow/baseline_test.go, template/root/design/workflows/board.md, template/CHANGELOG.md, flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/fields_test.go, flai/internal/upgrade/upgrade_test.go, flai/internal/check/check_workflow.go, flai/internal/check/workflow_test.go, design/system/workflow.md, design/system/work-hierarchy.md, design/system/flai-cli.md, docs/operators/runbooks/change-the-board.md, docs/users/flai.md, flai/internal/workflow/fields.go, flai/internal/workflow/fields_test.go, design/issues/I-0388-flai-serve-s-auto-upgrade-restart-interrupts-the-close-out-jobs-it-was-running.md, design/issues/summary.md]
date: 2026-11-18
status: absorbed
absorbed_by: "design/system/workflow.md#suites-a-step-needs-and-runs-adr-0312adrs0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transitionmd"
asked_again: []
cited: []
---

# D-0666 S-0806: older board.md files would fail flai check --strict on workflow.suite-unnamed: Recommended: agent-S-0806's first option, with one condition

The operator, alex, answered TH-1289 at 2026-11-18T04:05:38Z.

What was asked:

```text
**Recommended:** report `workflow.suite-unnamed` only for a workflow file that uses ADR-0312's keys somewhere (any `verify`, `runs`, or `verify_at`). A file that uses none of them predates ADR-0312 and gets no warning.
```

The answer:

```text
**Recommended:** agent-S-0806's first option, with one condition. `flai check` reports `workflow.suite-unnamed` only for a workflow file that uses at least one of ADR-0312's keys (`verify`, `runs`, or `verify_at`). The condition: a file with none of those keys is read as giving the baseline's keys, which are `verify: all` at `warn` on the move to review, `verify: all` on the acceptance, `batch: all`, `ci: all`, and `runs: {kluster: [deploy, checks]}` on `review`. Write that rule into workflow.md and the runbook in this story, so that S-0807, S-0808, S-0809, and S-0812 build on it.

**Why:**

- ADR-0312 §8 says a project that edits nothing runs what it ran before. A project whose `board.md` came from `flai upgrade` (ADR-0307 §4) edited nothing, so it should keep running every suite.
- ADR-0312 §9 has the warning so that "no suite stops running unseen". If a file with no keys runs every suite, none stops running, so the warning there would be false. It would also fail `--strict` and every close-out in those projects.
- Without the condition, the narrower check is unsafe. Once S-0807 and S-0808 read the keys, a file with no keys would need no suite at review or acceptance, and the check would then hide exactly what §9 guards.

**Why I am asking you:** this narrows §9 and settles how §8 applies to a file `flai upgrade` wrote. Neither is recorded yet.

**Alternatives:**

- **Make the warning advisory everywhere** (`c.advise`). `--strict` passes, but a project that does name suites loses the guard §9 asks for.
- **Keep it as built.** Each upgraded project's strict check fails until a story of its own carries the template's keys over.

Source: design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md § Decision

Source: design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md § Decision
```

Source: TH-1289.
