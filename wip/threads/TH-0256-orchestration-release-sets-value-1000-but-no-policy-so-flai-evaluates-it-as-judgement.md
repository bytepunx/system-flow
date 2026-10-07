---
id: TH-0256
title: orchestration.release sets value 1000 but no policy, so flai evaluates it as judgement
anchor:
  path: system-flow.yaml
status: resolved
participants: [orchestrator, alex]
created: 2026-10-07T08:20:57Z
updated: 2026-10-07T08:29:11Z
---

# TH-0256 orchestration.release sets value 1000 but no policy, so flai evaluates it as judgement

On system-flow.yaml.

## Entries

### 2026-10-07T08:20:57Z orchestrator
Recommendation: add `policy: threshold` under `orchestration.release`, so that the value of 1000 you set takes effect and I publish when the accepted, unreleased stories add up to 1000 USD/week.

As it stands, `orchestration.release` has `value: 1000` and no `policy`. `flai release --evaluate` treats it as `judgement`, which is never met by itself.

Pending now: S-0212 (33.79 USD/week) and S-0274 (64.00 USD/week), 97.79 USD/week in all.

Until you answer, I read your value of 1000 as your intent to batch, and I do not publish under judgement. Answer with `judgement` instead if you want me to publish each coherent batch as it is accepted.

### 2026-10-07T08:27:51Z alex
judgement

### 2026-10-07T08:28:12Z orchestrator
Done. Published S-0212 and S-0274 under judgement and pushed:

| Component | From | To | Tag |
|-----------|------|----|-----|
| flai | 1.34.2 | 1.34.3 | flai/v1.34.3 |
| flaiover | 0.36.1 | 0.37.0 | flaiover/v0.37.0 |
| template | 1.0.67 | 1.0.68 | v1.0.68 (ffe11b9) |

From here I publish each accepted batch I judge coherent and complete, and never one `whole_epics` holds back. Resolve this thread when you are content.

### 2026-10-07T08:29:11Z alex
Resolved.
