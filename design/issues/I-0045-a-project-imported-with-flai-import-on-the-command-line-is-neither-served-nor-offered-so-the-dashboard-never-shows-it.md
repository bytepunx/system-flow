---
id: I-0045
title: A project imported with flai import on the command line is neither served nor offered, so the dashboard never shows it
class: defect
status: open
count: 1
cost: 30m
first_reported: 2026-09-26T05:17:35Z
last_reported: 2026-09-26T05:17:35Z
updated: 2026-09-26T05:17:35Z
---

# I-0045 A project imported with flai import on the command line is neither served nor offered, so the dashboard never shows it

## Description
A project imported with flai import on the command line is neither served nor offered, so the dashboard never shows it

## Instances

### 2026-09-26T05:17:35Z
flai import in ~/git/blog wrote system-flow.yaml but never registered the project with flai serve; with a manifest it is no longer an import candidate either, so it appeared nowhere in the dashboard's switcher, and restarting the dashboard, the host, and editing ~/.flai/config.json did not help. Worked around with flai dashboard in the blog.

## Remediation
