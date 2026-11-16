---
id: T-3402
type: task
nature: remediation
title: The convention, the guides, and the template's changelog say a close-out asked for a branch behind main syncs first, and I-0383 is closed
status: done
parent: S-0798
owner: alex
created: 2026-11-16T01:43:38Z
updated: 2026-11-16T02:19:03Z
transitions:
  - to: ready
    at: 2026-11-16T02:17:49Z
    by: system-flow
  - to: in-progress
    at: 2026-11-16T02:17:49Z
    by: system-flow
  - to: done
    at: 2026-11-16T02:19:03Z
    by: agent-S-0798
stream: S-0798
tags: [close-out, template]
touches: [design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, docs/users/flai.md, scripts/README.md, template/root/scripts/README.md, design/issues/I-0383-a-close-out-asked-for-a-story-whose-branch-is-behind-main-stops-at-sync-without.md, design/issues/summary.md, design/system/flai-cli.md]
after: [T-3399, T-3400, T-3401]
usage:
  source: log
  seconds: 68
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 20
      output: 69
      cache_read: 2414732
      cache_write: 13361
      cost: 0.92
---
# T-3402 The convention, the guides, and the template's changelog say a close-out asked for a branch behind main syncs first, and I-0383 is closed

## Work

Say what the close-out now does, in the words of T-3399 and T-3400, and cite the ADR T-3401 wrote.

- `design/conventions/work-management.md` and `template/root/design/conventions/work-management.md`, the close-out rule (line 35). Say that a close-out whose first `flai verify` stops at sync, because the branch was behind main when it was asked for, runs `flai stream sync` itself as the last check's stop does, within the same three rounds. Make the same edit in both baselines.
- `template/CHANGELOG.md`, under `## Unreleased`: one line for the template's `scripts/close-out.sh` and one for the convention, with the bump, the ADR, and I-0383, as S-0781's lines are written.
- `docs/users/flai.md`, the close-out paragraph (line 963). It says today the close-out "refuses a branch that does not yet contain the main branch and says to sync". Say that it syncs, prints the round line, and stops at `flai stream sync` on a conflict. Check § The close-out job and the `verify_behind_main` paragraph (line 3397), which say the close-out syncs the branch: they now hold.
- `scripts/README.md` and `template/root/scripts/README.md`, the `close-out.sh` row: add the sync of a branch behind main, and the rounds when main moves.
- Close I-0383 with `flai issue close I-0383 --reason "<what fixed it: the close-out syncs when its first flai verify stops at sync, in the scripts and a container story's close-out, S-0798>"`. It regenerates `design/issues/summary.md`.

It waits for T-3399 and T-3400, whose behaviour it describes, and for T-3401, whose ADR it cites.

## Done when

- The convention, its template copy, the guides, and the READMEs say what the close-out does, and the template's changelog has its lines.
- I-0383 is closed with its reason.
- `flai test` on the changed files passes, the markdown lint included.

## Notes
