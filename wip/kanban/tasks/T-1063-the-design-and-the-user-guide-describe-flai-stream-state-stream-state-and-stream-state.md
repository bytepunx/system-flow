---
id: T-1063
type: task
nature: improvement
title: The design and the user guide describe flai stream state, stream_state, and stream.state
status: backlog
parent: S-0271
owner: alex
created: 2026-10-06T22:50:17Z
updated: 2026-10-06T22:50:17Z
transitions: []
stream: S-0271
tags: [docs]
touches: [design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1055, T-1056, T-1058, T-1059]
---
# T-1063 The design and the user guide describe flai stream state, stream_state, and stream.state

## Work

Criterion 3, for the design and the user guide.

- `design/system/agent-narrative.md`: under `## Current state` and `## Next steps`, and in the steps that say to rewrite them, say they are written with `flai stream state`, which replaces them and appends nothing to the log. Name the check rule from T-1059.
- `design/system/flai-cli.md`:
  - Add `flai stream state` to the stream commands' row.
  - Add `stream.state` to the `flai hostapi` method list.
  - Add `stream_state` wherever the MCP tools are listed.
  - Add the new check rule wherever the check rules are listed.
- `docs/users/flai.md`: in the narrative guidance, use `flai stream state` in place of editing the file.
- Regenerate `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md` with `scripts/flai-reference.sh`. Do not edit them by hand.

## Done when

- [ ] None of these documents says to edit `## Current state` or `## Next steps` by hand.
- [ ] `scripts/flai-reference.sh` leaves no diff after a second run.
- [ ] `flai check --strict` and the markdown lint pass.

## Notes

Drafted by the planner for S-0271. `docs/operators/settings.md` is predicted only because the regeneration adds `--current` and `--next` to its flag index.
