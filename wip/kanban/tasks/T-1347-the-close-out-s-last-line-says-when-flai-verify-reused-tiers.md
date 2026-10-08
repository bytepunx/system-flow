---
id: T-1347
type: task
nature: improvement
title: The close-out's last line says when flai verify reused tiers
status: backlog
parent: S-0341
owner: alex
created: 2026-10-08T08:05:32Z
updated: 2026-10-08T08:05:32Z
transitions: []
stream: S-0341
tags: [cli]
touches: [scripts/close-out.sh]
after: [T-1344]
---
# T-1347 The close-out's last line says when flai verify reused tiers

## Work

Criterion 3. It waits for the command task (`flai verify takes --fresh and prints reused tiers…`), whose last line it reads.

- Keep running `flai verify "$story" --record-issues` as now, with no `--fresh`, so a second close-out at an unchanged head resumes.
- Capture verify's last line while still printing its whole output, without a pipe that hides its exit status (POSIX `sh`, `set -eu`), and carry its `(reused N tiers from <time>)` into the close-out's own last line, for a pass and for a stop at `flai verify`.
- A commit the close-out makes after verify changes the head, so the next run is full; say nothing more about it.

## Done when

- A close-out run twice at an unchanged head, the first stopped at a late tier, ends the second with a last line naming the reused tiers; a run with nothing reused ends as it does today. Check it by hand in a scratch story or with a test of the script, and say which in the task's log entry.
- `flai test scripts/close-out.sh` passes.

## Notes

The last line is printed by the `finish` trap; `hint` already adds to it when verify fails.
