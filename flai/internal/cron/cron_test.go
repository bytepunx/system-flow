package cron

import (
	"strings"
	"testing"
	"time"
)

func at(text string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", text)
	if err != nil {
		panic(err)
	}
	return t
}

func TestNext(t *testing.T) {
	tests := []struct {
		name  string
		spec  string
		after string
		want  []string // the next times in turn
	}{
		{"daily is midnight", "daily", "2026-10-04 13:07:00",
			[]string{"2026-10-05 00:00:00", "2026-10-06 00:00:00"}},
		{"every fifteen minutes", "*/15 * * * *", "2026-10-04 13:07:00",
			[]string{"2026-10-04 13:15:00", "2026-10-04 13:30:00", "2026-10-04 13:45:00", "2026-10-04 14:00:00"}},
		{"weekdays skip the weekend", "30 9 * * 1-5", "2026-10-02 10:00:00", // a Friday
			[]string{"2026-10-05 09:30:00", "2026-10-06 09:30:00"}},
		{"first of the month across a month", "0 0 1 * *", "2026-10-04 00:00:00",
			[]string{"2026-11-01 00:00:00", "2026-12-01 00:00:00", "2027-01-01 00:00:00"}},
		{"first of the month across a year end", "0 0 1 * *", "2026-12-31 23:59:00",
			[]string{"2027-01-01 00:00:00"}},
		{"day of month or day of week when both are restricted", "0 12 13 * 5", "2026-11-01 00:00:00",
			// Every Friday, and 13 December, a Sunday, by its day of month.
			[]string{"2026-11-06 12:00:00", "2026-11-13 12:00:00", "2026-11-20 12:00:00", "2026-11-27 12:00:00", "2026-12-04 12:00:00", "2026-12-11 12:00:00", "2026-12-13 12:00:00"}},
		{"day of month alone when day of week is *", "0 0 13 * *", "2026-11-01 00:00:00",
			[]string{"2026-11-13 00:00:00", "2026-12-13 00:00:00"}},
		{"Sunday as 7", "5 4 * * 7", "2026-10-04 04:05:00", // a Sunday, on the match
			[]string{"2026-10-11 04:05:00", "2026-10-18 04:05:00"}},
		{"Sunday as 0", "5 4 * * 0", "2026-10-05 00:00:00",
			[]string{"2026-10-11 04:05:00"}},
		{"a range with a step", "0 8-18/2 * * *", "2026-10-04 17:00:00",
			[]string{"2026-10-04 18:00:00", "2026-10-05 08:00:00", "2026-10-05 10:00:00"}},
		{"a list of ranges and values", "0,30 9-10,17 * * *", "2026-10-04 10:15:00",
			[]string{"2026-10-04 10:30:00", "2026-10-04 17:00:00", "2026-10-04 17:30:00", "2026-10-05 09:00:00"}},
		{"a start with a step", "45/5 * * * *", "2026-10-04 10:50:00",
			[]string{"2026-10-04 10:55:00", "2026-10-04 11:45:00"}},
		{"a month list", "0 6 1 1,7 *", "2026-02-01 00:00:00",
			[]string{"2026-07-01 06:00:00", "2027-01-01 06:00:00"}},
		{"strictly after a match", "30 9 * * *", "2026-10-04 09:30:00",
			[]string{"2026-10-05 09:30:00"}},
		{"seconds past a match", "30 9 * * *", "2026-10-04 09:30:20",
			[]string{"2026-10-05 09:30:00"}},
		{"seconds before a match", "30 9 * * *", "2026-10-04 09:29:59",
			[]string{"2026-10-04 09:30:00"}},
		{"February 29th across a century year that is not leap", "0 0 29 2 *", "2096-03-01 00:00:00",
			[]string{"2104-02-29 00:00:00"}},
		{"31st skips the short months", "0 0 31 * *", "2026-04-01 00:00:00",
			[]string{"2026-05-31 00:00:00", "2026-07-31 00:00:00"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := Parse(tc.spec)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tc.spec, err)
			}
			cur := at(tc.after)
			for _, w := range tc.want {
				got := s.Next(cur)
				if !got.Equal(at(w)) {
					t.Fatalf("Next(%s) = %s, want %s", cur.Format(time.DateTime), got.Format(time.DateTime), w)
				}
				if got.Location() != time.UTC {
					t.Fatalf("Next(%s) is in %s, want UTC", cur.Format(time.DateTime), got.Location())
				}
				cur = got
			}
		})
	}
}

