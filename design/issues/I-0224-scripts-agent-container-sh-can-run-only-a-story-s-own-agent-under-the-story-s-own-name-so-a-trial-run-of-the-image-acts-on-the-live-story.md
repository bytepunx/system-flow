---
id: I-0224
title: scripts/agent-container.sh can run only a story's own agent, under the story's own name, so a trial run of the image acts on the live story
class: impression
status: open
count: 1
cost: 5m
first_reported: 2026-10-28T23:10:55Z
last_reported: 2026-10-28T23:10:55Z
updated: 2026-10-28T23:10:55Z
---

# I-0224 scripts/agent-container.sh can run only a story's own agent, under the story's own name, so a trial run of the image acts on the live story

## Description
scripts/agent-container.sh can run only a story's own agent, under the story's own name, so a trial run of the image acts on the live story

## Instances

### 2026-10-28T23:10:55Z
Story: S-0507.
S-0507's T-2174 checked pilot by hand with scripts/agent-container.sh S-0507. pilot builds the prompt from the start request as flai serve does, so the only run the script can make is the story's agent working the story, named agent-S-0507. The capped trial (haiku, $0.02) called story_start, prime, inbox, and stream_state on the live story, overwriting its narrative's state, which the story's agent wrote again. A smoke run of an agent image needs a request that does no work on a live item, such as a probe prompt or a story of its own.

## Remediation
