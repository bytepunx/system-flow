---
id: T-2177
type: task
nature: research
title: The results document is finished with its recommendation and indexed, and agent-adapters.md and docker.md record what the spike settled
status: done
parent: S-0505
owner: alex
created: 2026-10-26T04:52:14Z
updated: 2026-10-28T12:34:42Z
transitions:
  - to: ready
    at: 2026-10-28T12:33:50Z
    by: agent-S-0505
  - to: in-progress
    at: 2026-10-28T12:33:52Z
    by: agent-S-0505
  - to: done
    at: 2026-10-28T12:34:42Z
    by: agent-S-0505
stream: S-0505
tags: [cli, docker]
touches: [design/experiments/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md, design/experiments/README.md, design/system/agent-adapters.md, design/tech/docker.md]
after: [T-2173]
usage:
  source: log
  seconds: 50
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 17
      output: 5686
      cache_read: 1074392
      cache_write: 24166
      cost: 0.4965
---
# T-2177 The results document is finished with its recommendation and indexed, and agent-adapters.md and docker.md record what the spike settled

## Work

Finish the results document from `design/experiments/template.md`: the hypothesis (ADR-0203's container harness is workable with Claude Code baked in: the login, the mount, the HTTP MCP server, and the guard each hold inside a container), the success measure (the story reaches review inside the container; each open fact has a measurement), what was done, the results from T-2159, T-2166, and T-2173, and the recommendation, adopt, adapt, or drop, naming for each limitation the story of E-0029 that carries it on (S-0506 for the HTTP server, S-0507 for pilot, S-0508 for the adapter, S-0510 for the tools the tiers need) and each issue filed. Add the document to `design/experiments/README.md`'s index. Add one paragraph to `design/system/agent-adapters.md` § A container harness of flai's own, after "What a spike would settle", saying what the spike settled with the results linked, and list the two images in `design/tech/docker.md` beside `flaiover`'s, marked local and unsigned. Waits for T-2173, the last measurement.

## Done when

`flai check --strict` passes with the results document validated as the experiment story's, the README indexes it, agent-adapters.md and docker.md carry their paragraph and rows, and the markdown lint passes.

## Notes

`flai check` validates an experiment story's results document (`flai/internal/experiment`); its front matter names the story.
