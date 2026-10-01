---
id: T-0636
type: task
nature: remediation
title: The design, conventions, and user guide say where flai lints, and the issues are closed
status: done
parent: S-0179
owner: arobson
created: 2026-10-01T08:23:14Z
updated: 2026-10-01T08:52:08Z
transitions:
  - to: ready
    at: 2026-10-01T08:23:32Z
    by: claude-opus-5-5
  - to: in-progress
    at: 2026-10-01T08:46:17Z
    by: claude-opus-5-5
  - to: done
    at: 2026-10-01T08:52:08Z
    by: claude-opus-5-5
stream: S-0179
tags: []
touches: [design/system/flai-cli.md, design/adrs, design/conventions/tooling.md, template/root/design/conventions/tooling.md, template/template.yaml, template/CHANGELOG.md, docs/users/flai.md, design/issues, wip/archive/agents]
usage:
  source: log
  seconds: 351
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 83
      output: 51995
      cache_read: 9077423
      cache_write: 98546
      cost: 3.6441
---
# T-0636 The design, conventions, and user guide say where flai lints, and the issues are closed

## Work

Record the decision to embed a subset of markdownlint in an ADR; describe the lint in `design/system/flai-cli.md`, the tooling convention (template first), and `docs/users/flai.md`; close I-0027 and I-0043 with what fixed them; confirm `scripts/lint-md.sh` passes on main.

## Done when

- [x] The ADR, design, conventions, and user guide describe where flai lints what it writes
- [x] I-0027 and I-0043 are closed
- [ ] `make test`, `make integration`, `make smoke`, and `flai check --strict` pass

## Notes

Part of S-0179.

`make test` and `make integration` pass, and `make smoke` renders and checks the template. Its repository check stops on warnings that are not this story's (done epics not archived, threads answered on archived stories, S-0185's overlap) and on the two archived narratives this story fixes, which clear at acceptance; the last box stays unticked for that reason.
