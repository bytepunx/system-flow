---
id: T-1886
type: task
nature: feature
title: file.save with no hash creates a file that is absent, so the first example of a component can be saved
status: backlog
parent: S-0443
owner: alex
created: 2026-10-16T21:38:25Z
updated: 2026-10-16T21:48:32Z
transitions: []
stream: S-0443
tags: [cli]
touches: [flai/internal/docedit/file.go, flai/internal/docedit/file_test.go]
after: [T-1884]
---
# T-1886 file.save with no hash creates a file that is absent, so the first example of a component can be saved

## Work

In S-0436's `FileSave`, create a file that is absent when the hash is empty, with its folders, and refuse an empty hash for a file that exists. The check runs before and after, as for any save, so T-1885's finding refuses a bad example. The reach is unchanged: `design/components/` is under the design folder.

Waits for T-1884, whose ADR records the create.

## Done when

Tests cover the create with new folders, the empty hash on an existing file, and a created file the check refuses, leaving nothing behind. `flai test flai/internal/docedit` passes.

## Notes

None.
