---
id: T-1530
type: task
nature: research
title: Report the options to the operator on a thread and record the decision in the finding and in ADRs
status: backlog
parent: S-0380
owner: alex
created: 2026-10-13T17:46:11Z
updated: 2026-10-13T17:46:11Z
transitions: []
stream: S-0380
tags: [dashboard, cli]
touches: [design/system/design-system-tool.md, design/adrs, design/adrs/README.md]
after: [T-1529]
---
# T-1530 Report the options to the operator on a thread and record the decision in the finding and in ADRs

## Work

Open a thread on S-0380 with the recommendation first and each choice the operator must make, and hold until it is answered. Record the answer in the finding's Decision section and in one ADR per decision, made with `flai adr new`, whose file names are known only then. It waits for T-1529, whose recommendation it reports.

## Done when

- The thread is answered.
- The finding's Decision section and the ADRs record what was decided, and `design/adrs/README.md` lists the ADRs.

## Notes

- Layer 4 of 5. `design/adrs` is a folder touch because the ADRs' names are not known until they are made.
