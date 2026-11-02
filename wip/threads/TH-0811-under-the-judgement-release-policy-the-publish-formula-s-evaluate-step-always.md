---
id: TH-0811
title: "Under the judgement release policy, the publish formula's evaluate step always stops, and from flai 1.80.2 release_publish refuses: nothing can publish"
anchor:
  path: wip/kanban/stories/S-0562-the-publish-formula-is-run-once-in-this-repository-its-instance-worked-to.md
  item: S-0562
status: open
participants: [orchestrator]
created: 2026-11-02T05:32:39Z
updated: 2026-11-02T05:32:39Z
---

# TH-0811 Under the judgement release policy, the publish formula's evaluate step always stops, and from flai 1.80.2 release_publish refuses: nothing can publish

On wip/kanban/stories/S-0562-the-publish-formula-is-run-once-in-this-repository-its-instance-worked-to.md.

## Entries

### 2026-11-02T05:32:39Z orchestrator
**Recommended:** give this project a release policy that `--require-met` can meet. One way is `policy: threshold` beside the `value: 1000` that `orchestration.release` already has, with a `count: 1` if you want each accepted story released as now. The other way is a story that lets the formula's evaluate step pass under judgement once the orchestrator has judged the batch, for example a `--judged "<reason>"` flag.

**What I found while publishing S-0423 (flai 1.80.2, template 1.1.19):**

- `release_evaluate` gives policy `judgement` for this project. `orchestration.release` has only `value: 1000`, with no `policy` key.
- Under judgement, `flai release --evaluate --require-met` always exits 2, because judgement "is never met by itself". The publish formula's step 2 is that command. So every instance stops at step 2.
- T-2014 makes `release_publish` refuse while a publish formula exists. That change ships in 1.80.2. Once 1.80.2 is installed, I have no way left to publish under judgement.
- I published S-0423 with `release_publish` this time. The installed flai, 1.80.1, has neither the refusal nor `--require-met`.

S-0562's live run would also only show the evaluation stopping, and would release nothing.

**Alternative:** keep judgement, and publish by hand until a story settles it.
