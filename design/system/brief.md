---
title: Original brief
updated: 2026-09-15
status: active
---

# Original brief

The brief that started system-flow, moved verbatim from the root CLAUDE.md on 2026-09-15. It is the source for the epics in `wip/kanban/epics`. Where the living design and this brief disagree, the living design wins and the difference is recorded in an ADR.

I want your assistance designing an agentic lean project management system. I would like for system design to use a hierarchy system that represents how work is specified through to delivery. I would like all written documentation to be expressed in markdown. All internal documentation in a mono-repo belongs under "design" with subfolders based on documentation type and all outward facing (user/operator/etc) belongs in a "docs" folder with subfolders per audience.

ADRs go in a "design/adrs" folder and capture a point-in-time architecture decision. They can be superseded in whole or in part with new ADRs. I would like active system design to go in "design/system" and to be captured in up-to-date documentation that receives edits anytime a conversation resolves to new details, new direction, new standards, or decisions. A "design/tech" folder should capture all active technology choices, their versions, and in which projects it is used in and why.

A top level "wip" folder should be where all work-in-process is captured in documents. The "wip/agents" subfolder should capture the agentic narrative across all active work streams such that, if an environment were to crash, there is high recoverability of the work instead of pushing the effort of trying to piece lost context together onto the human operator. The "wip/kanban" folder should carry the representation of all wip for human consumption. This is where the Lean system design starts to become apparent. A hierarchy of work type should help capture how the work is progressing and how it moves through the system in real time. All hierarchical work types should describe the nature of the deliverable to the system such as feature, improvment, remediation, research, or experiment. Epics represent deliverables that span multiple Stories. Stories represent an incremental deliverable. Tasks capture the different pieces of work required to deliver a Story. Each of these items should be captured with sufficient front-matter so that we can build a kanban visualization, cycle time chart, burn-up chart, and statistics that help make the system of work measurable and efficacy clear. Front matter need to include either timestamps and/or durations so we can determine what parts of the process are using the most time and so that we can attempt to optimize the process over time.

I would like the following outputs from the initial delivery of this project:

1. a polished go CLI we'll call "flai" that I can use to create new projects or import existing projects such that they conform to these standards and that automates toil involved in managing projects that conform to the system-flow standard
2. a master template for mono-repo projects (for now) that is used for new projects and existing projects that need to conform. the template should be managed in a separate git repository. the CLI should have a default repo it pulls this from but the CLI should read and store its config under ~/.flai/config.json and the template repo can be changed to any fork or fork and branch
3. the CLI should use the template to initialize a new project that includes a CLAUDE.md document baseline from the template
4. a web dashboard the uses the docuementation and wip folders to create an interactive visualization of a given project's progress. the dashboard should include many different ways of viewing specific data. All markdown documentation should be rendered from a documentation explorer. We should include search across both the design and wip folders. A kanban view and chart views of the system has to be included
5. The web dashboard should be a sveltekit SPA that uses tailwind for styling. we'll call it "flaiover" and it should have a built and published docker image that the CLI can download and run locally such that it's reachable from a port with a writeable mount of the target monorepo
5. the CLI should be able to install and run the web dashboard pointed to the current project given that it adheres to the conventions
6. the CLI should be capable of analyzing any monorepo and making the directory structure we've described in this standard and then asking about creating template files and giving options to replace default names. it should also give the user the opportunity to copy existing markdown into the new structure by prompting the user as necessary
7. User-facing documentation for the system-flow github repo that explains the conventions, the CLI tool, the template, and the dashboard project

We will be using this project as a way to build itself. A prototype of the template repo should be under "./template". Both flai and and flaiover should be subfolders in this monorepo that follow the template. We'll need to move through the decisions and design for the template (which will include some dev ex/dev ops components). Once that is in place, we can begin work on the cli and dashboard.
