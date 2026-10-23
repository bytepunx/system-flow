---
id: D-0255
title: The close-out's install smoke test failed once and passed when run alone, with no cause in its output
kind: cause
sources: [I-0086]
topics: [release, flai, cli, go, code, all]
paths: [install.sh, flai/cmd/installsh_test.go, flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md, design/issues/summary.md, design/issues/I-0122-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0255 The close-out's install smoke test failed once and passed when run alone, with no cause in its output

I-0086 was closed on 2026-10-08.

The cause:

```text
The close-out's install smoke test failed once and passed when run alone, with no cause in its output
```

The fix:

```text
S-0291. The cause was GitHub dropping the connection partway through install.sh's 1 MB `releases?per_page=50` listing, with curl (92) or (56), and nothing tried it again. install.sh now makes up to three attempts at a call the network drops, two seconds apart, and never retries an HTTP error. It resolves the latest release from pages of ten. flai self-upgrade's page reads and downloads do the same, with three attempts a second apart. Tests that cut the response off mid-body reproduce the drop: flai/cmd/installsh_test.go (TestInstallShRetriesAListingTheNetworkDrops) and flai/internal/selfupgrade/selfupgrade_test.go (TestDroppedAnswerIsAskedAgain). Each fails without the retry. The hidden install.sh output in scripts/install-test.sh is left to S-0340's rewrite of that script, as agreed on MS-0013.
```

Source: I-0086.
