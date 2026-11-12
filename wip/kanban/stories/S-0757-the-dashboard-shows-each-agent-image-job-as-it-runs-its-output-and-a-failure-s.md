---
id: S-0757
type: story
nature: feature
title: The dashboard shows each agent image job as it runs, its output, and a failure's cause, and offers the operator a retry
status: backlog
parent: E-0036
owner: alex
created: 2026-11-12T16:03:25Z
updated: 2026-11-12T16:04:18Z
transitions: []
tags: [cli, dashboard]
topics: [docker, dashboard]
touches: [flai/internal/hostapi/images.go, flai/internal/hostapi/images_test.go, flai/internal/serve/imagejob.go, flai/internal/serve/imagejob_test.go, flaiover/src/routes/api/images/+server.ts, flaiover/src/routes/api/images/retry/+server.ts, flaiover/src/routes/ops/images/+page.svelte, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/ImageJob.svelte, flaiover/src/lib/components/ImageJob.svelte.test.ts, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, design/components/flaiover/src/lib/components/ImageJob.examples.json, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0751]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0757 The dashboard shows each agent image job as it runs, its output, and a failure's cause, and offers the operator a retry

## Goal

When an agent start fails today, the dashboard shows the refusal's text on the story's agent (`StoryAgent.svelte`), and nothing more. After S-0751, a missing image is brought back by a host job that may take minutes. The operator should see it while it runs, read its output when it fails, and start it again without a shell.

After this story the dashboard shows, under Ops, the project's agent images and their state from `flai image status`, and each image job: running or ended, its way, its step, and its output as it is written. A failed job shows its cause, the failing step, and its output, with a retry. A run held on the job shows that it waits for its image, linked to the job, in place of a failure.

## Acceptance criteria

- [ ] A host API read gives the project's agent images with their state, and the image jobs, running and the last ended, each with its id, its images, its way, its current step, its start and end, and its outcome. A log read gives a job's output from an offset, as `kluster.log` does for a deploy's step
- [ ] A host API write starts a new image job for the project, behind the `agent` host action. It is refused while a job for the same images runs, and it says why
- [ ] An Ops page lists the images and the jobs, follows a running job's output live, and shows a failed job's cause, step, and output with a Retry button that calls the write
- [ ] A story's agent held on an image job shows "waiting for its image" with a link to the job, and a run failed by a job shows the cause and a link, in place of the bare refusal
- [ ] flaiover's tests cover the page, the component, and the held run, and the component's examples file passes `components.examples`
- [ ] `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` describe the page, the retry, and the host action that gates it

## Tasks
- T-3181 The host API reads the project's agent images and image jobs, reads a job's log from an offset, and starts a new job behind the agent host action
- T-3182 flaiover's API routes read the agent images, the jobs, and a job's log, and post a retry
- T-3183 An ImageJob component shows one image job's state, step, live output, and a failure's cause with a Retry button
- T-3184 A story's agent held on an image job shows that it waits for its image, and one failed by a job shows the cause, each linked to the job
- T-3185 An Ops › Agent images page lists the project's images and jobs, follows a running job, and retries a failed one
- T-3186 The dashboard's design and user guide describe the Agent images page, a held run, and the retry behind the agent host action

## Notes

Proposed slice 4b of E-0036. The log read follows `kluster.log` (ADR-0263), and the page's live log reuses `LogFeed.svelte`.
