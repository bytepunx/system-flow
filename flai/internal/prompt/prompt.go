// Package prompt asks a person questions one line at a time: a value, a yes
// or no, or one of a few options. It reads answers from any reader and writes
// questions to any writer, so a terminal and a test drive it alike.
//
// flai asked through charmbracelet/huh until S-0160. huh brings in
// atotto/clipboard, which searches PATH for clipboard programs when its
// package loads, and that cost every flai process up to 145 ms on a host with
// a long PATH, prompt or none.
package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ErrNoAnswer is returned when the input ends before a question is answered.
var ErrNoAnswer = errors.New("input ended before an answer")

// Prompter asks questions on out and reads the answers from in.
type Prompter struct {
	in  *bufio.Reader
	out io.Writer
}

// New returns a Prompter that reads answers from in and writes questions to out.
func New(in io.Reader, out io.Writer) *Prompter {
	return &Prompter{in: bufio.NewReader(in), out: out}
}

// Option is one choice of a Select: what is shown and what is returned.
type Option struct {
	Label string
	Value string
}

// Input asks for a value. An empty answer takes def; validate, when not nil,
// refuses an answer with its error, and the question is asked again.
func (p *Prompter) Input(title, description, def string, validate func(string) error) (string, error) {
	if description != "" {
		fmt.Fprintln(p.out, description)
	}
	q := title
	if def != "" {
		q += " [" + def + "]"
	}
	for {
		ans, err := p.ask(q + ": ")
		if err != nil {
			return "", fmt.Errorf("%s: %w", title, err)
		}
		if ans == "" {
			ans = def
		}
		if validate != nil {
			if err := validate(ans); err != nil {
				fmt.Fprintf(p.out, "  %v\n", err)
				continue
			}
		}
		return ans, nil
	}
}

// Confirm asks a yes or no question. An empty answer takes def.
func (p *Prompter) Confirm(title string, def bool) (bool, error) {
	hint := "[y/N]"
	if def {
		hint = "[Y/n]"
	}
	for {
		ans, err := p.ask(title + " " + hint + " ")
		if err != nil {
			return false, fmt.Errorf("%s: %w", title, err)
		}
		switch strings.ToLower(ans) {
		case "":
			return def, nil
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}
		fmt.Fprintln(p.out, "  answer y or n")
	}
}

// Select asks for one of options by its number and returns its Value. An
// empty answer takes the option whose Value is def.
func (p *Prompter) Select(title string, options []Option, def string) (string, error) {
	fmt.Fprintln(p.out, title)
	defNum := 0
	for i, o := range options {
		fmt.Fprintf(p.out, "  %d) %s\n", i+1, o.Label)
		if o.Value == def && defNum == 0 {
			defNum = i + 1
		}
	}
	q := fmt.Sprintf("Choose 1-%d", len(options))
	if defNum > 0 {
		q += fmt.Sprintf(" [%d]", defNum)
	}
	for {
		ans, err := p.ask(q + ": ")
		if err != nil {
			return "", fmt.Errorf("%s: %w", title, err)
		}
		n := defNum
		if ans != "" {
			n, _ = strconv.Atoi(ans)
		}
		if n >= 1 && n <= len(options) {
			return options[n-1].Value, nil
		}
		fmt.Fprintf(p.out, "  answer a number from 1 to %d\n", len(options))
	}
}

// ask writes q and reads one line, trimmed. A last line without a newline is
// an answer; input that ends with nothing on the line is ErrNoAnswer.
func (p *Prompter) ask(q string) (string, error) {
	fmt.Fprint(p.out, q)
	line, err := p.in.ReadString('\n')
	if err != nil && (!errors.Is(err, io.EOF) || line == "") {
		fmt.Fprintln(p.out)
		if errors.Is(err, io.EOF) {
			return "", ErrNoAnswer
		}
		return "", err
	}
	return strings.TrimSpace(line), nil
}
