---
id: T-2493
type: task
nature: improvement
title: An ADR decides that a story's scoped check leaves out the agent.image notes only the host clears
status: done
parent: S-0582
owner: alex
created: 2026-11-07T00:30:33Z
updated: 2026-11-07T00:34:45Z
transitions:
  - to: ready
    at: 2026-11-07T00:34:11Z
    by: agent-S-0582
  - to: in-progress
    at: 2026-11-07T00:34:12Z
    by: agent-S-0582
  - to: done
    at: 2026-11-07T00:34:45Z
    by: agent-S-0582
stream: S-0582
tags: [adr, check]
touches: [design/adrs/0255-a-check-scoped-to-a-story-leaves-out-the-agent-image-findings-only-the-host.md, design/adrs/README.md]
usage:
  source: log
  seconds: 33
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 2186
      cache_read: 1034704
      cache_write: 20678
      cost: 0.4079
  strategic:
    - kind: orchestrator
      seconds: 6
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 13
          cache_read: 82505
          cache_write: 1711
          cost: 0.0203
---
# T-2493 An ADR decides that a story's scoped check leaves out the agent.image notes only the host clears

## Work

Propose the fix I-0264 asks for from its instance, as the story's goal says, before any code. Write the ADR with `flai adr new`, from I-0264's instance and ADR-0171's precedent. It decides:

- which `agent.image` findings a check scoped to a story leaves out when they are outside it: the notes that docker is not on `PATH` or could not say, and the warning that the layer's image is not built, since only the host clears them (Docker installed, `make agent-images`);
- that each stays inside, and is reported, when the story changes `system-flow.yaml` or the Dockerfile `agent.image.layer` names;
- that the warning for a layer that is not a file is unchanged, since a story's branch can fix it;
- that `--record-issues` records none of the left-out findings, and `flai check` without `--story` still reports them all.

Check the number against the ADRs on main when you write it: the touch names 0255, which another story may take first. Regenerate the index with `flai adr index`. This task waits for nothing; the code and docs tasks build on its decision.

## Done when

- The ADR exists, accepted or proposed as the story's agent judges, and names I-0264, ADR-0171, and ADR-0253.
- `design/adrs/README.md` lists it, and `flai check --strict` reports no `adr.index`.

## Notes
