---
id: TH-0759
title: release_publish found nothing pending while release_evaluate lists S-0544 as unreleased
anchor:
  path: wip/archive/kanban/stories/S-0544-make-the-dashboard-and-server-versions-easy-to-see.md
  item: S-0544
status: open
participants: [orchestrator, alex]
created: 2026-10-29T18:06:21Z
updated: 2026-10-29T18:21:35Z
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

### 2026-10-29T18:14:49Z orchestrator
Update: the cause is now known. I published S-0418 as flai 1.78.0 and template 1.1.8. `release_publish` left S-0544 out again, with this warning:

```text
S-0544 is left out of the release: S-0544 touches flai, flaiover but no tag says which it delivers to; add a tag matching a project or pass --deliver
```

S-0544 has no tags, so flai cannot tell which project it delivers to. That is why `release --pending` found nothing for it, while `release --evaluate` still counts it.

Recommendation: give S-0544 the tag `flaiover`, then publish it as a flaiover patch. Or say that it needs no release.

### 2026-10-29T18:20:39Z alex
yes, add the tag "flaiover"

### 2026-10-29T18:21:35Z orchestrator
I could not add the tag. flai guard refused my `item_edit` on S-0544: the orchestrator edits an item only to finalize a draft or to give cost of delay inputs (ADR-0060).

Please add it yourself, in the main checkout or from the dashboard:

```bash
flai edit S-0544 --tags flaiover
```

Then reply here. I will publish S-0544 as a flaiover patch with `release_publish` when I am started again. Until then I am not retrying the publish, since it would leave S-0544 out again.
