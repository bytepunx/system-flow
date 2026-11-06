---
id: T-2469
type: task
nature: research
title: A proposed ADR decides the shared cluster, the story and default namespaces, the manifest's kluster block, the commands' variables, the image load, acceptance, and the kluster host action, and the operator accepts it on a thread
status: backlog
parent: S-0576
owner: alex
created: 2026-11-06T21:41:55Z
updated: 2026-11-06T21:41:55Z
transitions: []
stream: S-0576
tags: [kluster, adr]
touches: [design/adrs/README.md]
---
# T-2469 A proposed ADR decides the shared cluster, the story and default namespaces, the manifest's kluster block, the commands' variables, the image load, acceptance, and the kluster host action, and the operator accepts it on a thread

## Work

Read S-0575's results document, `design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md`, and `design/system/project-manifest.md` § the manifest's blocks. Write a proposed ADR with `flai adr new`, with the topics `kluster`, `deployment`, and `manifest`. It decides each point of S-0576's first criterion, with the alternatives the experiment showed and why each was left. Settle these:

- how the story namespace is named from the serial;
- the `kluster:` keys, which are required, and their defaults;
- the variables given to the build, deploy, and teardown commands;
- whether flai loads images itself (`docker save` into `ctr`, or `k3d image import`) or waits for kluster to gain an import;
- what the review checks see in their environment.

List under Consequences each change the decision would ask of kluster's own repository, as a proposal for it.

Open a thread on S-0576 asking the operator to accept it, with the points they may want to change first. Set the ADR accepted with `flai adr accept` only once they do. The ADR index regenerates.

This is the first layer: it waits for no task.

## Done when

- The ADR exists, decides every point of the first criterion, and lists the kluster proposals under Consequences.
- The operator accepted it on a thread, and its status is accepted.
- `flai check --strict` passes on it and the index.

## Notes
