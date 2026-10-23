---
id: D-0231
title: "flai's MD034 does not report a bare `www.` literal, which markdownlint reports as a bare URL"
kind: cause
sources: [I-0110]
topics: [cli, flai, go, code, all]
paths: [flai/internal/mdlint/inline.go, flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/testdata/cases/www.md, flai/internal/mdlint/testdata/cases/expected.txt, design/issues/I-0110-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md, design/issues/summary.md, design/issues/I-0114-flai-verify-s-integration-tier-keeps-only-the-last-lines-of-go-test-s-output-so-the-failing-test-is-not-named.md, design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md, design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0231 flai's MD034 does not report a bare `www.` literal, which markdownlint reports as a bare URL

I-0110 was closed on 2026-10-08.

The cause:

```text
flai's MD034 does not report a bare `www.` literal, which markdownlint reports as a bare URL
```

The fix:

```text
S-0324: parseRange in flai/internal/mdlint/inline.go records a bare `www.` literal in out.urls through bareWww, as micromark takes it for markdownlint-cli2 0.20.0 (any case of `www.`, at the start or after a space or one of (*_~[], a domain with no underscore in its last two segments, outside link text and after no unclosed [), so md034 reports it. The www.md fixture settles the edges against markdownlint-cli2, and TestBareWwwOfI0110 reproduces the issue.
```

Source: I-0110.
