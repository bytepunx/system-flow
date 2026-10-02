---
id: E-0016
type: epic
nature: feature
title: Introduce Orchestration, Planning, and Analysis agents
status: backlog
owner: alex
created: 2026-10-02T10:51:49Z
updated: 2026-10-02T10:51:49Z
transitions: []
tags: [dashboard, cli]
topics: [orchestration, planning, analysis]
---
# E-0016 Introduce Orchestration, Planning, and Analysis agents

## Outcome

Three new agents get introduced to system-flow:

 - The planner:
   - capable of taking an epic and creating stories based on the initial understanding of the epic
   - examines stories to determine and enrich them with:
      - cost of delay estimates (if key information is supplied by the operator)
      - estimated time to completion and estimated delivery date (real world wall clock estimate given agent/model/configuration and story complexity, past timing data)
      - an accurate list of files anticipated in touch to help improve parallelism with safety
 - The orchestrator:
   - an agentic operator that manages projects with backlog and ready items
   - answers threads to keep work moving
   - decides when to publish releases
   - can optimize work through the system based on either maximized throughput or cost of delay
 - The analyzer:
   - an agentic analysis tool that:
      - identifies bottlenecks
      - identifies gaps between implementation and intent
      - identifies technical and security risks
      - creates stories for improvements and remediations with cost of delay pre-calculated based on impact

## Stories

## Notes
