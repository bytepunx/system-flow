package workitem

import "github.com/bytepunx/system-flow/flai/internal/manifest"

// AgentOf is the agent flai starts for the story it (ADR-0251): its own when
// it names a harness, otherwise the project's default with what its own does
// say laid over it; nil when neither sets anything.
//
// S-0573: a story made before the project had a default carries no agent, and
// setting the default later is enough for it to start. A story that names a
// harness, "none" included (ADR-0242), keeps its agent whole, so a later
// change to the default does not reach it (ADR-0037). Nothing is written to
// the story.
func (r *Repo) AgentOf(it *Item) *manifest.Agent {
	if it == nil {
		return nil
	}
	if it.Agent != nil && it.Agent.Harness != "" {
		return it.Agent
	}
	return r.Manifest.Agent.With(it.Agent)
}
