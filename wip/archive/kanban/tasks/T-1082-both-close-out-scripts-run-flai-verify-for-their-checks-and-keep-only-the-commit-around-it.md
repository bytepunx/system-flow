---
id: T-1082
type: task
nature: improvement
title: Both close-out scripts run flai verify for their checks and keep only the commit around it
status: done
parent: S-0270
owner: alex
created: 2026-10-06T22:52:46Z
updated: 2026-10-07T04:33:21Z
transitions:
  - to: ready
    at: 2026-10-07T03:49:38Z
    by: agent-S-0270
  - to: in-progress
    at: 2026-10-07T03:49:39Z
    by: agent-S-0270
  - to: done
    at: 2026-10-07T04:33:21Z
    by: agent-S-0270
stream: S-0270
tags: [flai, template]
touches: [scripts/close-out.sh, template/root/scripts/close-out.sh, scripts/README.md, template/root/scripts/README.md, system-flow.yaml, template/root/system-flow.yaml.tmpl]
after: [T-1075]
usage:
  source: log
  seconds: 2622
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 66
      output: 27465
      cache_read: 4754468
      cache_write: 137996
      cost: 2.3616
---
# T-1082 Both close-out scripts run flai verify for their checks and keep only the commit around it

## Work

Criterion 3: the close-out and `flai verify` never disagree, because the close-out calls it.

- In `scripts/close-out.sh`, replace the rebase check, the tier selection by `touched`, the markdown lint, `scripts/check.sh`, the flaiover tier, and the narrative check with one call, `flai verify "$story" --record-issues`, run through `scripts/flai.sh` as the rest of the script runs flai. Stop when it fails, and keep its last line as the close-out's.
- Keep what `flai verify` does not do: the commit with the options given, the commit of the issues it recorded, the sync check after the commit, and the clean worktree check. Keep the last line that names the outcome and the step it stopped at (S-0266).
- Do the same in `template/root/scripts/close-out.sh`, whose tiers come from the project's tier declaration (S-0273) instead of its fixed lint, test, integration, and smoke list.
- Update both `scripts/README.md` files.
- Waits for T-1075: the script calls the command. Runs alongside T-1089 and T-1096: no path in common.

## Done when

- [ ] Both close-out scripts call `flai verify` and run no check of their own that it also runs.
- [ ] A run on a story with a failing step stops at the step `flai verify` names, and a passing run commits and ends with the passing line.
- [ ] `scripts/template-test.sh` passes, and `scripts/close-out.sh` passes on this story's own branch.

## Notes

Drafted by the planner. A project made from the template needs the flai release that ships `flai verify`. T-1096's `template/CHANGELOG.md` line says so.
