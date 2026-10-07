---
title: Continuous improvement
updated: 2026-10-07
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
usage:                     # optional: what strategic agents spent on it
  strategic:
    - kind: analyzer
      seconds: 61
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1200
          output: 900
          cache_read: 40000
          cache_write: 3000
          cost: 0.25
---

# I-0001 golangci-lint on the host is v1 but the config is v2

## Description
What goes wrong, for whom, and what it costs.

## Instances
### 2026-09-15T17:00:00Z
Story: S-0008.
Where it happened and what was done about it.

### 2026-09-15T20:10:00Z
Report: design/analysis/2026-09-15-bottlenecks.md.
Occurred again.

## Impact
Evidence for the figures below, in any words.

- revenue_per_week: 1200
- penalty_per_week: 300
- time_lost_per_cycle: 4h

## Remediation
Proposed fix. Closed issues say what closed them.

Found by the analysis in [2026-09-15-bottlenecks.md](../analysis/2026-09-15-bottlenecks.md).

Story S-0009 remediates this issue, created from it at 2026-09-16T09:00:00Z.
```

`## Impact` is optional. `flai issue new` and `flai issue bump` write it from their impact flags (S-0224): `--revenue-per-week` and `--penalty-per-week`, each an amount of zero or more in `planning.currency`, `--time-lost-per-cycle`, a Go duration longer than zero, and `--evidence`, the words behind them. The section holds the evidence, then one `- key: value` line per figure, and is added before `## Remediation`. A bump replaces the figures it gives, keeps the others, and adds evidence the section does not hold already. A figure that does not parse is refused and nothing is written. `flai issue story` reads the section back: a line `revenue_per_week: <amount>`, `penalty_per_week: <amount>`, or `time_lost_per_cycle: <duration>`, as a list item or not, gives a cost of delay input. The first value of each key that parses is used. Every other line is evidence and is not read.

An issue names the analysis reports that found it (S-0224). `--report design/analysis/<file>.md`, on `flai issue new` or `bump`, puts a `Report: <path>.` line in the instance it records and adds `Found by the analysis in [<file>](../analysis/<file>).` to `## Remediation`, once per report, adding the section when the issue has none. The path is from the project root and must be a markdown file under `design/analysis/` other than its `README.md`; any other is refused and nothing is written. The report need not be written yet, since the analyzer files its issues while it writes it.

