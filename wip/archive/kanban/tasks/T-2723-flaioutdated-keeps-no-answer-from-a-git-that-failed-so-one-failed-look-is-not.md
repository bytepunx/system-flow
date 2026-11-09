---
id: T-2723
type: task
nature: remediation
title: FlaiOutdated keeps no answer from a git that failed, so one failed look is not served as not outdated for a minute
status: done
parent: S-0617
owner: alex
created: 2026-11-09T10:34:43Z
updated: 2026-11-09T10:47:00Z
transitions:
  - to: ready
    at: 2026-11-09T10:40:20Z
    by: agent-S-0617
  - to: in-progress
    at: 2026-11-09T10:40:22Z
    by: agent-S-0617
  - to: done
    at: 2026-11-09T10:47:00Z
    by: agent-S-0617
stream: S-0617
tags: [flai]
touches: [flai/internal/release/outdated.go, flai/internal/release/outdated_test.go]
usage:
  source: log
  seconds: 398
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 22
      output: 3112
      cache_read: 2537898
      cache_write: 40779
      cost: 0.8901
---
# T-2723 FlaiOutdated keeps no answer from a git that failed, so one failed look is not served as not outdated for a minute

## Work

`release.FlaiOutdated` (`flai/internal/release/outdated.go`) keeps each answer for a minute in `outdatedSeen`. `flaiOutdated` answers nil both when the flai is not outdated and when `git tag --merged HEAD --list flai/v*` fails. So one failed git call on a loaded host, such as a fork refused, is kept and served as "not outdated" for a minute: `project.info`'s `flai_outdated`, the inbox's, `flai task done`'s, and `flai verify`'s all read it. That is the one way left for `project.info` to answer the right project with `flai_outdated` nil while the project's history holds a newer flai tag, the symptom I-0252 records.

- Make `flaiOutdated` say whether git answered, and have `FlaiOutdated` store an answer only when it did. A failed look answers nil, as now, and the next call asks git again.
- Leave `tagRules` as it is: a tag with no rules file is a real answer of 0 (ADR-0180).
- Add a test in `flai/internal/release/outdated_test.go` with a runner whose first `git tag` call fails and whose second answers a newer tag: the first call answers nil, the second answers the newer release. It fails on the code before this change.

Waits for nothing: it changes only the release package.

## Done when

- A failed `git tag` look is not kept, and the next `FlaiOutdated` call for the same root and version asks git again.
- The new test fails before the change and passes after it; `flai test flai/internal/release` passes.

## Notes
