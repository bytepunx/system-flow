// Package cron reads the planner's schedule, planning.schedule in
// system-flow.yaml: a five-field cron expression in UTC or daily, and finds
// when it next comes round (S-0211, ADR-0084).
package cron

import (
	"fmt"
	"strings"
	"time"
)

// Daily is the one named schedule, midnight UTC every day.
const Daily = "daily"

// hint ends every refusal with a schedule that parses.
const hint = `write one such as "0 6 * * 1-5" (06:00 UTC on weekdays) or daily`

// horizon is how many years Next looks ahead. The longest gap between two
// times a schedule can name is February 29th across a century year that
// is not a leap year, eight years, so nine reach every time that comes round.
const horizon = 9

// field is one of the five fields of an expression and the values it takes.
type field struct {
	name string
	min  int
	// max is where * and a/n end.
	max int
	// limit is the largest value accepted, above max for Sunday as 7.
	limit int
	// accepted says the values accepted, for a refusal.
	accepted string
}

var fields = [5]field{
	{name: "minute", min: 0, max: 59, limit: 59, accepted: "0-59"},
	{name: "hour", min: 0, max: 23, limit: 23, accepted: "0-23"},
	{name: "day of month", min: 1, max: 31, limit: 31, accepted: "1-31"},
	{name: "month", min: 1, max: 12, limit: 12, accepted: "1-12"},
	{name: "day of week", min: 0, max: 6, limit: 7, accepted: "0-7, Sunday as 0 or 7"},
}

// Schedule is a parsed expression: the times it names, in UTC.
type Schedule struct {
	spec string
	// The sets hold bit v for each value v the field matches.
	minute, hour, dom, month, dow uint64
	// domStar and dowStar say the day fields were written *, so the day is
	// matched by the other field alone.
	domStar, dowStar bool
}

// Parse reads a five-field cron expression or daily.
func Parse(spec string) (Schedule, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return Schedule{}, fmt.Errorf("the schedule is empty; %s", hint)
	}
	expr := spec
	if spec == Daily {
		expr = "0 0 * * *"
	}
	parts := strings.Fields(expr)
	if len(parts) != len(fields) {
		return Schedule{}, fmt.Errorf("schedule %q has %d fields, not five (minute, hour, day of month, month, day of week); %s", spec, len(parts), hint)
	}
	var sets [5]uint64
	for i, part := range parts {
		set, err := parseField(fields[i], part)
		if err != nil {
			return Schedule{}, fmt.Errorf("schedule %q: %w; %s", spec, err, hint)
		}
		sets[i] = set
	}
	// Sunday is 0 and 7.
	if sets[4]&(1<<7) != 0 {
		sets[4] = sets[4]&^(1<<7) | 1
	}
	s := Schedule{
		spec:    spec,
		minute:  sets[0],
		hour:    sets[1],
		dom:     sets[2],
		month:   sets[3],
		dow:     sets[4],
		domStar: parts[2] == "*",
		dowStar: parts[4] == "*",
	}
	if !s.possible() {
		return Schedule{}, fmt.Errorf("schedule %q never comes round: no month it names has a day it names, as with February 31st; %s", spec, hint)
	}
	return s, nil
}

// parseField reads one field's comma list into its set.
func parseField(f field, text string) (uint64, error) {
	var set uint64
	for _, item := range strings.Split(text, ",") {
		lo, hi, step, err := parseItem(f, item)
		if err != nil {
			return 0, err
		}
		for v := lo; v <= hi; v += step {
			set |= 1 << uint(v)
		}
	}
	return set, nil
}

// parseItem reads one item of a list: *, a, a-b, each with an optional /n.
func parseItem(f field, item string) (lo, hi, step int, err error) {
	if item == "" {
		return 0, 0, 0, fmt.Errorf("the %s has an empty item in its comma list", f.name)
	}
	base, stepText, stepped := strings.Cut(item, "/")
	step = 1
	if stepped {
		n, ok := number(stepText)
		if !ok {
			return 0, 0, 0, fmt.Errorf("the %s step %q is not a number", f.name, item)
		}
		if n < 1 || n > f.max {
			return 0, 0, 0, fmt.Errorf("the %s step %q is outside 1-%d", f.name, item, f.max)
		}
		step = n
	}
	if base == "*" {
		return f.min, f.max, step, nil
	}
	loText, hiText, ranged := strings.Cut(base, "-")
	lo, err = value(f, loText, item)
	if err != nil {
		return 0, 0, 0, err
	}
	switch {
	case ranged:
		hi, err = value(f, hiText, item)
		if err != nil {
			return 0, 0, 0, err
		}
		if hi < lo {
			return 0, 0, 0, fmt.Errorf("the %s range %q runs backwards, from high to low", f.name, item)
		}
	case stepped:
		hi = f.max
		if lo > hi {
			hi = lo
		}
	default:
		hi = lo
	}
	return lo, hi, step, nil
}

// value reads one number of a field and checks it is in range.
func value(f field, text, item string) (int, error) {
	n, ok := number(text)
	if !ok {
		return 0, fmt.Errorf("the %s %q is not a number, a range a-b, or a step */n, a-b/n, or a/n (names such as MON are not read)", f.name, item)
	}
	if n < f.min || n > f.limit {
		return 0, fmt.Errorf("the %s %q is outside %s", f.name, item, f.accepted)
	}
	return n, nil
}

// number reads a run of up to nine decimal digits.
func number(text string) (int, bool) {
	if text == "" || len(text) > 9 {
		return 0, false
	}
	n := 0
	for _, c := range text {
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
	}
	return n, true
}

// daysIn is the most days each month has, February in a leap year.
var daysIn = [13]int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// possible says some day of some month the schedule names exists. A day of
// week names a day in every month, so only a day of month alone can miss.
func (s Schedule) possible() bool {
	if s.domStar || !s.dowStar {
		return true
	}
	for m := 1; m <= 12; m++ {
		if !has(s.month, m) {
			continue
		}
		for d := 1; d <= daysIn[m]; d++ {
			if has(s.dom, d) {
				return true
			}
		}
	}
	return false
}

// String returns the schedule as it was written, trimmed.
func (s Schedule) String() string {
	return s.spec
}

// Next returns the first whole minute strictly after after, in UTC, that
// the schedule names. It returns the zero time when none comes within nine
// years, which a schedule from Parse never does.
func (s Schedule) Next(after time.Time) time.Time {
	t := after.UTC().Truncate(time.Minute).Add(time.Minute)
	end := t.AddDate(horizon, 0, 0)
	for !t.After(end) {
		y, m, d := t.Date()
		switch {
		case !has(s.month, int(m)):
			t = time.Date(y, m+1, 1, 0, 0, 0, 0, time.UTC)
		case !s.day(t):
			t = time.Date(y, m, d+1, 0, 0, 0, 0, time.UTC)
		case !has(s.hour, t.Hour()):
			t = time.Date(y, m, d, t.Hour()+1, 0, 0, 0, time.UTC)
		case !has(s.minute, t.Minute()):
			t = t.Add(time.Minute)
		default:
			return t
		}
	}
	return time.Time{}
}

// day says the schedule names t's day. When both day fields are restricted
// either may match, as in cron; otherwise both must.
func (s Schedule) day(t time.Time) bool {
	dom := has(s.dom, t.Day())
	dow := has(s.dow, int(t.Weekday()))
	if s.domStar || s.dowStar {
		return dom && dow
	}
	return dom || dow
}

// has says set holds v.
func has(set uint64, v int) bool {
	return set&(1<<uint(v)) != 0
}
