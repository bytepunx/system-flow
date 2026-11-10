---
id: I-0280
title: scripts/close-out.sh takes no --fresh, so at a commit the verify record holds a pass for it runs no tier, and a second close-out in another place cannot check the same commit again
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-11-07T06:58:50Z
last_reported: 2026-11-07T06:58:50Z
updated: 2026-11-10T08:34:48Z
---

# I-0280 scripts/close-out.sh takes no --fresh, so at a commit the verify record holds a pass for it runs no tier, and a second close-out in another place cannot check the same commit again

## Description
scripts/close-out.sh takes no --fresh, so at a commit the verify record holds a pass for it runs no tier, and a second close-out in another place cannot check the same commit again

## Instances

### 2026-11-07T06:58:50Z
Story: S-0565.
Met on S-0565's T-2452, 2026-11-07. S-0556's close-out passed on the host as a flai serve job at `45b5126`. A second close-out inside the layered container at the same clean commit would answer from the verify record and run nothing (ADR-0229 §5). `scripts/close-out.sh` passes `--record-issues` alone to `flai verify`, and its remaining arguments go to `git commit`, so there is no way to ask it for `--fresh`. The run inside used `scripts/flai.sh verify S-0556 --fresh --record-issues` and then `--sync-only`: the close-out's two verify steps, without its commit step, which had nothing to commit. Remedy: an environment variable or a leading flag that `close-out.sh` passes on to `flai verify` as `--fresh`.

## Remediation

Story S-0671 remediates this issue, created from it at 2026-11-10T08:34:48Z.
