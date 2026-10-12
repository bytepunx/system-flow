---
id: I-0134
title: "flai check finds `issues.no-story` outside the story at close-out"
class: efficiency
status: open
count: 2
first_reported: 2026-10-09T18:33:15Z
last_reported: 2026-10-12T14:25:45Z
updated: 2026-10-12T14:25:45Z
---

# I-0134 flai check finds `issues.no-story` outside the story at close-out

## Description
flai check finds `issues.no-story` outside the story at close-out

## Instances

### 2026-10-09T18:33:15Z
Story: S-0289.
flai check found outside the story:
`design/issues/I-0052-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md`: I-0052 has been open since 2026-10-01 with no open story linking it (issues.story_after is unset, so 168h): make one with flai issue story I-0052

### 2026-10-12T14:25:45Z
Story: S-0361.
flai check found outside the story:
`design/issues/I-0052-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md`: I-0052 has been open since 2026-10-01 with no open story linking it (issues.story_after is unset, so 168h): make one with flai issue story I-0052

## Remediation

Story S-0365 remediates this issue, created from it at 2026-10-12T14:12:33Z.
