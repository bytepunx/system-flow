---
id: S-0060
type: story
nature: feature
title: Operators create ADRs from the dashboard
status: in-progress
parent: E-0006
owner: alex
created: 2026-09-19T05:34:10Z
updated: 2026-09-19T09:39:14Z
transitions:
  - to: ready
    at: 2026-09-19T06:37:41Z
    by: alex
  - to: in-progress
    at: 2026-09-19T09:38:09Z
    by: system-flow
tags: [dashboard, cli]
touches: [flaiover, flai/cmd, design/adrs]
---

# S-0060 Operators create ADRs from the dashboard

## Goal
An operator records a new architecture decision from the ADRs page by writing the decision in markdown. The number, file name, front matter, date, and the row in the ADR index are supplied behind the scenes, so the result is a correctly formed ADR without copying the template by hand.

## Acceptance criteria
- [ ] The ADRs page has a "new ADR" action that opens a form: a title, and the four sections the ADR template has (context, decision, consequences, alternatives considered) as markdown with the explorer's preview beside it. The form takes its sections from `design/adrs/0000-template.md`, so a project that changed its template gets its own
- [ ] The operator may mark which existing ADRs the new one supersedes or refines, chosen from a list. A superseded ADR gets its `superseded_by` set, which is the one edit the conventions allow to an accepted ADR; a refinement is noted in the index row the way the index already does it
- [ ] Creation goes through a new `flai adr new` (there is no ADR command today): it takes the next number from the files present, not from a counter in a document, writes `NNNN-slug.md` with `id`, `title`, `status`, `date`, `supersedes`, and `superseded_by`, adds the row to `design/adrs/README.md`, and takes the body on standard input so a script can do the same
- [ ] The status is the operator's choice between `proposed` and `accepted`, defaulting to `proposed`; a second action on a proposed ADR's page accepts it, setting the status and the date, through flai. An accepted ADR stays immutable: neither the editor (S-0040) nor this form changes its body afterwards
- [ ] The result is validated before it is kept: `flai check` runs with the new file and the index row in place, and anything it would report is shown in the form with the text kept and nothing left behind
- [ ] A created ADR is committed on its own unless `dashboard.autocommit: false`: the new file, the index, and any `superseded_by` edits in one commit, `docs:` with the ADR's number, the operator's identity as author and the dashboard's trailer. Nothing is pushed
- [ ] After creating, the operator lands on the new ADR in the explorer, and the ADRs list shows it with its status and its place in the supersession chain without a reload
- [ ] `flai check` gains what it needs to keep the index honest: a warning when an ADR file has no row in `design/adrs/README.md` or a row has no file. The stale "next is 0015" line in this repository's `decisions.md` project additions is removed, since the number now comes from the files
- [ ] Tests: `flai adr new` with real git (numbering with gaps, slug, index row, supersedes, check refusal leaving nothing behind, commit), the endpoint, the form as a component; `design/system/flaiover-dashboard.md`, `flai-cli.md`, `documentation-standard.md` or wherever ADR format is described, the decisions convention (template baseline first), and `docs/users` updated

## Tasks
- T-0233 flai adr new and flai adr accept: number from the files, slug, front matter, index row, supersedes and refines, body on standard input, checked before kept, committed on its own
- T-0234 flai check keeps the ADR index honest, and an accepted ADR's body cannot be edited
- T-0235 The decisions convention, the ADR index text, and the designs say ADRs are made with flai adr new
- T-0236 Dashboard: endpoints, the new ADR form on the ADRs page, and accepting a proposed ADR
- T-0237 Try it in a browser against a scratch project, user documentation, criteria

## Notes
Raised by the operator on 2026-09-19 with S-0059: "In the ADRs, operators should be able to create new ADRs", with the same aim: boilerplate behind the scenes, the operator writes the guidance in markdown, and eventually all project activity is driven from the dashboard.

What exists. ADRs are made by copying `design/adrs/0000-template.md` by hand, taking the next number by looking, and adding a row to the table in `design/adrs/README.md` by hand; agents did this for ADR-0021 to ADR-0024 in this session. flai has no `adr` command. `flai check` validates an ADR's front matter (`adr.id` against the file name, `adr.front-matter` for id, title, status, date) and does not look at the index. The dashboard's `/adrs` page lists ADRs with status, date, and the supersession chain (S-0012) and has no write action. The S-0040 editor saves existing documents only, and deliberately does not create them.

Because there is no flai command yet, this story is larger on the flai side than S-0059: the command, the index maintenance, and the check rule come first, and the dashboard form is a client of them (ADR-0016, ADR-0023).

The editor from S-0040 does not refuse edits to an accepted ADR today; that was left as an open question there. The fourth criterion settles it for ADRs: accepted means immutable, in the editor as well. If that is built here, say so in S-0040's follow-up notes; if the operator prefers a separate story, split it out.

Whether agents should also use `flai adr new` from now on is a convention change (`decisions.md` says to copy the template): yes, and the baseline should say so in this story, template first.

Sits with S-0059; whichever is built second reuses the first's form and endpoint shape.
