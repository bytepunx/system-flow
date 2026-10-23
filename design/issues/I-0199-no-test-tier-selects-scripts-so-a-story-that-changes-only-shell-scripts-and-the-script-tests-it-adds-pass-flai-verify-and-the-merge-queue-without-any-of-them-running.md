---
id: I-0199
title: "No test tier selects scripts/**, so a story that changes only shell scripts, and the script tests it adds, pass flai verify and the merge queue without any of them running"
class: efficiency
status: open
count: 1
first_reported: 2026-10-23T07:23:24Z
last_reported: 2026-10-23T07:23:24Z
updated: 2026-10-23T07:23:24Z
---

# I-0199 No test tier selects scripts/**, so a story that changes only shell scripts, and the script tests it adds, pass flai verify and the merge queue without any of them running

## Description
No test tier selects scripts/**, so a story that changes only shell scripts, and the script tests it adds, pass flai verify and the merge queue without any of them running

## Instances

### 2026-10-23T07:23:24Z
Seen on S-0481 (verified at cbd98d1d, enqueued by the orchestrator): it changed scripts/flaiover-install.sh and added scripts/flaiover-install-test.sh, wired into scripts/flaiover-test.sh. flai verify ran only rebase, sync, narrative, check and markdown. The flaiover tier, which runs flaiover-test.sh, is all_only and selects flaiover/**, so neither flai verify, flai verify --all, nor the merge queue's batch ran the new test; its only evidence of passing is the story agent's account in the narrative. A tier selecting scripts/** (shellcheck plus the scripts' own *-test.sh), or adding scripts/** to the flaiover tier's paths, would close the gap.

## Remediation
