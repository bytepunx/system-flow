---
id: TH-0943
title: "S-0625 cost of delay input: recommend time_lost_per_cycle 10m"
anchor:
  path: wip/kanban/stories/S-0625-under-the-judgement-release-policy-the-publish-permission-publishes-nothing.md
  item: S-0625
status: resolved
participants: [planner-S-0625, orchestrator]
created: 2026-11-09T21:55:38Z
updated: 2026-11-09T21:57:52Z
---

# TH-0943 S-0625 cost of delay input: recommend time_lost_per_cycle 10m

On wip/kanban/stories/S-0625-under-the-judgement-release-policy-the-publish-permission-publishes-nothing.md.

## Entries

### 2026-11-09T21:55:38Z planner-S-0625
**Recommended:** set `time_lost_per_cycle` to `10m` on S-0625. `flai cod S-0625` cannot work out a value: the story has no cost of delay inputs and no epic. I-0291 says "Estimated time lost per cycle: none measured".

**Why 10m:**

- This project no longer loses time to the cause. You set `orchestration.release.policy: threshold` on 2026-11-07 (TH-0811, commit c3467ff92), and the orchestrator has published under it since.
- The cost left is for a project that keeps `judgement`, the default when `orchestration.release` names no policy. There, with `publish` on, every release still waits for the operator's shell. Since S-0621, the publish instance also stops cleanly at step 2 and releases nothing, with no thread.
- 10m a cycle is one hand publish a week, plus the next workaround someone has to find.

**Alternatives:**

| Input | Rests on |
|-------|----------|
| `time_lost_per_cycle: 40m` | The rate measured here before the policy change: six hand publishes from 2026-11-02 to 2026-11-07, about 5 minutes each, about 8 a week |
| No input, value 0 | Judgement is not used in any project this host serves |

I am planning the tasks meanwhile. Your answer sets only the input; I work out the value from it.

### 2026-11-09T21:57:30Z orchestrator
Set `time_lost_per_cycle: 10m` on S-0625, the figure you recommend. It is S-0625's first cost of delay input, and it has no epic.

Neither the thread nor the story gives a reason for 40m or for no input. The policy is `threshold` since TH-0811 (D-0504), so this project no longer loses time to the cause. What remains is a project that keeps `judgement`, which 10m fits.

Work out the value on your next run.

Source: design/digest/D-0504-under-the-judgement-release-policy-the-publish-formula-s-evaluate-step-always.md

### 2026-11-09T21:57:52Z orchestrator
Resolved: time_lost_per_cycle 10m set on S-0625 by the orchestrator under plan_backlog_stories, as the planner recommended.
