package forge

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// story is a story's pull request input with two criteria, one ticked, and
// a passed verify result.
func story(template string) StoryPullRequest {
	return StoryPullRequest{
		ID:    "S-0393",
		Title: " Push a story's branch on its move to review ",
		Criteria: []Criterion{
			{Text: "The branch is pushed.", Ticked: true},
			{Text: "The request is opened.", Ticked: false},
		},
		Verify: &VerifyResult{
			Passed: true,
			Commit: "0123456789abcdef0123456789abcdef01234567",
			RanAt:  time.Date(2026, 10, 17, 8, 0, 0, 0, time.FixedZone("CEST", 2*3600)),
		},
		Template: template,
	}
}

const projectTemplate = `Thanks for the change.

## Story

<!-- flai:story -->
S-
<!-- /flai:story -->

## Criteria

<!--flai:criteria-->
<!-- /flai:criteria -->

## Our definition of done

<!-- flai:definition-of-done -->
- [ ] Reviewed by a human
<!-- /flai:definition-of-done -->

## Verified

<!-- flai:verify --><!-- /flai:verify -->

Signed off by the team.
`

func TestPullRequestTitle(t *testing.T) {
	if got, want := story("").PullRequestTitle(), "[S-0393] Push a story's branch on its move to review"; got != want {
		t.Errorf("PullRequestTitle() = %q, want %q", got, want)
	}
}

func TestPullRequestBodyFillsTheProjectsTemplate(t *testing.T) {
	got, err := story(projectTemplate).PullRequestBody()
	if err != nil {
		t.Fatalf("PullRequestBody(): %v", err)
	}
	want := `Thanks for the change.

## Story

<!-- flai:story -->
S-0393
<!-- /flai:story -->

## Criteria

<!--flai:criteria-->
- [x] The branch is pushed.
- [ ] The request is opened.
<!-- /flai:criteria -->

## Our definition of done

<!-- flai:definition-of-done -->
- [ ] Reviewed by a human
<!-- /flai:definition-of-done -->

## Verified

<!-- flai:verify -->
Passed at commit 0123456789abcdef0123456789abcdef01234567, run 2026-10-17T06:00:00Z.
<!-- /flai:verify -->

Signed off by the team.
`
	if got != want {
		t.Errorf("PullRequestBody() =\n%s\nwant\n%s", got, want)
	}
}

func TestPullRequestBodyWithoutATemplateHasTheFourPartsUnderHeadings(t *testing.T) {
	for _, template := range []string{"", " \n\n"} {
		got, err := story(template).PullRequestBody()
		if err != nil {
			t.Fatalf("PullRequestBody() with template %q: %v", template, err)
		}
		want := "## Story\n\n<!-- flai:story -->\nS-0393\n<!-- /flai:story -->\n\n" +
			"## Acceptance criteria\n\n<!-- flai:criteria -->\n- [x] The branch is pushed.\n- [ ] The request is opened.\n<!-- /flai:criteria -->\n\n" +
			"## Definition of done\n\n<!-- flai:definition-of-done -->\n" + defaultDefinitionOfDone + "\n<!-- /flai:definition-of-done -->\n\n" +
			"## Last flai verify\n\n<!-- flai:verify -->\nPassed at commit 0123456789abcdef0123456789abcdef01234567, run 2026-10-17T06:00:00Z.\n<!-- /flai:verify -->\n"
		if got != want {
			t.Errorf("PullRequestBody() with template %q =\n%s\nwant\n%s", template, got, want)
		}
	}
}

func TestPullRequestBodyWithNoVerifyResultSaysSo(t *testing.T) {
	s := story("")
	s.Verify = nil
	s.Criteria = nil
	got, err := s.PullRequestBody()
	if err != nil {
		t.Fatalf("PullRequestBody(): %v", err)
	}
	for _, want := range []string{
		"<!-- flai:verify -->\nNo `flai verify` result is stored for S-0393; run `flai verify S-0393`.\n<!-- /flai:verify -->",
		"<!-- flai:criteria -->\nThe story lists no acceptance criteria.\n<!-- /flai:criteria -->",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("PullRequestBody() =\n%s\nwant it to contain\n%s", got, want)
		}
	}
}

func TestPullRequestBodyNamesTheStepAVerifyStoppedAt(t *testing.T) {
	s := story("")
	s.Verify = &VerifyResult{Passed: false, StoppedAt: "go-test", Commit: "abc"}
	got, err := s.PullRequestBody()
	if err != nil {
		t.Fatalf("PullRequestBody(): %v", err)
	}
	if want := "<!-- flai:verify -->\nStopped at step `go-test` at commit abc.\n<!-- /flai:verify -->"; !strings.Contains(got, want) {
		t.Errorf("PullRequestBody() =\n%s\nwant it to contain\n%s", got, want)
	}
}

