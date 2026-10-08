---
id: I-0107
title: Publishing raises the manifest's flai minimum the moment it commits, before the release is built, so the host's flai drops the project until the binaries exist
class: defect
status: closed
count: 2
cost: 30m
first_reported: 2026-10-07T07:39:30Z
last_reported: 2026-10-07T15:09:19Z
updated: 2026-10-08T08:11:37Z
---

# I-0107 Publishing raises the manifest's flai minimum the moment it commits, before the release is built, so the host's flai drops the project until the binaries exist

## Description

`flai release --pending` (the board's Publish) commits in one step the release's changelog, the template version, and, when a flai release changed the front-matter fields flai reads, a raised `flai.minimum` in `system-flow.yaml`. The raised minimum takes effect on the host the instant the commit lands: `flai serve` re-reads the manifest every second, refuses to load one that needs a newer flai than the one running, drops the project, and cancels every request in flight for it, the publish itself included. The binaries for the new minimum do not exist yet: the `release flai` workflow starts when the tag reaches GitHub and takes minutes, and it can fail. Until `flai host upgrade` finds the release, the dashboard answers 503 for the project, its MCP server over HTTP is released, and no agent can be started for it. The operator sees the dashboard and the host stop at the moment of pressing Publish and reads it as a crash.

The same flai that commits the minimum is the one serving the project, so the publish removes the ground it stands on. A minimum should bind once the release it names can be installed, or a flai below the minimum should keep serving a project whose fields it already read until the operator upgrades, warning rather than dropping.

## Instances

### 2026-10-07T07:39:30Z
On 2026-10-07 the operator pressed Publish at 07:35:59Z. flai release committed the bump of flai.minimum to 1.34.1 and pushed the tags flai/v1.34.1 and flaiover/v0.36.1. At 07:36:01Z, two seconds later, flai serve (the installed flai 1.34.0) re-read system-flow.yaml, found it needs 1.34.1, dropped the project, and killed its own in-flight publish.run, which the journal records as failed with 'flai exited with -1'. The dashboard answered 503 for the project from then on and the host released its MCP server. The release flai workflow then failed at 07:38Z, so flai host upgrade finds nothing newer than 1.34.0 and the project stays unserved. The operator reported it as flai host having crashed.

### 2026-10-07T15:09:19Z
Again on 2026-10-07 at 15:04Z: a publish raised flai.minimum to 1.35.1 the moment it committed, with the release build for 1.35.1 still running on GitHub. The newest built flai, 1.35.0, which flai self-upgrade installed at 15:05Z, refuses every command in the project, and the running flai serve (1.34.6) dropped system-flow from the dashboard, which went on serving the other seven projects. flai dashboard check and upgrade had to be run from outside the project. The project stays unserved until the build finishes and flai host upgrade runs.

## Remediation

Directions to weigh: have `flai release` raise `flai.minimum` only after the release's binaries are published, in a later commit, or have the publish leave the bump for the first run of the upgraded flai; or have `flai serve` keep serving a project whose manifest newly demands a minimum above its own version, with the warning it already logs for a flai older than the project, and refuse only the writes that touch fields it does not know. Either way the publish should not cancel its own `publish.run`: the journal recorded it as failed with `flai exited with -1` although the tags had been pushed. A test that publishes a release raising the minimum against a running `flai serve` of the older version, and expects the project to stay served, would pin it. The release build's own failure that day (a duplicated `## 1.0.67` heading in `template/CHANGELOG.md`, which `TestRepositoryLintsClean` rejects, since `flai release` prepended a second section with the heading the stories had already written) is a separate defect and made the outage last longer.

Story S-0321 remediates this issue, created from it at 2026-10-07T18:59:52Z.
Closed 2026-10-08T08:11:37Z: S-0321: a publish raises flai.minimum only to a flai release no newer than the flai that publishes, so the minimum always names a release whose binaries exist and a newer raise waits for a publish from the upgraded flai (release.RaiseMinimum, TestRaiseMinimum reproduces the instance). A served project whose manifest stops loading now refuses new requests and lets those in flight answer, for up to five minutes, before flai serve drops it, so a publish is no longer killed mid-run (TestAProjectThatStopsLoadingAnswersItsRequestInFlightBeforeItIsDropped).
