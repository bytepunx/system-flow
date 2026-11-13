package verify

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// checksTimedOut is the outcome of a checks run that ran out of time, as the
// namespace's record keeps it beside passed, failed, and cancelled.
const checksTimedOut = "timed-out"

// klusterStep is the kluster step, which fills rep's Kluster part, and true,
// while the main checkout's manifest has a kluster: block; false, with no
// step, for a project without one, which has no namespace to speak of. A
// main checkout's manifest that cannot be read gives a step that passes
// saying so.
func (v *storyRun) klusterStep(rep *Report) (storyStep, bool) {
	m, err := v.mainManifest()
	if err == nil && m.Kluster == nil {
		return storyStep{}, false
	}
	return storyStep{StepKluster, func() []Finding {
		if err != nil {
			v.note = fmt.Sprintf("no kluster evidence is read for %s: the main checkout's manifest cannot be read: %v; fix it, then verify again", v.story, err)
			return nil
		}
		part, note, found := klusterEvidence(v.opts.Git, mainRoot(v.opts.Project), m, v.story, v.opts.KlusterOn, v.opts.ChecksOn)
		rep.Kluster, v.note = part, note
		return found
	}}, true
}

// mainManifest is the main checkout's manifest, whose kluster: block flai
// serve deploys the story's namespace by and whose checkout holds the
// namespace's record.
func (v *storyRun) mainManifest() (manifest.Manifest, error) {
	root := mainRoot(v.opts.Project)
	if root == v.opts.Project.Root {
		return v.opts.Project.Manifest, nil
	}
	return manifest.Load(filepath.Join(root, manifest.File))
}

// klusterEvidence is what the kluster step of story's close-out finds in its
// namespace's record under root, the main checkout, judged against the head
// of its branch there with git as kluster.ReadEvidence judges it (ADR-0292
// §6): the evidence for the report, nil while none applies; the step's note
// when it passes; and why it fails. It runs nothing. It passes, with no
// evidence, while m's kluster: block names no deploy, kluster.agent_namespaces
// is 0, the record's folder is not there to read, as in a story's container
// (ADR-0274), or the kluster host action is off (klusterOn); and while a
// record cannot be read, saying why. Otherwise it reports the deploy and
// checks run the record holds, and fails only when the deploy failed at the
// head, or, while the checks host action is on (checksOn), the checks run
// failed or timed out there: not deploying, or a record at an older commit,
// is no failure.
func klusterEvidence(git execx.Runner, root string, m manifest.Manifest, story string, klusterOn, checksOn bool) (part *kluster.Evidence, note string, found []Finding) {
	none := "no kluster evidence applies to " + story + ": "
	switch {
	case m.Kluster != nil && len(m.Kluster.Deploy) > 0 && m.KlusterAgentNamespaces() == 0:
		return nil, none + "kluster.agent_namespaces is 0, so no story's agent deploys before review; the review deploys it", nil
	case m.Kluster != nil && len(m.Kluster.Deploy) > 0 && !isDir(kluster.Dir(root)):
		return nil, fmt.Sprintf("no kluster evidence is read for %s here: the main checkout's %s is absent, as in a story's container, which does not mount it (ADR-0274), or because no kluster job has run in the project; item_get and the review page read the namespace's record on the host", story, filepath.Join(".flai-cache", "kluster")), nil
	}
	e, err := kluster.ReadEvidence(git, root, m, story, klusterOn, checksOn)
	if err != nil {
		return nil, fmt.Sprintf("no kluster evidence is read for %s: %v", story, err), nil
	}
	if !e.Applies {
		return nil, e.Summary(), nil
	}
	if why := failedAtHead(e); why != "" {
		return &e, "", []Finding{{Name: StepKluster, Message: why}}
	}
	note = e.Summary()
	if next := e.Next(); next != "" {
		note += "; " + next
	}
	return &e, note, nil
}

// failedAtHead is why e fails the close-out: its deploy failed at the head,
// or, while the checks must pass, its checks run failed or timed out there,
// with what to read and run next; "" when neither did.
func failedAtHead(e kluster.Evidence) string {
	if d := e.Deploy; d != nil && d.Head.AtHead && d.Job == "deploy" && d.Outcome == kluster.OutcomeFailed {
		why := e.Summary()
		if d.Error != "" {
			why += " (" + d.Error + ")"
		}
		if d.FailedStep != "" {
			why += fmt.Sprintf("; read its log with flai kluster log %s --step %s", e.Story, d.FailedStep)
		}
		return why + fmt.Sprintf("; fix the change, commit it, and deploy it again with %s, or the MCP tool %s", e.Commands.Deploy, e.Commands.DeployTool)
	}
	if c := e.Checks; e.ChecksOn && c != nil && c.Head.AtHead && (c.Outcome == kluster.OutcomeFailed || c.Outcome == checksTimedOut) {
		why := e.Summary()
		if c.Log != "" {
			why += "; read its log " + c.Log
		}
		return why + fmt.Sprintf("; fix the change, commit it, and deploy it again with %s, or the MCP tool %s, which runs the checks once the deploy passes", e.Commands.Deploy, e.Commands.DeployTool)
	}
	return ""
}

// isDir reports whether path is a directory that can be read.
func isDir(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}
