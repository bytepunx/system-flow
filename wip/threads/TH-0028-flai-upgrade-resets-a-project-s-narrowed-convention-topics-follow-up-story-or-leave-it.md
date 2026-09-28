---
id: TH-0028
title: "flai upgrade resets a project's narrowed convention topics: follow-up story, or leave it?"
anchor:
  path: wip/kanban/stories/S-0134-conventions-design-tech-files-and-adrs-carry-topics-on-the-file-and-on-headings-and-flai-check-keeps-them-honest.md
  item: S-0134
status: resolved
participants: [agent-S-0134, alex]
created: 2026-09-27T04:09:50Z
updated: 2026-09-28T22:04:56Z
---

# TH-0028 flai upgrade resets a project's narrowed convention topics: follow-up story, or leave it?

On wip/kanban/stories/S-0134-conventions-design-tech-files-and-adrs-carry-topics-on-the-file-and-on-headings-and-flai-check-keeps-them-honest.md.

## Entries

### 2026-09-27T04:09:50Z agent-S-0134
S-0134 is otherwise ready for review: every criterion is met, and all tiers and `flai check --strict` pass. One open question holds the move.

`flai upgrade` replaces everything above a convention's marker, front matter included. So when a project narrows a baseline convention's `topics` (say `code-quality.md` to `[code]`), or adds a heading comment above the marker, the next upgrade puts back the template's `[all]`. This repository doesn't hit it: here the baseline is edited in `template/` first and copied, so narrowing happens in the template. Every other project would hit it. `design/system/conventions.md` now states this as a caveat.

Recommended: **A.** Create a backlog story under E-0010: the upgrade merge keeps the project's `topics` key when it has one, and the template's `topics` otherwise. Heading comments above the marker stay the template's; a project narrows those in its additions. Nothing changes in S-0134.
**B.** Leave it as a documented limitation; no story.
**C.** Fix it inside S-0134 before review (touches `flai/internal/upgrade`, about an hour).

### 2026-09-27T04:11:59Z alex
what is the role of the flai upgrade command and why does it involve replacing everything above the marker?

### 2026-09-28T21:19:12Z agent-S-0134
**What `flai upgrade` is for.** A project made from the template (`flai new`) starts with the template's files: `CLAUDE.md`, the baseline conventions, scripts, the design skeleton. When the template publishes a new version (a new or changed rule, a new script), `flai upgrade` brings an existing project to it without overwriting what the project made its own (ADR-0015). It renders the new template in memory and sorts every file into one of five classes (`flai/internal/upgrade`):

- **add**: new in the template, missing in the project: written.
- **same**: already identical: nothing to do.
- **replace**: the project never touched it since it was applied (its hash still matches `system-flow.lock.yaml`): overwritten with the new version.
- **merge**: a marker file (the conventions and `CLAUDE.md`): see below.
- **conflict**: the project edited a non-marker file: reported, the operator chooses keep, replace, or diff.

**Why the marker.** Conventions are the one kind of file both sides need to write: the template keeps improving the baseline rules, and every project adds its own rules. Without a split, any project that added a single rule would turn every later baseline change into a conflict, and projects would drift off the baseline. So ADR-0013 split each convention in two with the marker comment: above it is the template's baseline, below it the project's additions. The ownership rule is simple and needs no diffing: upgrade takes the template's text above the marker and the project's text from the marker down (`merge`, a few lines). Precedence mirrors it: project additions below the marker win over the baseline above it, so a project never needs to edit the baseline to change a rule; it overrides it below.

**Where topics fall.** Front matter sits at the top of the file, above the marker, so it counts as the template's. That was fine while front matter was only `title` and `order`, both the template's to decide. `topics` is the first front-matter key a project might reasonably want to set for itself on a baseline file (narrowing `code-quality.md` to `[code]`, say), and upgrade would put the template's `[all]` back. Same for a heading comment on a baseline heading.

So the choice is still A, B, or C from the first entry. My recommendation stays **A**: a backlog story under E-0010 where merge keeps the project's `topics` key when it has one (a narrowly scoped exception to "above the marker is the template's"); until then a project can narrow in its additions. Say which, or ask more.

### 2026-09-28T21:19:33Z alex
C

### 2026-09-28T22:04:56Z agent-S-0134
Done as C, in S-0134 (ec591b6, T-0494). `system-flow.lock.yaml` now records the topics the template gave each convention. On upgrade, a project's `topics` that differ from that record count as the project's and are kept. Topics that match the record, or no topics at all, take the new template's. So a template that later narrows a convention still reaches projects that left it alone. When nothing is recorded (an older lock, or no lock), a project's topics are kept. Heading comments above the marker stay the template's; a project narrows a heading in its additions. Docs are updated. `make flai-test` passes.

Resolved: Answer C: fixed in S-0134 (T-0494, ec591b6)
