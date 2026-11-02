---
id: I-0251
title: A story branch rebased across a release brings back the changelog's old Unreleased heading and the lines the release moved, since union merge keeps both sides
class: defect
status: open
count: 1
cost: 15m
first_reported: 2026-11-02T03:57:28Z
last_reported: 2026-11-02T03:57:28Z
updated: 2026-11-02T03:57:28Z
---

# I-0251 A story branch rebased across a release brings back the changelog's old Unreleased heading and the lines the release moved, since union merge keeps both sides

## Description
A story branch rebased across a release brings back the changelog's old Unreleased heading and the lines the release moved, since union merge keeps both sides

## Instances

### 2026-11-02T03:57:28Z
Story: S-0518.
S-0518 added its line under `template/CHANGELOG.md`'s `## Unreleased`, which then sat below `## 1.1.14`. Release 1.1.17 on main renamed that section to `## 1.1.17` and put a new empty `## Unreleased` on top. When `flai stream sync` rebased story/S-0518, `merge=union` (ADR-0228) kept both sides, so the branch had two `## Unreleased` headings: the old one came back with S-0518's line and four lines already released under 1.1.17. flai check, task done, and the close-out at the branch head all passed, because the close-out had run before the release. Merge queue batch 20261102-035614 then failed at `markdown` with `template/CHANGELOG.md:26 MD024/no-duplicate-heading` (TH-0802), and S-0518 was sent back. Fixed by hand: the old block was dropped and S-0518's line moved under the top `## Unreleased`. A sync that finds a union-merged `CHANGELOG.md` with more than one `## Unreleased` could fold the branch's own lines into the top one, as it regenerates the ADR index and the issues summary.

## Remediation
