---
id: T-0702
type: task
nature: feature
title: I-0057 counts close-outs stopped by flai check findings outside the story
status: done
parent: S-0191
owner: arobson
created: 2026-10-02T16:13:44Z
updated: 2026-10-02T16:13:52Z
transitions:
  - to: ready
    at: 2026-10-02T16:13:52Z
    by: agent-S-0191
  - to: in-progress
    at: 2026-10-02T16:13:52Z
    by: agent-S-0191
  - to: done
    at: 2026-10-02T16:13:52Z
    by: agent-S-0191
stream: S-0191
tags: [dashboard]
touches: [design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md, design/issues/summary.md]
after: [T-0701]
usage:
  source: log
  seconds: 0
  models: []
---
# T-0702 I-0057 counts close-outs stopped by flai check findings outside the story

## Work

The fresh verifier's close-out stopped at `flai check --strict` on one warning outside the story: `threads.archived` on TH-0067, answered on archived S-0231, which still asks the operator a question and so is not this agent's to resolve. TH-0056's answer makes such a finding a note and asks for it to be recorded as an issue, or counted on the one that exists; none does, so record I-0057 with this occurrence and the earlier ones (S-0176, S-0181, S-0187). It waits for T-0701, whose verifier run found it.

## Done when

- [x] I-0057 records the occurrence and the earlier ones, and passes the markdown lint.

## Notes
