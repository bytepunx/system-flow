---
id: I-0128
title: flai test vets and tests only the changed Go packages, so a change that breaks a package importing them is found only by the integration tier
class: efficiency
status: closed
count: 1
cost: 10m
first_reported: 2026-10-08T10:50:48Z
last_reported: 2026-10-08T10:50:48Z
updated: 2026-10-12T15:36:43Z
---

# I-0128 flai test vets and tests only the changed Go packages, so a change that breaks a package importing them is found only by the integration tier

## Description
flai test vets and tests only the changed Go packages, so a change that breaks a package importing them is found only by the integration tier

## Instances

### 2026-10-08T10:50:48Z
Story: S-0338.
S-0338 added a slice field to workitem.Hold, which made it uncomparable. flai test on flai/internal/workitem passed: vet, lint, and the short tests run only ./internal/workitem and ./internal/messages. flai/cmd/story_start_test.go compares two holds with !=, so the cmd package stopped building. Only the close-out's integration tier found it, after 3m24s, and its output tail named no package (as S-0313 records). Finding the failing package took one more filtered full run. Remediation: have the vet tier, or the short tests, also build the packages that import the changed ones (go list -deps reversed, or go vet ./... since it is fast).

## Remediation

Story S-0364 remediates this issue, created from it at 2026-10-12T14:12:31Z.
Closed 2026-10-12T15:36:43Z: Fixed by S-0364: a test tier's command may hold {importers}, the changed Go packages and every package of the module that imports one of them, in its own files or its test files. This repository's vet tier runs go vet on {importers}, so a change that breaks a package importing the changed one, as S-0338's change to internal/workitem broke cmd/story_start_test.go, fails flai test and the close-out's vet tier rather than only the integration tier. flai test flai/internal/workitem now vets ./cmd and every other importer, in 2.5s.
