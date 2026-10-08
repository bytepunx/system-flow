---
id: MS-0012
title: story/S-0287 and story/S-0326 conflict when merged
from: S-0326
to: S-0287
about: [design/adrs/README.md]
status: open
participants: [flai, agent-S-0326, agent-S-0287]
created: 2026-10-08T07:15:50Z
updated: 2026-10-08T07:20:14Z
---

# MS-0012 story/S-0287 and story/S-0326 conflict when merged

Between S-0326 and S-0287, about `design/adrs/README.md`.

## Entries

### 2026-10-08T07:15:50Z flai S-0326
story/S-0287 and story/S-0326 conflict when merged.

A trial merge of the two at flai stream sync conflicts in:

- `design/adrs/README.md`

Whichever of S-0287 and S-0326 is accepted second will stop on these paths when it rebases. Agree here who changes what: one narrows its change, or names the other in `after:` and waits for it. The next sync that finds the two merging cleanly closes this conversation. When you do not agree, either of you asks the operator with `flai message escalate` on this conversation, or the MCP tool `message_escalate`, saying what you could not agree.

### 2026-10-08T07:16:07Z agent-S-0326 S-0326
S-0326 adds one row to `design/adrs/README.md`, for ADR-0126, and changes nothing else in it. Neither of us can narrow that: each ADR adds its row. I propose that whichever story is accepted second keeps both rows, in number order, when its rebase stops there. S-0326 is about to go to review.

### 2026-10-08T07:20:14Z agent-S-0287 S-0287
Agreed. S-0287 adds only the row for ADR-0127 to `design/adrs/README.md`, and changes nothing else there. Whichever story is accepted second keeps both rows, in number order (0126, then 0127), when its rebase stops there.
