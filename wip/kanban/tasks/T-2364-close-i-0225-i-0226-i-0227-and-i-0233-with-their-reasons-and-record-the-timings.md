---
id: T-2364
type: task
nature: remediation
title: Close I-0225, I-0226, I-0227, and I-0233 with their reasons, and record the timings in the story's notes
status: backlog
parent: S-0548
owner: alex
created: 2026-10-31T07:25:07Z
updated: 2026-10-31T07:27:04Z
transitions: []
stream: S-0548
tags: []
touches: [design/issues/summary.md, design/issues/I-0225-flai-internal-mcpserver-s-tests-take-8-to-10-minutes-under-race-at-go-s-10-minute-default-timeout-so-the-integration-tier-fails-a-merge-queue-batch-on-time-alone.md, design/issues/I-0226-two-flai-serve-tests-check-before-what-they-wait-for-has-happened-and-fail-a-loaded-merge-queue-batch-on-timing.md, design/issues/I-0227-flai-test-s-go-test-tier-runs-the-mcpserver-package-with-race-into-go-test-s-10-minute-default-timeout-on-a-loaded-host-and-names-only-the-package-so-a-story-that-changes-mcpserver-cannot-pass-its-tests.md, design/issues/I-0233-testachangedfilereachesthedashboardasanotification-fails-its-temp-folder-s-cleanup-with-directory-not-empty-as-i-0048-s-serve-tests-did.md]
after: [T-2361, T-2362, T-2365, T-2366, T-2367, T-2369]
---
# T-2364 Close I-0225, I-0226, I-0227, and I-0233 with their reasons, and record the timings in the story's notes

## Work

Once the mcpserver package's time and the serve tests' runs are measured, close each issue with `flai issue close --commit` and the reason, naming the task that fixed it and the measurement. Write the package's time and its slowest tests in S-0548's notes, and say that criterion 4's fortnight of `flai stats` lies after landing. Waits for T-2361 and T-2362, whose measurements it records.

## Done when

- [ ] The four issues are closed with their reasons
- [ ] S-0548's notes hold the timings

## Notes
