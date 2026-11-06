---
id: T-2482
type: task
nature: feature
title: flai serve's kluster pass tears down a story that leaves review or is cancelled, and on acceptance tears it down and promotes main to the default namespace
status: backlog
parent: S-0579
owner: alex
created: 2026-11-06T21:44:23Z
updated: 2026-11-06T21:44:23Z
transitions: []
stream: S-0579
tags: [kluster, serve]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go]
after: [T-2481]
---
# T-2482 flai serve's kluster pass tears down a story that leaves review or is cancelled, and on acceptance tears it down and promotes main to the default namespace

## Work

It waits for T-2481, whose pass it extends, in the same files.

In `flai/internal/serve/kluster.go`, on each look, for each deployment record whose story is no longer in review:

- back to in progress, or cancelled: tear it down;
- done: tear it down, then promote main to the default namespace, once per acceptance, with acceptances that land together promoting once.

Teardowns and promotions take their turn behind the deploys of the same project. The journal logs each. A failed teardown or promotion opens one thread on the story, or on the epic for a promotion of several, naming the step and the log.

Test each case in `kluster_test.go`: a story sent back, a cancelled story, an accepted story, two acceptances landing together promoting once, and a failed promotion's one thread.

## Done when

- Leaving review, cancellation, and acceptance tear down and promote as above.
- Every case above is tested, and `flai test` on the two paths passes.

## Notes
