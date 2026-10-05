---
id: T-0847
type: task
nature: improvement
title: S-0249's own close-out ends clean past main's findings, and I-0057 is closed
status: backlog
parent: S-0249
owner: alex
created: 2026-10-05T00:18:47Z
updated: 2026-10-05T00:18:47Z
transitions: []
stream: S-0249
tags: [issues]
touches: [design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md, design/issues/summary.md]
after: [T-0844, T-0845]
---
# T-0847 S-0249's own close-out ends clean past main's findings, and I-0057 is closed

## Work

Run `scripts/close-out.sh S-0249` with nothing committed yet beyond the work, as the dress rehearsal for the second criterion. Confirm that it passes the check step while main still carries findings outside the story, and that it records those findings in issues. Confirm that the issues it opens or bumps are the ones T-0841 names, not I-0057, which this task closes.

Then close I-0057 with `flai issue close I-0057 --reason`. Name the ADR T-0845 wrote and the flags: `flai check --story` reports findings outside the story as notes, `--record-issues` records them, and the close-out scripts pass both. Closing regenerates `design/issues/summary.md`.

It waits for T-0844 and T-0845, so that what it verifies and the reason it gives are complete.

## Done when

- The close-out's output shows main's findings outside the story as notes and names the issues they went to, and the check step does not stop it.
- I-0057 has status closed, with a reason naming the fix, and `design/issues/summary.md` no longer lists it as open.
- `flai check --strict --story S-0249` passes.

## Notes
