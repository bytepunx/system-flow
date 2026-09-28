package workitem

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

// S-0135: stories and epics carry topics; tasks do not.
func TestCreateWithTopics(t *testing.T) {
	r := newProject(t)
	e, err := r.Create(NewOptions{Type: Epic, Title: "E", Topics: []string{"release"}, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	s, err := r.Create(NewOptions{Type: Story, Title: "S", Parent: e.ID, Topics: []string{" logging", "release", "logging", ""}, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"logging", "release"}; !reflect.DeepEqual(s.Topics, want) {
		t.Errorf("story topics = %v, want %v", s.Topics, want)
	}
	data, _ := os.ReadFile(s.Path)
	if !strings.Contains(string(data), "\ntags: []\ntopics: [logging, release]\n") {
		t.Errorf("topics follow tags in the front matter:\n%s", data)
	}
	back, err := ReadItem(s.Path)
	if err != nil {
		t.Fatal(err)
	}
	if back.Marshal() != string(data) || !reflect.DeepEqual(back.Topics, s.Topics) {
		t.Errorf("a story with topics does not round-trip:\n%s", back.Marshal())
	}
	if e2, _ := r.Get(e.ID); !reflect.DeepEqual(e2.Topics, []string{"release"}) {
		t.Errorf("epic topics = %v", e2.Topics)
	}
	if _, err := r.Create(NewOptions{Type: Task, Title: "T", Parent: s.ID, Topics: []string{"logging"}, Now: t0}); err == nil || !strings.Contains(err.Error(), "stories and epics") {
		t.Errorf("a task with topics must be refused: %v", err)
	}
	if _, err := r.Create(NewOptions{Type: Story, Title: "S2", Topics: []string{"two words"}, Now: t0}); err == nil || !strings.Contains(err.Error(), "one word") {
		t.Errorf("a topic of two words must be refused: %v", err)
	}
}

func TestValidateTopics(t *testing.T) {
	it := &Item{ID: "T-0001", Type: Task, Nature: "feature", Title: "x", Status: Backlog, Parent: "S-0001", Owner: "a",
		Created: "2026-09-15T10:00:00Z", Updated: "2026-09-15T10:00:00Z", Topics: []string{"logging"}}
	if err := it.Validate(); err == nil || !strings.Contains(err.Error(), "topics are for stories and epics") {
		t.Errorf("a task's topics must be refused: %v", err)
	}
	it.Type, it.ID, it.Parent, it.Topics = Story, "S-0002", "", []string{"-dash"}
	if err := it.Validate(); err == nil || !strings.Contains(err.Error(), `topics[0] "-dash"`) {
		t.Errorf("a topic that is not a word must be refused: %v", err)
	}
	it.Topics = []string{"logging", "v1.2", "a_b"}
	if err := it.Validate(); err != nil {
		t.Errorf("words are topics: %v", err)
	}
}

func TestStrictParsingAcceptsTopics(t *testing.T) {
	doc := "---\nid: E-0001\ntype: epic\nnature: feature\ntitle: E\nstatus: backlog\nowner: a\ncreated: 2026-09-15T10:00:00Z\nupdated: 2026-09-15T10:00:00Z\ntransitions: []\ntags: []\ntopics: [logging]\n---\n# E-0001 E\n"
	it, err := ParseItem(doc)
	if err != nil {
		t.Fatalf("strict parsing accepts topics: %v", err)
	}
	if it.Marshal() != doc {
		t.Errorf("round trip:\n%s", it.Marshal())
	}
}
