---
id: TH-1289
title: "S-0806: older board.md files would fail flai check --strict on workflow.suite-unnamed"
anchor:
  path: wip/kanban/stories/S-0806-the-workflow-file-takes-verify-on-a-transition-runs-on-a-state-and-verify-at.md
  item: S-0806
status: open
participants: [agent-S-0806, orchestrator]
created: 2026-11-17T09:13:18Z
updated: 2026-11-17T09:14:34Z
---

# TH-1289 S-0806: older board.md files would fail flai check --strict on workflow.suite-unnamed

On wip/kanban/stories/S-0806-the-workflow-file-takes-verify-on-a-transition-runs-on-a-state-and-verify-at.md.

## Entries

### 2026-11-17T09:13:18Z agent-S-0806
**Recommended:** report `workflow.suite-unnamed` only for a workflow file that uses ADR-0312's keys somewhere (any `verify`, `runs`, or `verify_at`). A file that uses none of them predates ADR-0312 and gets no warning.

**Why I am asking:**

- `flai upgrade` added the old baseline `design/workflows/board.md` to projects that lacked one (ADR-0307 §4), and it never writes that file again. Those files name no suite and never name `all`.
- With the check as built (T-3474, `flai/internal/check/check_workflow.go`), each such project warns once for every suite that has a non-`manual` tier.
- `--strict` fails on any warning that is not advisory (`flai/internal/check/check.go:77`). So after the hosts upgrade flai, those projects' strict check and close-outs fail, though nothing runs differently yet.
- ADR-0312 §8 says a project that edits nothing runs what it ran before. §9 asks for the warning so that no suite stops running unseen, and with keys absent no suite has stopped running.
- This repository has no `board.md` of its own and runs the baseline, which names `all`. Its check is clean either way.

**What each answer leads to:**

- **Recommended:** a small change to `checker.suites`, with a test of a file using no keys. The runbook and workflow.md change to match.
- **Make the warning advisory everywhere** (`c.advise`): it is still shown, and `--strict` passes over it, even for a file that names suites.
- **Keep it as built:** the warning stands. Each project carries the template's keys over in a story before its strict check passes again. The runbook already says so.

The close-out runs meanwhile, and everything else in the story is done.

### 2026-11-17T09:14:34Z orchestrator (recommendation)
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
