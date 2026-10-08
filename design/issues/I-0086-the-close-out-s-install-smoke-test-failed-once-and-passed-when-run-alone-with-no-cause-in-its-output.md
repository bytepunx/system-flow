---
id: I-0086
title: The close-out's install smoke test failed once and passed when run alone, with no cause in its output
class: efficiency
status: open
count: 6
cost: 7m
first_reported: 2026-10-06T10:22:10Z
last_reported: 2026-10-08T06:55:19Z
updated: 2026-10-08T06:55:19Z
---

# I-0086 The close-out's install smoke test failed once and passed when run alone, with no cause in its output

## Description
The close-out's install smoke test failed once and passed when run alone, with no cause in its output

## Instances

### 2026-10-06T10:22:10Z
Story: S-0283.
S-0283's close-out stopped at smoke's last sub-step, scripts/install-test.sh, after gofmt, vet, lint, the full go test -race, the template render and check, and the markdown lint passed; its output showed self-upgrade --check reaching GitHub for release 1.31.4 but no failure reason. The verifier ran scripts/install-test.sh alone once: exit 0 after over five minutes. Likely the network or a timeout; the cause is not confirmed.

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

## Remediation

Story S-0291 remediates this issue, created from it at 2026-10-06T10:31:55Z.
