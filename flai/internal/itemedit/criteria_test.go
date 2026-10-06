package itemedit

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

const storyBody = `## Goal

Tick what was verified.

## Acceptance criteria

- [ ] flai criteria list numbers the boxes
- [x] flai criteria tick ticks them
- [ ] the dashboard ticks them too

## Tasks

- [ ] T-0001 A task
- [x] T-0002 Another

## Notes

- [ ] a box in the notes
`

// S-0282: the acceptance criteria are the boxes under their heading, nested
// ones included, numbered in document order; boxes elsewhere are not theirs.
func TestCriteriaListsTheBoxesUnderAcceptanceCriteria(t *testing.T) {
	for _, c := range []struct {
		name string
		body string
		want []Criterion
	}{
		{"a typical story", storyBody, []Criterion{
			{N: 1, Text: "flai criteria list numbers the boxes"},
			{N: 2, Text: "flai criteria tick ticks them", Ticked: true},
			{N: 3, Text: "the dashboard ticks them too"},
		}},
		{"nested boxes and a capital X", "## Acceptance criteria  \n\n- [ ] outer\n  - [X] inner\n\tnot a box\n- [x] trailing spaces  \n- [x]  two spaces is not a box\n- [ ]\n### Detail\n- [ ] under a subheading\n", []Criterion{
			{N: 1, Text: "outer"},
			{N: 2, Text: "inner", Ticked: true},
			{N: 3, Text: "trailing spaces", Ticked: true},
			{N: 4, Text: "under a subheading"},
		}},
		{"no section", "## Goal\n\n- [ ] not a criterion\n", []Criterion{}},
		{"a section without a box", "## Acceptance criteria\n\nSomething true.\n\n## Notes\n- [ ] note\n", []Criterion{}},
		{"the last section", "## Notes\n\n- [x] note\n\n## Acceptance criteria\n- [x] last", []Criterion{{N: 1, Text: "last", Ticked: true}}},
	} {
		t.Run(c.name, func(t *testing.T) {
			got := Criteria(c.body)
			if got == nil || !reflect.DeepEqual(got, c.want) {
				t.Errorf("Criteria = %#v, want %#v", got, c.want)
			}
		})
	}
}

// S-0282: a tick or an untick changes the marks it names and no other byte.
func TestTickSetsTheNamedBoxesAndNothingElse(t *testing.T) {
	for _, c := range []struct {
		name         string
		body         string
		tick, untick []int
		want         string
	}{
		{"tick", storyBody, []int{1, 3}, nil,
			strings.Replace(strings.Replace(storyBody, "- [ ] flai criteria list", "- [x] flai criteria list", 1), "- [ ] the dashboard", "- [x] the dashboard", 1)},
		{"untick", storyBody, nil, []int{2},
			strings.Replace(storyBody, "- [x] flai criteria tick", "- [ ] flai criteria tick", 1)},
		{"tick and untick at once, a number twice", storyBody, []int{3, 3}, []int{2},
			strings.Replace(strings.Replace(storyBody, "- [x] flai criteria tick", "- [ ] flai criteria tick", 1), "- [ ] the dashboard", "- [x] the dashboard", 1)},
		{"ticking a ticked box leaves the body as it was", storyBody, []int{2}, nil, storyBody},
		{"a capital X ticked again stays", "## Acceptance criteria\n- [X] done\n", []int{1}, nil, "## Acceptance criteria\n- [X] done\n"},
		{"a capital X is unticked", "## Acceptance criteria\n- [X] done\n", nil, []int{1}, "## Acceptance criteria\n- [ ] done\n"},
		{"a nested box", "## Acceptance criteria\n- [ ] outer\n  - [ ] inner\n", []int{2}, nil, "## Acceptance criteria\n- [ ] outer\n  - [x] inner\n"},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := Tick(c.body, c.tick, c.untick)
			if err != nil {
				t.Fatal(err)
			}
			if got != c.want {
				t.Errorf("Tick =\n%s\nwant\n%s", got, c.want)
			}
		})
	}
}

// S-0282: what cannot be ticked is refused as the caller's mistake, with
// what to do instead.
func TestTickRefusesWhatCannotBeTicked(t *testing.T) {
	for _, c := range []struct {
		name         string
		body         string
		tick, untick []int
		says         string
	}{
		{"no criteria", "## Goal\n\n- [ ] not a criterion\n", []int{1}, nil, "no acceptance criteria with a checkbox"},
		{"nothing to tick or untick", storyBody, nil, nil, "flai criteria list"},
		{"zero", storyBody, []int{0}, nil, "there are 3, numbered 1 to 3 as flai criteria list"},
		{"past the last", storyBody, nil, []int{4}, "no acceptance criterion 4: there are 3"},
		{"in both lists", storyBody, []int{1}, []int{1}, "both to tick and to untick"},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := Tick(c.body, c.tick, c.untick)
			var inv *InvalidError
			if !errors.As(err, &inv) {
				t.Fatalf("want an InvalidError, got %q, %v", got, err)
			}
			if !strings.Contains(err.Error(), c.says) {
				t.Errorf("message %q does not say %q", err, c.says)
			}
		})
	}
}
