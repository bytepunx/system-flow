---
id: I-0047
title: A project imported with flai import on the command line is neither served nor offered, so the dashboard never shows it
class: defect
status: closed
count: 1
cost: 30m
first_reported: 2026-09-26T05:17:35Z
last_reported: 2026-09-26T05:17:35Z
updated: 2026-09-26T05:37:51Z
---

# I-0047 A project imported with flai import on the command line is neither served nor offered, so the dashboard never shows it

## Description
A project imported with flai import on the command line is neither served nor offered, so the dashboard never shows it

## Instances

### 2026-09-26T05:17:35Z
flai import in ~/git/blog wrote system-flow.yaml but never registered the project with flai serve; with a manifest it is no longer an import candidate either, so it appeared nowhere in the dashboard's switcher, and restarting the dashboard, the host, and editing ~/.flai/config.json did not help. Worked around with flai dashboard in the blog.

## Remediation

S-0120. `flai import` registers the project with `flai serve` when a `flai host` runs, and names `flai dashboard` when none does. `flai serve` serves the git repositories with a `system-flow.yaml` below the folders named for import that are not registered ([ADR-0045](../adrs/0045-a-system-flow-repository-under-a-folder-named-for-import-is-served-whether-or.md)), and `flai serve status` says why any there is not served. The `repo_url` prompt defaults to the origin remote as an https URL and refuses a value that is not a URL.
Closed 2026-09-26T05:37:51Z: S-0120: flai import registers with the host flai, flai serve serves the projects below import_roots, and flai serve status says why any is not served
