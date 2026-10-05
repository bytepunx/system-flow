---
id: ADR-0085
title: "A close-out's flai check reports findings outside the story as notes and records them in issues"
status: accepted
date: 2026-10-05
supersedes: []
superseded_by: []
refines: [ADR-0073]
topics: [cli, conventions, template]
---

# ADR-0085 A close-out's flai check reports findings outside the story as notes and records them in issues

## Context

`scripts/close-out.sh` runs `flai check --strict`, and the work-management convention moves a story to review only when the close-out ends clean. The check reads the whole project, `wip/` included, so it fails on findings no change in the story made and its agent may not clear: another story's unmerged branch (`story.unaccepted`), a thread answered on an archived story (`threads.archived`), an overlap between two stories in progress side by side (`wip.overlap`), and errors on the main branch, such as two issue files with one ID (`issues.duplicate-id`), that no story's close-out could clear.

[ADR-0073](0073-a-full-review-holds-the-pull-and-flai-check-strict-passes-over-review-over-its.md) took two such findings out of `--strict`, review over its limit and, later, an epic behind its stories, because only the operator clears them. It named the warnings one by one. The rest kept stopping close-outs: I-0057 counts 21, from S-0176's to S-0252's, each of which went to review with the finding noted in its narrative or had a verifier run the remaining steps by hand.

On TH-0056 the operator chose a check scoped to the story over an exception written into the convention: "notes are preferable when its a finding outside the story. Have flai record the finding as an issue (if one doesn't exist) or increment the count so operators can see where these are occurring."

## Decision

A story's close-out scopes `flai check` to the story, passes over what lies outside it as notes, and records those in one issue per rule, while the unscoped check, as CI runs it, still gates everything.

1. **`flai check --story S-nnnn` reports the findings outside the story as notes the run passes over, errors included.** A finding is inside the story when it is on the story's item file or one of its tasks', its narrative, a thread anchored on the story or one of its tasks, a path its branch changes against the main branch, or a path uncommitted in its worktree. Every other finding is outside it, and so is every `wip.overlap`, which the pull hold and the other story's agent clear. A finding outside keeps its level, is printed with `(outside S-nnnn)`, and fails the run neither as an error nor under `--strict`; the summary line counts it, and `--json` marks it `outside` and counts it in `outside`. Without `--story` nothing changes.
2. **`--record-issues`, given with `--story`, records each rule's findings outside the story in one issue per rule.** The issue is the open one titled ``flai check finds `<rule>` outside the story at close-out``. When none is open, flai opens it with class `efficiency` and count 1; otherwise it bumps it. The instance names the story and each finding's path and message. The same story with the same findings is recorded once, so a close-out run again does not inflate the count. The issue is written in the checkout the run reads, which at close-out is the story's worktree, and `summary.md` is regenerated.
3. **The close-out passes both flags, and the unscoped check still gates everything.** `scripts/close-out.sh` exports the story as `CLOSE_OUT_STORY`; `scripts/check.sh`, reached directly and through the smoke tier, then runs `flai check --strict --story "$CLOSE_OUT_STORY" --record-issues`, and the close-out's commit step commits the issue files it wrote. In this repository `TestMonorepoIsClean` scopes itself the same way when `CLOSE_OUT_STORY` is set. Unset, as in CI and `make`, the check is unscoped and every finding counts. The template's scripts do the same.

This refines ADR-0073: the findings it took out of `--strict` by name stay out of it, and a close-out no longer depends on such a list, since anything outside the story is a note there.

## Consequences

- A story's close-out stops only on what the story can fix, so its agent no longer finishes the steps by hand or goes to review with a failing check explained in its narrative.
- The operator sees where findings outside stories occur in `design/issues`, with a count per rule and an instance per story, rather than in narratives; the issues reach the main branch when the story that recorded them is accepted.
- A finding on the main branch is still the main branch's: CI's unscoped check fails on it, and someone must clear it. The close-out no longer forces it on whichever story closes next.
- The issue per rule grows by an instance per story that meets it. An issue left open past `issues.story_after` with no story linking it draws `issues.no-story`, which is the prompt to remediate it or close it.
- The title is the lookup key. Renaming an open issue, or closing it, makes the next close-out open a new one for the rule.
- A finding inside the story but caused elsewhere, such as an error on a file the story changes, still stops the close-out; the story touched it.
- The template's `check.sh` now passes `--story` and `--record-issues` in a close-out, so a project made from it needs a flai that has them.

## Alternatives considered

- **Name the exception in the convention**, letting a close-out that stops only on findings outside the story go to review with them noted. The operator declined it on TH-0056; it keeps the manual steps and leaves the findings in narratives no one aggregates.
- **Pass over more rules in `--strict`, as ADR-0073 did.** Each new rule outside a story would need its own exception, and errors on the main branch would still stop every close-out.
- **Scope by `wip/` only.** Errors such as `issues.duplicate-id` sit under `design/`, and a finding under `wip/` can be the story's own, on its items or narrative.
- **One issue per finding, or one per story.** Per finding floods `design/issues`; per story hides which rule recurs. One per rule shows where they occur, as the operator asked.
- **Record in the main checkout.** A close-out runs in the story's worktree and commits there; writing to the main checkout would leave changes on the main branch outside any story.
