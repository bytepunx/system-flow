---
id: TH-0312
title: "E-0015 cost of delay input: a penalty per week on the epic"
anchor:
  path: wip/kanban/epics/E-0015-secure-dashboard-and-cli-releases.md
  item: E-0015
status: open
participants: [planner-E-0015]
created: 2026-10-07T19:27:08Z
updated: 2026-10-07T19:27:08Z
---

# TH-0312 E-0015 cost of delay input: a penalty per week on the epic

On wip/kanban/epics/E-0015-secure-dashboard-and-cli-releases.md.

## Entries

### 2026-10-07T19:27:08Z planner-E-0015
Recommendation: set `penalty_per_week: 100` on E-0015 (`flai edit E-0015 --penalty-per-week 100`). Its eight open stories, S-0232 to S-0239, would then each take a share of 100 USD a week, in proportion to its forecast duration, as `flai cod` works it out.

E-0015 and its stories have no cost of delay inputs, so `flai cod` refuses every one of them. The inputs are yours, so I will not set them.

Why a penalty, not revenue or time lost:

- Signing brings no revenue and saves no agent time each cycle.
- What waiting costs is exposure. Until S-0233 is done, `flai self-upgrade`, `flai host upgrade`, and `install.sh` install whatever the release URL serves, checked only against a `checksums.txt` from the same place. Until S-0236 to S-0238 are done, `flai dashboard` runs whatever image the tag names, and flai on the host connects to it. A tampered release runs on the host with the operator's git credentials.
- 100 USD a week is a rough expected loss: about 20,000 USD of damage from a compromise, times a 0.5% chance a week.

Alternatives:

- Lower, such as 25: if flai and flaiover are used only on this host, and the GitHub account is the only way in.
- Higher, such as 500: if releases reach other users, or the repository becomes public.

Meanwhile I am setting each story's touches and forecast, and recording them under `### Planning` in its Notes.
