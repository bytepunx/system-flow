package issues

import (
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/analysis"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Impact is what an issue costs while it stays open, as the analyzer measures
// it (S-0224): the cost of delay inputs its "## Impact" section gives, which
// impact reads back for the story made from it, and the evidence for them.
// An empty field is not given. An amount is a number of zero or more in
// planning.currency; a duration is a Go duration longer than zero.
type Impact struct {
	RevenuePerWeek, PenaltyPerWeek, TimeLostPerCycle string
	Evidence                                         string
}

// impactKeys are the Impact section's inputs, in the order they are written.
var impactKeys = []string{"revenue_per_week", "penalty_per_week", "time_lost_per_cycle"}

// IsZero reports whether nothing is given.
func (im Impact) IsZero() bool {
	return strings.TrimSpace(im.RevenuePerWeek+im.PenaltyPerWeek+im.TimeLostPerCycle+im.Evidence) == ""
}

// Validate refuses an amount that is not a number of zero or more and a
// duration that is not a Go duration longer than zero.
func (im Impact) Validate() error {
	_, err := im.entries()
	return err
}

// entries are the inputs given, as key and value in impactKeys order, each
// value written as impact reads it back: 1200, 4h, 1h30m.
func (im Impact) entries() ([][2]string, error) {
	var out [][2]string
	var errs []string
	for i, v := range []string{im.RevenuePerWeek, im.PenaltyPerWeek, im.TimeLostPerCycle} {
		key, v := impactKeys[i], strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if key == "time_lost_per_cycle" {
			d, err := time.ParseDuration(v)
			if err != nil || d <= 0 {
				errs = append(errs, fmt.Sprintf("impact %s %q is not a duration longer than zero, like 4h", key, v))
				continue
			}
			out = append(out, [2]string{key, normalise(d.String())})
			continue
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil || math.IsNaN(f) || math.IsInf(f, 0) || f < 0 {
			errs = append(errs, fmt.Sprintf("impact %s %q is not an amount of zero or more, like 1200", key, v))
			continue
		}
		out = append(out, [2]string{key, strconv.FormatFloat(f, 'f', -1, 64)})
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return out, nil
}

// setImpact writes the Impact section of body, which it adds before
// Remediation, or at the end, when body has none. Each input given replaces
// the lines that give that key; the lines of the others and the evidence stay.
// Evidence given is added as a paragraph, unless the section holds it already.
// The inputs follow the evidence as "- key: value" lines. im must be valid.
func setImpact(body string, im Impact) string {
	entries, _ := im.entries()
	given := map[string]bool{}
	var keyed [][2]string // key and line
	for _, e := range entries {
		given[e[0]] = true
		keyed = append(keyed, [2]string{e[0], "- " + e[0] + ": " + e[1]})
	}
	var evidence []string
	start, end, ok := section(body, "## Impact")
	if ok {
		for _, line := range strings.Split(body[start:end], "\n") {
			if m := impactLineRe.FindStringSubmatch(strings.TrimSpace(line)); m != nil {
				if !given[m[1]] {
					keyed = append(keyed, [2]string{m[1], line})
				}
				continue
			}
			evidence = append(evidence, line)
		}
	}
	text := strings.TrimSpace(strings.Join(evidence, "\n"))
	if ev := strings.TrimSpace(im.Evidence); ev != "" && !strings.Contains(text, ev) {
		if text != "" {
			text += "\n\n"
		}
		text += ev
	}
	var lines []string
	for _, key := range impactKeys {
		for _, k := range keyed {
			if k[0] == key {
				lines = append(lines, k[1])
			}
		}
	}
	if len(lines) > 0 {
		if text != "" {
			text += "\n\n"
		}
		text += strings.Join(lines, "\n")
	}
	text += "\n"
	if ok {
		if end < len(body) {
			text += "\n"
		}
		return body[:start] + text + body[end:]
	}
	if idx := strings.Index(body, "\n## Remediation"); idx >= 0 {
		return body[:idx+1] + "## Impact\n" + text + "\n" + body[idx+1:]
	}
	return strings.TrimRight(body, "\n") + "\n\n## Impact\n" + text
}

// reportPath checks that report names a markdown file under the analysis
// reports folder, design/analysis, other than its index, and returns its path
// relative to the project root, slash separated, or "" for none. A relative
// path is taken from the project root. The report need not be written yet.
func reportPath(r *workitem.Repo, report string) (string, error) {
	report = strings.TrimSpace(report)
	if report == "" {
		return "", nil
	}
	rel := report
	if filepath.IsAbs(rel) {
		var err error
		if rel, err = filepath.Rel(r.Root, rel); err != nil {
			rel = report
		}
	}
	rel = path.Clean(filepath.ToSlash(rel))
	dir := analysis.Dir(r.Manifest)
	if !strings.HasPrefix(rel, dir+"/") || path.Ext(rel) != ".md" || analysis.Skipped(path.Base(rel)) {
		return "", fmt.Errorf("report %q is not an analysis report: give the markdown file under %s/ that found it, like %s/2026-10-06-risk.md", report, dir, dir)
	}
	if fi, err := os.Stat(filepath.Join(r.Root, filepath.FromSlash(rel))); err == nil && fi.IsDir() {
		return "", fmt.Errorf("report %q is a folder, not an analysis report", report)
	}
	return rel, nil
}

// reportLine is the line naming an instance's report, with its newline, or ""
// for none.
func reportLine(report string) string {
	if report == "" {
		return ""
	}
	return "Report: " + report + ".\n"
}

// hasReport reports whether an instance of the issue names the report.
func hasReport(is *Issue, report string) bool {
	start, end, ok := section(is.Body, "## Instances")
	return ok && report != "" && contains(strings.Split(is.Body[start:end], "\n"), strings.TrimSuffix(reportLine(report), "\n"))
}

// reportLinkRe is the line linkReport writes in the Remediation section,
// capturing the link's target.
var reportLinkRe = regexp.MustCompile(`(?m)^Found by the analysis in \[[^\]]*\]\(([^)]+)\)\.$`)

// ReportLinks are the targets of the report links in the issue's Remediation
// section, in order, relative to the issue's folder: ../analysis/<file>.md.
func ReportLinks(is *Issue) []string {
	start, end, ok := section(is.Body, "## Remediation")
	if !ok {
		return nil
	}
	var out []string
	for _, m := range reportLinkRe.FindAllStringSubmatch(is.Body[start:end], -1) {
		out = append(out, m[1])
	}
	return out
}

// NamingReport is the IDs of the issues, open or closed, that name the
// analysis report, a path from the project root, as Reports gives them, in
// ID order; nil for an empty report. An analyzer run's end is charged to them
// (S-0227).
func NamingReport(r *workitem.Repo, report string) ([]string, error) {
	report = strings.TrimSpace(report)
	if report == "" {
		return nil, nil
	}
	report = path.Clean(filepath.ToSlash(report))
	list, err := List(r)
	if err != nil {
		return nil, fmt.Errorf("issues naming report %s: %w", report, err)
	}
	var ids []string
	for _, is := range list {
		if contains(Reports(is), report) {
			ids = append(ids, is.ID)
		}
	}
	return ids, nil
}

// linkReport adds a line linking the report, a path from the project root
// root, to the issue's Remediation section, as its last paragraph, adding the
// section when it has none. A report it links already is not linked again.
func linkReport(is *Issue, root, report string) {
	target, err := filepath.Rel(filepath.Dir(is.Path), filepath.Join(root, filepath.FromSlash(report)))
	if err != nil {
		target = report
	}
	target = filepath.ToSlash(target)
	if contains(ReportLinks(is), target) {
		return
	}
	line := fmt.Sprintf("Found by the analysis in [%s](%s).\n", path.Base(report), target)
	body := strings.TrimRight(is.Body, "\n") + "\n"
	if start, end, ok := section(body, "## Remediation"); ok {
		text, rest := strings.TrimSpace(body[start:end]), body[end:]
		if text != "" {
			text += "\n\n"
		}
		if rest != "" {
			rest = "\n" + rest
		}
		body = body[:start] + "\n" + text + line + rest
	} else {
		body += "\n## Remediation\n\n" + line
	}
	is.Body = body
}
