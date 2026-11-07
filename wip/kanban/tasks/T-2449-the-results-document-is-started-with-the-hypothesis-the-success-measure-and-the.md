---
id: T-2449
type: task
nature: experiment
title: The results document is started with the hypothesis, the success measure, and the host's preconditions, and the two proof stories carry harness container
status: done
parent: S-0565
owner: alex
created: 2026-11-03T01:24:13Z
updated: 2026-11-07T04:20:07Z
transitions:
  - to: ready
    at: 2026-11-07T03:43:19Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T03:43:20Z
    by: system-flow
  - to: done
    at: 2026-11-07T04:20:07Z
    by: agent-S-0565
stream: S-0565
tags: [docker]
touches: [design/experiments/S-0565-a-story-worked-end-to-end-by-flai-serve-on-harness-container.md]
usage:
  source: log
  seconds: 146
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 4413
      cache_read: 3406102
      cache_write: 29307
      cost: 1.0041
---
# T-2449 The results document is started with the hypothesis, the success measure, and the host's preconditions, and the two proof stories carry harness container

## Work

Start `design/experiments/S-0565-a-story-worked-end-to-end-by-flai-serve-on-harness-container.md` from `design/experiments/template.md`. Write the hypothesis (the `container` harness of S-0508, with S-0510's image layer and S-0517's and S-0520's mounts, runs a story's agent from `story_start` to review under `flai serve` as a host run does) and the success measure (the story's three criteria).

Check each precondition on the host and record it in the document, with what you ran and what it answered:

- the installed `flai` on `PATH`, which `flai serve` runs, is a release that holds S-0508, S-0510, S-0517, and S-0520;
- the images are built: `flai-agent-bedrock:local`, `flai-agent-claude-code:local`, and this repository's layer from `flai/images/system-flow/Dockerfile` (`make agent-images`, `docker image inspect`);
- the host's container harness settings are set (`flai serve agent harness container` with `--login`, the image, and the network). They are the operator's: `flai guard` refuses an agent `flai serve` started every command that writes the host configuration (ADR-0220);
- a container on this host reaches the host's HTTP MCP listener (I-0254).

Ask the operator on one thread on S-0565 for each precondition that does not hold, and wait for the answer before the next task.

Then take the two proof stories TH-0821 settled: a story that changes documentation or shell only, for criteria 1 and 2, and a story that changes Go, for criterion 3. Set `harness: container` on each with `flai edit S-nnnn --harness container`. Moving them to ready is the operator's or the orchestrator's, not yours.

## Done when

- The results document has its hypothesis, its success measure, and a table of the preconditions, each with the command and what it answered.
- Every precondition holds, or a thread on S-0565 asks the operator for it and is answered.
- Both proof stories carry `harness: container` on their agent, and the document names them by ID.

## Notes
