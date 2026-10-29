---
id: TH-0734
title: S-0417 renames the published template 1.1.4 heading to 1.1.5, so S-0416's released entry will read as 1.1.5
anchor:
  path: template/CHANGELOG.md
status: open
participants: [orchestrator]
created: 2026-10-29T05:49:10Z
updated: 2026-10-29T05:49:10Z
---

# TH-0734 S-0417 renames the published template 1.1.4 heading to 1.1.5, so S-0416's released entry will read as 1.1.5

On template/CHANGELOG.md.

## Entries

### 2026-10-29T05:49:10Z orchestrator
**Recommended:** once S-0417 lands, restore the `## 1.1.4 - 2026-10-29` heading over S-0416's one-line entry in a small follow-up. Put S-0417's entry under the heading the next publish writes. Until that fix lands, I hold the next release, which would add a second `## 1.1.5` heading.

## What changes

- The publish of flai 1.75.0 and template 1.1.4 (0c2e06cb3) wrote `## 1.1.4 - 2026-10-29` over S-0416's entry. template.yaml is at 1.1.4.
- S-0417's branch renames that heading to `## 1.1.5 - 2026-10-29` and adds its own detailed entry above S-0416's. template.yaml is not bumped.
- Once it lands, S-0416's released entry sits under a version not yet released. The next publish adds a `## 1.1.5` section of its own, so the file gets two 1.1.5 headings and no 1.1.4.

## Why I enqueued S-0417 anyway

- All six criteria are met, flai verify passed, and flai reports no blocker. Criterion 6 asks only that `template/CHANGELOG.md` says how a defined role starts, which it does.
- I accept on flai's blockers and the story's criteria. A heading defect in the changelog is yours to weigh.

## What fixes it

- A one-line restore of the 1.1.4 heading. It goes on a story of yours, or on the next story that touches `template/CHANGELOG.md`. The same file already holds an older duplicate `## 1.1.3` heading (2026-10-28) with S-0416's pre-written entry.
- If you would rather publish without waiting, answer "publish", and I release with the next acceptance.
