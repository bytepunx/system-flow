---
id: T-2474
type: task
nature: feature
title: A kluster package runs kluster up, down, and status from the manifest's block, parses status into data, and says how to install kluster when it is missing
status: done
parent: S-0577
owner: alex
created: 2026-11-06T21:42:47Z
updated: 2026-11-07T01:11:51Z
transitions:
  - to: ready
    at: 2026-11-07T01:09:09Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T01:09:10Z
    by: system-flow
  - to: done
    at: 2026-11-07T01:11:51Z
    by: agent-S-0577
stream: S-0577
tags: [kluster, flai]
touches: [flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go]
after: [T-2472]
usage:
  source: log
  seconds: 161
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 6437
      cache_read: 2875456
      cache_write: 32146
      cost: 0.9306
---
# T-2474 A kluster package runs kluster up, down, and status from the manifest's block, parses status into data, and says how to install kluster when it is missing

## Work

It waits for T-2472, whose manifest type for the `kluster:` block it takes.

Add `flai/internal/kluster`. Given the project root and the block, it:

- finds `kluster` on `PATH`, and when it is missing, returns an error whose sentence names kluster's `install.sh` and `go install github.com/bytepunx/kluster/kluster@latest`;
- runs every command in the project's main checkout, as `design/system/kluster.md` § The project's cluster says, with the config path from the repository root;
- runs `kluster --config <path> up`, with the provider the config names, streaming its output to a writer it is given, and once it exits 0 writes the cluster's kubeconfig to `.flai-cache/kluster/kubeconfig` with `kluster --config <path> kubeconfig --name <name> --output <file>`, never merging it into `~/.kube/config`;
- runs `kluster --config <path> down --yes`;
- runs `kluster --config <path> status` and parses its table (NAME, RUNNING, AGE) into a status value for the block's cluster, with the profile and provider read from the kluster config, and a cluster absent from the table reported as not created.

Run commands through the project's process helper (`flai/internal/execx`), so tests swap in a fake. Test in `kluster_test.go`, with a fake `kluster` script on `PATH` written into a temporary folder: up, down, a running cluster, a stopped one, an absent one, a status table kluster's next version might widen, a failing exit, and a missing binary.

## Done when

- The package brings a cluster up and down and returns its status as data, with every case above tested.
- No test needs a real kluster, cluster, or Docker.
- `flai test` on the two paths passes.

## Notes
