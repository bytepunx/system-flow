---
id: T-2921
type: task
nature: improvement
title: The template's scripts/close-out.sh takes --fresh and --all and reads CLOSE_OUT_FRESH and CLOSE_OUT_ALL as this repository's does
status: done
parent: S-0671
owner: alex
created: 2026-11-10T08:48:22Z
updated: 2026-11-10T09:41:41Z
transitions:
  - to: ready
    at: 2026-11-10T09:41:32Z
    by: agent-S-0671
  - to: in-progress
    at: 2026-11-10T09:41:34Z
    by: agent-S-0671
  - to: done
    at: 2026-11-10T09:41:41Z
    by: agent-S-0671
stream: S-0671
tags: [template, close-out, scripts]
touches: [template/root/scripts/close-out.sh, template/root/scripts/README.md, template/CHANGELOG.md]
after: [T-2918]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 415
      cache_read: 470809
      cache_write: 4273
      cost: 0.1367
---
# T-2921 The template's scripts/close-out.sh takes --fresh and --all and reads CLOSE_OUT_FRESH and CLOSE_OUT_ALL as this repository's does

## Work

- Port T-2918's change to `template/root/scripts/close-out.sh`: the leading `--fresh` and `--all` options, the two variables, and their pass to each round's `flai verify --record-issues`. Keep the template copy's own wording, without this repository's story and ADR references.
- Update the script's row in `template/root/scripts/README.md`.
- Add an entry under `## Unreleased` in `template/CHANGELOG.md`: S-0671, the two options and variables, and that flai serve's close-out job now passes `fresh` and `all` on (patch).
- Waits for T-2918, since it copies that task's code and wording rather than writing its own.

## Done when

- `flai test template/root/scripts/close-out.sh` passes, ShellCheck and the template tier included.
- The template's script and this repository's take the same options and variables.

## Notes
