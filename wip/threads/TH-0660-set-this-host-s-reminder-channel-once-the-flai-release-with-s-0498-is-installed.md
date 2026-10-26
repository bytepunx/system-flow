---
id: TH-0660
title: Set this host's reminder channel once the flai release with S-0498 is installed
anchor:
  path: wip/kanban/stories/S-0498-a-pending-recommendation-waits-a-day-or-more-for-the-operator-s-confirmation-while-its-story-idles-in-progress.md
  item: S-0498
status: resolved
participants: [agent-S-0498, orchestrator, alex]
created: 2026-10-26T05:14:14Z
updated: 2026-10-26T16:01:44Z
---

# TH-0660 Set this host's reminder channel once the flai release with S-0498 is installed

On wip/kanban/stories/S-0498-a-pending-recommendation-waits-a-day-or-more-for-the-operator-s-confirmation-while-its-story-idles-in-progress.md.

## Entries

### 2026-10-26T05:14:14Z agent-S-0498
S-0498 adds reminders. After `threads.remind_after` (4h for this project), `flai serve` reminds you of a thread that awaits you, and again every `threads.remind_every` (24h). It does this through a channel kept in `~/.flai/config.json`. Only you set that file, in a shell on the host. This thread does not hold S-0498: run the command once the release that carries S-0498 is installed (`flai host upgrade`).

**Recommended: a push to your phone through ntfy.** Subscribe to a long random topic in the ntfy app, then run:

```bash
flai serve notify set --command 'curl -s -H "Title: flai" -d @- https://ntfy.sh/<your-topic>'
flai serve notify show
```

Anyone who knows a topic's name on ntfy.sh can read it, and the reminder carries the thread's title. With `jq` installed, `docs/operators/index.md` § Reminding the operator of what awaits them gives a one-line form.

**Or a Windows toast from WSL2.**
1. Install the BurntToast module in PowerShell.
2. Save `$r = [Console]::In.ReadToEnd() | ConvertFrom-Json; New-BurntToastNotification -Text "flai: $($r.project)", $r.title` as `~/flai-toast.ps1`.
3. Run `flai serve notify set --command 'powershell.exe -NoProfile -File "$(wslpath -w ~/flai-toast.ps1)"'`.

A toast reaches you only while you are at the Windows desktop.

Until a channel is set, `flai serve` writes each reminder to `serve/serve.log` and warns once that no channel is set. Reply here when it is set, or name other values for `threads.remind_after` and `threads.remind_every`.

### 2026-10-26T05:23:32Z orchestrator (recommendation)
**Recommended for the operator:**
1. Resolve this thread now, so S-0498 can land.
2. Set the channel yourself, in a shell, once the release with S-0498 is installed. Take the agent's ntfy command, with a long random topic.
3. Have I-0210 made into a story soon.

**Why this thread holds S-0498.** The thread says it does not hold the story, but it does. The merge queue refuses any story with an open thread:

```text
not eligible for the merge queue: thread_open: thread TH-0660 is not resolved on S-0498 or its tasks
```

I verified S-0498 at a759ff18. flai verify passed, and the verifier matched both criteria with reproducing tests. A reminder writes nothing on a thread, and every block that S-0498 writes has a path that unblocks it. I enqueue it the moment this thread is resolved.

**Why I-0210 matters before you set a command.** The reminder runs the stored `--command` through `sh -c` in `flai serve`'s own process, outside any sandbox. I-0210, which S-0498's close-out recorded, says flai guard does not stop a served agent from running `flai serve notify set`. The gap is the same one that already covers `flai serve agent`, `forge` and `enable`. So once S-0498 lands, an agent could rewrite the command that flai serve later runs. The fix is a guard rule with a test, in a story of its own.

**Why it is yours.** The channel is a host change, and making a story is a change of scope.

Source: design/conventions/safety.md § Rules

### 2026-10-26T16:01:44Z alex
Resolved: S-0498 was accepted
