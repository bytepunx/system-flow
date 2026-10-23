---
id: D-0262
title: flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version
kind: cause
sources: [I-0085]
topics: [flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src/lib/server/notify.test.ts, design/issues/I-0085-flaiover-s-notify-test-ts-fails-now-and-then-under-the-full-vitest-run-because-project-info-reads-a-system-flow-yaml-with-no-version.md, design/issues/summary.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0262 flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version

I-0085 was closed on 2026-10-08.

The cause:

```text
flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version
```

The fix:

```text
S-0288 fixed it in the test's setup, not in Repo. notify.test.ts's setManifest truncated system-flow.yaml with writeFile and never reported the change, so a manifest read that an earlier repo.changed had started could see the file empty ("version is required") or cache it without notify_url (startNotifier returned null). setManifest now writes a temporary file, renames it over the manifest, and calls repo.changed('system-flow.yaml') as flai on the host does. A new test, "reads the notify_url set after a change was reported and the manifest read", reproduces the stale cached read and failed against the old setManifest; the whole flaiover vitest suite then passed five runs in a row.
```

Source: I-0085.