func TestNextComputesInUTC(t *testing.T) {
	s, err := Parse("0 6 * * *")
	if err != nil {
		t.Fatal(err)
	}
	// 05:30 in New York on 4 October 2026 is 09:30 UTC, past 06:00 UTC.
	ny := time.FixedZone("EDT", -4*60*60)
	got := s.Next(time.Date(2026, 10, 4, 5, 30, 0, 0, ny))
	want := at("2026-10-05 06:00:00")
	if !got.Equal(want) || got.Location() != time.UTC {
		t.Fatalf("Next = %s, want %s UTC", got, want.Format(time.DateTime))
	}
}

func TestString(t *testing.T) {
	for spec, want := range map[string]string{
		"  daily\n":        "daily",
		" 0  6 * * 1-5 ":   "0  6 * * 1-5",
		"*/15 * * * *":     "*/15 * * * *",
		"\t5 4 * * 7\t":    "5 4 * * 7",
		"0 8-18/2 * * *  ": "0 8-18/2 * * *",
	} {
		s, err := Parse(spec)
		if err != nil {
			t.Fatalf("Parse(%q): %v", spec, err)
		}
		if got := s.String(); got != want {
			t.Errorf("Parse(%q).String() = %q, want %q", spec, got, want)
		}
	}
}

func TestParseRefuses(t *testing.T) {
	tests := []struct {
		spec string
		want []string // what the error must say
	}{
		{"", []string{"empty"}},
		{"   ", []string{"empty"}},
		{"* * * *", []string{"4 fields", "not five", "daily"}},
		{"* * * * * *", []string{"6 fields", "not five"}},
		{"weekly", []string{"1 fields", "not five", "daily"}},
		{"60 * * * *", []string{"minute", `"60"`, "0-59"}},
		{"* 24 * * *", []string{"hour", `"24"`, "0-23"}},
		{"* * 0 * *", []string{"day of month", `"0"`, "1-31"}},
		{"* * 32 * *", []string{"day of month", `"32"`, "1-31"}},
		{"* * * 13 *", []string{"month", `"13"`, "1-12"}},
		{"* * * 0 *", []string{"month", `"0"`, "1-12"}},
		{"* * * * 8", []string{"day of week", `"8"`, "0-7"}},
		{"*/0 * * * *", []string{"minute step", `"*/0"`, "1-59"}},
		{"*/60 * * * *", []string{"minute step", `"*/60"`}},
		{"*/x * * * *", []string{"minute step", "not a number"}},
		{"5-3 * * * *", []string{"minute range", `"5-3"`, "backwards"}},
		{"a * * * *", []string{"minute", `"a"`, "not a number"}},
		{"* * * JAN *", []string{"month", `"JAN"`, "names"}},
		{"* * * * MON", []string{"day of week", `"MON"`}},
		{"-5 * * * *", []string{"minute", `"-5"`, "not a number"}},
		{"1,,2 * * * *", []string{"minute", "empty item"}},
		{"1-2-3 * * * *", []string{"minute", `"1-2-3"`}},
		{"@daily", []string{"1 fields"}},
		{"0 0 31 2 *", []string{"never comes round", "February 31st"}},
		{"0 0 30,31 2 *", []string{"never comes round"}},
		{"0 0 31 4,6,9,11 *", []string{"never comes round"}},
	}
	for _, tc := range tests {
		t.Run(tc.spec, func(t *testing.T) {
			_, err := Parse(tc.spec)
			if err == nil {
				t.Fatalf("Parse(%q) took it, want a refusal", tc.spec)
			}
			msg := err.Error()
			for _, w := range append(tc.want, `"0 6 * * 1-5"`) {
				if !strings.Contains(msg, w) {
					t.Errorf("Parse(%q) = %q, want it to say %q", tc.spec, msg, w)
				}
			}
		})
	}
}

func TestParseTakesDayFieldsThatComeRound(t *testing.T) {
	// A restricted day of week makes any day of month possible, as either
	// may match.
	for _, spec := range []string{"0 0 31 2 1", "0 0 29 2 *", "0 0 30 1-2 *", "0 0 * 2 *"} {
		if _, err := Parse(spec); err != nil {
			t.Errorf("Parse(%q): %v", spec, err)
		}
	}
}

func TestNextOfZeroScheduleIsZero(t *testing.T) {
	var s Schedule
	if got := s.Next(at("2026-10-04 00:00:00")); !got.IsZero() {
		t.Fatalf("Next of an unparsed schedule = %s, want the zero time", got)
	}
}
