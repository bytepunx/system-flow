---
id: T-3547
type: task
nature: remediation
title: flai serve copies a story namespace's kubeconfig, record, and step logs into the story's agent cache folder, and a teardown removes them
status: done
parent: S-0834
owner: alex
created: 2026-11-17T08:32:18Z
updated: 2026-11-18T07:40:19Z
transitions:
  - to: ready
    at: 2026-11-18T07:11:51Z
    by: agent-S-0834
  - to: in-progress
    at: 2026-11-18T07:11:52Z
    by: agent-S-0834
  - to: done
    at: 2026-11-18T07:40:19Z
    by: agent-S-0834
stream: S-0834
tags: [kluster, go]
touches: [flai/internal/kluster/agentcopy.go, flai/internal/kluster/agentcopy_test.go, flai/internal/kluster/agentconfig.go, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go]
after: [T-3546]
usage:
  source: log
  seconds: 1707
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 60
      output: 20135
      cache_read: 5983212
      cache_write: 222617
      cost: 2.9697
---
# T-3547 flai serve copies a story namespace's kubeconfig, record, and step logs into the story's agent cache folder, and a teardown removes them

## Work

This task builds the copy that T-3546's ADR decides, so it waits for that ADR. If the operator changes the decision, follow the ADR as accepted rather than these words.

- In `flai/internal/kluster/agentcopy.go`, add the copy's folder, `AgentCopyDir(root, story)`, which is `<root>/.flai-cache/agents/<S-nnnn>/kluster`. Add a function that writes three things there, each whole with `atomicfile`:
  - the namespace's kubeconfig, mode 0600, with its certificate authority inlined as data, since the file `agentconfig.go` makes absolute is not mounted in a container;
  - the deployment record;
  - each step log of the namespace.
- The copy is the same whatever the container's network. Rewriting the server for `bridge` is T-3549's work, at read time, from the variable T-3548 gives the container.
- Call the copy in `deploy.go` at the end of each deploy and checks run into a story's namespace, whether it passed or failed. A teardown that passes removes the copy. A failed teardown keeps it, as it keeps the record.
- The default namespace gets no copy.
- A failure to copy is logged in the step's log. It fails nothing, because the main checkout's files stay the evidence.

T-3548 changes the harness at the same time; this task shares no file with it.

## Done when

- Tests in `agentcopy_test.go` and `deploy_test.go` show the copy after a deploy and after a checks run, with the certificate authority inlined, and gone after a teardown that passes.
- `flai test` passes on the changed paths.

## Notes
