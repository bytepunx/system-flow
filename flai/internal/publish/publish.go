// Package publish pushes a locally developed template to its git remote:
// clone the branch, replace the contents, commit, tag, push. Git failures
// are returned verbatim; nothing is retried or forced unless asked.
package publish

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/template"
)

// Options control a push.
type Options struct {
	Dir      string // local template directory (has template.yaml)
	Remote   string // git URL
	Ref      string // branch
	CacheDir string // working clones live under CacheDir/publish
	Tag      bool   // also create and push v<version>
	Force    bool   // force push
	DryRun   bool   // report and leave the remote untouched
	Message  string // commit subject override
	Source   string // where the template came from, for the commit body
}

// Result reports what happened.
type Result struct {
	Remote  string   `json:"remote"`
	Ref     string   `json:"ref"`
	Version string   `json:"version"`
	Files   []string `json:"files"`   // git status --porcelain lines
	Commit  string   `json:"commit"`  // short hash, empty when nothing to push or dry run
	Message string   `json:"message"` // commit message used or proposed
	Tag     string   `json:"tag,omitempty"`
	Pushed  bool     `json:"pushed"`
	Nothing bool     `json:"nothing_to_push"`
	Created bool     `json:"branch_created"` // branch did not exist on the remote
}

// Run performs the publish.
func Run(r execx.Runner, opt Options) (*Result, error) {
	m, err := template.LoadManifest(opt.Dir)
	if err != nil {
		return nil, err
	}
	if opt.Remote == "" {
		return nil, fmt.Errorf("no remote: pass --remote or set publish.repo in %s", template.ManifestFile)
	}
	if opt.Ref == "" {
		return nil, fmt.Errorf("no branch: pass --ref or set publish.ref in %s", template.ManifestFile)
	}
	if err := execx.Require(r, "git", "Publishing needs git."); err != nil {
		return nil, err
	}
	res := &Result{Remote: opt.Remote, Ref: opt.Ref, Version: m.Version}
	sum := sha256.Sum256([]byte(opt.Remote))
	work := filepath.Join(opt.CacheDir, "publish", hex.EncodeToString(sum[:])[:12])
	_ = os.RemoveAll(work)
	defer func() {
		if opt.DryRun {
			_ = os.RemoveAll(work)
		}
	}()
	if err := os.MkdirAll(filepath.Dir(work), 0o755); err != nil {
		return nil, err
	}
	if _, err := r.Run("", "git", "clone", "--quiet", opt.Remote, work); err != nil {
		return nil, err
	}
	if _, err := r.Run(work, "git", "checkout", "--quiet", opt.Ref); err != nil {
		// branch does not exist on the remote: create it from the default branch
		if _, err2 := r.Run(work, "git", "checkout", "--quiet", "-b", opt.Ref); err2 != nil {
			return nil, err2
		}
		res.Created = true
	}
	if err := replaceContents(work, opt.Dir); err != nil {
		return nil, err
	}
	if _, err := r.Run(work, "git", "add", "-A"); err != nil {
		return nil, err
	}
	status, err := r.Run(work, "git", "status", "--porcelain")
	if err != nil {
		return nil, err
	}
	for _, l := range strings.Split(strings.TrimSpace(status), "\n") {
		if l != "" {
			res.Files = append(res.Files, l)
		}
	}
	res.Message = opt.Message
	if res.Message == "" {
		res.Message = fmt.Sprintf("template %s", m.Version)
	}
	if opt.Source != "" {
		res.Message += "\n\nSynced from " + opt.Source + " by flai template push."
	}
	if opt.Tag {
		res.Tag = "v" + m.Version
		out, err := r.Run(work, "git", "ls-remote", "--tags", "origin", res.Tag)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(out) != "" && !opt.Force {
			return nil, fmt.Errorf("tag %s already exists on %s; bump the version in %s or pass --force", res.Tag, opt.Remote, template.ManifestFile)
		}
	}
	if len(res.Files) == 0 && !(opt.Tag && res.Tag != "") {
		res.Nothing = true
		return res, nil
	}
	if opt.DryRun {
		return res, nil
	}
	if len(res.Files) > 0 {
		if _, err := r.Run(work, "git", "-c", "user.name=flai", "-c", "user.email=flai@localhost", "commit", "--quiet", "-m", res.Message); err != nil {
			return nil, err
		}
	}
	short, _ := r.Run(work, "git", "rev-parse", "--short", "HEAD")
	res.Commit = strings.TrimSpace(short)
	if opt.Tag {
		args := []string{"tag", "-a", res.Tag, "-m", "template " + m.Version}
		if opt.Force {
			args = append(args, "-f")
		}
		if _, err := r.Run(work, "git", append([]string{"-c", "user.name=flai", "-c", "user.email=flai@localhost"}, args...)...); err != nil {
			return nil, err
		}
	}
	push := []string{"push", "--quiet"}
	if opt.Force {
		push = append(push, "--force")
	}
	push = append(push, "origin", opt.Ref)
	if _, err := r.Run(work, "git", push...); err != nil {
		return nil, err
	}
	if opt.Tag {
		tp := []string{"push", "--quiet"}
		if opt.Force {
			tp = append(tp, "--force")
		}
		tp = append(tp, "origin", res.Tag)
		if _, err := r.Run(work, "git", tp...); err != nil {
			return nil, err
		}
	}
	res.Pushed = true
	return res, nil
}

// replaceContents makes work mirror src except for .git.
func replaceContents(work, src string) error {
	entries, err := os.ReadDir(work)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Name() == ".git" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(work, e.Name())); err != nil {
			return err
		}
	}
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		if rel == "." {
			return nil
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return os.MkdirAll(filepath.Join(work, rel), 0o755)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(work, rel), data, info.Mode().Perm())
	})
}
