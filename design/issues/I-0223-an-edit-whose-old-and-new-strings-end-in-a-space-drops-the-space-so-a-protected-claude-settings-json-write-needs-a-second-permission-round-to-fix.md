---
id: I-0223
title: An Edit whose old and new strings end in a space drops the space, so a protected .claude/settings.json write needs a second permission round to fix
class: efficiency
status: open
count: 1
cost: 15m
first_reported: 2026-10-28T22:40:11Z
last_reported: 2026-10-28T22:40:11Z
updated: 2026-10-29T02:01:44Z
---

# I-0223 An Edit whose old and new strings end in a space drops the space, so a protected .claude/settings.json write needs a second permission round to fix

## Description
An Edit whose old and new strings end in a space drops the space, so a protected .claude/settings.json write needs a second permission round to fix

## Instances

### 2026-10-28T22:40:11Z
Story: S-0416.
S-0416 (T-1733): the story's agent edited the Edit hook's condition in both .claude/settings.json files with old and new strings ending in "exit 0; ". The written files read "exit 0;out=$(...)", the trailing space gone, and the hook tests failed. Each fix is another protected write through permission_prompt (TH-0712/TH-0713, then TH-0720). End such an edit's strings on a non-space character.

## Remediation

Story S-0539 remediates this issue, created from it at 2026-10-29T02:01:44Z.
