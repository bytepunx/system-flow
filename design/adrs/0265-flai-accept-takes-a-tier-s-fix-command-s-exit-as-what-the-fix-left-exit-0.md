---
id: ADR-0265
title: "flai accept takes a tier's fix command's exit as what the fix left: exit 0 passes its files, a non-zero exit refuses with the findings its output gives, and the check runs after the fix only when the fix cannot say"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0139, ADR-0140]
---

# ADR-0265 flai accept takes a tier's fix command's exit as what the fix left: exit 0 passes its files, a non-zero exit refuses with the findings its output gives, and the check runs after the fix only when the fix cannot say

## Context

[ADR-0139](0139-a-test-tier-s-fix-command-rewrites-the-files-a-story-changes-before-flai-task.md) runs a tier's `fix_command` before flai commits and ignores its exit status. [ADR-0140](0140-flai-accept-refuses-an-acceptance-whose-files-a-tier-with-a-fix-command-still.md) then has `flai accept` run the tier's `command` over the same files to learn what the fix left. For this repository's `markdown` tier that is two runs of `scripts/lint-md.sh` per set where one says everything: `markdownlint-cli2 --fix` exits 1 and prints each error it could not fix, as the check does. An acceptance ran the tier five times, each through `npx`, about 1 s apiece (cause 19 in [server-performance.md](../system/server-performance.md#where-the-time-goes), S-0609).

## Decision

`flai accept` takes a tier's fix command's exit as what the fix left, in the main checkout at preflight and on the story's branch after its rebase (`verify.FixRemaining`).

1. **Exit 0** passes the files the fix ran over: nothing remains, and no check runs.
2. **A non-zero exit** is lint that remains: the acceptance is refused with the findings the fix's output gives, parsed in the tier's `format` as its `command`'s output is, with the same text, findings, and `--json` refusal ADR-0140 and S-0378 give.
3. **A fix that cannot say** runs the check after it, as ADR-0140 decides: a command that cannot start, a run cut short, or a structured format whose output gives no finding of its own. A branch whose fix did not run is checked the same way.
4. **The fixes before the acceptance's commits** run only over the files that are new or changed since the preflight fix passed them, by content, and not at all when none are.

`flai task done` and `flai stream sync` keep ADR-0139's rule: their fix stops nothing, and its exit is logged at warn. The merge queue's preflight keeps its fix and check.

## Consequences

- An acceptance lints the uncommitted set once and the branch's set once.
- A `fix_command` must exit non-zero while lint remains, as `markdownlint-cli2 --fix` and `eslint --fix` do. One that exits 0 and leaves lint, such as `true`, lets that lint through an acceptance. `project-manifest.md` says so.
- A fix that exits non-zero for a reason other than lint, such as a missing tool, refuses the acceptance with its last lines of output, as ADR-0140's check already did for a check that cannot start.

## Alternatives considered

- **Keep the check after the fix.** Correct for any fix command, but it doubles the lint of each set, which is what S-0609 removes.
- **A manifest flag saying a fix's exit can be trusted.** More to configure for every tier, when every fix command in use already exits as its check does.
