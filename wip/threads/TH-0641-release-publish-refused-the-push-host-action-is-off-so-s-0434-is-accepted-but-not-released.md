---
id: TH-0641
title: "release_publish refused: the push host action is off, so S-0434 is accepted but not released"
anchor:
  path: wip/archive/kanban/stories/S-0434-flai-stats-reports-per-week-the-threads-opened-to-the-operator-per-completed-story-and-how-many-asked-what-a-digest-entry-already-answered-and-the-dashboard-charts-it.md
  item: S-0434
status: open
participants: [orchestrator]
created: 2026-10-24T22:56:44Z
updated: 2026-10-24T23:51:12Z
---

# TH-0641 release_publish refused: the push host action is off, so S-0434 is accepted but not released

On wip/archive/kanban/stories/S-0434-flai-stats-reports-per-week-the-threads-opened-to-the-operator-per-completed-story-and-how-many-asked-what-a-digest-entry-already-answered-and-the-dashboard-charts-it.md.

## Entries

### 2026-10-24T22:56:44Z orchestrator
**Recommended:** if you want S-0434 released, run this on the host, in the project. I will publish on the next acceptance, or when you answer here:

```sh
flai serve enable push
```

If you turned `push` off on purpose, say so here. I will then leave publishing to you until you turn it back on.

## What happened

S-0434 landed through merge queue batch 20261024-224219 at 22:56:21Z, with integration and smoke passing. release_evaluate lists it alone as pending, at 92.98 USD/week, under the judgement policy. I judged it complete and called release_publish. It refused:

```text
release_publish publishes as publish.run does, which needs the push host action: the host action "push" is not enabled for this project. On the host, in the project, run: flai serve enable push
```

The same call published flai 1.70.0 and template 1.0.115 at 21:38Z, so `push` was turned off after that. S-0435, the digest experiment, reads the measure S-0434 adds, so it is better measured once S-0434 is released.

### 2026-10-24T23:51:12Z orchestrator
Still refused after S-0435 landed at 23:50:56Z, which completes E-0025. The experiment releases nothing, so S-0434 is the only pending story (92.98 USD/week). release_publish gave the same answer:

```text
the host action "push" is not enabled for this project. On the host, in the project, run: flai serve enable push
```

I will not try again until you answer here or another story is accepted.
