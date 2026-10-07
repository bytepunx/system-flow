package verify

import (
	"strconv"
	"strings"
)

// vitestReport is the part of vitest's JSON reporter output read here.
type vitestReport struct {
	TestResults []struct {
		Name             string `json:"name"`
		Status           string `json:"status"`
		Message          string `json:"message"`
		AssertionResults []struct {
			FullName        string   `json:"fullName"`
			Status          string   `json:"status"`
			FailureMessages []string `json:"failureMessages"`
			Location        *struct {
				Line int `json:"line"`
			} `json:"location"`
		} `json:"assertionResults"`
	} `json:"testResults"`
}

// parseVitest is a finding for each failed test in vitest's JSON reporter
// output, and one for each test file that failed with no test failing, as
// one that does not load does.
func parseVitest(stdout string, e env) []Finding {
	var r vitestReport
	if !decodeFirst(stdout, "testResults", &r) {
		return nil
	}
	var got []Finding
	for _, file := range r.TestResults {
		p := e.rel(file.Name)
		failed := 0
		for _, a := range file.AssertionResults {
			if a.Status != "failed" {
				continue
			}
			failed++
			f := Finding{Name: a.FullName, Path: p}
			var msg string
			if len(a.FailureMessages) > 0 {
				msg = a.FailureMessages[0]
				f.Message = vitestMessage(msg)
			}
			if a.Location != nil && a.Location.Line > 0 {
				f.Line = a.Location.Line
			} else {
				f.Line = stackLine(msg, file.Name)
			}
			got = append(got, f)
		}
		if failed == 0 && file.Status == "failed" {
			got = append(got, Finding{Name: "suite", Path: p, Message: clip(file.Message)})
		}
	}
	return got
}

// vitestMessage is a failure message without its stack.
func vitestMessage(msg string) string {
	var keep []string
	for _, line := range strings.Split(msg, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "at ") {
			break
		}
		keep = append(keep, line)
	}
	return clip(strings.Join(keep, "\n"))
}

// stackLine is the line of file that a stack names first, 0 when it names
// none.
func stackLine(stack, file string) int {
	if file == "" {
		return 0
	}
	i := strings.Index(stack, file+":")
	if i < 0 {
		return 0
	}
	rest := stack[i+len(file)+1:]
	end := strings.IndexFunc(rest, func(r rune) bool { return r < '0' || r > '9' })
	if end < 0 {
		end = len(rest)
	}
	n, _ := strconv.Atoi(rest[:end])
	return n
}
