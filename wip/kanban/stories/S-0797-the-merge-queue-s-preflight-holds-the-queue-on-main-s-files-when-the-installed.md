---
id: S-0797
type: story
nature: remediation
title: The merge queue's preflight holds the queue on main's files when the installed flai is older than the flai main's files are written for, and the hold cannot clear until a release is installed
status: backlog
owner: alex
created: 2026-11-16T00:41:48Z
updated: 2026-11-16T00:41:48Z
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
    time_lost_per_cycle: 3h
    by: flai
    at: 2026-11-16T00:41:48Z
---
# S-0797 The merge queue's preflight holds the queue on main's files when the installed flai is older than the flai main's files are written for, and the hold cannot clear until a release is installed

## Goal

This story remediates [I-0382](../../../design/issues/I-0382-the-merge-queue-s-preflight-holds-the-queue-on-main-s-files-when-the-installed.md), "The merge queue's preflight holds the queue on main's files when the installed flai is older than the flai main's files are written for, and the hold cannot clear until a release is installed". The issue recommends this solution:

Proposed, for the story to settle in an ADR: the preflight tells a finding the installed flai cannot judge from one in the files. When main's own flai sub-project is ahead of the running flai's commit, the preflight either runs its `flai check` with a flai built from main, as `scripts/flai.sh` does, or holds with words that name the cause, the installed flai's version against main's unreleased commits, and the remedy, publish and install the release, in place of telling the operator to fix files that are right. Consider also holding an acceptance that changes a format the running flai reads (role files, the workflow file, the manifest) until a flai that reads it is installed.

## Acceptance criteria
- [ ] The cause I-0382 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0382 is closed with `flai issue close I-0382 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0382. time_lost_per_cycle 3h carried over from I-0382's Impact section. Its Impact time_lost_per_cycle was taken rather than the 3h derived from its cost and count.
