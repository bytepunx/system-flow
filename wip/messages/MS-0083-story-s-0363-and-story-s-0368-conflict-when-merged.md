---
id: MS-0083
title: story/S-0363 and story/S-0368 conflict when merged
from: S-0363
to: S-0368
about: [design/adrs/README.md]
status: closed
participants: [flai, agent-S-0363, alex]
created: 2026-10-12T18:15:08Z
updated: 2026-10-12T18:20:02Z
---

# MS-0083 story/S-0363 and story/S-0368 conflict when merged

Between S-0363 and S-0368, about `design/adrs/README.md`.

## Entries

### 2026-10-12T18:15:08Z flai S-0363
story/S-0363 and story/S-0368 conflict when merged.

A trial merge of the two at flai stream sync conflicts in:

- `design/adrs/README.md`

Whichever of S-0363 and S-0368 is accepted second will stop on these paths when it rebases. Agree here who changes what: one narrows its change, or names the other in `after:` and waits for it. The next sync that finds the two merging cleanly closes this conversation. When you do not agree, either of you asks the operator with `flai message escalate` on this conversation, or the MCP tool `message_escalate`, saying what you could not agree.

### 2026-10-12T18:15:14Z agent-S-0363 S-0363
S-0368 is in review and goes first. S-0363 only appends ADR-0142's row to `design/adrs/README.md`; when S-0368 is accepted, S-0363's next sync resolves the index by keeping both rows. Nothing for S-0368 to change.

### 2026-10-12T18:20:02Z alex
Closed: S-0368 was accepted
