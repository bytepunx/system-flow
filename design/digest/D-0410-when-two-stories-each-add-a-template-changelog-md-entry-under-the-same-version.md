---
id: D-0410
title: When two stories each add a template/CHANGELOG.md entry under the same version, the one accepted second renumbers its heading to the next version at its sync and keeps both entries
kind: decision
sources: [S-0255, S-0467, S-0498, S-0502]
topics: [template, git]
paths: [template/CHANGELOG.md, template/template.yaml]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: [S-0418]
---

# D-0410 When two stories each add a template/CHANGELOG.md entry under the same version, the one accepted second renumbers its heading to the next version at its sync and keeps both entries

Stories that change `template/` each add an entry to `template/CHANGELOG.md`, and two in flight at once often head theirs with the same version. They settle it the same way each time: the story accepted second takes the next version.

- S-0255: "template/CHANGELOG.md conflicted (main took 1.0.42), so S-0255's entry and template.yaml are now 1.0.43."
- S-0467, on MS-0177 with S-0468: "S-0467 keeps both when it rebases after S-0468 lands."
- S-0498, on MS-0195 with S-0504: "Agreed that the story accepted second renumbers its heading and `template/template.yaml` to 1.0.119 at its sync."
- S-0502, on MS-0201 with S-0414: "whichever story lands second takes the next version"; main was 1.1.0, so S-0502's entry became 1.1.1.

Whether a story also bumps `template/template.yaml` is not settled the same way: S-0443 and S-0491 left it alone, since the release tooling treats the commit that sets the version as the release, while S-0255 and S-0498 bumped it with the entry. Renumber the version line too only when the story bumped it.

Sources: the Decisions of wip/archive/agents/S-0255.md, S-0467.md, S-0498.md, and S-0502.md.
