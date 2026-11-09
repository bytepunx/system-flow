---
id: I-0242
title: The merge queue's send-back thread quotes go test's summary line with its tabs, the markdown lint refuses the thread (MD010), and the queue retries forever with the batch wedged
class: defect
status: closed
count: 1
cost: 5h
first_reported: 2026-10-31T16:43:40Z
last_reported: 2026-10-31T16:43:40Z
updated: 2026-11-09T10:08:07Z
---

# I-0242 The merge queue's send-back thread quotes go test's summary line with its tabs, the markdown lint refuses the thread (MD010), and the queue retries forever with the batch wedged

## Description
The merge queue's send-back thread quotes go test's summary line with its tabs, the markdown lint refuses the thread (MD010), and the queue retries forever with the batch wedged

## Instances

### 2026-10-31T16:43:40Z
Batch 20261031-110330 failed at smoke (the OpenRouter gateway test, a 402 on the account's balance) and was bisected to S-0546 at 11:15:56Z. SendBack opens a thread on the story quoting the failing tier's findings in a code block; the findings carry go test's summary line, FAIL<tab>github.com/bytepunx/system-flow/flai/internal/serve<tab>131.960s, and the thread writer refuses a body the markdown lint fails: MD010 hard tabs. So nothing was written, the story was not sent back, and flai serve retried at every look: 639 identical errors in the journal by 16:37Z, one every 30 s, with the batch at step bisected, S-0546 enqueued and in review, and no other story batched. Every review story's hand acceptance meanwhile was refused by the ADR-0221 gate, since no batch could run the queue tiers. Fixed in fenced (bisect.go) by writing tabs as four spaces, with TestTheFailingTiersOutputIsQuotedWithoutTabs; the queue was cleared by hand with flai queue remove S-0546 and held.

## Remediation

Story S-0614 remediates this issue, created from it at 2026-11-07T23:39:04Z.
Closed 2026-11-09T10:08:07Z: S-0614: send-back threads now quote tool output through mdlint.Fenced, which makes line ends, tabs, ANSI escapes and other control characters, trailing spaces, and blank runs lint-safe; the merge queue's failingText and the pull request send-back's quote both use it. When the lint still refuses the queue's thread, SendBack logs a warning and opens it again with the findings left out and the log folder named, so the story is sent back and the batch goes on instead of being retried at every look. TestARefusedSendBackThreadLeavesTheFindingsOut and TestAGoTestFailureIsSentBackUnderTheProjectsLint reproduce the case.
