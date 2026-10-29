---
id: I-0223
title: An Edit whose old and new strings end in a space drops the space, so a protected .claude/settings.json write needs a second permission round to fix
class: efficiency
status: closed
count: 1
cost: 15m
first_reported: 2026-10-28T22:40:11Z
last_reported: 2026-10-28T22:40:11Z
updated: 2026-10-29T04:10:03Z
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
Closed 2026-10-29T04:10:03Z: Fixed by S-0539: permission_prompt refuses at once, before any thread and under auto-approve too, an Edit or a MultiEdit edit whose old string ends in a space or a tab while its new string does not, since Claude Code drops that space from the new string before flai sees the call. The refusal tells the agent to end both strings on another character, so it retries at once instead of writing the wrong text and asking the operator a second time. TestPermissionPromptRefusesAtOnceAnEditWhoseNewStringLostTheTrailingSpace reproduces S-0416's edit from TH-0712.
