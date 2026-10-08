---
id: T-1375
type: task
nature: improvement
title: Record in an ADR and the CLI design that the sync check passes over commits on main that change only wip the branch does not change
status: done
parent: S-0347
owner: alex
created: 2026-10-08T08:43:52Z
updated: 2026-10-08T10:31:07Z
transitions:
  - to: ready
    at: 2026-10-08T10:29:59Z
    by: agent-S-0347
  - to: in-progress
    at: 2026-10-08T10:30:00Z
    by: agent-S-0347
  - to: done
    at: 2026-10-08T10:31:07Z
    by: agent-S-0347
stream: S-0347
tags: [cli, docs]
touches: [design/adrs/0133-the-close-out-s-sync-check-passes-over-commits-on-main-that-change-only-wip.md, design/adrs/README.md, design/system/flai-cli.md, design/system/devex.md]
usage:
  source: log
  seconds: 67
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 17
      output: 5995
      cache_read: 979793
      cache_write: 31548
      cost: 0.5224
---
# T-1375 Record in an ADR and the CLI design that the sync check passes over commits on main that change only wip the branch does not change

## Work

I-0119's three instances: a close-out passed every verify step, then stopped at its last check that the branch contains main, because main moved during the four-to-ten-minute run. In two of them every commit main gained was flai's own commit of `wip/`: a replan after a reorder or a cancel, a forecast replan, or an edit of another story's criteria. Those commits cannot change what the branch's tiers verified: `wip/` stays on main ([ADR-0019](../../../design/adrs/0019-story-branches-and-touches.md)), the branch does not change it, the tiers that read the repository scope themselves to the story under `CLOSE_OUT_STORY` ([ADR-0085](../../../design/adrs/0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md)), and acceptance merges the branch, so the merge brings them.

Propose the fix and record it before it is built:

- Write an ADR with `flai adr new` (the number and slug in `touches` are a prediction; take what flai gives): the `sync` step of `flai verify` and the close-out's last check pass when every commit of the main branch that the story's branch lacks changes only paths under the manifest's `wip` folder, none of which the branch changes. They name those commits in a note. Any other commit the branch lacks, such as an acceptance or a release, still fails the step, as now.
- Name in the ADR the alternatives weighed: an automatic sync followed by a full run, and S-0341's resume, which a rebase defeats because it changes the head.
- Name the way the close-out reaches the rule: `flai verify S-nnnn --sync-only`, which runs the `rebase` and `sync` steps alone and stores no record, in place of `git merge-base --is-ancestor`. The template's close-out then needs the flai that has it.
- Update `design/system/flai-cli.md`'s `flai verify` row (the `sync` step, `--sync-only`, the close-out's last check) and the close-out row of `design/system/devex.md`. Add the ADR to `design/adrs/README.md`.

Runs alone in the first layer: every other task builds or documents what it decides.

## Done when

- The ADR is accepted and listed in `design/adrs/README.md`.
- `design/system/flai-cli.md` and `design/system/devex.md` describe the `sync` step's rule, `--sync-only`, and the close-out's last check as the ADR decides them.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes

Drafted by the planner for S-0347.
