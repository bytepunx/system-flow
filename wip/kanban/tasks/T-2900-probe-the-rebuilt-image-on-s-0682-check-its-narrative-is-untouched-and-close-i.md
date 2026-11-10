---
id: T-2900
type: task
nature: remediation
title: Probe the rebuilt image on S-0682, check its narrative is untouched, and close I-0224
status: in-progress
parent: S-0682
owner: alex
created: 2026-11-10T08:46:31Z
updated: 2026-11-10T09:08:44Z
transitions:
  - to: ready
    at: 2026-11-10T09:08:43Z
    by: agent-S-0682
  - to: in-progress
    at: 2026-11-10T09:08:44Z
    by: agent-S-0682
stream: S-0682
tags: [containers, issues]
touches: [design/issues/I-0224-scripts-agent-container-sh-can-run-only-a-story-s-own-agent-under-the-story-s-own-name-so-a-trial-run-of-the-image-acts-on-the-live-story.md, design/issues/summary.md]
after: [T-2897, T-2898]
---
# T-2900 Probe the rebuilt image on S-0682, check its narrative is untouched, and close I-0224

## Work

Waits for T-2897 and T-2898: the trial runs the probe pilot and the script now make, in an image built from them.

- Rebuild the images from the story's worktree with `make agent-images`, so pilot carries the probe.
- Run the free check first: `scripts/agent-container.sh --probe S-0682 -- env`, and record in this task's Notes that `FLAI_AGENT=probe-S-0682` is set and `FLAI_STORY` is not.
- The capped trial calls a paid model. Ask on a thread on S-0682 before it, recommending one run with `FLAI_AGENT_MODEL=haiku FLAI_AGENT_MAX_BUDGET_USD=0.02 scripts/agent-container.sh --probe S-0682`. With the operator's yes, run it, and check that S-0682's narrative's Current state, Next steps, and Log are what they were before it and that no item changed (`flai stream state` and `inbox` show nothing by `probe-S-0682`). Record what it reported. Without docker or a Claude login on the host, record that instead and rest criterion 1 on T-2897's test.
- Tick criterion 1 with `flai criteria tick S-0682 1` once T-2897's test passes and the checks above hold.
- Run `flai issue close I-0224 --reason "<reason>"`: the reason names `scripts/agent-container.sh --probe` and the start request's `probe`, which runs the image as `probe-S-nnnn` with a prompt that works no item and no `FLAI_STORY`, held by T-2897's test. Then tick criterion 2.

## Done when

- The free check, and the capped trial or the reason it did not run, are recorded in this task's Notes.
- I-0224 is closed with that reason, and `design/issues/summary.md` no longer lists it as open.
- Both of S-0682's criteria are ticked.

## Notes
