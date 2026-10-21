---
id: E-0027
type: epic
nature: feature
title: Rework main layout
status: in-progress
owner: alex
created: 2026-10-21T06:54:07Z
updated: 2026-10-21T07:55:01Z
transitions:
  - to: ready
    at: 2026-10-21T06:54:13Z
    by: alex
  - to: in-progress
    at: 2026-10-21T07:55:01Z
    by: agent-S-0458
tags: []
usage:
  source: sum
  seconds: 2810
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-21
      ceremony: 7
      hand_edits: 3
      work: 63
  models:
    - model: claude-opus-5-5
      input: 458
      output: 142004
      cache_read: 23560091
      cache_write: 888923
      cost: 12.4458
  strategic:
    - kind: orchestrator
      seconds: 932
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 122
          output: 1847
          cache_read: 11706505
          cache_write: 46523
          cost: 2.9468
        - model: claude-sonnet-5-5
          input: 12
          output: 76
          cache_read: 205218
          cache_write: 57700
          cost: 0.2152
---
# E-0027 Rework main layout

## Outcome

Currently, the top nav orientation creates some usability issues in how the operator navigates between sections and content in flaiover.

Lets work on incorporating the following layout patterns in flaiover:

## Left-aligned Nav Rail

- The top bar should show flaiover, the project picker, the theme setting, and status pills along the top
- Navigation rail separated into sections grouped by concern:
  - Flow
    - Board
    - Inbox
    - Charts
    - Workflow (coming soon)
  - Agents
    - Orchestration
    - Planning
    - Analysis
    - Activity
    - Threads
    - Messages
  - Design
    - System
    - Components
  - Documentation
    - System
    - Conventions
    - ADRs
    - Issues
    - Search
  - Configuration
    - Features
    - Agents
    - File System
    - Projects
    - Harnesses

Each grouping should expand or collapse to show the icons beneath it. Highlight the icon for each of the expanded groups, dim the icon for collapsed groups.

## Side bars

For many of these selections, the page should contain a left aligned sidebar to allow the operator to navigate within the page. Where possible clicking the sidebar should either display a pane/page or, if appropriate, scroll to the active area. We want to minimize the amount of scrolling necessary for a user.

For example - in the story's page, the user has to scroll through a lot of content to find what they're looking for. A side bar lets them find the content they want quickly without the additional noise that exists when all information about the story is visible and requires scrolling.

## Easier to Navigate Board

### Mobile

When in especially constrained widths, each lane's heading should remain visible, but the lane should be collapsed where clicking/tapping the lane heading should toggle the lane's collapsed state.

### Desktop

When a wider layout is possible, do not stack lanes vertically! Instead, allow the user to horizontally scroll between lanes.

## Stories
- S-0458 The dashboard's pages are reached from a left nav rail of five collapsible groups, and the top bar holds only the name, the project picker, the theme, and the status pills
- S-0459 The settings page becomes five Configuration pages, Features, Agents, File System, Projects, and Harnesses, each reached from the rail
- S-0460 A page may carry a left sidebar of its sections, and the item page shows one section at a time from it instead of one long scroll
- S-0461 The review, charts, and documents pages navigate within themselves from the left sidebar
- S-0462 The board's lanes sit side by side and scroll sideways on a wide screen, and collapse to their headings, toggled by a tap, on a narrow one

## Notes
- 2026-10-21T07:55:01Z: moved to in-progress: follows S-0458, which moved to in-progress
