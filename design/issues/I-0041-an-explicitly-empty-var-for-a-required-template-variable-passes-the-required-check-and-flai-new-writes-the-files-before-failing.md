---
id: I-0041
title: An explicitly empty --var for a required template variable passes the required check, and flai new writes the files before failing
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-09-24T08:16:45Z
last_reported: 2026-09-24T08:16:45Z
updated: 2026-10-01T08:53:26Z
---

# I-0041 An explicitly empty --var for a required template variable passes the required check, and flai new writes the files before failing

## Description
An explicitly empty --var for a required template variable passes the required check, and flai new writes the files before failing

## Instances

### 2026-09-24T08:16:45Z
S-0018: flai new DIR --defaults --var project_name= wrote CLAUDE.md, design/, docs/ into DIR and then failed with rendered project has an invalid manifest: name is required. collectVars in cmd/new.go continues past the required check for a value given with --var.

## Remediation
Closed 2026-10-01T08:53:26Z: S-0185: flai new and flai import apply the required check to a value given with --var too, and stop before writing anything (TestNewRefusesAnEmptyRequiredVarBeforeWriting)
