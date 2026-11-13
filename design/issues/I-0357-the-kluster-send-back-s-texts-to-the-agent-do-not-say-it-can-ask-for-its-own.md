---
id: I-0357
title: The kluster send-back's texts to the agent do not say it can ask for its own deploy before review
class: impression
status: open
count: 1
cost: 5m
first_reported: 2026-11-13T08:42:12Z
last_reported: 2026-11-13T08:42:12Z
updated: 2026-11-13T08:42:12Z
---

# I-0357 The kluster send-back's texts to the agent do not say it can ask for its own deploy before review

## Description
The kluster send-back's texts to the agent do not say it can ask for its own deploy before review

## Instances

### 2026-11-13T08:42:12Z
Story: S-0660.
Found in S-0660's docs review: harness.SentBack's kluster text (flai/internal/harness/harness.go) and heldText in flai/internal/serve/kluster.go still say only that flai serve deploys the story again once it is back in review. Since S-0660 the agent of a story sent back can ask kluster_deploy (flai kluster deploy S-nnnn --job) before review, and a deploy and checks run it asks for that pass at the head resolve the send-back's threads. The agent is not told.

## Remediation
