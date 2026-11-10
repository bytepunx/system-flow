---
id: I-0261
title: A story's prime pack came in 22 parts, most of each part the same topic sources list, so reading every part as the start prompt asks would take about 590 KB of context
class: efficiency
status: open
count: 2
cost: 10m
first_reported: 2026-11-02T19:44:22Z
last_reported: 2026-11-07T09:46:28Z
updated: 2026-11-10T08:55:29Z
---

# I-0261 A story's prime pack came in 22 parts, most of each part the same topic sources list, so reading every part as the start prompt asks would take about 590 KB of context

## Description
A story's prime pack came in 22 parts, most of each part the same topic sources list, so reading every part as the start prompt asks would take about 590 KB of context

## Instances

### 2026-11-02T19:44:22Z
Story: S-0425.
S-0425: story_start answered a pack of 591,634 bytes in 22 parts, exceeded at conventions. Each part repeats the `topics` list with its `sources`, one entry per claimed file for each of the story's and tasks' 50 touches, under every topic (cli, flai, go, dashboard, flaiover, sveltekit...), about 25 KB of each 40 KB part. Part 1 held two conventions and the README. The agent read part 1, then read the conventions from their files instead of the other 21 parts.

### 2026-11-07T09:46:28Z
Story: S-0586.
S-0586's pack came in 14 parts, 359 KB, each part repeating the topic sources list (about 20 KB of claim paths per topic). I read it once as text with flai prime --story into a git-ignored file instead of 14 prime calls.

## Remediation

Story S-0699 remediates this issue, created from it at 2026-11-10T08:55:29Z.
