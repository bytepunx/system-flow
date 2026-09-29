package prompt

import (
	"errors"
	"strings"
	"testing"
)

func prompter(in string) (*Prompter, *strings.Builder) {
	var out strings.Builder
	return New(strings.NewReader(in), &out), &out
}

func TestInput(t *testing.T) {
	p, out := prompter("\n  given  \n")
	got, err := p.Input("Project name", "Shown above", "flai", nil)
	if err != nil || got != "flai" {
		t.Fatalf("empty line: %q, %v", got, err)
	}
	if !strings.Contains(out.String(), "Shown above\nProject name [flai]: ") {
		t.Errorf("question: %q", out.String())
	}
	if got, err := p.Input("Project name", "", "flai", nil); err != nil || got != "given" {
		t.Errorf("an answer: %q, %v", got, err)
	}
}

func TestInputAsksAgainUntilValid(t *testing.T) {
	p, out := prompter("\nbad\ngood")
	got, err := p.Input("URL", "", "", func(s string) error {
		if s != "good" {
			return errors.New("not good")
		}
		return nil
	})
	if err != nil || got != "good" {
		t.Fatalf("got %q, %v", got, err)
	}
	if n := strings.Count(out.String(), "  not good\n"); n != 2 {
		t.Errorf("refused %d times, want 2:\n%s", n, out)
	}
}

func TestConfirm(t *testing.T) {
	p, out := prompter("\nmaybe\nYes\nn\n")
	if ok, err := p.Confirm("Apply?", true); err != nil || !ok {
		t.Errorf("default yes: %v, %v", ok, err)
	}
	if ok, err := p.Confirm("Apply?", false); err != nil || !ok {
		t.Errorf("maybe then Yes: %v, %v", ok, err)
	}
	if ok, err := p.Confirm("Apply?", true); err != nil || ok {
		t.Errorf("n: %v, %v", ok, err)
	}
	s := out.String()
	if !strings.Contains(s, "Apply? [Y/n] ") || !strings.Contains(s, "Apply? [y/N] ") || !strings.Contains(s, "  answer y or n\n") {
		t.Errorf("questions:\n%s", s)
	}
}

func TestSelect(t *testing.T) {
	opts := []Option{{"keep the project's version", "keep"}, {"replace it", "replace"}, {"show diff", "diff"}}
	p, out := prompter("\n0\nthree\n3\n")
	if got, err := p.Select("Conflict", opts, "replace"); err != nil || got != "replace" {
		t.Errorf("default: %q, %v", got, err)
	}
	if got, err := p.Select("Conflict", opts, "keep"); err != nil || got != "diff" {
		t.Errorf("3 after two refusals: %q, %v", got, err)
	}
	s := out.String()
	if !strings.Contains(s, "Conflict\n  1) keep the project's version\n  2) replace it\n  3) show diff\nChoose 1-3 [2]: ") {
		t.Errorf("question:\n%s", s)
	}
	if n := strings.Count(s, "  answer a number from 1 to 3\n"); n != 2 {
		t.Errorf("refused %d times, want 2:\n%s", n, s)
	}
}

func TestSelectWithoutDefaultNeedsANumber(t *testing.T) {
	p, _ := prompter("\n1\n")
	if got, err := p.Select("Where?", []Option{{"leave", "leave"}}, "none"); err != nil || got != "leave" {
		t.Errorf("got %q, %v", got, err)
	}
}

func TestEndOfInputIsNoAnswer(t *testing.T) {
	for _, c := range []struct {
		name, in string
		ask      func(*Prompter) error
	}{
		{"input", "", func(p *Prompter) error { _, err := p.Input("Name", "", "x", nil); return err }},
		{"confirm", "", func(p *Prompter) error { _, err := p.Confirm("Go?", true); return err }},
		{"select", "", func(p *Prompter) error { _, err := p.Select("Pick", []Option{{"a", "a"}}, "a"); return err }},
		{"after a refusal", "maybe\n", func(p *Prompter) error { _, err := p.Confirm("Go?", true); return err }},
	} {
		p, _ := prompter(c.in)
		if err := c.ask(p); !errors.Is(err, ErrNoAnswer) {
			t.Errorf("%s: %v, want ErrNoAnswer", c.name, err)
		}
	}
}
