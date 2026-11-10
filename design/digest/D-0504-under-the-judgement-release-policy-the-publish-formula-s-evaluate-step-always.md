---
id: D-0504
title: "Under the judgement release policy, the publish formula's evaluate step always stops, and from flai 1.80.2 release_publish refuses: nothing can publish: I changed the policy to threshold instead of…"
kind: answer
sources: [TH-0811, TH-0256, I-0291, S-0625]
topics: [orchestration, cli, planning, template, conventions, release, all]
paths: [wip/kanban/stories/S-0562-the-publish-formula-is-run-once-in-this-repository-its-instance-worked-to.md, design/issues, design/issues/summary.md, design/issues/I-0283-a-publish-instance-with-nothing-to-release-stops-at-its-evaluate-step-down-the.md, design/issues/I-0291-under-the-judgement-release-policy-the-publish-permission-publishes-nothing.md, design/adrs/0273-under-judgement-the-orchestrator-s-judgement-reaches-the-publish-formula-as-an.md, flai/internal/release/evaluate.go]
date: 2026-11-07
status: absorbed
absorbed_by: "design/adrs/0273-under-judgement-the-orchestrator-s-judgement-reaches-the-publish-formula-as-an.md#context"
asked_again: []
cited: [S-0625, TH-0839, TH-0943]
---

# D-0504 Under the judgement release policy, the publish formula's evaluate step always stops, and from flai 1.80.2 release_publish refuses: nothing can publish: I changed the policy to threshold instead of…

This project's release policy is `threshold`, not `judgement`: `orchestration.release` in `system-flow.yaml` now has `policy: threshold` beside `value: 1000` and `count: 1`. This reverses the operator's answer on TH-0256 (D-0216), which kept `judgement`.

The operator, alex, answered TH-0811 at 2026-11-07T09:22:57Z: "I changed the policy to threshold instead of judgement".

Why it was asked: under `judgement`, `flai release --evaluate --require-met` always exits 2, since judgement is never met by itself, so the publish formula's step 2 stopped every instance. From flai 1.80.2 (T-2014) `release_publish` refuses while a publish formula exists, which left no way to publish under judgement. The orchestrator recommended `policy: threshold`, with `count: 1` to release each accepted story as before.

TH-0839 (D-0508) waited on this answer to word S-0562's criterion 1: with threshold, the instance can release for real.

## Since S-0625 (I-0291, ADR-0273)

The cause behind the question is fixed: judgement can publish again. I-0291, "Under the judgement release policy the publish permission publishes nothing", was closed on 2026-11-10 by S-0625. From its Remediation: "flai release --evaluate takes --judged "<reason>", which under judgement meets the policy while a story is accepted and not yet released and whole_epics holds nothing back; the publish formula declares the input judgement, which its evaluate step passes as --judged; and the orchestrator, under judgement, runs flai formula run publish --ready --autocommit --input judgement="<its reasoning>"." ADR-0273 records it.

The operator's answer stands: this project's policy stays `threshold` until the operator changes it. A later question whether to go back to `judgement` no longer needs the threshold workaround as its reason. `release_evaluate` and `release_publish` over MCP take no judgement and still answer as before under judgement (S-0625, T-2790, T-2793).

Sources: TH-0811, TH-0256, I-0291, S-0625.
