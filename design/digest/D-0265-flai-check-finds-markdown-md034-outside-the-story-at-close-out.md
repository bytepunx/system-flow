---
id: D-0265
title: "flai check finds `markdown.MD034` outside the story at close-out"
kind: cause
sources: [I-0118]
topics: [cli, flai, go, code, all]
paths: [flai/internal/mdlint/inline.go, flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/testdata/cases/www.md, flai/internal/mdlint/testdata/cases/expected.txt, design/issues/I-0110-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md, design/issues/summary.md, design/issues/I-0114-flai-verify-s-integration-tier-keeps-only-the-last-lines-of-go-test-s-output-so-the-failing-test-is-not-named.md, design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md, design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0265 flai check finds `markdown.MD034` outside the story at close-out

I-0118 was closed on 2026-10-08.

The cause:

```text
flai check finds `markdown.MD034` outside the story at close-out
```

The fix:

```text
All five instances were bare `www.` literals that an installed flai older than 1.39.4 wrote in wip, before S-0324 gave flai's MD034 the bare `www.` rule (I-0110). Every write an agent makes to wip now passes through LintGuard, which refuses one. S-0346 fixed the one write flai makes from text no agent checks: a strategic run's end entry now puts the bare URLs in its summary in a code span (mdlint.QuoteBareURLs, called from logRunEndSaying). TestARunEndQuotesABareURLInItsSummary reproduces the refusal. flai check --strict on main finds no markdown.MD034.
```

Source: I-0118.
