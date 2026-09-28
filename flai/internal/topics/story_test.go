package topics

import (
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var projects = []manifest.Project{
	{Name: "flai", Path: "flai", Kind: "go", Tags: []string{"cli"}},
	{Name: "flaiover", Path: "flaiover/", Kind: "sveltekit", Tags: []string{"dashboard"}},
	{Name: "template", Path: "template", Kind: "template", Tags: []string{"template", "conventions"}},
}

// render is a story's topics as "topic: source; source" lines.
func render(list []StoryTopic) string {
	var b strings.Builder
	for _, t := range list {
		srcs := make([]string, len(t.Sources))
		for i, s := range t.Sources {
			srcs[i] = s.String()
		}
		b.WriteString(t.Topic + ": " + strings.Join(srcs, "; ") + "\n")
	}
	return b.String()
}

func story(tags, touches, topics []string) *workitem.Item {
	return &workitem.Item{ID: "S-0001", Type: workitem.Story, Parent: "E-0001", Tags: tags, Touches: touches, Topics: topics}
}

// S-0135, ADR-0047: each source a story's topics come from.
func TestOfStorySources(t *testing.T) {
	epic := &workitem.Item{ID: "E-0001", Type: workitem.Epic, Topics: []string{"release", "logging"}}
	tasks := []*workitem.Item{
		{ID: "T-0001", Type: workitem.Task, Parent: "S-0001", Status: workitem.InProgress, Touches: []string{"flaiover/src/lib"}},
		{ID: "T-0002", Type: workitem.Task, Parent: "S-0001", Status: workitem.Done, Touches: []string{"template/root"}},
		{ID: "T-0003", Type: workitem.Task, Parent: "S-0009", Status: workitem.Ready, Touches: []string{"template"}},
	}
	got := render(OfStory(story([]string{"cli", "research"}, []string{"./flai/cmd/", "docs/users"}, []string{"logging"}), epic, tasks, projects))
	want := `logging: own; epic E-0001
release: epic E-0001
flai: tag cli reaches flai; S-0001 touches flai/cmd, in flai
cli: tag cli reaches flai; S-0001 touches flai/cmd, in flai
go: tag cli reaches flai; S-0001 touches flai/cmd, in flai
flaiover: T-0001 touches flaiover/src/lib, in flaiover
dashboard: T-0001 touches flaiover/src/lib, in flaiover
sveltekit: T-0001 touches flaiover/src/lib, in flaiover
code: flai is code; flaiover is code
all: every story
`
	if got != want {
		t.Errorf("topics:\n%s\nwant:\n%s", got, want)
	}
}

// A story with no tags, touches, topics, or epic has all alone; a touch
// outside every sub-project, and a tag that names none, add nothing.
func TestOfStoryWithNothingIsAll(t *testing.T) {
	if got := Names(OfStory(story(nil, nil, nil), nil, nil, projects)); !reflect.DeepEqual(got, []string{"all"}) {
		t.Errorf("no tags or touches: %v", got)
	}
	if got := Names(OfStory(story([]string{"research"}, []string{"design/system/work-hierarchy.md", "flaiover-old"}, nil), nil, nil, projects)); !reflect.DeepEqual(got, []string{"all"}) {
		t.Errorf("outside every sub-project: %v", got)
	}
}

// A story that reaches only the template has its words but not code; one
// that names a sub-project's tag in touches reaches it as its path would.
func TestOfStoryTemplateIsNotCode(t *testing.T) {
	got := render(OfStory(story(nil, []string{"conventions"}, nil), nil, nil, projects))
	want := `template: S-0001 touches conventions, in template
conventions: S-0001 touches conventions, in template
all: every story
`
	if got != want {
		t.Errorf("topics:\n%s\nwant:\n%s", got, want)
	}
}

// A touch above a sub-project reaches it too: overlap works both ways.
func TestOfStoryTouchAboveASubProject(t *testing.T) {
	ps := []manifest.Project{{Name: "api", Path: "services/api", Kind: "go"}}
	if got := Names(OfStory(story(nil, []string{"services"}, nil), nil, nil, ps)); !reflect.DeepEqual(got, []string{"api", "go", "code", "all"}) {
		t.Errorf("topics: %v", got)
	}
}

// Summary names each sub-project once, however many touches reach it.
func TestStoryTopicSummary(t *testing.T) {
	epic := &workitem.Item{ID: "E-0001", Type: workitem.Epic, Topics: []string{"logging"}}
	tasks := []*workitem.Item{{ID: "T-0001", Type: workitem.Task, Parent: "S-0001", Status: workitem.Ready, Touches: []string{"flai/cmd", "flaiover"}}}
	list := OfStory(story([]string{"cli"}, []string{"flai/cmd", "flai/internal"}, []string{"logging"}), epic, tasks, projects)
	var got []string
	for _, t := range list {
		got = append(got, t.Topic+": "+t.Summary())
	}
	want := []string{
		"logging: own; epic E-0001",
		"flai: flai by tag cli and 2 touches (S-0001, T-0001)",
		"cli: flai by tag cli and 2 touches (S-0001, T-0001)",
		"go: flai by tag cli and 2 touches (S-0001, T-0001)",
		"flaiover: flaiover by touches flaiover (T-0001)",
		"dashboard: flaiover by touches flaiover (T-0001)",
		"sveltekit: flaiover by touches flaiover (T-0001)",
		"code: flai, flaiover are code",
		"all: every story",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("summaries:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}
