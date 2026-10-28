---
id: D-0229
title: The dashboard upgrade stops the old container and cannot start the new one, because the name flaiover is still in use
kind: cause
sources: [I-0094]
topics: [flai, cli, go, code, all]
paths: [flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_test.go, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/runbooks/update.md, design/system/flai-cli.md, design/issues/I-0094-the-dashboard-upgrade-stops-the-old-container-and-cannot-start-the-new-one-because-the-name-flaiover-is-still-in-use.md, design/issues/summary.md, design/issues/I-0116-flai-dashboard-restart-starts-the-container-s-tag-again-so-a-newer-image-a-check-pulled-under-a-floating-tag-is-started-without-an-upgrade.md, design/issues/I-0109-flai-check-finds-narrative-state-outside-the-story-at-close-out.md, design/issues/I-0110-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md, design/issues/I-0118-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md]
date: 2026-10-08
status: absorbed
absorbed_by: "design/system/flai-cli.md#listing-and-installing-a-published-release"
asked_again: []
cited: []
---

# D-0229 The dashboard upgrade stops the old container and cannot start the new one, because the name flaiover is still in use

I-0094 was closed on 2026-10-08.

The cause:

```text
`dashboard.upgrade` replaces the running flaiover container with one on the new image. It stops the old container and at once runs the new one under the same name. Docker refused the new one: the name was still held by the old container. The likely reason, not confirmed, is that the old container runs with `--rm`, as the new one does, and Docker had not finished removing it when the new `docker run` asked for the name. The upgrade then reports failure with the dashboard down.

flai host's own watch recovered it: it found the dashboard gone and started it about a minute later. To the operator, in the dashboard, the page they had just used to upgrade stopped answering.
```

The fix:

```text
S-0316 fixed it (commit 66f2113e): after docker stop, flai dashboard upgrade and restart wait until Docker no longer lists the container's name, up to 10 tries 500 ms apart, then remove it with docker rm -f before starting under it; and when the new image still fails to start, the upgrade starts the previous image again by its image ID and fails saying so, so the operator is not left without a dashboard. Tests in flai/cmd/dashboard_test.go reproduce the lingering name with a stand-in for docker.
```

Source: I-0094.
