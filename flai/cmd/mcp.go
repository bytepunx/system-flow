package cmd

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
)

func newMCPCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "Serve this repository to agents over the Model Context Protocol on stdio",
		Long: `An MCP server for agent sessions (ADR-0020): the inbox of threads awaiting
the agent, replies, work item and document reads, transitions with the
workflow rules, who is touching a path, and wait_for_events, which blocks
until the designer or another agent changes a thread, an item, or a
narrative. Every write goes through the same code as the CLI. The agent is
FLAI_AGENT. Register it in .mcp.json:

  { "mcpServers": { "flai": { "command": "flai", "args": ["mcp"] } } }

Standard output is the protocol channel; log events go to standard error.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			agent, _ := agentIdentity()
			srv := mcpserver.New(mcpserver.Options{Repo: repo, Agent: agent, Version: buildinfo.Version, Now: a.now})
			a.logger().Info("mcp server started", "component", "mcp", "agent", agent, "root", repo.Root)
			return srv.Run(cmd.Context(), &mcp.StdioTransport{})
		},
	}
}
