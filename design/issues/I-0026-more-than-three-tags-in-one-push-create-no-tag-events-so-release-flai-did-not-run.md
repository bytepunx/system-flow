---
id: I-0026
title: More than three tags in one push create no tag events, so release flai did not run
class: defect
status: open
count: 1
cost: 15m
first_reported: 2026-09-20T12:14:09Z
last_reported: 2026-09-20T12:14:09Z
updated: 2026-09-20T12:14:09Z
---

# I-0026 More than three tags in one push create no tag events, so release flai did not run

## Description
More than three tags in one push create no tag events, so release flai did not run

## Instances

### 2026-09-20T12:14:09Z
2026-09-20: S-0073, S-0074 and S-0075 were accepted together and one git push carried main and six tags (flai 1.5.1 to 1.5.3, flaiover 0.20.0 to 0.22.0). GitHub creates no events for tags when more than three are pushed at once, so release flai, which triggers on the tag, never ran: no GitHub release exists for flai 1.5.1, 1.5.2 or 1.5.3, while the flaiover 0.22.0 image, which needs the host methods of 1.5.3, was published because its workflow triggers on main. flai push sends every pending ref in one git push (cmd/push.go). Remedy: push the branch, then the tags three at a time.

## Remediation
