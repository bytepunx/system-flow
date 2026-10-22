---
id: I-0164
title: The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501
class: impression
status: closed
count: 1
cost: 0m
first_reported: 2026-10-19T05:38:11Z
last_reported: 2026-10-19T05:38:11Z
updated: 2026-10-22T22:08:36Z
---

# I-0164 The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501

## Description
The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501

## Instances

### 2026-10-19T05:38:11Z
Story: S-0449.
S-0449's headless pass (167 renders, 182 clicks) saw no 500 once the harness answered a framed component's requests to the app's own server with 501 (ADR-0172). Fourteen components ask the app's server: AcceptConfirm, ActivityView, AnalyzerPanel, CancelConfirm, ImportPrompt, InboxView, Messages, NewAdrForm, NewItemForm, OrchestratorPanel, PlannerPanel, PublishBanner, Threads, and WorktreeUncommitted. Each now shows the harness's 501 error, or an empty list, where the app shows the server's answer, and an example cannot give it one. A later story could let an example name the answers to such requests, for the harness to give in place of the 501.

## Remediation

Story S-0476 remediates this issue, created from it at 2026-10-21T11:59:18Z.
Closed 2026-10-22T22:08:36Z: S-0476: an example may carry answers (ADR-0187), which the component harness gives a framed component's fetch and EventSource to the app's own server in place of the 501, and the Components page sends, edits, and saves them. The fourteen components carry examples with answers, and a headless pass rendered their 17 answered examples with no error and no 501. integration_test.go's harness test reproduces it: an answered fetch and EventSource, and an unmatched request still 501.
