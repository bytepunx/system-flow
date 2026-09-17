---
id: T-090
type: task
nature: feature
title: "ADR-0016: dashboard writes and metrics delegate to the flai binary; server wrapper with JSON and error mapping"
status: done
parent: S-013
owner: alex
created: 2026-09-17T04:32:26Z
updated: 2026-09-17T04:38:32Z
transitions:
  - to: ready
    at: 2026-09-17T04:38:31Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:38:31Z
    by: agent
  - to: done
    at: 2026-09-17T04:38:32Z
    by: agent
stream: S-013
tags: [dashboard]
---

# T-090 ADR-0016: dashboard writes and metrics delegate to the flai binary; server wrapper with JSON and error mapping

## Work
Decide and record (ADR-0016) that flaiover performs every write (move, block, unblock, stream log) and every metric by invoking the flai binary with --json, rather than porting the rules to TypeScript; src/lib/server/flai.ts locates the binary (FLAI_BIN, then PATH), runs it in PROJECT_DIR with FLAI_CONFIG and FLAI_CACHE_DIR inside the project's .flai-cache and FLAI_AGENT=flaiover, parses stdout JSON, and maps a fatal stderr event to a RepoError whose message is the rule text.

## Done when
Wrapper tested against the built flai on a temp project; a missing binary yields a clear 503.

## Notes
