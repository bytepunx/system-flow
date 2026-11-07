---
id: TH-0863
title: "S-0564: remove the scratch repository /home/alex/git/flai-scratch-import now that S-0001 is accepted?"
anchor:
  path: wip/kanban/stories/S-0564-the-import-and-analysis-cadence-formulas-are-each-run-once-in-this-repository.md
  item: S-0564
status: open
participants: [agent-S-0564]
created: 2026-11-07T08:52:21Z
updated: 2026-11-07T08:52:21Z
---

# TH-0863 S-0564: remove the scratch repository /home/alex/git/flai-scratch-import now that S-0001 is accepted?

On wip/kanban/stories/S-0564-the-import-and-analysis-cadence-formulas-are-each-run-once-in-this-repository.md.

## Entries

### 2026-11-07T08:52:21Z agent-S-0564
**Both runs are done and accepted:** S-0001 (import) at `049cf36` in flai-scratch-import, and S-0584 (cadence) at `07fa45be7`. S-0564's Notes record them.

**Question:** may the scratch repository `/home/alex/git/flai-scratch-import` be removed? Its Notes say to remove it only on your word.

**Recommended:** remove it yourself, since it is outside this repository and flai serve serves it as a project. First turn off its host actions, then delete the folder:

```sh
cd /home/alex/git/flai-scratch-import && flai serve disable agent
rm -rf /home/alex/git/flai-scratch-import
```

**Alternative:** keep it as a fixture for the next import run, such as a check of I-0300 or I-0302 once fixed.

S-0564 does not wait on this. It goes to review either way.
