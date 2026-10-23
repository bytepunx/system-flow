---
id: D-0234
title: "S-0324's close-out waits for main to commit four wip lines its fix now flags: Recommendation: the operator commits the two files on main with the command on TH-0355, which covers the same two files…"
kind: answer
sources: [TH-0354]
topics: [cli, flai, go, code, all]
paths: [flai/internal/mdlint/inline.go, flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/testdata/cases/www.md, flai/internal/mdlint/testdata/cases/expected.txt, design/issues/I-0110-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md, design/issues/summary.md, design/issues/I-0114-flai-verify-s-integration-tier-keeps-only-the-last-lines-of-go-test-s-output-so-the-failing-test-is-not-named.md, design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md, design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0234 S-0324's close-out waits for main to commit four wip lines its fix now flags: Recommendation: the operator commits the two files on main with the command on TH-0355, which covers the same two files…

The operator, alex, answered TH-0354 at 2026-10-08T04:20:12Z.

What was asked:

```text
Recommendation: commit the four fixed `wip/` files on main now (option 1), so I can sync and run the close-out again.
```

The answer:

```text
Recommendation: the operator commits the two files on main with the command on TH-0355, which covers the same two files as yours. One commit clears both S-0316 and S-0324. No flai acceptance can land them: review is empty, and I accept only a verified story (ADR-0093). I will tell you here when it lands.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md
```

Source: TH-0354.
