---
id: S-0056
type: story
nature: improvement
title: Estimates are written by the agent that pulls the story, in agent time
status: in-progress
parent: E-0001
owner: alex
created: 2026-09-19T02:05:29Z
updated: 2026-09-19T09:57:19Z
transitions:
  - to: ready
    at: 2026-09-19T06:37:43Z
    by: alex
  - to: in-progress
    at: 2026-09-19T09:57:19Z
    by: system-flow
tags: [cli, conventions]
touches: [flai/cmd, design/conventions, design/system, template]
---

# S-0056 Estimates are written by the agent that pulls the story, in agent time

## Goal
Stories and tasks carry estimates of how long an agent will take to implement them, written by the agent that pulls the work and calibrated against how long agents have actually taken here, so estimate versus actual means something and nobody reads a human-effort number as a forecast.

## Acceptance criteria
- [ ] The standard says what an estimate is: the expected time in `in-progress` for an agent doing the work, as a Go duration, written by the agent that pulls the story when it writes the tasks (ADR-0021), for the story and optionally for each task; `design/system/work-hierarchy.md` and the work-management convention, template baseline first, say so
- [ ] `flai` can set it without a hand edit of front matter: a flag on `story new` and `task new` and a command to set or change it on an existing item; `flai show` and the item page display it
- [ ] The agent has something to calibrate against: `flai stats` (or a small addition to it) reports actual cycle time for done stories by nature, median and 85th percentile, and the convention tells the agent to read it before estimating
- [ ] Estimates already on items, if any are found, are reviewed and either rewritten in agent time or removed; the story's notes record what was found
- [ ] The estimate versus actual chart and the `estimate error` metric work with the new estimates; if the definition in `design/system/metrics.md` needs a word changed to say the estimate is agent time, that change comes with an ADR, as that file requires
- [ ] Tests for the flag, the command, and the stats output; `docs/users/flai.md` and the conventions guide updated

## Tasks

## Notes
Raised by the operator on 2026-09-18: "I would like the estimates to be based on an agent implementing them, their current estimates are likely based off some idealized human effort."

What was found when queuing this, which does not match the premise and should be settled with the operator before work starts: no work item in this repository, active or archived, has an `estimate` field. The schema allows one (`estimate: 4h`, optional, "used for estimate-vs-actual" in `work-hierarchy.md`), `metrics.md` defines estimate error, and the dashboard has an estimate versus actual chart and shows the field on the item page, but no `flai` command sets it and no convention asks anyone to write it. The only sizing guidance is that a story is "one to a few sessions". So there may be no current estimates to correct; the operator may be looking at the chart, at a forecast line, or at numbers from somewhere else. The first thing the pulling agent does is ask where the operator sees them.

For calibration: the stories accepted on 2026-09-18 (S-0048 to S-0051) each spent roughly ten to forty-five minutes in progress, with an agent doing the work, which is the scale an agent-time estimate lives on.

`design/system/metrics.md` is the contract between `flai stats` and the dashboard and changes only with an ADR.
