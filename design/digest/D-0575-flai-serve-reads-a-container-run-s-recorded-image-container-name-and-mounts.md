---
id: D-0575
title: flai serve reads a container run's recorded image, container name, and mounts from the docker run argv it built, so the record and the argv cannot disagree, and the mount type is harness.ContainerMount, apart from storygit.Mount
kind: decision
sources: [S-0665]
topics: [docker, harness, go, code]
paths: [flai/internal/harness/container.go, flai/internal/harness/harness.go, flai/internal/serve/agents_start.go, flai/internal/serve/role.go]
date: 2026-11-10
status: absorbed
absorbed_by: "design/system/flai-cli.md#the-container-harness-s-0508"
asked_again: []
cited: []
---

# D-0575 flai serve reads a container run's recorded image, container name, and mounts from the docker run argv it built, so the record and the argv cannot disagree, and the mount type is harness.ContainerMount, apart from storygit.Mount

S-0665 fixed I-0274: `flai serve` recorded no image, container name, or mounts for a run on the `container` harness, so nobody could show afterwards that a run used the project's layer image. It now records them in `agents.json` (`image`, `container`, `mounts`) and names the image and container on the serve journal's start line.

How it was built, from S-0665's Decisions (`wip/archive/agents/S-0665.md`):

```text
T-2942 reads the record from the docker run argv (`containerRun`), so the record and the argv cannot disagree; the mount type is `ContainerMount`, apart from `storygit.Mount`.
T-2945 also changed `flai/internal/serve/role.go`, so one helper, `startedSaid`, writes every start line; a defined role cannot run on the container harness, so its line is unchanged.
No test of a container start without a login: the only inner harness, claude-code, refuses one.
```

So a later change to what a container run records derives it from `containerRun`'s argv, not from a second description of the mounts, and a change to a start line goes through `startedSaid`. The reproducing test is `TestAContainerRunRecordsItsImageContainerAndMounts`.

Sources: S-0665 (T-2942, T-2943, T-2945), I-0274.
