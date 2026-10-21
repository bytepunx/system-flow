---
id: E-0027
type: epic
nature: feature
title: Rework main layout
status: backlog
owner: alex
created: 2026-10-21T06:54:07Z
updated: 2026-10-21T06:54:07Z
transitions: []
tags: []
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

## Notes
