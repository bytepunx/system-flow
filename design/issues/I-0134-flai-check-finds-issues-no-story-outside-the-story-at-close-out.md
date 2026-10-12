---
id: I-0134
title: "flai check finds `issues.no-story` outside the story at close-out"
class: efficiency
status: closed
count: 6
first_reported: 2026-10-09T16:51:57Z
last_reported: 2026-10-12T15:18:29Z
updated: 2026-10-12T15:18:29Z
---

# I-0134 flai check finds `issues.no-story` outside the story at close-out

## Description
flai check finds `issues.no-story` outside the story at close-out

## Instances

### 2026-10-09T16:51:57Z
Story: S-0343.
flai check found outside the story:
`design/issues/I-0052-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md`: I-0052 has been open since 2026-10-01 with no open story linking it (issues.story_after is unset, so 168h): make one with flai issue story I-0052

### 2026-10-09T18:33:15Z
Story: S-0289.
flai check found outside the story:
`design/issues/I-0052-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md`: I-0052 has been open since 2026-10-01 with no open story linking it (issues.story_after is unset, so 168h): make one with flai issue story I-0052

### 2026-10-12T14:23:39Z
Story: S-0342.
flai check found I-0052's issues.no-story outside the story at S-0342's close-out on 2026-10-08T21:20:32Z, first recorded as I-0130 on the story's branch and folded in here when the branch was rebased onto main

### 2026-10-12T14:25:45Z
Story: S-0361.
flai check found outside the story:
`design/issues/I-0052-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md`: I-0052 has been open since 2026-10-01 with no open story linking it (issues.story_after is unset, so 168h): make one with flai issue story I-0052

### 2026-10-12T14:46:12Z
Story: S-0362.
flai check found outside the story:
`design/issues/I-0052-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md`: I-0052 has been open since 2026-10-01 with no open story linking it (issues.story_after is unset, so 168h): make one with flai issue story I-0052

### 2026-10-12T15:18:29Z
Story: S-0232.
flai check found outside the story:
`design/issues/I-0052-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md`: I-0052 has been open since 2026-10-01 with no open story linking it (issues.story_after is unset, so 168h): make one with flai issue story I-0052

## Remediation

Story S-0365 remediates this issue, created from it at 2026-10-12T14:12:33Z.
Closed 2026-10-12T15:06:12Z: Fixed by S-0365 (ADR-0136): a check scoped to a story, as the close-out runs it, leaves out every issues.no-story outside the story and takes back its counts, so --record-issues no longer records an open issue no story links, such as I-0052; one on an issue file the story changes stays inside, and the main checkout's flai check still warns on every one. TestScopedCheckLeavesOutAnIssuesNoStoryOutsideTheStory and TestCheckRecordIssuesLeavesOutAnIssuesNoStoryOutsideTheStory reproduce it.
