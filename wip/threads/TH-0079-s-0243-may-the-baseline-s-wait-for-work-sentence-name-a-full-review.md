---
id: TH-0079
title: "S-0243: may the baseline's wait_for_work sentence name a full review?"
anchor:
  path: design/conventions/work-management.md
status: resolved
participants: [agent-S-0243, alex]
created: 2026-10-03T03:19:05Z
updated: 2026-10-03T03:26:35Z
---

# TH-0079 S-0243: may the baseline's wait_for_work sentence name a full review?

On design/conventions/work-management.md.

## Entries

### 2026-10-03T03:19:05Z agent-S-0243
From S-0243 (option A on TH-0078). A baseline sentence in `work-management.md` is now false. It says `wait_for_work` "names the story to pull as soon as one is ready and the in-progress limit leaves room", and that on a timeout it says "whether it is waiting for room or for a story to be ready".

With a full review holding the pull (ADR-0073), I'd make the same edit here and in `template/root/design/conventions/work-management.md`:
- "…as soon as one is ready, the in-progress limit leaves room, and review is under its limit"
- "…whether it is waiting for room, for a story in review to be accepted or sent back, or for a story to be ready"

It only describes what the tool now does; no rule changes. The template copy would be 1.0.35, with a changelog line.

Recommended: yes. If no, I'll take the edit out of both copies and record it in the story's Open questions.

### 2026-10-03T03:26:02Z alex
yes

### 2026-10-03T03:26:35Z agent-S-0243
Resolved: The edit is made here and in the template (1.0.35).
