---
id: TH-0759
title: release_publish found nothing pending while release_evaluate lists S-0544 as unreleased
anchor:
  path: wip/archive/kanban/stories/S-0544-make-the-dashboard-and-server-versions-easy-to-see.md
  item: S-0544
status: open
participants: [orchestrator]
created: 2026-10-29T18:06:21Z
updated: 2026-10-29T18:06:21Z
---

# TH-0759 release_publish found nothing pending while release_evaluate lists S-0544 as unreleased

On wip/archive/kanban/stories/S-0544-make-the-dashboard-and-server-versions-easy-to-see.md.

## Entries

### 2026-10-29T18:06:21Z orchestrator
Recommendation: check why `flai release --pending` and `flai release --evaluate` disagree about S-0544, then publish it, or say it needs no release.

Under the `judgement` policy I judged S-0544 coherent and complete on its own. It is a self-contained flaiover improvement with no epic, both criteria are met, and you accepted it. So I called `release_publish`, and it refused:

```text
flai release --pending published nothing: nothing was pending
```

`release_evaluate` still answers `count: 1`, `pending: [S-0544]`, `released: false`, `value: 0` (unvalued). One of two things looks true: the change touches only `flaiover/src` and docs, which no release unit picks up as pending; or the two commands read pending work differently, which would be a defect in flai. To fix it, either release flaiover by hand, or record that S-0544 needs no release so the evaluation stops listing it. If the commands disagree for no reason, file it as an issue.

I will not try again until you answer this thread or another story is accepted.
