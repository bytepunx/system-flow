---
id: I-0148
title: A close-out whose run main outruns takes two rounds of the integration and smoke tiers, about 24 minutes, past the Bash tool's 10-minute limit, so the agent's call is killed mid-round
class: efficiency
status: open
count: 1
cost: 25m
first_reported: 2026-10-15T00:07:49Z
last_reported: 2026-10-15T00:07:49Z
updated: 2026-10-15T00:07:49Z
---

# I-0148 A close-out whose run main outruns takes two rounds of the integration and smoke tiers, about 24 minutes, past the Bash tool's 10-minute limit, so the agent's call is killed mid-round

## Description
A close-out whose run main outruns takes two rounds of the integration and smoke tiers, about 24 minutes, past the Bash tool's 10-minute limit, so the agent's call is killed mid-round

## Instances

### 2026-10-15T00:07:49Z
Story: S-0306.
S-0306's close-out, run as one Bash call with the 600000 ms timeout the conventions give, passed every step of round 1 (integration 5m17s, smoke 4m30s). Main moved during it (S-0389 accepted), so it synced and verified again, and the call was moved to the background at 10 minutes and killed at the background limit in round 2's verify. Running it again in the background with a longer timeout passed, but main moved once more (S-0388), so round 2 ran integration (5m42s) and smoke (6m10s) again: about 24 minutes in all. ADR-0156's reuse does not help when main's new paths are under flai/, which select both tiers. The conventions' single 600000 ms call cannot hold a close-out that needs a second round of these tiers.

## Remediation
