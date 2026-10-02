---
id: I-0057
title: flai check --strict stops a story's close-out on wip/ findings outside the story
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-10-02T16:13:36Z
last_reported: 2026-10-02T16:13:36Z
updated: 2026-10-02T16:13:36Z
---

# I-0057 flai check --strict stops a story's close-out on wip/ findings outside the story

## Description
flai check --strict stops a story's close-out on wip/ findings outside the story

## Instances

### 2026-10-02T16:13:36Z
2026-10-02: S-0191's close-out stopped at flai check --strict on one warning, threads.archived on TH-0067 (answered, its story S-0231 archived), which the story did not change and its agent may not resolve: the thread still asks the operator a question (point 5). S-0176 (TH-0032), S-0181 (14 warnings, none its own), and S-0187 (TH-0032) stopped the same way, and each went to review with the finding noted. TH-0056's answer: a finding outside the story is a note, and flai records it as an issue or bumps its count; until flai does, the story's agent records it here.

## Remediation
