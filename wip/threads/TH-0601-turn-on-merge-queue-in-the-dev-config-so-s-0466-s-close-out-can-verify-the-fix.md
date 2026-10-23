---
id: TH-0601
title: Turn on merge-queue in the dev config so S-0466's close-out can verify the fix
anchor:
  path: wip/kanban/stories/S-0466-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md
  item: S-0466
status: resolved
participants: [agent-S-0466, alex]
created: 2026-10-23T03:33:53Z
updated: 2026-10-23T04:45:25Z
---

# TH-0601 Turn on merge-queue in the dev config so S-0466's close-out can verify the fix

On wip/kanban/stories/S-0466-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md.

## Entries

### 2026-10-23T03:33:53Z agent-S-0466
The code and docs for option 1 are done and committed on story/S-0466 (T-2061, T-2062). One step is yours, as the recommendation on TH-0594 said. Run this in a shell from the main checkout:

```sh
scripts/flai.sh serve enable merge-queue
```

It adds this project to `merge-queue` in `.flai-cache/config.json`, the configuration `scripts/flai.sh` and the close-out read. `flai serve enable` says that turning a host action on is the operator's, so I haven't run it.

Once it is on, S-0466's close-out leaves integration and smoke to the queue, and I-0148 (T-2063) closes on that run. Reply here once it's done. If you'd rather allow me to run it, reply "allow" and I will.

If you'd rather keep it off, I'll run the close-out as it is. It will run both tiers, about 12 minutes, longer than the Bash call's 10-minute limit, so I'd run it in the background. Its output will show the new line naming the setting that is off, and I'd close I-0148 on the verify line and its test instead.

### 2026-10-23T04:44:59Z alex
I ran the command

### 2026-10-23T04:45:25Z agent-S-0466
Resolved: The operator turned on merge-queue in .flai-cache/config.json. I-0148 is closed (T-2063), and S-0466's close-out runs next.
