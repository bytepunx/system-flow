---
id: S-0559
type: story
nature: improvement
title: PageSidebar's link entries have no page using them since the charts page dropped its sidebar
status: backlog
owner: alex
created: 2026-11-02T00:04:42Z
updated: 2026-11-02T00:04:42Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4m
    by: flai
    at: 2026-11-02T00:04:42Z
---
# S-0559 PageSidebar's link entries have no page using them since the charts page dropped its sidebar

## Goal

This story remediates [I-0190](../../../design/issues/I-0190-pagesidebar-s-link-entries-have-no-page-using-them-since-the-charts-page-dropped-its-sidebar.md), "PageSidebar's link entries have no page using them since the charts page dropped its sidebar". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0190 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0190 is closed with `flai issue close I-0190 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0190. time_lost_per_cycle 4m: 5m per occurrence × 1 occurrence ÷ 1.4 cycles of 168h (first reported 2026-10-22T20:35:40Z, 10.1 days before this story).
