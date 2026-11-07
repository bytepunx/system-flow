package check

import (
	"os"
	"path/filepath"
)

// instructionFiles are the files at a project's root that an agent's harness
// reads first as the project's map: OpenCode reads AGENTS.md and falls back
// to CLAUDE.md, and Claude Code reads CLAUDE.md (S-0406).
var instructionFiles = []string{"AGENTS.md", "CLAUDE.md"}

// instructions warns, as instructions.missing, when the project's root holds
// none of instructionFiles, so no harness gives its agents the map. The
// warning is on CLAUDE.md, the file a project made from the template has, so
// a check scoped to a story is outside it unless the story changes that file
// (ScopeToStories).
func (c *checker) instructions() {
	for _, name := range instructionFiles {
		if st, err := os.Stat(filepath.Join(c.repo.Root, name)); err == nil && !st.IsDir() {
			return
		}
	}
	c.add(Warning, "instructions.missing", filepath.Join(c.repo.Root, "CLAUDE.md"), 1,
		"the project root has neither AGENTS.md nor CLAUDE.md, the map a harness reads first: add one there (Claude Code reads CLAUDE.md; OpenCode reads AGENTS.md, else CLAUDE.md)")
}
