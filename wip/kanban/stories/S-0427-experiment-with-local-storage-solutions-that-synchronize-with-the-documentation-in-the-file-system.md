---
id: S-0427
type: story
nature: experiment
title: Experiment with local storage solutions that synchronize with the documentation in the file system
status: backlog
owner: alex
created: 2026-10-16T05:55:07Z
updated: 2026-10-16T05:55:07Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0427 Experiment with local storage solutions that synchronize with the documentation in the file system

## Goal

As of now, flai relies heavily on reading front-matter and documents from across the repository in order to present critical information for both the operator and agents.

As projects grow larger, this overhead starts to show up in both agents and tools like the board and chart.

Lets look at what it might take to leverage a database instance that gets deterministically populated from the file system and kept up-to-date with the changes on the local system such that the agents and flai can perform more demanding read operations over a SQLite connection vs. reading from hundreds or thousands of files.

Don't worry about the overhead of introducing a pass-through storage adapter in addition to the file operations, focus primarily on the potential time savings and look at how this might reduce both operator wait time and agent overhead.

## Acceptance criteria
- [ ] Side-by-side comparisons of chart data (load a sql lite db then make comparative calls from flai's api vs. a throw-away script that reads the same broad information from Sqlite)
- [ ] Exploration of storage, RAM, and CPU overhead that might be added with a SQLite instance
- [ ] If there might be additional advantages to a postgres instance, spin one up via docker and perform the same kind of test and present the data alongside the other comparisons

## Tasks

## Notes
