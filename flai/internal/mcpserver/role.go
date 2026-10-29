package mcpserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/harness"
)

// Starting a role the project defines on the host (S-0417, ADR-0207): an
// operator's own agent can have flai do what flai role run does, under the
// host action the role's file names. flai guard refuses the tool to a
// sub-agent, as it refuses every tool that is not a read (ADR-0060).

// RoleStart starts name, a role the project at root defines, on item, or on
// none when item is empty, with inputs in the order given, as flai role run
// does, on by's word. A refusal is an error that says why.
type RoleStart func(ctx context.Context, root, name, item string, inputs []harness.Input, by string) (RoleStarted, error)

// RoleStarted is the run of the role started.
type RoleStarted struct {
	Role    string `json:"role" jsonschema:"the role's name, as its file gives it"`
	Item    string `json:"item,omitempty" jsonschema:"the item the run is on, empty for none"`
	Action  string `json:"action" jsonschema:"the host action the role's file names, which its run is journalled under"`
	Agent   string `json:"agent"`
	Harness string `json:"harness,omitempty"`
	Command string `json:"command,omitempty"`
	PID     int    `json:"pid,omitempty"`
	Log     string `json:"log,omitempty" jsonschema:"the role's output, on the host"`
	Session string `json:"session,omitempty"`
	Started string `json:"started,omitempty"`
	Trigger string `json:"trigger,omitempty" jsonschema:"what started the run: asked, for the operator's asking"`
}

// RoleRunIn names the role to start, its item, and its inputs.
type RoleRunIn struct {
	Project string          `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Name    string          `json:"name" jsonschema:"the role's name, as its file under design/workflows/roles gives it, such as reviewer"`
	Item    string          `json:"item,omitempty" jsonschema:"the epic or story ID to run the role on, such as S-0042 (any zero padding), when its file's item takes one; none when it takes none"`
	Inputs  []harness.Input `json:"inputs,omitempty" jsonschema:"the run's inputs, each a name and a value, in the order its prompt names them"`
}

func (in RoleRunIn) project() string { return in.Project }

const roleRunDescription = "Start a role the project defines in its role file under design/workflows/roles on this host now, as flai role run does (S-0417, ADR-0207): in the project's main checkout, with the role's agent (roles.<name>.agent over agent in system-flow.yaml), as its file's harness definition, its body the prompt, on the item given when its file's item takes one, an epic or a story, and on none when it takes none, with inputs, each a name and a value, in the order its prompt names them. It does not wait for the role: the run is recorded where flai serve tracks agents, under roles by its name and item, and flai serve settles it once it ends, journals it under the role's action, and logs its activity in wip/agents/<name>.md. It is no story's agent and holds no story back. For the operator's own agent: refused to an agent flai serve started, as plan is. Needs the host action the role's file names on for the project (flai serve enable <action>). Refused, saying why, for a built-in role (the tools plan, analyze, and curate and the orchestrate host action start those) and a name no role file gives, while the role's host action is off, for an item when its file takes none, without one when it takes some, for an item of a kind its file does not take or one archived, done, or cancelled, while a run of the role on the same item, or on none, runs, and when its agent names no harness and no command is set on the host."

func (s *server) roleRun(ctx context.Context, _ *mcp.CallToolRequest, in RoleRunIn) (*mcp.CallToolResult, RoleStarted, error) {
	name, item := strings.TrimSpace(in.Name), strings.TrimSpace(in.Item)
	if name == "" {
		return nil, RoleStarted{}, fmt.Errorf("give name, the role's name as its file under design/workflows/roles gives it; flai role list lists them")
	}
	how := strings.TrimSpace("flai role run " + name + " " + item)
	for _, input := range in.Inputs {
		if strings.TrimSpace(input.Name) == "" {
			return nil, RoleStarted{}, fmt.Errorf("an input with the value %q has no name: give each input its name, which the role's prompt names, and its value", input.Value)
		}
		how += " --input " + input.Name + "=" + input.Value
	}
	if s.roles == nil {
		return nil, RoleStarted{}, fmt.Errorf("this flai mcp cannot start a role; on the host run %s", how)
	}
	out, err := s.roles(ctx, projectRoot(s.repo), name, item, in.Inputs, s.agent)
	if err != nil {
		return nil, RoleStarted{}, err
	}
	return nil, out, nil
}
