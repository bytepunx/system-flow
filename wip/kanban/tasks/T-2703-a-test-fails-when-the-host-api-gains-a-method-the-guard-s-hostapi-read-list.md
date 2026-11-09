---
id: T-2703
type: task
nature: improvement
title: A test fails when the host API gains a method the guard's hostapi read list does not class
status: in-progress
parent: S-0612
owner: alex
created: 2026-11-09T07:57:22Z
updated: 2026-11-09T08:38:44Z
transitions:
  - to: ready
    at: 2026-11-09T08:38:42Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T08:38:44Z
    by: system-flow
stream: S-0612
tags: [flai, guard, hostapi]
touches: [flai/internal/hostapi/guardreads_test.go]
after: [T-2702]
---
# T-2703 A test fails when the host API gains a method the guard's hostapi read list does not class

## Work

Waits for T-2702, which writes the guard's list of the host API's read methods that this test checks.

The guard keeps its list apart from the method table, since it cannot import `flai/internal/hostapi`. Keep them in step with a test in a new `flai/internal/hostapi/guardreads_test.go`, which may import the guard: every name on the guard's list is a method of `MethodsFor`'s table, and none of them is a method `writeMethods` builds or any other method that writes (the `*.save`, `settings.*`, `components.start` and `stop`, and `kluster` writes among them). Class every other method of the table as a write in a list in the test, so that a new method that is on neither list fails the test with a message saying to add it to one, and why.

A change under `flai/internal/hostapi/` selects the `flaiover-flai` tier; run it with the Go tests.

## Done when

- The test passes, and fails, with a message naming the method, when a method is added to the table on neither list, or when a write method is put on the guard's list (checked by hand once and said in this task's Notes).
- `flai test flai/internal/hostapi/guardreads_test.go` passes.

## Notes
