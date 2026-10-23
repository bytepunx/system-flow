---
id: D-0221
title: "Accept ADR-0115, the remedy for I-0076: Recommendation: accept ADR-0115, with one clarification first"
kind: answer
sources: [TH-0269]
topics: [cli, conventions, template, flai, go, code, all]
paths: [design/adrs, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/workflow.md, design/system/agent-coordination.md, design/system/continuous-improvement.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md, design/issues/summary.md, design/system/work-hierarchy.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0221 Accept ADR-0115, the remedy for I-0076: Recommendation: accept ADR-0115, with one clarification first

The operator, alex, answered TH-0269 at 2026-10-07T14:34:55Z.

What was asked:

```text
ADR-0115 (`design/adrs/0115-a-close-out-records-no-wip-overlap-and-notes-only-an-overlap-naming-its-story.md`, on story/S-0279) proposes the remedy for I-0076, from its 14 instances. Recommended: accept it as written.
```

The answer:

```text
Recommendation: accept ADR-0115, with one clarification first.

Why accept: it brings `wip.overlap` in line with ADR-0096. ADR-0096 §3 defines a story's claim, with its tasks narrowing a folder touch, and §2 already gives `wip.overlap` the shared-paths exception. Comparing claims, one finding per pair, removes the duplicates I-0076's instances show (story against task, folder against narrowed files). Recording no `wip.overlap` at close-out changes ADR-0085 §1–2. It does not lose the signal: the pull hold, the `overlapped` change, and the trial merge still report it.

The clarification: ADR-0096 §1 says a story in review keeps its claim for `flai check`'s `wip.overlap`. ADR-0115 as summarised compares "the claims of the stories in progress". Either keep stories in review in that comparison, or have ADR-0115 say it refines ADR-0096 §1 as well as ADR-0085.

This is a new decision, so the call is the operator's; I recommend and do not answer it.

Source: design/adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md § Decision
```

Source: TH-0269.
