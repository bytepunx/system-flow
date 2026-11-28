---
id: T-3600
type: task
nature: improvement
title: Accept ADR-0316 as superseding ADR-0240 and ADR-0273, and ask on S-0811 for its publish formula criterion to name release_publish
status: done
parent: S-0837
owner: alex
created: 2026-11-17T12:37:23Z
updated: 2026-11-27T23:57:47Z
transitions:
  - to: ready
    at: 2026-11-27T23:57:37Z
    by: agent-S-0837
  - to: in-progress
    at: 2026-11-27T23:57:38Z
    by: agent-S-0837
  - to: done
    at: 2026-11-27T23:57:47Z
    by: agent-S-0837
stream: S-0837
tags: [docs]
touches: [design/adrs/0316-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md, design/adrs/0240-under-publish-the-orchestrator-starts-the-baseline-publish-formula-whose-main.md, design/adrs/0273-under-judgement-the-orchestrator-s-judgement-reaches-the-publish-formula-as-an.md, design/adrs/README.md]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 1
      output: 183
      cache_read: 47285
      cache_write: 1437
      cost: 0.0217
---
# T-3600 Accept ADR-0316 as superseding ADR-0240 and ADR-0273, and ask on S-0811 for its publish formula criterion to name release_publish

## Work

This is the story's first task, as its Goal asks. Every other task waits for it, because they build what the accepted ADR says.

- Ask the operator on a thread on S-0837 whether ADR-0316 stands as written. Say which open point you settle in it: whether `release_publish` keeps `--judged` under `judgement` (Consequences, third point). Record any change the operator makes in the ADR while it is still `proposed`.
- ADR-0316 lists ADR-0240 and ADR-0273 under `refines`. The story's first criterion asks that both name ADR-0316 as what supersedes them. Move them from `refines` to `supersedes` in ADR-0316, keep ADR-0094 under `refines`, and accept it with `flai adr accept 316`.
- Set `superseded_by: [ADR-0316]` on ADR-0240 and ADR-0273. That field is the one edit an accepted ADR allows. `flai adr accept` does not set it, so write the field by hand and change nothing else in either file. `flai adr index` regenerates `design/adrs/README.md`.
- S-0811 is finalized and in ready, and its second criterion names the publish formula. Open a thread on S-0811 to the operator, or to the orchestrator if it may edit it. Propose the new words: "`release_publish` and the dashboard's Publish refuse on that code, and the orchestrator's activity log says which suite failed". Do not edit S-0811 yourself.

## Done when

- ADR-0316 is `accepted`, with ADR-0240 and ADR-0273 under `supersedes`.
- ADR-0240 and ADR-0273 carry `superseded_by: [ADR-0316]` and are otherwise unchanged.
- The ADR index is regenerated, and `flai check --strict` passes.
- A thread on S-0811 asks for its criterion to be rewritten.

## Notes

- Drafted by planner-S-0837.
