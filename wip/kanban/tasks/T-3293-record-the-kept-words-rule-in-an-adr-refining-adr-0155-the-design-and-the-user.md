---
id: T-3293
type: task
nature: remediation
title: Record the kept-words rule in an ADR refining ADR-0155, the design, and the user guide, and close I-0360
status: in-progress
parent: S-0767
owner: alex
created: 2026-11-13T11:48:34Z
updated: 2026-11-13T12:14:44Z
transitions:
  - to: ready
    at: 2026-11-13T12:14:44Z
    by: agent-S-0767
  - to: in-progress
    at: 2026-11-13T12:14:44Z
    by: agent-S-0767
stream: S-0767
tags: [docs, adr]
touches: [design/adrs/0297-a-markdown-line-a-rebased-branch-changes-keeping-every-word-of-main-s-line-in.md, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md]
after: [T-3291]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 8
          cache_read: 20199
          cache_write: 757
          cost: 0.0061
---
# T-3293 Record the kept-words rule in an ADR refining ADR-0155, the design, and the user guide, and close I-0360

## Work

Write the rule T-3291 built down where its readers look. It waits for T-3291, since it describes the rule as that task built it.

- Write an ADR with `flai adr new`, refining ADR-0155: in a markdown path, a line main added that the branch's net change removes is not undone when a line the branch adds keeps every word of it in order, each added line answering for one of main's lines; other paths keep the multiset of whole lines. Give I-0360's instance as its context and the commented-out line as why code paths are left out. Accept it with `flai adr accept`. The touch names the predicted number and slug; `flai task done` corrects them.
- In `design/system/flai-cli.md`, in the `flai stream sync` row, where it says the compare is "a multiset of line text", say what the new rule keeps and link the ADR.
- In `design/system/workflow.md` § Branches and collisions, in the S-0371 paragraph, say the same in a sentence.
- In `docs/users/flai.md`, in the "Main's change undone" bullet, say that a markdown line the branch extends, keeping every word of main's, is not reported.
- Close the issue: `flai issue close I-0360 --reason "<what fixed it, naming S-0767 and the ADR>"`.

## Done when

- [ ] The ADR is accepted, refines ADR-0155, and the ADR index lists it
- [ ] The three documents describe the rule and link the ADR
- [ ] I-0360 is closed with a reason naming what fixed it
- [ ] `flai check --strict` and the markdown lint pass on the changed files

## Notes
