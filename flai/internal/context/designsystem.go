package context

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// KindTokens is the kind of the design token file, the one document a pack
// holds that is not markdown (ADR-0152).
const KindTokens = "tokens"

// DesignSystemReason starts the reason the design system is in a story's
// pack, followed by the guidance's topics that reached the story: "design
// system: topics flaiover".
const DesignSystemReason = "design system: topics "

// DesignSystem is a project's design system as prime gives it to the stories
// it reaches (ADR-0152): its guidance, one of the design documents, and its
// token file, a DTCG JSON file read as its raw text.
type DesignSystem struct {
	Guidance *Doc // nil when the project has none, and then it has no design system
	Tokens   *Doc // nil when the token file is not there
}

// loadDesignSystem finds the guidance among docs at the path the manifest
// names, and reads the token file at the path it names. Without the
// guidance there is no design system and the token file is not read.
func loadDesignSystem(root string, m manifest.Manifest, docs []*Doc) (*DesignSystem, error) {
	ds := &DesignSystem{}
	for _, d := range docs {
		if d.Path == m.DesignGuidance() {
			ds.Guidance = d
		}
	}
	if ds.Guidance == nil {
		return ds, nil
	}
	rel := m.DesignTokens()
	raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return ds, nil
	case err != nil:
		return nil, fmt.Errorf("reading the design token file %s: %w", rel, err)
	}
	ds.Tokens = tokensDoc(rel, string(raw))
	return ds, nil
}

// tokensDoc is the token file as a document: its raw text in one section
// with no heading, so it loads whole or is briefed, and is never ranked.
func tokensDoc(rel, raw string) *Doc {
	return &Doc{Path: rel, Kind: KindTokens, Title: "Design tokens", Raw: raw, Sections: []topics.Section{{Line: 1, Text: raw}}}
}

// Reaches is the guidance's front-matter topics that reach a story with
// these topics, as ByTopics matches them: all, or those the story has. It
// is empty when the project has no design system.
func (ds *DesignSystem) Reaches(story []string) []string {
	if ds == nil || ds.Guidance == nil {
		return nil
	}
	return matched(ds.Guidance.Topics, story)
}

// HasDesignSystem reports whether the project has a design system: a
// guidance document at the path the manifest names (ADR-0146). Without one
// no agent reads the design-system convention (ADR-0152).
func HasDesignSystem(repo *workitem.Repo) (bool, error) {
	rel := repo.Manifest.DesignGuidance()
	_, err := os.Stat(filepath.Join(repo.Root, filepath.FromSlash(rel)))
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, fs.ErrNotExist):
		return false, nil
	}
	return false, fmt.Errorf("looking for the design system's guidance %s: %w", rel, err)
}

// DesignSystem chooses the design system for a story with these topics when
// the guidance's topics reach them (ADR-0152): the guidance and the token
// file, whole, as if the story named them, with the reason
// DesignSystemReason and the topics that reached it. Each larger than
// BriefOver is briefed instead (ADR-0050), and with BriefNamed each is
// briefed, as a sub-agent's pack briefs what the story names. Run before
// ByTopics, so the guidance's topics list their reason on this choice
// rather than brief it again.
func (s *Selection) DesignSystem(ds *DesignSystem, story []string) {
	reach := ds.Reaches(story)
	if len(reach) == 0 {
		return
	}
	reason := DesignSystemReason + strings.Join(reach, ", ")
	for _, d := range []*Doc{ds.Guidance, ds.Tokens} {
		if d == nil {
			continue
		}
		if s.byPath[d.Path] == nil {
			s.add(d)
		}
		switch {
		case s.Whole(d):
			s.choose(d, d.all(), "", reason)
		case s.BriefNamed || (s.BriefOver > 0 && len(d.Raw) > s.BriefOver):
			s.addBrief(d, reason, nil)
		default:
			s.choose(d, d.all(), StepNamed, reason)
		}
	}
}

// tokensBrief is the text of the token file's brief, which has no outline.
const tokensBrief = "A DTCG token file in JSON, not markdown: it has no outline, and doc_get and flai doc show do not read it. Read the file itself, or every token with its value and CSS variable with flai design show --json, before building UI from it.\n"

// designSystemBrief follows the line naming a brief of the design system's
// guidance or token file in a story's pack (ADR-0152).
const designSystemBrief = "The design system reaches this story; it is briefed for its size. Read it whole before building or changing UI.\n"
