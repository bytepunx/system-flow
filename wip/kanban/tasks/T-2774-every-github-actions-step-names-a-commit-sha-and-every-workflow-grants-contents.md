---
id: T-2774
type: task
nature: remediation
title: Every GitHub Actions step names a commit SHA, and every workflow grants contents read at the top level and writes only on the jobs that need them
status: backlog
parent: S-0630
owner: alex
created: 2026-11-09T21:58:15Z
updated: 2026-11-09T21:58:15Z
transitions: []
stream: S-0630
tags: [ci, security]
touches: [".github/workflows/flai.yml", ".github/workflows/flaiover.yml", ".github/workflows/system-flow-check.yml", ".github/workflows/release-flai.yml", ".github/workflows/release-flaiover.yml", ".github/workflows/install-published.yml", template/root/.github/workflows/system-flow-check.yml]
---
# T-2774 Every GitHub Actions step names a commit SHA, and every workflow grants contents read at the top level and writes only on the jobs that need them

## Work

- In each of the six files under `.github/workflows/`, replace every `uses: <owner>/<action>@<tag>` with the full 40-character commit SHA that the tag names today, and keep the version as a trailing comment: `uses: actions/checkout@<sha> # v7.x.y`. Look each SHA up with `git ls-remote https://github.com/<owner>/<action> 'refs/tags/<tag>*'`. For an annotated tag, take the peeled `^{}` commit. `.github/dependabot.yml`'s `github-actions` entry updates SHA pins and their comments, so this needs no other change.
- `release-flai.yml:81`: replace GoReleaser's `version: "~> v2"` with `v2.18.1`, the version `scripts/install-tools.sh` and `system-flow-check.yml` pin.
- `flai.yml`, `flaiover.yml`, `system-flow-check.yml`: add a top-level `permissions:` block with `contents: read`. None of them writes.
- `release-flai.yml:8-11` and `release-flaiover.yml:9-11`: make the top-level block `contents: read`. Move each write grant (`contents`, `packages`, `id-token`, `attestations`) to the jobs whose steps use it. `release-flaiover.yml` already gives its jobs blocks at lines 28, 88, and 168, so check that each one names every write that job's steps need.
- `install-published.yml`: pin its actions. Keep its `permissions:` block, and check that it grants no write.
- `template/root/.github/workflows/system-flow-check.yml`: add the same top-level `permissions: contents: read`. Keep its version tags: a project made from the template has no `dependabot.yml` to move a SHA, and its pins are the project's to choose.

It waits for no task: no other task of S-0630 changes these files.

## Done when

- `grep -nE 'uses: [^ ]+@' .github/workflows/*.yml` lists only `@` followed by 40 hexadecimal characters.
- Each of the six workflows has a top-level `permissions:` block, and no top-level block grants `write`.
- `release-flai.yml` pins GoReleaser at `v2.18.1`.
- The template's `system-flow-check.yml` has a top-level `permissions: contents: read`.
- Each changed file parses as YAML, and `scripts/template-test.sh` passes.

## Notes
