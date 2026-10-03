---
title: Continuous improvement
updated: 2026-10-03
status: active
topics: [cli, conventions]
---

# Continuous improvement

`design/issues/` records the friction, defects, blockers, and inefficiencies that agents and operators hit while working, with a count and an average cost per occurrence, so remediation can be prioritised by evidence. Decided in [ADR-0014](../adrs/0014-design-issues.md); the agent-facing rules are in [conventions/continuous-improvement.md](../conventions/continuous-improvement.md).

## Layout

```text
design/issues/
├── README.md              # what the folder is
├── summary.md             # generated table of open issues
└── I-0001-slug.md          # one file per issue
```

## Issue file

```yaml
---
id: I-0001
title: golangci-lint on the host is v1 but the config is v2
class: efficiency          # defect | blocker | efficiency | impression
status: open               # open | closed
count: 3
cost: 5m                   # average wall-clock per occurrence, Go duration
first_reported: 2026-09-15T17:00:00Z
last_reported: 2026-09-15T20:10:00Z
updated: 2026-09-15T20:10:00Z
---

# I-0001 golangci-lint on the host is v1 but the config is v2

## Description
What goes wrong, for whom, and what it costs.

## Instances
### 2026-09-15T17:00:00Z
Story: S-0008.
Where it happened and what was done about it.

## Remediation
Proposed fix. Closed issues say what closed them.
```

| Class | Meaning |
|-------|---------|
| `defect` | Something is wrong and produced a wrong result |
| `blocker` | Progress stopped until someone acted |
| `efficiency` | It works but costs time every occurrence |
| `impression` | A suspected problem without a measured instance yet |

## Summary

`summary.md` is a table of open issues: ID, class, title, count, average cost per occurrence, total cost (count times average), last reported. Sorted by total cost descending so the most expensive friction is first. Every `flai issue` command regenerates it; `flai check` warns when it is stale, and `flai prime` prints it after the conventions when anything is open.

## Cadence

- Record an occurrence when it happens, not at the end of the story.
- Each instance names the story it was recorded for, in a `Story: S-nnnn.` line under its heading (S-0198). `flai issue new` and `bump` take it from `--story`, else `FLAI_STORY`, else `FLAI_AGENT` of the form `agent-S-nnnn`, else the `story/S-nnnn` branch checked out; outside a story the instance names none. `flai issue list --story S-nnnn` lists a story's issues. Issues live under the checkout's `design/`, so a story's own issues are committed on its branch and are in its worktree, not the main checkout, until it is accepted; `--story` reads them there.
- A story links an issue when its body names the issue's ID. That link is the only record of an issue's story: the issue file is not changed, so making a story from the dashboard leaves nothing uncommitted outside `wip/`. `flai issue list --json` gives each issue's `stories` and `story`, the open story that links it.
- The story's agent makes no story for an issue. The operator chooses at acceptance (decided 2026-10-02, replacing the summary presented at an epic's end; refined on TH-0074): a story's review page lists the open issues it recorded or bumped, checked, then every other open issue no open story links, unchecked. With any checked, Accept reads `Accept and Create Stories`; once the acceptance succeeds, a backlog story is made, and committed, for each checked issue, and none when it fails. The page names each story made and each issue refused.
- `flai issue story I-nnnn [--epic] [--story]`, or the MCP tool `issue_story`, makes that story: the issue's title; nature `remediation` for a `defect` or `blocker`, `improvement` for `efficiency` or `impression`; a goal that links the issue's document and carries its `## Remediation` text as the recommended solution, or asks for one; and a last criterion closing the issue with `flai issue close`. A closed issue, or one an open story links, is refused. An agent runs it when the operator asks, or for an issue `flai check` warns about.
- `flai check` warns with `issues.no-story` about an open issue first reported longer ago than `issues.story_after` in `system-flow.yaml` (a Go duration, `168h` when unset, `0` turns it off) that no open story links. Age, rather than the recording story, finds the issues from before instances named their story.
- Remediation is a story with nature `remediation` or `improvement` that closes the issue.

## Metrics

Issue counts and costs are inputs to the process-improvement loop in [workflow.md](workflow.md) alongside the flow metrics: they show where time goes outside the value stream.
