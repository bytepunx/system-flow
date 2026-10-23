package digest

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Match is text with the shape of a secret: the 1-based line it is on, the
// kind of secret it looks like, and a description that names the kind and
// shows the value's first characters only.
type Match struct {
	Line        int    `json:"line"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
}

// shape is one kind of secret: what it looks like, and a test the text it
// matched must also pass, when the pattern alone would match prose.
type shape struct {
	kind, name string
	re         *regexp.Regexp
	value      int // the submatch that is the secret's value
	also       func(string) bool
}

// assigned is a long value right after a word naming a key, token, secret,
// or password and an = or a colon.
var assigned = regexp.MustCompile(`(?i)\b([A-Za-z0-9_.-]*(?:key|token|secret|passwd|password)[A-Za-z0-9_.-]*)["']?\s*[:=]\s*["']?([A-Za-z0-9+/_-]{20,}={0,2})(?:[^A-Za-z0-9+/_.=-]|$)`)

var shapes = []shape{
	{kind: "private-key", name: "a private key block", re: regexp.MustCompile(`-----BEGIN (?:[A-Z0-9]+ )*PRIVATE KEY-----`)},
	{kind: "api-key", name: "an API key", re: regexp.MustCompile(`\b(sk-(?:ant-|proj-)?[A-Za-z0-9_-]{20,})`), value: 1, also: random},
	{kind: "github-token", name: "a GitHub token", re: regexp.MustCompile(`\b((?:ghp|gho|ghs|ghu|ghr)_[A-Za-z0-9]{36})\b`), value: 1},
	{kind: "github-token", name: "a GitHub token", re: regexp.MustCompile(`\b(github_pat_[A-Za-z0-9_]{22,})`), value: 1},
	{kind: "slack-token", name: "a Slack token", re: regexp.MustCompile(`\b(xox[abposr]-[A-Za-z0-9-]{10,})`), value: 1, also: hasDigit},
	{kind: "aws-access-key", name: "an AWS access key", re: regexp.MustCompile(`\b((?:AKIA|ASIA)[A-Z0-9]{16})\b`), value: 1},
	{kind: "jwt", name: "a JSON web token", re: regexp.MustCompile(`\b(eyJ[A-Za-z0-9_-]{10,}\.eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,})`), value: 1},
}

// Secrets finds the text with the shape of a key or token: a private key
// block, an sk- API key, a GitHub, Slack, or AWS token, a JSON web token, or
// a long random value given to a word naming a key, token, secret, or
// password. A word such as token in prose, an environment variable's name, a
// path, a short value, or a commit hash given to no such word is not one.
func Secrets(text string) []Match {
	var out []Match
	for i, line := range strings.Split(text, "\n") {
		var taken [][2]int
		for _, s := range shapes {
			for _, at := range s.re.FindAllStringSubmatchIndex(line, -1) {
				v := line[at[2*s.value]:at[2*s.value+1]]
				if s.also != nil && !s.also(v) {
					continue
				}
				taken = append(taken, [2]int{at[0], at[1]})
				desc := s.name
				if s.value > 0 {
					desc = fmt.Sprintf("%s (%s)", s.name, redact(v))
				}
				out = append(out, Match{Line: i + 1, Kind: s.kind, Description: desc})
			}
		}
		for _, at := range assigned.FindAllStringSubmatchIndex(line, -1) {
			word, v := line[at[2]:at[3]], line[at[4]:at[5]]
			if overlaps(taken, at[4], at[5]) || !hex(v) && !random(v) {
				continue
			}
			out = append(out, Match{Line: i + 1, Kind: "assigned-secret", Description: fmt.Sprintf("a value given to %s (%s)", word, redact(v))})
		}
	}
	return out
}

// redact shows a value's first four characters and its length, never more.
func redact(v string) string {
	n := min(4, len(v))
	return fmt.Sprintf("%s… of %d characters", v[:n], len(v))
}

func overlaps(spans [][2]int, from, to int) bool {
	for _, s := range spans {
		if from < s[1] && s[0] < to {
			return true
		}
	}
	return false
}

// hex reports a run of 32 or more hexadecimal digits.
func hex(v string) bool {
	if len(v) < 32 {
		return false
	}
	for _, r := range v {
		if !unicode.Is(unicode.ASCII_Hex_Digit, r) {
			return false
		}
	}
	return true
}

func hasDigit(v string) bool { return strings.ContainsFunc(v, unicode.IsDigit) }

// random reports a value that looks generated rather than written: it has a
// digit and a letter, and either both cases or no dash or underscore, which a
// slug, a path, or an identifier has.
func random(v string) bool {
	if !hasDigit(v) || !strings.ContainsFunc(v, unicode.IsLetter) {
		return false
	}
	mixed := strings.ContainsFunc(v, unicode.IsUpper) && strings.ContainsFunc(v, unicode.IsLower)
	return mixed || !strings.ContainsAny(v, "-_")
}
