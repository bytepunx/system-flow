---
id: T-1338
type: task
nature: remediation
title: install.sh retries a GitHub API call the network drops and resolves the latest release from a small page
status: done
parent: S-0291
owner: alex
created: 2026-10-08T07:59:51Z
updated: 2026-10-08T08:24:31Z
transitions:
  - to: ready
    at: 2026-10-08T08:14:00Z
    by: agent-S-0291
  - to: in-progress
    at: 2026-10-08T08:14:01Z
    by: agent-S-0291
  - to: done
    at: 2026-10-08T08:24:31Z
    by: agent-S-0291
stream: S-0291
tags: [install, network]
touches: [install.sh, flai/cmd/installsh_test.go]
usage:
  source: log
  seconds: 630
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 45
      output: 181
      cache_read: 1279857
      cache_write: 70404
      cost: 0.6066
---
# T-1338 install.sh retries a GitHub API call the network drops and resolves the latest release from a small page

## Work

I-0086's instances show the same fault each time: while `install.sh` lists `releases?per_page=50` (about 1 MB), GitHub drops the connection with curl (92) `HTTP/2 stream 1 was not closed cleanly` or curl (56) `unexpected eof while reading`, and `install.sh` stops at "could not list releases". TH-0365 measured 2 of 5 fetches of that listing dropped and 3 of 3 of `per_page=5` coming back.

- In `install.sh`, make `api()` retry a call that fails on the network: up to three attempts with a short pause, in a POSIX `sh` loop, not `--retry-all-errors`, which older curl lacks. Do not retry an HTTP error (curl exit 22): a 401 or 404 still fails at once with its hint. Print one line per retry naming curl's exit code.
- Resolve the latest release from a small page, such as `per_page=10`, and ask for the next page only when that page holds no `flai/v*` tag. Dashboard releases are tags without a GitHub release, so the first page holds a flai release in practice.
- Add no environment variable. If one is unavoidable, it needs a row in `docs/operators/settings.md` § Read by install.sh, which `flai/cmd/settings_doc_test.go` checks, and that file joins this task's touches.
- Add `flai/cmd/installsh_test.go`. It runs `sh ../../install.sh` with `FLAI_API` set to an `httptest` server that answers as the releases API does. The server cuts the first releases listing off mid-body by hijacking the connection, then serves a small archive and its `checksums.txt`. A Go test is used because the story's close-out selects the Go tiers by `flai/**`, and `install.sh` alone selects none of them. Skip the test where `sh` or `curl` is missing.

Waits for nothing: no other task of this story changes these paths.

## Done when

- `flai/cmd/installsh_test.go` fails against the current `install.sh`, which stops at "could not list releases", and passes with the change.
- A 404 from the stand-in still makes `install.sh` fail at once with its "private repository?" hint, which the test also checks.
- `flai test install.sh flai/cmd/installsh_test.go` passes.

## Notes
