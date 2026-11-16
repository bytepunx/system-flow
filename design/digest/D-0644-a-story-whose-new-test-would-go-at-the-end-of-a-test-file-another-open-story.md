---
id: D-0644
title: A story whose new test would go at the end of a test file another open story also appends to puts it in a test file of its own, since the two appends conflict at the trial merge
kind: decision
sources: [S-0748]
topics: [go, git, code]
paths: [flai/internal/hostapi/people_test.go, flai/internal/hostapi/settings_roles_test.go]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0644 A story whose new test would go at the end of a test file another open story also appends to puts it in a test file of its own, since the two appends conflict at the trial merge

Two open stories that each append a test at the end of the same Go test file conflict when flai stream sync trial-merges their branches, even though neither test touches the other's lines: both change the file's last lines. The story that finds the conflict moves its own test into a new file in the same package, which removes the overlap without either story waiting on the other in `after:`.

- S-0748, Decisions: "The hostapi to-condition test lives in its own file, settings_roles_test.go, not at the end of people_test.go, where S-0747 also appends and a trial merge conflicted (MS-0580)."
- MS-0580: the sync of S-0748 reported that `story/S-0747` and `story/S-0748` conflicted in `flai/internal/hostapi/people_test.go`. S-0748 committed 0178058, `test: [S-0748] the settings page's to-condition test moves to its own file`, and the next sync closed the conversation: "story/S-0747 and story/S-0748 merge cleanly".

Go puts every `_test.go` file in a directory into one package, so a new file needs no wiring. Keep its helper names distinct from the other files', since two files that each declare the same helper collide when both land.

Sources: the Decisions of wip/archive/agents/S-0748.md; wip/messages/MS-0580-t-3238-of-s-0748-changed-paths-s-0747-s-claim-covers.md.
