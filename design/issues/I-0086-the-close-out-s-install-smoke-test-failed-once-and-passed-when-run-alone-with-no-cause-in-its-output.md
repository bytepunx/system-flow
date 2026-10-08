---
id: I-0086
title: The close-out's install smoke test failed once and passed when run alone, with no cause in its output
class: efficiency
status: closed
count: 8
cost: 7m
first_reported: 2026-10-06T10:22:10Z
last_reported: 2026-10-08T08:22:57Z
updated: 2026-10-08T08:25:32Z
---

# I-0086 The close-out's install smoke test failed once and passed when run alone, with no cause in its output

## Description
The close-out's install smoke test failed once and passed when run alone, with no cause in its output

## Instances

### 2026-10-06T10:22:10Z
Story: S-0283.
S-0283's close-out stopped at smoke's last sub-step, scripts/install-test.sh, after gofmt, vet, lint, the full go test -race, the template render and check, and the markdown lint passed; its output showed self-upgrade --check reaching GitHub for release 1.31.4 but no failure reason. The verifier ran scripts/install-test.sh alone once: exit 0 after over five minutes. Likely the network or a timeout; the cause is not confirmed.

### 2026-10-08T05:36:43Z
Story: S-0320.
S-0320's close-out failed smoke at `install-test: install.sh with an explicit FLAI_INSTALL_DIR`. This time the output names the cause: `curl: (56) OpenSSL SSL_read: ... unexpected eof while reading`, then `could not list releases of bytepunx/system-flow`. A transient network failure reaching GitHub's API fails the whole close-out after the integration tier has passed.

### 2026-10-08T05:40:46Z
Story: S-0318.
S-0318's close-out failed its smoke tier in install-test's default-path case. The tier's output gave no cause; `.flai-cache/install-test-home.log` did: curl (92) "HTTP/2 stream 1 was not closed cleanly: CANCEL (err 8)" while install.sh listed the releases, which it reports as "could not list releases (private repository?)". A transient network fault, not a change of the story's; the same tier passed at 04:52Z.

### 2026-10-08T05:59:18Z
Story: S-0319.
S-0319's close-out failed its smoke tier in install-test's explicit FLAI_INSTALL_DIR case: curl (92) "HTTP/2 stream 1 was not closed cleanly: CANCEL (err 8)" while install.sh listed the releases, reported as "could not list releases of bytepunx/system-flow (private repository?)". Same transient network fault as S-0318's; S-0319 changes nothing install.sh or the smoke tier runs. The close-out was run again.

### 2026-10-08T06:04:47Z
Story: S-0319.
S-0319's second close-out failed its smoke tier in install-test's default-path case, with no cause in the tier's output. `.flai-cache/install-test-home.log` held it: the same curl (92) "HTTP/2 stream 1 was not closed cleanly: CANCEL (err 8)" while listing releases. Two runs in a row each lost one of install-test's GitHub calls. That points at curl's HTTP/2 to the GitHub API on this host, which install.sh could retry or ask for with --http1.1.

### 2026-10-08T06:40:55Z
Story: S-0336.
S-0336's close-out stopped at the smoke tier in `install-test: install.sh with no FLAI_INSTALL_DIR installs under a fresh HOME/.flai/bin, without sudo`, with `exit status 1` and no cause printed, after the explicit FLAI_INSTALL_DIR case had installed flai 1.39.8. S-0336 changes no installer or script; every earlier tier, integration included, passed.

### 2026-10-08T06:55:19Z
Story: S-0336.
S-0336's second close-out stopped at the same step. Cause found: run by hand with the same fresh HOME, `install.sh` printed `curl: (56) OpenSSL SSL_read: OpenSSL/3.5.4: error:0A000126:SSL routines::unexpected eof while reading` and `x could not list releases of bytepunx/system-flow`, exit 1; a run straight after installed flai 1.39.8. So the step fails on a dropped connection to GitHub. `scripts/install-test.sh` writes install.sh's output to `.flai-cache/install-test-home.log` and, under `set -e`, exits as install.sh fails, before the `cat "$OUT"` that would show it: the remedy is to print the log when install.sh itself exits non-zero, and to retry the release listing once on a network error.

### 2026-10-08T07:15:27Z
Story: S-0326.
S-0326's close-out passed every tier through integration, then smoke failed in install-test: `curl: (92) HTTP/2 stream 1 was not closed cleanly: CANCEL (err 8)` while resolving the latest release, reported as "could not list releases of bytepunx/system-flow (private repository?)". A network fault, not the code; the close-out was run again and failed the same way. Five fetches of the listing `install.sh` asks for (`releases?per_page=50`, about 1 MB) on this host: two ended with `curl: (56) OpenSSL SSL_read: ... unexpected eof while reading`, three returned 200; `per_page=5` returned 200 three times out of three. A smaller listing, or a retry on a network error, would avoid it. A third close-out passed the first install case and failed the fresh-HOME case silently; `.flai-cache/install-test-home.log` held the same `unexpected eof while reading`.

### 2026-10-08T08:22:57Z
Story: S-0321.
S-0321's close-out smoke tier failed at install.sh with an explicit FLAI_INSTALL_DIR: "curl: (92) HTTP/2 stream 1 was not closed cleanly: CANCEL (err 8)" and "could not list releases of bytepunx/system-flow". Every other tier passed.

## Remediation

Story S-0291 remediates this issue, created from it at 2026-10-06T10:31:55Z.
Closed 2026-10-08T08:25:32Z: S-0291. The cause was GitHub dropping the connection partway through install.sh's 1 MB `releases?per_page=50` listing, with curl (92) or (56), and nothing tried it again. install.sh now makes up to three attempts at a call the network drops, two seconds apart, and never retries an HTTP error. It resolves the latest release from pages of ten. flai self-upgrade's page reads and downloads do the same, with three attempts a second apart. Tests that cut the response off mid-body reproduce the drop: flai/cmd/installsh_test.go (TestInstallShRetriesAListingTheNetworkDrops) and flai/internal/selfupgrade/selfupgrade_test.go (TestDroppedAnswerIsAskedAgain). Each fails without the retry. The hidden install.sh output in scripts/install-test.sh is left to S-0340's rewrite of that script, as agreed on MS-0013.

S-0340 (2026-10-08) takes the smoke tier off GitHub: `scripts/install-test.sh` installs and self-upgrades from releases built from the tree with `scripts/flai-snapshot.sh --local` and served on 127.0.0.1 by `scripts/release-server.sh`, with no token and no request to GitHub, and prints a failing step's output and the server's request log. The check that the latest published release installs from GitHub is `scripts/install-published-test.sh`, run by `install-published.yml` on push to main, daily, and by hand, and by no close-out.
