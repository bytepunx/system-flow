---
id: S-0231
type: story
nature: feature
title: Author LICENSE.md file
status: backlog
owner: alex
created: 2026-10-02T12:24:10Z
updated: 2026-10-02T12:24:10Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: high
---
# S-0231 Author LICENSE.md file

## Goal

Write an adaptation of polycore shield that grants free usage to:
 - 501(c)(3)s
 - Educators and students at State-funded educational institutions
 - Security researchers analyzing the system
 - Start-ups earning under $1M USD in annual revenue and having less than $15M in cumulative funding

That also prohibits any class of persons (whether granted free use or not) from using this software as a basis for hosted products or services. The intent being to shield system-flow or a custom fork of it from being offered as SaaS or licensed software by another entity.

## Acceptance criteria
- [ ] The license grants free use for the classes mentioned above
- [ ] The license makes it clear that any commercial use falling outside of the permitted criteria requires a paid license
- [ ] The license makes it clear that offering licenses to system-flow or a fork of it, or offering system-flow or a fork of it as a service is prohibited
- [ ] The license requires that a copy of it be distributed in all forms
- [ ] The license file is shipped with the docker image and is visible under Host in its own submenu and in the CLI binary and is printed when the command `flai license` is issued.

## Tasks

## Notes
