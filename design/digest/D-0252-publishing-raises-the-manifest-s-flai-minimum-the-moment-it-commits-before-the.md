---
id: D-0252
title: Publishing raises the manifest's flai minimum the moment it commits, before the release is built, so the host's flai drops the project until the binaries exist
kind: cause
sources: [I-0107]
topics: [release, flai, cli, go, code, all]
paths: [flai/internal/release/release.go, flai/internal/release/release_test.go, flai/cmd/release.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flai/internal/channel/channel.go, design/system/flai-cli.md, design/system/project-manifest.md, docs/users/flai.md, docs/operators/index.md, design/issues/I-0107-publishing-raises-the-manifest-s-flai-minimum-the-moment-it-commits-before-the-release-is-built-so-the-host-s-flai-drops-the-project-until-the-binaries-exist.md, design/issues/summary.md, design/system/work-hierarchy.md, flai/internal/manifest/manifest.go, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, design/issues/I-0113-a-failed-integration-tier-in-the-close-out-shows-only-the-last-lines-of-go-test-so-the-failing-test-is-not-named.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md, design/issues/I-0123-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0252 Publishing raises the manifest's flai minimum the moment it commits, before the release is built, so the host's flai drops the project until the binaries exist

I-0107 was closed on 2026-10-08.

The cause:

```text
`flai release --pending` (the board's Publish) commits in one step the release's changelog, the template version, and, when a flai release changed the front-matter fields flai reads, a raised `flai.minimum` in `system-flow.yaml`. The raised minimum takes effect on the host the instant the commit lands: `flai serve` re-reads the manifest every second, refuses to load one that needs a newer flai than the one running, drops the project, and cancels every request in flight for it, the publish itself included. The binaries for the new minimum do not exist yet: the `release flai` workflow starts when the tag reaches GitHub and takes minutes, and it can fail. Until `flai host upgrade` finds the release, the dashboard answers 503 for the project, its MCP server over HTTP is released, and no agent can be started for it. The operator sees the dashboard and the host stop at the moment of pressing Publish and reads it as a crash.

The same flai that commits the minimum is the one serving the project, so the publish removes the ground it stands on. A minimum should bind once the release it names can be installed, or a flai below the minimum should keep serving a project whose fields it already read until the operator upgrades, warning rather than dropping.
```

The fix:

```text
S-0321: a publish raises flai.minimum only to a flai release no newer than the flai that publishes, so the minimum always names a release whose binaries exist and a newer raise waits for a publish from the upgraded flai (release.RaiseMinimum, TestRaiseMinimum reproduces the instance). A served project whose manifest stops loading now refuses new requests and lets those in flight answer, for up to five minutes, before flai serve drops it, so a publish is no longer killed mid-run (TestAProjectThatStopsLoadingAnswersItsRequestInFlightBeforeItIsDropped).
```

Source: I-0107.
