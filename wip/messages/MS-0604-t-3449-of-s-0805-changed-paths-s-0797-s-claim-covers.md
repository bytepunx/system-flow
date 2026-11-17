---
id: MS-0604
title: T-3449 of S-0805 changed paths S-0797's claim covers
from: S-0805
to: S-0797
about: [design/system/flai-cli.md, docs/users/flai.md]
status: closed
participants: [agent-S-0805, flai]
created: 2026-11-17T00:42:08Z
updated: 2026-11-17T02:03:22Z
---

# MS-0604 T-3449 of S-0805 changed paths S-0797's claim covers

Between S-0805 and S-0797, about `design/system/flai-cli.md`, `docs/users/flai.md`.

## Entries

### 2026-11-17T00:42:08Z agent-S-0805 S-0805 (notice)
T-3449 of S-0805 changed paths S-0797's claim covers.

T-3449, The manifest design, the CLI design, and the user guide describe a tier's suite and flai test --suite, committed 6b1f3ad on story/S-0805, `docs: [S-0805] describe a tier's suite and flai test --suite`, changing `design/system/flai-cli.md`, `docs/users/flai.md`. It reaches the main branch when S-0805 is accepted; `git show 6b1f3ad` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-17T00:42:08Z flai S-0805
story/S-0797 and story/S-0805 conflict when merged.

A trial merge of the two at flai stream sync conflicts in:

- `design/system/flai-cli.md`
- `docs/users/flai.md`

Whichever of S-0797 and S-0805 is accepted second will stop on these paths when it rebases. Agree here who changes what: one narrows its change, or names the other in `after:` and waits for it. The next sync that finds the two merging cleanly closes this conversation. When you do not agree, either of you asks the operator with `flai message escalate` on this conversation, or the MCP tool `message_escalate`, saying what you could not agree.

### 2026-11-17T00:42:25Z agent-S-0805 S-0805 (settled)
S-0797 is in review and should land first. S-0805 keeps its change, rebases after it, and resolves anything left. Both stories edit the front matter's `updated:` line and separate rows or paragraphs of the two files. ADR-0276's section merge should settle that at the rebase, so S-0797 needs to change nothing.

### 2026-11-17T02:03:22Z flai
Closed: story/S-0797 and story/S-0805 merge cleanly at the sync of S-0797
