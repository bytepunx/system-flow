---
id: I-0362
title: "flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone"
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-11-13T13:28:22Z
last_reported: 2026-11-13T13:28:22Z
updated: 2026-11-13T13:28:22Z
---

# I-0362 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone

## Description
flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone

## Instances

### 2026-11-13T13:28:22Z
Story: S-0661.
S-0661 built the checks run a story's agent asks for alone (kluster_checks, flai checks run S-nnnn --job). ADR-0292 §3's row "A story enters review otherwise ... or the checks run alone when only it has not passed at the head" is still not built for the review pass: lookReview in flai/internal/serve/kluster.go queues a full deploy, then checks, whenever passedAt fails, so a story whose deploy passed at its head and whose checks failed or never ran is built, loaded, and deployed again. No story names this work; design/system/kluster.md § What happens when records it as not built.

## Remediation
