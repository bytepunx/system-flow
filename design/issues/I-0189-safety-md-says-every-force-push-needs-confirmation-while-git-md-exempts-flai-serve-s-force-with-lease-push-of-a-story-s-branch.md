---
id: I-0189
title: safety.md says every force push needs confirmation, while git.md exempts flai serve's force-with-lease push of a story's branch
class: defect
status: closed
count: 1
first_reported: 2026-10-21T16:01:40Z
last_reported: 2026-10-21T16:01:40Z
updated: 2026-11-07T04:38:55Z
---

# I-0189 safety.md says every force push needs confirmation, while git.md exempts flai serve's force-with-lease push of a story's branch

## Description
safety.md says every force push needs confirmation, while git.md exempts flai serve's force-with-lease push of a story's branch

## Instances

### 2026-10-21T16:01:40Z
Raised on TH-0496 and TH-0498; the operator answered draft it on TH-0498. S-0393 added the exception to both copies of git.md (template 1.0.97, ADR-0166), but safety.md § Rules, in design/conventions and template/root/design/conventions, still says force pushes always need confirmation, per action and per context. An agent that reads both meets a contradiction. The fix is one line in each copy of safety.md naming the exception, with a template release. The orchestrator cannot create a story and no epic holds this, so it is recorded here for the issue step to make a draft story.

## Remediation

Story S-0558 remediates this issue, created from it at 2026-11-02T00:04:34Z.
Closed 2026-11-07T04:38:55Z: S-0558 rewrote the clause in both copies of safety.md § Rules: an agent's force push still always needs confirmation, and flai serve's --force-with-lease push of a story's own branch under pull-request acceptance is named as the one force push that needs none, linking git.md, which states the exception.
