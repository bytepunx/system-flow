---
id: T-0937
type: task
nature: feature
title: Design and user guide say how the analyzer files findings as issues and how the issue step carries their impact and report
status: done
parent: S-0224
owner: alex
created: 2026-10-05T05:45:21Z
updated: 2026-10-06T21:24:49Z
transitions:
  - to: ready
    at: 2026-10-06T21:20:04Z
    by: agent-S-0224
  - to: in-progress
    at: 2026-10-06T21:20:05Z
    by: agent-S-0224
  - to: done
    at: 2026-10-06T21:24:49Z
    by: agent-S-0224
stream: S-0224
tags: [flai, docs]
touches: [design/system/continuous-improvement.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, docs/operators/index.md, template/CHANGELOG.md]
after: [T-0922, T-0930]
usage:
  source: log
  seconds: 284
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 91
      output: 41961
      cache_read: 4139792
      cache_write: 166601
      cost: 2.6116
---
# T-0937 Design and user guide say how the analyzer files findings as issues and how the issue step carries their impact and report

## Work

Waits for T-0922 and T-0930, and through them T-0918 and T-0925: it describes what they built, as built.

- `design/system/continuous-improvement.md`: the `## Impact` section is written by `flai issue new` and `bump` with the impact flags, an issue names the reports that found it in its instances and `## Remediation`, `--report` bumps an open issue of the same title, and the story `flai issue story` makes links the report.
- `design/system/strategic-agents.md`: the analyzer's section says how it files a finding, deduplicates it, and links each issue from its report, and that its guard refuses story creation.
- `design/system/flai-cli.md` and `docs/users/flai-reference.md`: the new flags of `flai issue new` and `bump`, and the MCP tools `issue_new` and `issue_bump`.
- `docs/users/flai.md`: how an operator turns an analyzer's issue into a draft story with its cost of delay in place.

## Done when

- Each document says what the code does, with no stale sentence such as the analyzer of S-0224 being yet to write the Impact format
- `flai check --strict` and the markdown lint pass on the changed documents

## Notes
