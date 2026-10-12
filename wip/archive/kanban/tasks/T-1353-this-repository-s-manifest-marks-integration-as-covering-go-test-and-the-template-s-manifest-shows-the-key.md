---
id: T-1353
type: task
nature: improvement
title: This repository's manifest marks integration as covering go-test, and the template's manifest shows the key
status: done
parent: S-0342
owner: alex
created: 2026-10-08T08:06:00Z
updated: 2026-10-08T21:07:36Z
transitions:
  - to: ready
    at: 2026-10-08T21:03:05Z
    by: agent-S-0342
  - to: in-progress
    at: 2026-10-08T21:03:05Z
    by: agent-S-0342
  - to: done
    at: 2026-10-08T21:07:36Z
    by: agent-S-0342
stream: S-0342
tags: [template]
touches: [system-flow.yaml, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md]
after: [T-1345]
usage:
  source: log
  seconds: 271
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 22
      output: 6551
      cache_read: 1073193
      cache_write: 43027
      cost: 0.6139
---
# T-1353 This repository's manifest marks integration as covering go-test, and the template's manifest shows the key

## Work

- In `system-flow.yaml`, give the `integration` tier `covers: [go-test]`: `scripts/integration.sh` runs `go test -race -count=1 ./...`, every package and every test, so it holds everything `go-test` (`-race -short` on the changed packages) runs. Edit the `tests` list by hand, keeping its comments, or with `flai manifest set tests='<JSON>'`.
- In `template/root/system-flow.yaml.tmpl`, the template's `test` and `integration` scripts are placeholders that run different things, so mark nothing as covered. Add a comment on the `integration` tier saying that a project whose integration script runs its behaviour tests too writes `covers: [test]` there, so that `flai verify` runs them once. The file the story declares, `template/root/system-flow.yaml`, does not exist; this is the template's manifest.
- Add a `template/CHANGELOG.md` entry for S-0342 under the next version, in the form the entries above it use.
- Run `scripts/flai.sh check --strict` with the tree's flai, built after T-1345, so the `covers` key validates, and `scripts/template-test.sh`.

It waits for T-1345: before it, `flai manifest set` refuses the key and `flai check` cannot validate it. An older flai reads the manifest leniently, so the installed flai ignores the key and runs both tiers until it is released.

## Done when

- `system-flow.yaml`'s `integration` tier has `covers: [go-test]`, and `scripts/flai.sh check --strict` reports no `manifest.tests` finding.
- The template's manifest carries the comment, `scripts/template-test.sh` passes, and `template/CHANGELOG.md` has the entry.

## Notes

Layer 2 of S-0342's plan, beside T-1349.
