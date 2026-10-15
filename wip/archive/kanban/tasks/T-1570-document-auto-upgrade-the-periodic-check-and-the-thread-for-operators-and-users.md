---
id: T-1570
type: task
nature: improvement
title: Document auto-upgrade, the periodic check, and the thread for operators and users
status: done
parent: S-0387
owner: alex
created: 2026-10-14T20:54:56Z
updated: 2026-10-15T02:09:47Z
transitions:
  - to: ready
    at: 2026-10-15T02:09:26Z
    by: agent-S-0387
  - to: in-progress
    at: 2026-10-15T02:09:26Z
    by: agent-S-0387
  - to: done
    at: 2026-10-15T02:09:47Z
    by: agent-S-0387
stream: S-0387
tags: [docs]
touches: [docs/operators/index.md, docs/operators/settings.md, docs/operators/runbooks/update.md, docs/users/flai.md]
after: [T-1568, T-1569]
usage:
  source: log
  seconds: 21
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 5313
      cache_read: 777349
      cache_write: 21980
      cost: 0.402
---
# T-1570 Document auto-upgrade, the periodic check, and the thread for operators and users

## Work

Describe what T-1568 and T-1569 built, as they built it; it waits for both so the documents say what the code does.

- `docs/operators/index.md` § Keeping the host's flai current: `flai serve` notices a newer release while it runs, upgrades the host with `auto-upgrade` on, and otherwise asks on one thread; a story's agent keeps running through the host's restart.
- `docs/operators/settings.md`: a row for the `auto-upgrade` host action in the index.
- `docs/operators/runbooks/update.md`: turning `auto-upgrade` on in a shell, and what to do with the thread when it is off.
- `docs/users/flai.md`: the host actions row, and the refusal `flai test` gives when the running flai is older than the project.

## Done when

- The four documents describe `auto-upgrade`, the periodic check, the thread, and the version-aware refusal, with `updated` bumped.
- The settings index test and the markdown lint pass (`flai test` on the four paths).

## Notes
