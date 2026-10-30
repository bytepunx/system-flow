---
id: D-0075
title: A story created without a component tag cannot be accepted when it touches two components
kind: cause
sources: [I-0024, TH-0759]
topics: [release, flai, cli, go, flaiover, dashboard, sveltekit, code, all, orchestration]
paths: [flai/internal/release, flai/cmd/release.go, flai/internal/hostapi, flai/internal/check, flaiover/src, flai/cmd/hostapi_reads_test.go, flai/cmd/push_test.go, flai/cmd/release_pending_test.go, flai/internal/pending/pending.go, flai/internal/preview/push.go, flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users, design/issues]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0075 A story created without a component tag cannot be accepted when it touches two components

I-0024 was closed on 2026-10-01.

The cause:

```text
A story created without a component tag cannot be accepted when it touches two components
```

The fix:

```text
S-0174: flai check warns on the untagged story (story.component-tag), and the pending plan names what it cannot plan
```

The same gap stops a release. On 2026-10-29 S-0544, an untagged flaiover story whose change touched both flai and flaiover, was accepted. `release_publish` left it out each time while `release_evaluate` still listed it as pending. The orchestrator found the reason in the publish warning (TH-0759):

```text
S-0544 is left out of the release: S-0544 touches flai, flaiover but no tag says which it delivers to; add a tag matching a project or pass --deliver
```

The operator answered: "yes, add the tag \"flaiover\"". The fix for such a story is a tag naming the project it delivers to, set with `flai edit S-nnnn --tags <project>` by the operator or a story's agent. `flai guard` refuses the orchestrator's `item_edit` of an item's tags (ADR-0060), so the orchestrator asks for the tag on a thread and does not set it itself.

Source: I-0024; TH-0759 (D-0428 records the operator's answer).
