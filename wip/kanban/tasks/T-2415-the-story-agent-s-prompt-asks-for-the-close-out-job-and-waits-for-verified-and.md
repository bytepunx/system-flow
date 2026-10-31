---
id: T-2415
type: task
nature: improvement
title: The story agent's prompt asks for the close-out job and waits for verified, and scripts/close-out.sh and the template's copy say they run as the job, not in a turn
status: backlog
parent: S-0546
owner: alex
created: 2026-10-31T07:44:50Z
updated: 2026-10-31T07:45:23Z
transitions: []
stream: S-0546
tags: [flai]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, scripts/close-out.sh, template/root/scripts/close-out.sh]
after: [T-2407]
---
# T-2415 The story agent's prompt asks for the close-out job and waits for verified, and scripts/close-out.sh and the template's copy say they run as the job, not in a turn

## Work

Rewrite the close-out line of `delegation(r)` in `flai/internal/harness/harness.go`, which today says to run `scripts/close-out.sh` as one Bash call with a 600000 ms timeout. It should say to call the MCP `verify` tool (or `flai verify --job`) with the commit message, then hold `wait_for_events` until `verified` comes, then read the verdict, fix what failed, and ask again. It should also say never to run the script in a turn. Update `harness_test.go`, which quotes the old line.

In `scripts/close-out.sh` and `template/root/scripts/close-out.sh`, say in the header that the script is what the serve job runs, as T-2396's ADR decides, and that an agent asks for the job and does not call the script itself. If the ADR has the job run the steps in Go, the script becomes the call to `flai verify --job`. The two copies differ today (the template's calls `flai verify` directly), so keep each one's own call.

Waits for T-2407: the prompt names the tool and flag, which must exist. Shares no path with T-2417, so the two run together.

## Done when

- [ ] `harness_test.go` asserts the new close-out line and that the 600000 ms Bash instruction is gone
- [ ] Both scripts pass ShellCheck, and the template's passes `scripts/template-test.sh`
- [ ] `go test -race ./internal/harness` passes
- [ ] Story criterion 3 (the prompt and the script)

## Notes
