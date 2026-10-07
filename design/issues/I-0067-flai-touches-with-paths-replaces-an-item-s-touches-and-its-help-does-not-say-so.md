---
id: I-0067
title: flai touches with paths replaces an item's touches, and its help does not say so
class: impression
status: closed
count: 1
cost: 2m
first_reported: 2026-10-03T18:47:44Z
last_reported: 2026-10-03T18:47:44Z
updated: 2026-10-07T00:43:45Z
---

# I-0067 flai touches with paths replaces an item's touches, and its help does not say so

## Description
flai touches with paths replaces an item's touches, and its help does not say so

## Instances

### 2026-10-03T18:47:44Z
Story: S-0206.
Adding the paths a story's later tasks touch with flai touches S-0206 <paths> replaced the story's five original paths instead of adding to them; the help's example reads as either. Restored by naming every path again.

## Remediation

S-0254: `flai touches` gains `--add`, which adds the paths given and keeps the rest, and `--remove`, which takes them out and refuses one the item does not touch. Its help says that paths given alone replace the list, and replacement stays the default for scripts that rely on it. `flai stream sync`'s widen hint is now `flai touches <id> --add <paths outside>`, so following it no longer means naming every path again. `TestTouchesAddsRemovesAndReplaces` covers the three modes and the refusals.
Closed 2026-10-07T00:43:45Z: S-0254: flai touches adds with --add and removes with --remove, its help says paths given alone replace the list, and flai stream sync's widen hint uses --add
