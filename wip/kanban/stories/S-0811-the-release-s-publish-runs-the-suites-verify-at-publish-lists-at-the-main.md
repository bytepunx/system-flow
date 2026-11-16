---
id: S-0811
type: story
nature: feature
title: The release's publish runs the suites verify_at.publish lists at the main branch's head before it tags, and publishes nothing when one fails
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:07Z
updated: 2026-11-16T23:42:07Z
transitions: []
tags: [flai]
topics: [verify, workflow]
after: [S-0806]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0811 The release's publish runs the suites verify_at.publish lists at the main branch's head before it tags, and publishes nothing when one fails

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §6. Today publishing runs no tier.

## Acceptance criteria

- [ ] With `verify_at.publish` set, the publish of what is pending runs those suites at the main branch's head before any tag or version file is written, with the output kept, and exits with its own code, naming the tier, when one fails.
- [ ] The publish formula stops on that code and its story says which suite failed.
- [ ] With no `publish` in the file, as on the baseline, publishing runs no tier, as before this story.
- [ ] The release design and the runbook describe it.

## Tasks

## Notes

- After S-0806: the publish reads verify_at.publish, which S-0806 parses.
