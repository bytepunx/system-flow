---
id: I-0110
title: "flai's MD034 does not report a bare `www.` literal, which markdownlint reports as a bare URL"
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-10-07T08:55:00Z
last_reported: 2026-10-07T08:55:00Z
updated: 2026-10-08T00:33:50Z
---

# I-0110 flai's MD034 does not report a bare `www.` literal, which markdownlint reports as a bare URL

## Description
flai's MD034 does not report a bare `www.` literal, which markdownlint reports as a bare URL

## Instances

### 2026-10-07T08:55:00Z
Story: S-0265.
Found while making MD034 take bare email addresses (S-0265): flai/internal/mdlint/inline.go's parseRange skips a GFM `www.` literal (the c == 'w' case) without adding it to out.urls, so md034 never reports it, while markdownlint-cli2 0.20.0's micromark takes it as a literalAutolink and MD034 reports it. A thread entry or work item with a bare `www.example.com` would pass flai check and fail a close-out's markdown lint, as I-0056 did for email addresses. Left out of S-0265 to keep it to I-0056; the fix is to record the literal's position as the http case does, gated the same way on link text and an unclosed [, with a `www.` line in the email.md fixture or one of its own.

## Remediation

Story S-0324 remediates this issue, created from it at 2026-10-07T18:59:55Z.
Closed 2026-10-08T00:33:50Z: S-0324: parseRange in flai/internal/mdlint/inline.go records a bare `www.` literal in out.urls through bareWww, as micromark takes it for markdownlint-cli2 0.20.0 (any case of `www.`, at the start or after a space or one of (*_~[], a domain with no underscore in its last two segments, outside link text and after no unclosed [), so md034 reports it. The www.md fixture settles the edges against markdownlint-cli2, and TestBareWwwOfI0110 reproduces the issue.
