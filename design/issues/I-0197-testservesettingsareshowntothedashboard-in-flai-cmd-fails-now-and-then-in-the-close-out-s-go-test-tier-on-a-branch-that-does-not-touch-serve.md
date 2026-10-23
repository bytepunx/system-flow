---
id: I-0197
title: TestServeSettingsAreShownToTheDashboard in flai/cmd fails now and then in the close-out's go-test tier on a branch that does not touch serve
class: defect
status: closed
count: 1
cost: 7m
first_reported: 2026-10-23T05:12:51Z
last_reported: 2026-10-23T05:12:51Z
updated: 2026-10-23T19:47:43Z
---

# I-0197 TestServeSettingsAreShownToTheDashboard in flai/cmd fails now and then in the close-out's go-test tier on a branch that does not touch serve

## Description
TestServeSettingsAreShownToTheDashboard in flai/cmd fails now and then in the close-out's go-test tier on a branch that does not touch serve

## Instances

### 2026-10-23T05:12:51Z
Story: S-0468.
S-0468's first close-out stopped at go-test on flai/cmd/serve_actions_test.go:981, "no token in the settings". S-0468 changes nothing the test reads. The second run, on the same commit, passed. The finding's message is cut off before the string it matched, so the cause is not confirmed. Likely: the settings JSON carries the test's random t.TempDir() paths (the project root and import_roots), and the test refuses any "4242" in it, which a random temp-directory suffix can contain. A rerun of the go-test tier, about 6 minutes, was the cost.

## Remediation

Story S-0488 remediates this issue, created from it at 2026-10-23T17:53:08Z.
Closed 2026-10-23T19:47:43Z: Fixed by S-0488: TestServeSettingsAreShownToTheDashboard refused any 4242 in the settings JSON, which also holds two t.TempDir() paths whose random numbers sometimes held 4242. The check is narrowed to the registered entry's URL, 127.0.0.1:4242, and the import root is now a folder named harbour-4242, which failed the old check on every run and passes the new one.
