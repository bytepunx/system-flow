---
id: D-0504
title: "Under the judgement release policy, the publish formula's evaluate step always stops, and from flai 1.80.2 release_publish refuses: nothing can publish: I changed the policy to threshold instead of…"
kind: answer
sources: [TH-0811]
topics: [orchestration, cli, planning, template, conventions, all]
paths: [wip/kanban/stories/S-0562-the-publish-formula-is-run-once-in-this-repository-its-instance-worked-to.md, design/issues, design/issues/summary.md, design/issues/I-0283-a-publish-instance-with-nothing-to-release-stops-at-its-evaluate-step-down-the.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [S-0625, TH-0839, TH-0943]
---

# D-0504 Under the judgement release policy, the publish formula's evaluate step always stops, and from flai 1.80.2 release_publish refuses: nothing can publish: I changed the policy to threshold instead of…

This project's release policy is `threshold`, not `judgement`: `orchestration.release` in `system-flow.yaml` now has `policy: threshold` beside `value: 1000` and `count: 1`. This reverses the operator's answer on TH-0256 (D-0216), which kept `judgement`.

The operator, alex, answered TH-0811 at 2026-11-07T09:22:57Z: "I changed the policy to threshold instead of judgement".

Why it was asked: under `judgement`, `flai release --evaluate --require-met` always exits 2, since judgement is never met by itself, so the publish formula's step 2 stopped every instance. From flai 1.80.2 (T-2014) `release_publish` refuses while a publish formula exists, which left no way to publish under judgement. The orchestrator recommended `policy: threshold`, with `count: 1` to release each accepted story as before.

TH-0839 (D-0508) waited on this answer to word S-0562's criterion 1: with threshold, the instance can release for real.

Sources: TH-0811, TH-0256.
