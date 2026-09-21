---
id: S-0089
type: story
nature: remediation
title: Inbox should not continue showing open items for stories that moved to review
status: backlog
parent: E-0003
owner: alex
created: 2026-09-21T03:59:59Z
updated: 2026-09-21T03:59:59Z
transitions: []
tags: [dashboard]
touches: [flaiover/src]
---
# S-0089 Inbox should not continue showing open items for stories that moved to review

## Goal

Once open questions have been resolved via agent chat or by completing work, they should no longer appear in the inbox - this applies especially for stories that have entered review. If there were still open questions waiting on the operator, the story should remain in "in progress" and not be moved to review. Once a question is answered, it should get removed from the inbox.

## Acceptance criteria
- [ ] If a story is moved to the review step, it should no longer show open questions in the inbox
- [ ] Stories with unresolved questions in the Inbox cannot be moved to review since the agent can't complete work if the operator has not resolved open items.
- [ ] Once a question is answered, it should not linger in the Open Questions section of the Inbox.

## Tasks

## Notes
