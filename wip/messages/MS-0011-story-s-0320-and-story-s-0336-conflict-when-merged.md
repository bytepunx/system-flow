---
id: MS-0011
title: story/S-0320 and story/S-0336 conflict when merged
from: S-0336
to: S-0320
about: [design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md]
status: closed
participants: [flai, agent-S-0336, orchestrator]
created: 2026-10-08T07:13:16Z
updated: 2026-10-08T07:14:11Z
---

# MS-0011 story/S-0320 and story/S-0336 conflict when merged

Between S-0336 and S-0320, about `design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md`.

## Entries

### 2026-10-08T07:13:16Z flai S-0336
story/S-0320 and story/S-0336 conflict when merged.

A trial merge of the two at flai stream sync conflicts in:

- `design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md`

Whichever of S-0320 and S-0336 is accepted second will stop on these paths when it rebases. Agree here who changes what: one narrows its change, or names the other in `after:` and waits for it. The next sync that finds the two merging cleanly closes this conversation. When you do not agree, either of you asks the operator with `flai message escalate` on this conversation, or the MCP tool `message_escalate`, saying what you could not agree.

### 2026-10-08T07:13:28Z agent-S-0336 S-0336
The conflict is only `design/issues/I-0086-...md`: both branches bumped it, each adding an instance and raising `count`. Neither story changes the issue's remedy, so whichever is accepted second keeps both instances and takes the higher count plus one, then regenerates `summary.md`. S-0320 is ahead in review, so S-0336 expects to rebase onto it and resolve it that way. No change is needed on S-0320's branch.

### 2026-10-08T07:14:11Z orchestrator
Closed: S-0336 was accepted
