---
id: T-2474
type: task
nature: feature
title: A kluster package runs kluster up, down, and status from the manifest's block, parses status into data, and says how to install kluster when it is missing
status: backlog
parent: S-0577
owner: alex
created: 2026-11-06T21:42:47Z
updated: 2026-11-06T21:42:47Z
transitions: []
stream: S-0577
tags: [kluster, flai]
touches: [flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go]
after: [T-2472]
---
# T-2474 A kluster package runs kluster up, down, and status from the manifest's block, parses status into data, and says how to install kluster when it is missing

## Work

It waits for T-2472, whose manifest type for the `kluster:` block it takes.

Add `flai/internal/kluster`. Given the project root and the block, it:

- finds `kluster` on `PATH`, and when it is missing, returns an error whose sentence names kluster's `install.sh` and `go install github.com/bytepunx/kluster/kluster@latest`;
- runs `kluster --config <path> up`, with the provider the config names, streaming its output to a writer it is given;
- runs `kluster --config <path> down --yes`;
- runs `kluster --config <path> status` and parses its table (NAME, RUNNING, AGE) into a status value for the block's cluster, with the profile and provider read from the kluster config, and a cluster absent from the table reported as not created.

Run commands through the project's process helper (`flai/internal/execx`), so tests swap in a fake. Test in `kluster_test.go`, with a fake `kluster` script on `PATH` written into a temporary folder: up, down, a running cluster, a stopped one, an absent one, a status table kluster's next version might widen, a failing exit, and a missing binary.

## Done when

- The package brings a cluster up and down and returns its status as data, with every case above tested.
- No test needs a real kluster, cluster, or Docker.
- `flai test` on the two paths passes.

## Notes
