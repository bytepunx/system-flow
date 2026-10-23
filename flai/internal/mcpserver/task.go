package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/taskdone"
)

const taskDoneDescription = "Close a task in one call from its story's worktree (ADR-0107), with the answer flai task done --json prints. In order: commit with message the paths changed in the worktree that the task's touches cover and those no other open task of the story covers (ADR-0128), leaving uncommitted, and listed in left, those only another open task covers, for its close (nothing to commit is no failure, and needs no message), message each other story in progress or in review whose claim covers a path the commit changed, the shared paths included, naming the task, the commit, its subject, and those paths, and answer them as told (a message that cannot be sent is logged and stops nothing), sync the story's branch as flai stream sync does (a refusal for the paths left alone stops nothing: the close that leaves none syncs it), move the task to done with what moves with it, log message's subject line, or log, or with neither 'Closed T-nnnn: <title>', to the story's narrative, add the paths the commit changed that the story's branch still changes against the main branch, never a path left, to the task's and the story's touches where they do not cover them, listing the committed paths the branch no longer changes in touches.unchanged (ADR-0144), run flai check --strict scoped to the story, listing every finding in the story and the first 10 notes outside it, counting the rest in check.notes_omitted (ADR-0198; flai check --story lists them), and read your inbox, as inbox does but listing only the threads on the story and its tasks and counting the others in inbox.threads_omitted (call inbox to list them) (ADR-0179). The whole answer is at most 40,000 bytes of JSON: the inbox's oldest changes are left out first, counted in changes_omitted, then each listed thread's last entry and pending recommendation is shortened, counted in entries_shortened (read them whole with thread_get). It stops at the first step that fails, and answers stopped naming it, error saying why and what to do, and what the steps before it did; a stop is an answer, not a tool error. After a stopped sync, resolve each conflicting path in the worktree, git add it, run git rebase --continue there, and call again; after a failed check, fix its findings and call again, which commits the fix. A task already done is not moved again. When this server's flai is older than the newest flai release in the project's history and that release raised the task-done rules number or the major version (or, for a release without a number, the minor version), which refuses the close (ADR-0142, ADR-0180), the flai installed at the path the server started from closes the task instead once it has caught up by the same rule, answering its result with closed_by naming its version and path; while it is behind too, the close is refused, naming its version: call again once the host is upgraded (ADR-0160). Running the tests and ticking criteria stay yours."

// taskDoneInstructions sends the agent to task_done at every task
// transition (ADR-0107).
const taskDoneInstructions = "At every task transition, close the task with task_done, in place of committing, syncing, moving the task, logging, widening touches, checking, and calling inbox one by one: it does all seven in that order, stops at the first that fails with what to do, and is called again once that is done."

// TaskDoneIn closes a task (ADR-0107).
type TaskDoneIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Task    string `json:"task" jsonschema:"the task to close, T-nnnn, in any zero padding"`
	Message string `json:"message,omitempty" jsonschema:"the commit message, needed only when there is something to commit; its subject line is the narrative's log entry unless log is given"`
	Log     string `json:"log,omitempty" jsonschema:"the narrative's log entry, when it is not to be the message's subject line"`
}

func (in TaskDoneIn) project() string { return in.Project }

// TaskDoneOut is taskdone's result itself, so that task_done answers field
// for field what flai task done --json prints.
type TaskDoneOut = taskdone.Result

// taskDone closes the task as the server's agent. A step that fails is the
// answer's stopped and error; only a close that cannot start is the tool's
// error. A close this server's outdated flai refuses goes to the installed
// flai (ADR-0160).
func (s *server) taskDone(ctx context.Context, _ *mcp.CallToolRequest, in TaskDoneIn) (*mcp.CallToolResult, TaskDoneOut, error) {
	res, err := taskdone.Run(ctx, taskdone.Options{
		Repo: s.repo, Runner: s.runner, Task: in.Task, Message: in.Message, LogEntry: in.Log,
		Agent: s.agent, Session: s.session, Now: s.now, Version: s.version, Logger: s.logger,
	})
	var outdated *taskdone.OutdatedError
	if errors.As(err, &outdated) {
		return s.closeWithInstalled(ctx, in, outdated)
	}
	if err != nil {
		return nil, TaskDoneOut{}, err
	}
	return nil, res, nil
}

// closeWithInstalled closes the task that refusal refused with the flai
// installed at the path this server was started from, when that flai has
// caught up with the newest release by the refusal's rule, BehindRules
// (ADR-0180); otherwise it refuses, naming the installed flai's version and
// task-done rules number or why they could not be read (ADR-0160).
func (s *server) closeWithInstalled(ctx context.Context, in TaskDoneIn, refusal *taskdone.OutdatedError) (*mcp.CallToolResult, TaskDoneOut, error) {
	if s.installed == "" {
		return nil, TaskDoneOut{}, fmt.Errorf("%w; no installed flai is known to close it with (ADR-0160)", refusal)
	}
	f := installedFlai{path: s.installed, dir: projectRoot(s.repo), env: append(os.Environ(), "FLAI_AGENT="+s.agent, "FLAI_SESSION="+s.session)}
	version, rules, err := f.version(ctx)
	if err != nil {
		return nil, TaskDoneOut{}, fmt.Errorf("%w; the installed flai at %s could not be read (%w), so it cannot close it either; call task_done again once the host is upgraded (ADR-0160)", refusal, s.installed, err)
	}
	if taskdone.BehindRules(version, rules, refusal.Newest, refusal.NewestRules) != "" {
		return nil, TaskDoneOut{}, fmt.Errorf("%w; the installed flai at %s is %s, behind too; call task_done again once the host is upgraded (ADR-0160, ADR-0180)", refusal, s.installed, installedIs(version, rules))
	}
	res, err := f.taskDone(ctx, in)
	if err != nil {
		return nil, TaskDoneOut{}, err
	}
	res.ClosedBy = fmt.Sprintf("flai %s at %s", version, s.installed)
	// closed_by adds to what the installed flai fitted, and an installed flai
	// from before ADR-0179 fitted nothing
	res.Fit()
	if s.logger != nil {
		s.logger.Info("task closed with the installed flai", "component", "mcp", "agent", s.agent, "task", res.Task, "running", refusal.Running, "installed", version, "path", s.installed, "stopped", res.Stopped)
	}
	return nil, res, nil
}
