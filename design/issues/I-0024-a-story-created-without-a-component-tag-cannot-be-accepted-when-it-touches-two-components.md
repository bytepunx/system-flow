---
id: I-0024
title: A story created without a component tag cannot be accepted when it touches two components
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-09-19T09:00:02Z
last_reported: 2026-09-19T09:00:02Z
updated: 2026-10-01T08:04:45Z
---

# I-0024 A story created without a component tag cannot be accepted when it touches two components

## Description
A story created without a component tag cannot be accepted when it touches two components

## Instances

### 2026-09-19T09:00:02Z
S-0062: the agent created S-0062, S-0063 and S-0064 with flai story new and no --tag; E-0006 has no tags either, so the release plan could not tell which component S-0062 delivers to and the operator's acceptance from the board was refused in the preview. Fixed by hand-adding tags. flai story new could warn when the parent epic has no component tag and none is given, or flai check could flag a story in review whose plan is ambiguous.

## Remediation

S-0174: `flai check` warns (`story.component-tag`) on an open story whose touches, and its open tasks', reach two or more components while no tag of its own or its epic's names one of them, and says which tag to add. The pending plan no longer drops such a story silently: `flai release --pending`, `publish.preview`, and the board's Publish banner name each accepted item they could not plan, with the reason.
Closed 2026-10-01T08:04:45Z: S-0174: flai check warns on the untagged story (story.component-tag), and the pending plan names what it cannot plan
