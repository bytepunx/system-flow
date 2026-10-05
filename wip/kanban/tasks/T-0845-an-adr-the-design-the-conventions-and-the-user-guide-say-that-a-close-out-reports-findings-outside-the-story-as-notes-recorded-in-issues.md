---
id: T-0845
type: task
nature: improvement
title: An ADR, the design, the conventions, and the user guide say that a close-out reports findings outside the story as notes recorded in issues
status: backlog
parent: S-0249
owner: alex
created: 2026-10-05T00:18:31Z
updated: 2026-10-05T00:18:31Z
transitions: []
stream: S-0249
tags: [flai, template, docs]
touches: [design/adrs, design/system/flai-cli.md, design/system/continuous-improvement.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, docs/users/flai.md]
after: [T-0841]
---
# T-0845 An ADR, the design, the conventions, and the user guide say that a close-out reports findings outside the story as notes recorded in issues

## Work

Write a new ADR, with `flai adr new`, that refines ADR-0073. It records three decisions:

- `flai check --story` reports the findings outside the story as notes the run passes over, errors included, and says which paths count as the story's.
- `--record-issues` opens or bumps one issue per rule.
- The close-out scripts pass both flags. The unscoped check, as CI runs it, still gates everything.

Say why: I-0057's 19 instances, and TH-0056's answer.

Update `design/system/flai-cli.md`, the `flai check` section, with the flags. Update `design/system/continuous-improvement.md` to say that flai itself records issues from a close-out's findings outside the story. In `docs/users/flai.md`, the section "Check the repository" gets the flags and an example.

The close-out rule in `design/conventions/work-management.md` says to move to review only when the close-out ends clean. Add that findings outside the story are printed as notes and recorded in issues, and do not stop it. Make the same edit to `template/root/design/conventions/work-management.md`. Record both template changes, the convention and T-0844's scripts, in `template/CHANGELOG.md`.

It waits for T-0841, so that it describes the flags as built. It runs beside T-0844: the two share no path.

## Done when

- The ADR is accepted in `design/adrs`, and `design/adrs/README.md` lists it.
- `flai-cli.md`, `continuous-improvement.md`, `docs/users/flai.md`, and both copies of `work-management.md` describe the scoped check. The two `work-management.md` baselines match.
- `template/CHANGELOG.md` has an entry for the scripts and the convention.
- `scripts/lint-md.sh` and `flai check --strict --story S-0249` pass.

## Notes
