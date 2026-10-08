---
id: T-1362
type: task
nature: remediation
title: Document that a dashboard restart keeps the image that runs, and close I-0116
status: done
parent: S-0344
owner: alex
created: 2026-10-08T08:35:52Z
updated: 2026-10-08T09:39:29Z
transitions:
  - to: ready
    at: 2026-10-08T09:38:11Z
    by: agent-S-0344
  - to: in-progress
    at: 2026-10-08T09:38:12Z
    by: agent-S-0344
  - to: done
    at: 2026-10-08T09:39:29Z
    by: agent-S-0344
stream: S-0344
tags: [flai, dashboard, docs]
touches: [docs/users/flai.md, docs/operators/index.md, design/system/flai-cli.md, design/issues/I-0116-flai-dashboard-restart-starts-the-container-s-tag-again-so-a-newer-image-a-check-pulled-under-a-floating-tag-is-started-without-an-upgrade.md, design/issues/summary.md]
after: [T-1360, T-1361]
usage:
  source: log
  seconds: 77
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 21
      output: 7530
      cache_read: 1267143
      cache_write: 48383
      cost: 0.7248
---
# T-1362 Document that a dashboard restart keeps the image that runs, and close I-0116

## Work

Third layer: it waits for T-1360 and T-1361, because it documents what they built and closes the issue only once both restarts are fixed.

- `docs/users/flai.md` § Run the dashboard: add `flai dashboard restart` to the command block, saying it starts the image that runs again, never a newer one a check pulled.
- `docs/operators/index.md` § The dashboard's watch: say the record keeps the image ID beside the reference, the watch restarts from the ID, and falls back to the reference when the image is gone.
- `design/system/flai-cli.md` § A chosen dashboard: say restart and the watch start from the image ID, and how the tag is still shown.
- Close the issue with `flai issue close I-0116 --reason "<what fixed it>"`, naming T-1360 and T-1361, which updates `design/issues/summary.md`.

## Done when

- The three documents say what restart and the watch now start from.
- I-0116 is closed with a reason naming the fix.
- `flai check --strict` is clean and the markdown lint passes on the changed files.

## Notes

Drafted by the planner for S-0344. It meets criterion 2.
