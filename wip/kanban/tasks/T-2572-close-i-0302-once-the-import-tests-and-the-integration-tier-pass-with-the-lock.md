---
id: T-2572
type: task
nature: remediation
title: Close I-0302 once the import tests and the integration tier pass with the lock written
status: backlog
parent: S-0590
owner: alex
created: 2026-11-07T10:03:08Z
updated: 2026-11-07T10:03:08Z
transitions: []
stream: S-0590
tags: [issues, import]
touches: [design/issues/I-0302-flai-import-writes-no-system-flow-lock-yaml-so-the-import-formula-s-flai-check.md, design/issues/summary.md]
after: [T-2570, T-2571]
---
# T-2572 Close I-0302 once the import tests and the integration tier pass with the lock written

## Work

Close I-0302 only once the fix and its documents are in, so this task waits for T-2570, the fix and its test, and T-2571, the documents.

First run `flai test flai/cmd/import.go flai/cmd/import_commit.go`. Then run the integration tier, `flai test integration`, which runs a real import. A failure there that is not this story's is recorded as the close-out says, and does not hold this task.

Then close the issue:

```bash
flai issue close I-0302 --reason "flai import writes system-flow.lock.yaml from the render's hashes and commits it with the import (S-0590, T-2570), so the import formula's flai check --strict step no longer fails on layout.lock"
```

`flai issue close` rewrites `design/issues/summary.md`. Commit both files.

## Done when

- [ ] `flai test` on the import files and the integration tier pass
- [ ] I-0302 is closed with a reason naming the fix
- [ ] Story criteria 1 and 2 are ticked with `flai criteria tick`

## Notes
