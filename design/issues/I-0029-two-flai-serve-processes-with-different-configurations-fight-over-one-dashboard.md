---
id: I-0029
title: Two flai serve processes with different configurations fight over one dashboard
class: defect
status: open
count: 1
cost: 25m
first_reported: 2026-09-20T13:50:06Z
last_reported: 2026-09-20T13:50:06Z
updated: 2026-09-20T13:50:06Z
---

# I-0029 Two flai serve processes with different configurations fight over one dashboard

## Description
Two flai serve processes with different configurations fight over one dashboard

## Instances

### 2026-09-20T13:50:06Z
2026-09-20: after releasing S-0078 the dashboard said push was not enabled and that push.run was missing, although the tree's flai serve had both. A second flai serve, the installed flai 1.6.1 under the home configuration, had had this project registered since 13:01Z with the same agent credential. A dashboard keeps only the newest proven connection, so each replaced the other every few seconds: 457 connections in 45 minutes in the installed one's log, and whichever held the connection at the moment answered, with its own version, methods, and host actions. A write in flight during a swap can be lost or retried on the other process. One flai serve per user is only kept per configuration folder. Resolved for now by stopping the installed one on the operator's word. Remedy to consider: a dashboard refuses a second flai for the same project while one is healthy instead of replacing it, or flai serve backs off for good when it is replaced by a connection it did not lose, and says so in flai serve status.

## Remediation
