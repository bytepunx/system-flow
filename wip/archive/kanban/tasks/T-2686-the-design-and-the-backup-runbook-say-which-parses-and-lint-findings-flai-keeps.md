---
id: T-2686
type: task
nature: improvement
title: The design and the backup runbook say which parses and lint findings flai keeps under .flai-cache
status: done
parent: S-0600
owner: alex
created: 2026-11-07T23:51:16Z
updated: 2026-11-08T01:15:46Z
transitions:
  - to: ready
    at: 2026-11-08T01:02:12Z
    by: agent-S-0600
  - to: in-progress
    at: 2026-11-08T01:02:14Z
    by: agent-S-0600
  - to: done
    at: 2026-11-08T01:15:46Z
    by: agent-S-0600
stream: S-0600
tags: [cli, docs]
touches: [design/system/flai-cli.md, docs/operators/runbooks/backup.md]
after: [T-2682, T-2685]
usage:
  source: log
  seconds: 812
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 17
      output: 5303
      cache_read: 2565103
      cache_write: 78619
      cost: 1.0636
---
# T-2686 The design and the backup runbook say which parses and lint findings flai keeps under .flai-cache

## Work

In `design/system/flai-cli.md`, under `## Internal structure`, say that a flai process parses items and lints a folder on a bounded pool of workers, and that it keeps parsed items and each file's lint findings in the two files T-2682 and T-2685 write under `.flai-cache`: their names, their keys, what makes an entry a miss, and that they are written only into a `.flai-cache` that exists. In `docs/operators/runbooks/backup.md`, add a row for each file to the table under `## What exists, and what is worth keeping`, beside `.flai-cache/stories-commits.json`: not worth keeping, rebuilt from the files when missing. Bump each document's `updated`.

Third layer: it waits for T-2682 and T-2685, which fix the file names and keys it describes.

## Done when

- [ ] `design/system/flai-cli.md` and `docs/operators/runbooks/backup.md` name both cache files as the code writes them.
- [ ] The markdown lint and `flai check --strict` pass on both documents.

## Notes
