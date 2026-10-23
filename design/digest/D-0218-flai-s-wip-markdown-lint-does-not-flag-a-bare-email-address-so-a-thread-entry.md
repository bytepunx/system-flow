---
id: D-0218
title: flai's wip markdown lint does not flag a bare email address, so a thread entry with one reached main and failed a story's close-out
kind: cause
sources: [I-0056]
topics: [cli, flai, go, code, all]
paths: [flai/internal/mdlint/inline.go, flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/testdata/cases, docs/users/flai.md, design/issues/I-0056-flai-s-wip-markdown-lint-does-not-flag-a-bare-email-address-so-a-thread-entry-with-one-reached-main-and-failed-a-story-s-close-out.md, design/issues/summary.md, design/issues/I-0110-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md, design/issues/I-0111-flai-check-finds-narrative-state-outside-the-story-at-close-out.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0218 flai's wip markdown lint does not flag a bare email address, so a thread entry with one reached main and failed a story's close-out

I-0056 was closed on 2026-10-07.

The cause:

```text
flai's wip markdown lint does not flag a bare email address, so a thread entry with one reached main and failed a story's close-out
```

The fix:

```text
flai's MD034 now reports GFM's extended email autolink, a bare address, as markdownlint-cli2 0.20.0 does, skipping it in link text, after an unclosed [, and where a text directive takes the name after a colon; tested by the email.md fixture against markdownlint and TestBareEmailOfI0056 (S-0265, T-0981).
```

Source: I-0056.
