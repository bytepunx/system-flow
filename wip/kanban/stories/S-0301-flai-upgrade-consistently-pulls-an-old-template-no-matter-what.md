---
id: S-0301
type: story
nature: remediation
title: flai upgrade consistently pulls an old template no matter what
status: backlog
owner: alex
created: 2026-10-06T22:46:13Z
updated: 2026-10-06T22:46:13Z
transitions: []
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0301 flai upgrade consistently pulls an old template no matter what

## Goal

### flai upgrade

flai upgrade should default to pulling the latest available version tag from the template repository it is configured for.

If there is a ref passed, it should overwrite the lock file.
If there is no ref passed and the operator either set a different version in system-flow.yaml or there is a newer version, flai should prompt the user about their intent rather than silently reverting back to the template version in the lockfile.

### Description

In another project, the template defaulted to 1.0.18 (not 1.0.60 which was the latest at the time of this writing). When I change the version manually in system-flow.yaml, it reverts (possibly due to the version in the lock file). But even when providing a version reference to the tag 1.0.60, I still see it revert to 1.0.18.

## Acceptance criteria
- [ ] flai correctly updates system-flow files to reflect the new version supplied in the --ref argument before completing the upgrade
- [ ] flai defaults to the latest available version of the template when creating a new project or upgrading an existing project not created with flai
- [ ] flai assumes that `flai upgrade` should find and use the latest available tag, update the system-flow files, and then perform the upgrade
- [ ] if there is a conflict, flai's CLI must prompt the operator with the version they intend before proceeding rather than reverting to the version found in the lock file.

## Tasks

## Notes
