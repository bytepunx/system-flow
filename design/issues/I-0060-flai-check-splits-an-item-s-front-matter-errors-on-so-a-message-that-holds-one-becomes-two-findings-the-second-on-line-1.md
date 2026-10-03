---
id: I-0060
title: "flai check splits an item's front-matter errors on \"; \", so a message that holds one becomes two findings, the second on line 1"
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-03T06:01:53Z
last_reported: 2026-10-03T06:01:53Z
updated: 2026-10-03T06:01:53Z
---

# I-0060 flai check splits an item's front-matter errors on "; ", so a message that holds one becomes two findings, the second on line 1

## Description
flai check splits an item's front-matter errors on "; ", so a message that holds one becomes two findings, the second on line 1

## Instances

### 2026-10-03T06:01:53Z
Story: S-0199.
check.oneItem splits Item.Validate's joined error on "; ". S-0199's planning messages held "; " and came out as two item.front-matter findings; they now use ": ". manifest/agent.go's role messages ("agent role %s sets nothing; ...") still split. Fix: Validate returns a list (a Problems type) that check ranges over.

## Remediation
