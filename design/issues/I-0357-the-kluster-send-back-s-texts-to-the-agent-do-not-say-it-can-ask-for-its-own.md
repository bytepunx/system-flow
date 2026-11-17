---
id: I-0357
title: The kluster send-back's texts to the agent do not say it can ask for its own deploy before review
class: impression
status: closed
count: 1
cost: 5m
first_reported: 2026-11-13T08:42:12Z
last_reported: 2026-11-13T08:42:12Z
updated: 2026-11-17T09:15:15Z
---

# I-0357 The kluster send-back's texts to the agent do not say it can ask for its own deploy before review

## Description
The kluster send-back's texts to the agent do not say it can ask for its own deploy before review

## Instances

### 2026-11-13T08:42:12Z
Story: S-0660.
Found in S-0660's docs review: harness.SentBack's kluster text (flai/internal/harness/harness.go) and heldText in flai/internal/serve/kluster.go still say only that flai serve deploys the story again once it is back in review. Since S-0660 the agent of a story sent back can ask kluster_deploy (flai kluster deploy S-nnnn --job) before review, and a deploy and checks run it asks for that pass at the head resolve the send-back's threads. The agent is not told.

## Remediation

Story S-0831 remediates this issue, created from it at 2026-11-17T08:17:25Z.
Closed 2026-11-17T09:15:15Z: S-0831: harness.SentBack.then for a kluster send-back, the agent's start prompt, and heldText's sent-back branch in flai/internal/serve/kluster.go, which the deploy and checks failure threads end with, now say that the agent may ask flai serve to deploy the story at its committed head before review with kluster_deploy or flai kluster deploy S-nnnn --job, and that back in review flai serve deploys it again only when its deploy, and its checks while checks is on, have not passed at its head. TestAnAgentItsKlusterSentBackIsToldItMayDeployBeforeReview and the failure-thread tests in serve hold the texts to it.
