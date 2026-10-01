---
id: I-0013
title: A git-backed test relied on the developer's global git identity and failed in CI
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-09-17T20:35:09Z
last_reported: 2026-09-17T20:35:09Z
updated: 2026-10-01T08:00:48Z
---

# I-0013 A git-backed test relied on the developer's global git identity and failed in CI

## Description
A git-backed test relied on the developer's global git identity and failed in CI

## Instances

### 2026-09-17T20:35:09Z
TestStoryBranchLifecycle passed locally because ~/.gitconfig has a user; CI has none, so flai's rebase inside the test repository failed with 'Committer identity unknown'. Fixed by setting user.name and user.email in the test repository. Reproduce locally with GIT_CONFIG_GLOBAL=/dev/null.

## Remediation
Closed 2026-10-01T08:00:48Z: Fixed by S-0037 (411ca5a): git-backed tests set their own identity. The inverse, the host's identity leaking into a test, is I-0045 (S-0180).
