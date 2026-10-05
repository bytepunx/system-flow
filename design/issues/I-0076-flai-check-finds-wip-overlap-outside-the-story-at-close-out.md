---
id: I-0076
title: "flai check finds `wip.overlap` outside the story at close-out"
class: efficiency
status: open
count: 4
first_reported: 2026-10-05T03:24:33Z
last_reported: 2026-10-05T04:30:02Z
updated: 2026-10-05T04:30:02Z
---

# I-0076 flai check finds `wip.overlap` outside the story at close-out

## Description
flai check finds `wip.overlap` outside the story at close-out

## Instances

### 2026-10-05T03:24:33Z
Story: S-0260.
flai check found outside the story:
`wip/kanban/stories/S-0253-an-acceptance-merge-committed-conflict-markers-to-a-design-document-on-main-and-nothing-caught-it.md`: S-0253 touches docs/users/flai.md, which S-0260 (in progress) also touches as docs/users/flai.md

### 2026-10-05T03:23:20Z
Story: S-0253.
flai check found outside the story:
`wip/kanban/stories/S-0253-an-acceptance-merge-committed-conflict-markers-to-a-design-document-on-main-and-nothing-caught-it.md`: S-0253 touches docs/users/flai.md, which S-0260 (in progress) also touches as docs/users/flai.md

### 2026-10-05T04:26:15Z
Story: S-0262.
flai check found outside the story:
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches design/issues, which S-0262 (in progress) also touches as design/issues
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches design/system/flai-cli.md, which S-0262 (in progress) also touches as design/system/flai-cli.md
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches docs/users/flai.md, which S-0262 (in progress) also touches as docs/users/flai.md

### 2026-10-05T04:30:02Z
Story: S-0262.
flai check found outside the story:
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches design/issues, which S-0262 (in progress) also touches as design/issues
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches design/system/flai-cli.md, which S-0262 (in progress) also touches as design/system/flai-cli.md
`wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md`: S-0257 touches docs/users/flai.md, which S-0262 (in progress) also touches as docs/users/flai.md
`wip/kanban/stories/S-0262-flai-s-wip-markdown-lint-does-not-flag-a-space-inside-a-code-span-so-a-task-body-reached-main-and-failed-a-story-s-close-out.md`: S-0262 touches design/issues, which T-0877 (in progress) also touches as design/issues
`wip/kanban/stories/S-0262-flai-s-wip-markdown-lint-does-not-flag-a-space-inside-a-code-span-so-a-task-body-reached-main-and-failed-a-story-s-close-out.md`: S-0262 touches design/system/flai-cli.md, which T-0877 (in progress) also touches as design/system/flai-cli.md
`wip/kanban/stories/S-0262-flai-s-wip-markdown-lint-does-not-flag-a-space-inside-a-code-span-so-a-task-body-reached-main-and-failed-a-story-s-close-out.md`: S-0262 touches docs/users/flai.md, which T-0877 (in progress) also touches as docs/users/flai.md

## Remediation
