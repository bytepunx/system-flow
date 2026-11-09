---
id: T-2725
type: task
nature: remediation
title: Close I-0252 naming S-0569's fix and T-2723's
status: done
parent: S-0617
owner: alex
created: 2026-11-09T10:34:56Z
updated: 2026-11-09T10:49:34Z
transitions:
  - to: ready
    at: 2026-11-09T10:49:25Z
    by: agent-S-0617
  - to: in-progress
    at: 2026-11-09T10:49:25Z
    by: agent-S-0617
  - to: done
    at: 2026-11-09T10:49:34Z
    by: agent-S-0617
stream: S-0617
tags: [flai, issues]
touches: [design/issues/I-0252-testservesaystheflaiisolderthantheproject-fails-under-load-project-info-is.md, design/issues/summary.md]
after: [T-2723, T-2724]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 658
      cache_read: 536543
      cache_write: 8621
      cost: 0.1882
---
# T-2725 Close I-0252 naming S-0569's fix and T-2723's

## Work

Close I-0252 in the story's worktree with `flai issue close I-0252 --reason "<reason>"`. The reason says what fixed it, in one or two sentences:

- its instance was I-0247's notification read as the reply to `project.info`, which S-0569 (`ea5280843`) fixed after the instance;
- T-2723 stopped `release.FlaiOutdated` keeping a failed git look as "not outdated" for a minute, the one way left for `project.info` to answer `flai_outdated` nil;
- T-2724's stress run, with its count, passed.

Correct the reason to what T-2723 and T-2724 found if it differs from this. `flai issue close` regenerates `design/issues/summary.md`.

Waits for T-2723 and T-2724: the reason names what each found.

## Done when

- I-0252's status is closed, with a reason naming S-0569, T-2723, and the stress run's result.
- `flai check --strict` passes on the story.

## Notes
