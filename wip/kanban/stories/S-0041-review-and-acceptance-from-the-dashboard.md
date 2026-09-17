---
id: S-0041
type: story
nature: feature
title: Review and acceptance from the dashboard
status: backlog
parent: E-0006
owner: alex
created: 2026-09-17T19:46:29Z
updated: 2026-09-17T19:46:29Z
transitions: []
tags: [dashboard, cli]
---

# S-0041 Review and acceptance from the dashboard

## Goal
Stories in review can be read, discussed, and accepted or sent back from the dashboard: the designer sees the diff of the story branch against main, the ticked acceptance criteria, the narrative, and open threads, then accepts (which runs `flai accept`) or moves the story back with a reason.

## Acceptance criteria
- [ ] A review page per story in review: acceptance criteria, narrative summary, open threads, the branch diff against main (files and hunks), and the release plan `flai release --dry-run` would produce
- [ ] Accept runs `flai accept <id> --by <designer>` with the token holder as `--by`, streams progress, and shows the resulting tags; send back runs `flai move <id> in-progress --reason`
- [ ] Failures (dirty tree, rebase conflict, release error) are shown verbatim from flai and leave the story in review
- [ ] Tests with a fake flai; docs/users and docs/operators updated

## Tasks

## Notes
Depends on S-0037 for the branch diff and S-0036 for authentication.