`usage` is optional and carries `strategic` alone: one entry per strategic agent kind, in the shape an item's has ([ADR-0083](../adrs/0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md)), and no agents' figures of its own (S-0227, [ADR-0100](../adrs/0100-an-analyzer-activity-s-usage-is-charged-evenly-to-the-issues-it-names-under.md)). flai writes it; an issue without it has no `usage` key. An analyzer activity's usage is split evenly between the issues it names and added under each one's `analyzer` entry: at the end of an analyzer run, the issues, open or closed, that name the report it wrote, by an instance's `Report:` line or a Remediation link; through `activity_log`, the issues the call names. An activity that names no issue is left in the analyzer's project strategic total ([strategic-agents.md](strategic-agents.md#what-the-analyzers-run-costs-s-0227)).

| Class | Meaning |
|-------|---------|
| `defect` | Something is wrong and produced a wrong result |
| `blocker` | Progress stopped until someone acted |
| `efficiency` | It works but costs time every occurrence |
| `impression` | A suspected problem without a measured instance yet |

## Summary

`summary.md` is a table of open issues: ID, class, title, count, average cost per occurrence, total cost (count times average), last reported. Sorted by total cost descending so the most expensive friction is first. Every `flai issue` command regenerates it; `flai check` warns when it is stale, and `flai prime` prints it after the conventions when anything is open.

The file is generated and committed, so two stories that each record or close an issue both change it, and git cannot merge them ([I-0074](../issues/I-0074-two-stories-that-each-record-or-close-an-issue-always-conflict-in-design-issues-summary-md-whose-updated-line-both-rewrite.md)). Sync and acceptance regenerate it rather than merge it ([ADR-0098](../adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md)). When the rebase `flai stream sync` or `flai accept` runs stops on `summary.md` alone, flai writes it again from the issue files in the worktree, stages it, and continues; its `updated:` line is then the time of the sync. When other paths conflict too, the rebase stops for the agent as for any conflict, `summary.md` among the paths, and `flai issue summary` in the worktree writes it once the others are resolved. The trial merge between open story branches leaves the file out, so no conflict is told over it. Two branches that bump the same issue still conflict in that issue's file, whose count and instances both sides add to (I-0092).

## Cadence

- Record an occurrence when it happens, not at the end of the story.
- Each instance names the story it was recorded for, in a `Story: S-nnnn.` line under its heading (S-0198). `flai issue new` and `bump` take it from `--story`, else `FLAI_STORY`, else `FLAI_AGENT` of the form `agent-S-nnnn`, else the `story/S-nnnn` branch checked out; outside a story the instance names none. `flai issue list --story S-nnnn` lists a story's issues. Issues live under the checkout's `design/`, so a story's own issues are committed on its branch and are in its worktree, not the main checkout, until it is accepted; `--story` reads them there. So that two stories worked at once do not take the same number, `flai issue new` numbers one past the highest issue on the main branch, in the main checkout, in every story worktree, committed or not, and on every story branch (S-0252).
- Recording from a story is one call (S-0275). In the story's worktree, `flai issue new`, `bump`, or `close` with `--commit`, or the MCP tool `issue_new`, `issue_bump`, or `issue_close` with `commit`, writes the issue, commits the issue's file and `summary.md` and nothing else on `story/S-nnnn` as `docs: [S-nnnn] record|bump|close I-nnnn <title>`, and adds them to the story's touches. The agent runs no `git add`, `git commit`, or `flai touches` after it. Off the story's branch, with no story, or mid-rebase it refuses and writes nothing. `flai adr new --commit` and the MCP tool `adr_new` with `commit` do the same for an ADR. The details are in [flai-cli.md](flai-cli.md#commands).
- The dashboard records through the host channel's `issue.new`, `issue.bump`, and `issue.close` (S-0275). They run the commands in the main checkout with `--autocommit`, which commits the issue's file and `summary.md` there on their own as `docs: record|bump|close I-nnnn <title>`, unless `dashboard.autocommit` is false, and widens no story's touches. An occurrence filed from the dashboard against a story names it in its instance but is committed on the main branch, not the story's.
- `flai issue new --report` records an analysis report's finding without counting it twice (S-0224). When an open issue has exactly the title given, flai bumps it, with the report, the impact, and the note, rather than open a second; when an instance of it already names that report, it leaves the issue as it is. The output says which it did, and `--json` gives it as `outcome`: `opened`, `bumped`, or `already recorded`. Without `--report`, `flai issue new` always opens a new issue. A finding the analyzer judges the same as an open issue under another title it bumps by ID with `flai issue bump <id> --report`. The MCP tools `issue_new` and `issue_bump` do the same as the commands, and write the issue in the recording story's worktree when it has one, else in the project; they regenerate `summary.md` and commit nothing unless given `commit`. How the analyzer uses them is in [strategic-agents.md](strategic-agents.md#how-the-analyzer-files-its-findings-s-0224).
- flai itself records the findings a story's close-out meets outside the story (S-0249, [ADR-0085](../adrs/0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md)). `scripts/close-out.sh` runs `flai check --strict --story S-nnnn --record-issues`, through `scripts/check.sh`: a finding outside the story is a note that does not stop the close-out, and each rule's are recorded in the open issue titled ``flai check finds `<rule>` outside the story at close-out``, opened with class `efficiency` when none is open and bumped otherwise, with an instance naming the story and the findings. The same story and findings are recorded once, so a rerun does not inflate the count. A `wip.overlap` is never recorded (S-0279, [ADR-0115](../adrs/0115-a-close-out-records-no-wip-overlap-and-notes-only-an-overlap-naming-its-story.md)): the close-out notes only one that names its story, and the pull hold, the `overlapped` change, and the trial merge report it to the agents who can act. The issue is written in the story's worktree and committed by the close-out, so it reaches the main branch with the story. The operator sees per rule how often, and in which stories, close-outs meet findings that are not theirs; the title is how flai finds the issue again, so renaming or closing it starts a new one.
- A story links an issue when its body names the issue's ID. That link is what ties them: `flai issue list --json` gives each issue's `stories` and `story`, the open story that links it. The line `flai issue story` adds to the issue's Remediation section is for a reader of the issue; flai does not link by it.
- The story's agent makes no story for an issue. The operator chooses at acceptance (decided 2026-10-02, replacing the summary presented at an epic's end; refined on TH-0074): a story's review page lists the open issues it recorded or bumped, checked, then every other open issue no open story links, unchecked. With any checked, Accept reads `Accept and Create Stories`; once the acceptance succeeds, a backlog story is made, and committed, for each checked issue, and none when it fails. The page names each story made and each issue refused.
- `flai issue story I-nnnn [--epic] [--story]`, or the MCP tool `issue_story`, makes that story: the issue's title; nature `remediation` for a `defect` or `blocker`, `improvement` for `efficiency` or `impression`; a goal that links the issue's document and carries its `## Remediation` text as the recommended solution, or asks for one; and a last criterion closing the issue with `flai issue close`. When the issue names analysis reports, in its instances' `Report:` lines or its Remediation's `Found by the analysis in` links, the goal links each from the story's file, "which the analysis in [<file>](…) found", the Notes say that the analyzer found the issue, and the report links are left out of the solution text (S-0224). A closed issue, or one an open story links, is refused. An agent runs it when the operator asks, or for an issue `flai check` warns about.
- That story is a draft (`draft: true`, S-0203, [ADR-0074](../adrs/0074-work-items-carry-planning-data-a-story-s-draft-flag-an-epic-s-or-story-s-cost.md)): an agent's words until the operator finalizes it, so it cannot go to ready before. It carries the cost of delay inputs the issue gives, in `cost_of_delay.inputs` with `by: flai` and `at` the time it was made:
  - `time_lost_per_cycle` is the issue's `cost` × `count` ÷ the cycles of `planning.cycle` (default `168h`) since `first_reported`, to a tenth and at least one, rounded to the minute (to the second under a minute). A `cost` that is not a duration longer than zero, or a `first_reported` that is not a timestamp (counted as one cycle), is said in the Notes.
  - `revenue_per_week`, `penalty_per_week`, and `time_lost_per_cycle` from the issue's `## Impact`, as the analyzer's `flai issue new` and `bump` wrote them, a bump's figure in place of the one it replaced. An Impact time lost takes precedence over the derived one. A value that does not parse is left out.
  - With no inputs, the story has no `cost_of_delay`. Its `## Notes` say how each input was set and what was left out.
- It carries the issue's strategic usage (S-0227): each of the issue's `usage.strategic` entries is charged to the story under its kind, and to every item above it, up to its epic, as a planner's charge is summed, and its Notes say what was carried. The issue keeps its own entries. `flai stats` counts them once: on the issue until a story is made from it, as the line below names it, then on the story ([metrics.md](metrics.md#strategic-usage)). A story that merely links the issue carries nothing.
- The issue's `## Remediation` then gets a last paragraph, `Story S-nnnn remediates this issue, created from it at <ts>.`, and its `updated` is set; the section is added when the issue has none. The story is named by ID, because its file moves to `wip/archive/` on acceptance. With `--autocommit` the issue is committed with the story and its epic, so the dashboard's acceptance flow leaves nothing uncommitted; with `--story` the issue is in that story's worktree and is that story's to commit. MCP `issue_story` commits neither.
- `flai check` warns with `issues.no-story` about an open issue first reported longer ago than `issues.story_after` in `system-flow.yaml` (a Go duration, `168h` when unset, `0` turns it off) that no open story links. Age, rather than the recording story, finds the issues from before instances named their story.
- Remediation is a story with nature `remediation` or `improvement` that closes the issue.

## Metrics

Issue counts and costs are inputs to the process-improvement loop in [workflow.md](workflow.md) alongside the flow metrics: they show where time goes outside the value stream.
