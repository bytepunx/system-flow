---
id: T-2495
type: task
nature: improvement
title: The docs and the close-out rule say a story's scoped check leaves out the agent.image notes only the host clears
status: done
parent: S-0582
owner: alex
created: 2026-11-07T00:30:57Z
updated: 2026-11-07T00:41:05Z
transitions:
  - to: ready
    at: 2026-11-07T00:35:14Z
    by: agent-S-0582
  - to: in-progress
    at: 2026-11-07T00:35:16Z
    by: agent-S-0582
  - to: done
    at: 2026-11-07T00:41:05Z
    by: agent-S-0582
stream: S-0582
tags: [docs, conventions, template]
touches: [docs/users/flai.md, design/system/flai-cli.md, design/system/workflow.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md]
after: [T-2493]
usage:
  source: log
  seconds: 349
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 40
      output: 9599
      cache_read: 4543024
      cache_write: 90790
      cost: 1.7911
  strategic:
    - kind: orchestrator
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 12
          cache_read: 82505
          cache_write: 1711
          cost: 0.0203
---
# T-2495 The docs and the close-out rule say a story's scoped check leaves out the agent.image notes only the host clears

## Work

Write T-2493's decision wherever the scoped check's leave-outs are listed, as S-0446 did for ADR-0171. It waits for T-2493, whose ADR it links; it shares no file with T-2494, so the two run together.

- `docs/users/flai.md`: the `--story` paragraph and the `--record-issues` paragraph under Check the repository, and the `agent.image` table under the container harness, which says when a finding is left out of a scoped check.
- `design/system/flai-cli.md`: the scoped check's leave-outs, and the project's layer image paragraph.
- `design/system/workflow.md`: the close-out's list of what the scoped check leaves out, where it names the others.
- `design/conventions/work-management.md` and `template/root/design/conventions/work-management.md`: the close-out rule's list of notes the scoped check leaves out, the same words in both.
- `template/CHANGELOG.md`: a patch entry for S-0582 naming the convention change and the ADR.

## Done when

- Each file above names the left-out `agent.image` findings and links the ADR.
- The two `work-management.md` files say the same, and `flai check --strict` and the markdown lint pass on them.

## Notes
