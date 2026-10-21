---
id: I-0189
title: safety.md says every force push needs confirmation, while git.md exempts flai serve's force-with-lease push of a story's branch
class: defect
status: open
count: 1
first_reported: 2026-10-21T16:01:40Z
last_reported: 2026-10-21T16:01:40Z
updated: 2026-10-21T16:01:40Z
---

# I-0189 safety.md says every force push needs confirmation, while git.md exempts flai serve's force-with-lease push of a story's branch

## Description
safety.md says every force push needs confirmation, while git.md exempts flai serve's force-with-lease push of a story's branch

## Instances

### 2026-10-21T16:01:40Z
Raised on TH-0496 and TH-0498; the operator answered draft it on TH-0498. S-0393 added the exception to both copies of git.md (template 1.0.97, ADR-0166), but safety.md § Rules, in design/conventions and template/root/design/conventions, still says force pushes always need confirmation, per action and per context. An agent that reads both meets a contradiction. The fix is one line in each copy of safety.md naming the exception, with a template release. The orchestrator cannot create a story and no epic holds this, so it is recorded here for the issue step to make a draft story.

## Remediation
