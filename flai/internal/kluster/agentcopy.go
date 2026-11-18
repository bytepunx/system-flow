package kluster

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
)

// AgentCopyDir is the copy of story's namespace in the story's agent cache
// folder, .flai-cache/agents/S-nnnn/kluster under root, which the story's
// container mounts: a Folder holding that namespace's record, step logs,
// state folder, and kubeconfig alone (ADR-0319 §1).
func AgentCopyDir(root, story string) string {
	return filepath.Join(root, ".flai-cache", "agents", strings.ToUpper(story), "kluster")
}

// copyForAgent writes the copy of r's namespace for its story's agent, as
// writeAgentCopy does, when r is a story's; the default namespace has none.
// A failure fails nothing, since the main checkout's files stay the evidence
// and the next write of the record writes the copy again: it is noted in the
// step log at log, when one is given and is there.
func copyForAgent(root string, r Record, log string) {
	if r.Story == "" {
		return
	}
	err := writeAgentCopy(root, r.Story)
	if err == nil || log == "" {
		return
	}
	f, ferr := os.OpenFile(log, os.O_WRONLY|os.O_APPEND, 0)
	if ferr != nil {
		return
	}
	defer func() { _ = f.Close() }()
	_, _ = fmt.Fprintf(f, "flai: the copy of the namespace %s for %s's agent, in %s, is not up to date: %v\n", r.Namespace, r.Story, AgentCopyDir(root, r.Story), err)
}

// writeAgentCopy writes into AgentCopyDir the copy of story's namespace as
// the project run from root keeps it: its record, its step logs, its state
// folder, and its kubeconfig, with the certificate authority it names
// inlined. Files are written whole and readable by their owner alone, in
// folders of mode 0700; what the main checkout no longer holds is removed
// from the copy. Each part is copied whatever another's failure, and the
// error joins theirs. Nothing of another namespace, nor the cluster's
// kubeconfig, is copied.
func writeAgentCopy(root, story string) error {
	ns, err := StoryNamespace(story)
	if err != nil {
		return err
	}
	src, dst := Folder(Dir(root)), Folder(AgentCopyDir(root, story))
	return errors.Join(
		wrapCopy("record", ns, mirrorFile(src.Record(ns), dst.Record(ns), nil)),
		wrapCopy("step logs", ns, mirrorDir(src.LogDir(ns), dst.LogDir(ns))),
		wrapCopy("state folder", ns, mirrorDir(src.State(ns), dst.State(ns))),
		wrapCopy("kubeconfig", ns, mirrorFile(src.AgentKubeconfig(ns), dst.AgentKubeconfig(ns), inlineAuthority)),
	)
}

// wrapCopy says which part of namespace ns's copy err failed to write.
func wrapCopy(part, ns string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("copy the namespace %s's %s for its story's agent: %w; the namespace's next deploy, checks run, or teardown writes it again", ns, part, err)
}

// mirrorFile writes the file src at dst, through transform when it is not
// nil, or removes dst when src is not there.
func mirrorFile(src, dst string, transform func(data []byte, path string) ([]byte, error)) error {
	data, err := os.ReadFile(src)
	if errors.Is(err, fs.ErrNotExist) {
		if err := os.Remove(dst); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("remove %s, which %s no longer is: %w", dst, src, err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", src, err)
	}
	if transform != nil {
		if data, err = transform(data, src); err != nil {
			return err
		}
	}
	return writeCopy(dst, data)
}

// writeCopy writes data at path whole, readable by its owner alone, in
// folders of mode 0700.
func writeCopy(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("make the folder %s: %w", filepath.Dir(path), err)
	}
	if err := atomicfile.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// mirrorDir makes the folder dst hold what the folder src holds: each of its
// folders and regular files, and nothing else, so that what src no longer
// holds leaves dst too. A src that is not there removes dst. Links and other
// special files in src are not copied.
func mirrorDir(src, dst string) error {
	if _, err := os.Stat(src); errors.Is(err, fs.ErrNotExist) {
		if err := os.RemoveAll(dst); err != nil {
			return fmt.Errorf("remove %s, which %s no longer is: %w", dst, src, err)
		}
		return nil
	}
	kinds := map[string]bool{}
	if err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Type().IsRegular() {
			rel, err := filepath.Rel(src, path)
			if err != nil {
				return err
			}
			kinds[rel] = d.IsDir()
		}
		return nil
	}); err != nil {
		return fmt.Errorf("list %s: %w", src, err)
	}
	if err := prune(dst, kinds); err != nil {
		return err
	}
	for rel, dir := range kinds {
		to := filepath.Join(dst, rel)
		if dir {
			if err := os.MkdirAll(to, 0o700); err != nil {
				return fmt.Errorf("make the folder %s: %w", to, err)
			}
			continue
		}
		if err := mirrorChanged(filepath.Join(src, rel), to); err != nil {
			return err
		}
	}
	return nil
}

// mirrorChanged writes the file src at dst with src's modification time,
// unless dst already has src's size and modification time, so that a record
// write copies again only the logs and state that changed since the last.
func mirrorChanged(src, dst string) error {
	from, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("read %s: %w", src, err)
	}
	if to, err := os.Stat(dst); err == nil && to.Size() == from.Size() && to.ModTime().Equal(from.ModTime()) {
		return nil
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read %s: %w", src, err)
	}
	if err := writeCopy(dst, data); err != nil {
		return err
	}
	if err := os.Chtimes(dst, from.ModTime(), from.ModTime()); err != nil {
		return fmt.Errorf("set the modification time of %s to its source's: %w", dst, err)
	}
	return nil
}

// prune removes from dst each entry that kinds, the folders (true) and
// files (false) the source holds by their path relative to it, does not
// hold as the same kind.
func prune(dst string, kinds map[string]bool) error {
	err := filepath.WalkDir(dst, func(path string, d fs.DirEntry, err error) error {
		if errors.Is(err, fs.ErrNotExist) && path == dst {
			return fs.SkipAll
		}
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dst, path)
		if err != nil {
			return err
		}
		if dir, ok := kinds[rel]; ok && dir == d.IsDir() && (dir || d.Type().IsRegular()) {
			return nil
		}
		if err := os.RemoveAll(path); err != nil {
			return err
		}
		if d.IsDir() {
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("remove from %s what its source no longer holds: %w", dst, err)
	}
	return nil
}