func TestPullRequestBodyAppendsThePartsATemplateDoesNotMark(t *testing.T) {
	template := "## Story\n\n<!-- flai:story -->\n<!-- /flai:story -->\n\n## Definition of done\n\n- [ ] Reviewed"
	got, err := story(template).PullRequestBody()
	if err != nil {
		t.Fatalf("PullRequestBody(): %v", err)
	}
	want := "## Story\n\n<!-- flai:story -->\nS-0393\n<!-- /flai:story -->\n\n## Definition of done\n\n- [ ] Reviewed\n" +
		"\n## Acceptance criteria\n\n<!-- flai:criteria -->\n- [x] The branch is pushed.\n- [ ] The request is opened.\n<!-- /flai:criteria -->\n" +
		"\n## Definition of done\n\n<!-- flai:definition-of-done -->\n" + defaultDefinitionOfDone + "\n<!-- /flai:definition-of-done -->\n" +
		"\n## Last flai verify\n\n<!-- flai:verify -->\nPassed at commit 0123456789abcdef0123456789abcdef01234567, run 2026-10-17T06:00:00Z.\n<!-- /flai:verify -->\n"
	if got != want {
		t.Errorf("PullRequestBody() =\n%s\nwant\n%s", got, want)
	}
}

func TestPullRequestBodyFillsEveryPairOfAPartsMarks(t *testing.T) {
	template := "<!-- flai:story -->\n<!-- /flai:story -->\nand again <!-- flai:story --> S- <!-- /flai:story -->\n"
	got, err := story(template).PullRequestBody()
	if err != nil {
		t.Fatalf("PullRequestBody(): %v", err)
	}
	if want := "<!-- flai:story -->\nS-0393\n<!-- /flai:story -->\nand again <!-- flai:story -->\nS-0393\n<!-- /flai:story -->\n"; !strings.HasPrefix(got, want) {
		t.Errorf("PullRequestBody() =\n%s\nwant it to begin\n%s", got, want)
	}
}

func TestPullRequestBodyIsTheSameBuiltFromItself(t *testing.T) {
	s := story(projectTemplate)
	first, err := s.PullRequestBody()
	if err != nil {
		t.Fatalf("PullRequestBody(): %v", err)
	}
	s.Template = first
	second, err := s.PullRequestBody()
	if err != nil {
		t.Fatalf("PullRequestBody() from its own body: %v", err)
	}
	if second != first {
		t.Errorf("PullRequestBody() from its own body =\n%s\nwant\n%s", second, first)
	}
}

func TestPullRequestBodyRefusesAnUnpairedMark(t *testing.T) {
	for _, tc := range []struct{ template, want string }{
		{"<!-- flai:criteria -->\n", "<!-- flai:criteria --> without <!-- /flai:criteria --> after it"},
		{"<!-- /flai:verify -->\n<!-- flai:verify -->\n", "<!-- /flai:verify --> without <!-- flai:verify --> before it"},
	} {
		_, err := story(tc.template).PullRequestBody()
		if err == nil {
			t.Errorf("PullRequestBody() with template %q: no error, want one", tc.template)
			continue
		}
		for _, want := range []string{tc.want, "S-0393", ".github/pull_request_template.md"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("PullRequestBody() with template %q: error %q does not say %q", tc.template, err, want)
			}
		}
	}
}

func TestTheTemplatesPullRequestTemplateMarksEveryPart(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "template", "root", ".github", "pull_request_template.md"))
	if err != nil {
		t.Fatalf("read the template's pull request template: %v", err)
	}
	got, err := story(string(data)).PullRequestBody()
	if err != nil {
		t.Fatalf("PullRequestBody() from the template's pull request template: %v", err)
	}
	for _, p := range templateParts {
		if n := strings.Count(got, "<!-- flai:"+p.name+" -->"); n != 1 {
			t.Errorf("the body built from the template's pull request template has <!-- flai:%s --> %d times, want once", p.name, n)
		}
	}
	if want := "<!-- flai:definition-of-done -->\n" + defaultDefinitionOfDone + "\n<!-- /flai:definition-of-done -->"; !strings.Contains(got, want) {
		t.Errorf("the template's definition of done differs from flai's default; want the body to contain\n%s\ngot\n%s", want, got)
	}
}
