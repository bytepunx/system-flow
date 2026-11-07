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
cited: []
---

# D-0504 Under the judgement release policy, the publish formula's evaluate step always stops, and from flai 1.80.2 release_publish refuses: nothing can publish: I changed the policy to threshold instead of…

The operator, alex, answered TH-0811 at 2026-11-07T09:22:57Z.

What was asked:

```text
**Recommended:** give this project a release policy that `--require-met` can meet. One way is `policy: threshold` beside the `value: 1000` that `orchestration.release` already has, with a `count: 1` if you want each accepted story released as now. The other way is a story that lets the formula's evaluate step pass under judgement once the orchestrator has judged the batch, for example a `--judged "<reason>"` flag.
```

The answer:

```text
I changed the policy to threshold instead of judgement
```

Source: TH-0811.
