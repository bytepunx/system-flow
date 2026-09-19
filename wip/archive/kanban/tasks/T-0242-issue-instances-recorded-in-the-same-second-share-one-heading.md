---
id: T-0242
type: task
nature: feature
title: Issue instances recorded in the same second share one heading
status: done
parent: S-0047
owner: alex
created: 2026-09-19T10:30:01Z
updated: 2026-09-19T10:30:50Z
transitions:
  - to: ready
    at: 2026-09-19T10:30:02Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T10:30:02Z
    by: system-flow
  - to: done
    at: 2026-09-19T10:30:50Z
    by: system-flow
stream: S-0047
tags: []
---

# T-0242 Issue instances recorded in the same second share one heading

## Work
`insertInstance` in `flai/internal/issues` writes a new `### <timestamp>` heading only when the last instance does not already carry that timestamp; otherwise the note goes under the existing heading as a further paragraph, the way narrative log entries do since I-0011. `flai issue new` followed by `flai issue bump` in the same second is the same case. A test bumps twice at one instant and checks one heading, both notes, the count, and that the repository check's duplicate-heading rule passes on the file.

## Done when
- One heading for same-second instances, tested, and the test fails without the change

## Notes
