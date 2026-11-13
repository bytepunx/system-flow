---
id: T-3062
type: task
nature: remediation
title: Document that an acceptance waits for a verify over main's new paths
status: done
parent: S-0683
owner: alex
created: 2026-11-10T10:35:37Z
updated: 2026-11-12T16:34:23Z
transitions:
  - to: ready
    at: 2026-11-12T16:18:49Z
    by: agent-S-0683
  - to: in-progress
    at: 2026-11-12T16:18:50Z
    by: agent-S-0683
  - to: done
    at: 2026-11-12T16:34:23Z
    by: agent-S-0683
stream: S-0683
tags: [docs, acceptance]
touches: [design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-3051, T-3058]
usage:
  source: log
  seconds: 933
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 63
      output: 18165
      cache_read: 6048449
      cache_write: 114098
      cost: 2.3889
---
# T-3062 Document that an acceptance waits for a verify over main's new paths

## Work

It waits for T-3051, whose ADR it links, and for T-3058, whose blocker and fields it names.

Say what the gate does, link the ADR, and keep the queue-tier wording as it is:

- `design/system/workflow.md` § What waits on an acceptance (S-0296): an acceptance by hand runs no tier. So every acceptance but the queue's landing waits while main's paths since the story's last verify send back a tier.
- `design/system/flai-cli.md` § The acceptance gate: the blocker `verify_stale`, its fields, and what clears it.
- `docs/users/flai.md`: the `flai accept` section, for the operator. Say what the refusal says and to ask for the close-out again.

## Done when

- The three documents describe the blocker and what clears it, and link the ADR.
- Each changed document's `updated` is bumped.
- `flai test` passes the markdown lint on the three files.

## Notes
