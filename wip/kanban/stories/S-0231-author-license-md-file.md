---
id: S-0231
type: story
nature: feature
title: Author LICENSE.md file
status: review
owner: alex
created: 2026-10-02T12:24:10Z
updated: 2026-10-02T16:04:15Z
transitions:
  - to: ready
    at: 2026-10-02T12:24:10Z
    by: alex
  - to: in-progress
    at: 2026-10-02T12:46:56Z
    by: agent-S-0231
  - to: review
    at: 2026-10-02T13:15:04Z
    by: agent-S-0231
tags: [flai]
touches: [LICENSE.md, ".dockerignore", flai/cmd, flai/internal/license, flai/.goreleaser.yaml, flaiover/Dockerfile, flaiover/package.json, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/SiteMenu.svelte, flaiover/src/lib/components/SiteMenu.svelte.test.ts, flaiover/src/lib/server/license.ts, flaiover/src/lib/server/license.test.ts, flaiover/src/lib/server/no-project-files.test.ts, flaiover/src/routes/license, flaiover/src/routes/api/license, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/operators/index.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/tech/docker.md, README.md]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: high
usage:
  source: log
  seconds: 1812
  models:
    - model: claude-fable-5-1
      input: 1416
      output: 65564
      cache_read: 5871573
      cache_write: 180744
      cost: 8.3751
    - model: claude-haiku-4-5-20251001
      input: 564
      output: 19728
      cache_read: 2210625
      cache_write: 112171
      cost: 0.4605
    - model: claude-sonnet-5-5
      input: 38
      output: 8196
      cache_read: 633335
      cache_write: 91039
      cost: 0.4363
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
- [x] The license grants free use for the classes mentioned above
- [x] The license makes it clear that any commercial use falling outside of the permitted criteria requires a paid license
- [x] The license makes it clear that offering licenses to system-flow or a fork of it, or offering system-flow or a fork of it as a service is prohibited
- [x] The license requires that a copy of it be distributed in all forms
- [x] The license file is shipped with the docker image and is visible under Host in its own submenu and in the CLI binary and is printed when the command `flai license` is issued.

## Tasks
- T-0696 LICENSE.md adapts PolyForm Shield with free use for the named classes and bars offering it as a service or licensing it
- T-0697 flai license prints the license embedded in the binary
- T-0698 The docker image ships LICENSE.md and the dashboard shows it under Host in its own submenu
- T-0699 The design, user docs, and README describe the license and where it is shown

## Notes
- The license is the system-flow Shield License 1.0, an adaptation of PolyForm Shield 1.0.0 without the PolyForm name (the PolyForm Project asks that of adaptations). The free classes are exactly the four the goal lists; the hosted-service and licensing prohibition binds every licensee, paid or free. Licensor name, paid-license contact, and two scope questions are in TH-0067, answered with the defaults stated there until the operator says otherwise.
- The binary embeds a copy of `LICENSE.md` at `flai/internal/license/LICENSE.md`, kept equal to the root by a test, because Go embeds nothing above its module. The dashboard reads the image's own `/app/LICENSE.md` rather than asking flai: the license is the image's, not the project's.
