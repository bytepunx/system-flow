---
id: I-0224
title: scripts/agent-container.sh can run only a story's own agent, under the story's own name, so a trial run of the image acts on the live story
class: impression
status: closed
count: 1
cost: 5m
first_reported: 2026-10-28T23:10:55Z
last_reported: 2026-10-28T23:10:55Z
updated: 2026-11-10T09:25:18Z
---

# I-0224 scripts/agent-container.sh can run only a story's own agent, under the story's own name, so a trial run of the image acts on the live story

## Description
scripts/agent-container.sh can run only a story's own agent, under the story's own name, so a trial run of the image acts on the live story

## Instances

### 2026-10-28T23:10:55Z
Story: S-0507.
S-0507's T-2174 checked pilot by hand with scripts/agent-container.sh S-0507. pilot builds the prompt from the start request as flai serve does, so the only run the script can make is the story's agent working the story, named agent-S-0507. The capped trial (haiku, $0.02) called story_start, prime, inbox, and stream_state on the live story, overwriting its narrative's state, which the story's agent wrote again. A smoke run of an agent image needs a request that does no work on a live item, such as a probe prompt or a story of its own.

## Remediation

Story S-0682 remediates this issue, created from it at 2026-11-10T08:37:06Z.
Closed 2026-11-10T09:25:18Z: S-0682: `scripts/agent-container.sh --probe S-nnnn` runs the agent image as `probe-S-nnnn` in a new session, with `"probe":true` in pilot's start request. harness.Prompt answers a probe with a prompt that only reads and reports and names every writing tool as never to be called, and pilot starts it without FLAI_STORY. pilot's TestRunStartsAProbeWithoutTheStory replays S-0507's trial request against the same request marked probe. A capped probe of the rebuilt image on S-0682 (haiku, 0.0098 USD) called only read tools and left S-0682's narrative and items unchanged.
