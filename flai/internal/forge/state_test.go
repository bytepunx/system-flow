package forge

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// prAt is S-0394's pull request with head at commit.
func prAt(commit string) PullRequest {
	return PullRequest{Number: 42, URL: "https://github.com/o/r/pull/42", Head: "story/S-0394", Base: "main", HeadSHA: commit, State: PullRequestOpen}
}

func TestWriteStatusThenReadStatusRoundTrips(t *testing.T) {
	dir := t.TempDir()
	want := Status{
		Story: "S-0394", PullRequest: 42, URL: "https://github.com/o/r/pull/42", HeadSHA: "aaa", ReadAt: "2026-10-17T08:00:00Z",
		Checks:   []CheckStatus{{Name: "test", State: CheckFailed, URL: "https://ci/1", Summary: "2 failed", Since: "2026-10-17T07:55:00Z"}},
		Reviews:  []ReviewStatus{{Author: "alex", State: ReviewChangesRequested, URL: "https://github.com/o/r/pull/42#r1", Body: "Rename it.", SubmittedAt: "2026-10-17T07:50:00Z"}},
		Events:   []StatusEvent{{At: "2026-10-17T07:55:00Z", Kind: EventCheck, Name: "test", From: "pending", To: "failed", Commit: "aaa"}},
		SentBack: []SendBack{{Key: CheckKey("test", "aaa"), Kind: EventCheck, Name: "test", Commit: "aaa", URL: "https://ci/1", Thread: "M-0001", At: "2026-10-17T08:00:00Z"}},
	}
	if err := WriteStatus(dir, want); err != nil {
		t.Fatalf("WriteStatus: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "forge", "status", "S-0394.json"))
	if err != nil {
		t.Fatalf("the status file is not at its path: %v", err)
	}
	if !strings.HasSuffix(string(data), "}\n") {
		t.Errorf("the status file does not end in a newline: %q", data)
	}
	got, ok, err := ReadStatus(dir, "S-0394")
	if err != nil || !ok {
		t.Fatalf("ReadStatus = %v, %v; want the status written", ok, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReadStatus =\n%+v\nwant\n%+v", got, want)
	}
}

func TestReadStatusOfAMissingFileIsNoStatus(t *testing.T) {
	got, ok, err := ReadStatus(t.TempDir(), "S-0394")
	if err != nil || ok || !reflect.DeepEqual(got, Status{}) {
		t.Errorf("ReadStatus = %+v, %v, %v; want no status and no error", got, ok, err)
	}
}

func TestReadStatusOfACorruptFileSaysToDeleteIt(t *testing.T) {
	dir := t.TempDir()
	path := StatusPath(dir, "S-0394")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, ok, err := ReadStatus(dir, "S-0394")
	if err == nil || ok {
		t.Fatalf("ReadStatus = %v, %v; want an error", ok, err)
	}
	for _, want := range []string{"S-0394", path, "delete it"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not say %q", err, want)
		}
	}
}

func TestRecordAFirstReadAppendsTheHeadEachRequiredCheckAndEachReview(t *testing.T) {
	now := time.Date(2026, 10, 17, 10, 0, 0, 0, time.FixedZone("CEST", 2*3600))
	checks := []Check{
		{Name: "lint", State: CheckPassed, Required: true, URL: "https://ci/lint"},
		{Name: "optional", State: CheckFailed, Required: false},
		{Name: "test", State: CheckPending, Required: true, Summary: "has not reported yet"},
	}
	reviews := []Review{{Author: "alex", State: ReviewCommented, URL: "https://github.com/o/r/pull/42#r1", Body: "Looks fine.", SubmittedAt: "2026-10-17T07:50:00Z"}}
	got := Status{Story: "S-0394"}.Record(prAt("aaa"), checks, reviews, now)

	at := "2026-10-17T08:00:00Z"
	want := Status{
		Story: "S-0394", PullRequest: 42, URL: "https://github.com/o/r/pull/42", HeadSHA: "aaa", ReadAt: at,
		Checks: []CheckStatus{
			{Name: "lint", State: CheckPassed, URL: "https://ci/lint", Since: at},
			{Name: "test", State: CheckPending, Summary: "has not reported yet", Since: at},
		},
		Reviews: []ReviewStatus{{Author: "alex", State: ReviewCommented, URL: "https://github.com/o/r/pull/42#r1", Body: "Looks fine.", SubmittedAt: "2026-10-17T07:50:00Z"}},
		Events: []StatusEvent{
			{At: at, Kind: EventHead, To: "aaa", Commit: "aaa"},
			{At: at, Kind: EventCheck, Name: "lint", To: "passed", Commit: "aaa"},
			{At: at, Kind: EventCheck, Name: "test", To: "pending", Commit: "aaa"},
			{At: at, Kind: EventReview, Name: "alex", To: "commented", Commit: "aaa"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Record =\n%+v\nwant\n%+v", got, want)
	}
}

func TestRecordASecondReadAppendsOnlyTheChangesAndCarriesOverWhatItKept(t *testing.T) {
	first := time.Date(2026, 10, 17, 8, 0, 0, 0, time.UTC)
	second := first.Add(5 * time.Minute)
	sent := SendBack{Key: ReviewKey("https://github.com/o/r/pull/42#r0"), Kind: EventReview, Name: "kim", Commit: "000", Thread: "M-0001", At: "2026-10-17T07:00:00Z"}
	old := StatusEvent{At: "2026-10-17T07:00:00Z", Kind: EventReview, Name: "kim", To: "changes-requested", Commit: "000"}
	was := Status{Story: "S-0394", Events: []StatusEvent{old}, SentBack: []SendBack{sent}}
	review := Review{Author: "alex", State: ReviewCommented, URL: "https://github.com/o/r/pull/42#r1"}
	was = was.Record(prAt("aaa"), []Check{
		{Name: "lint", State: CheckPassed, Required: true},
		{Name: "test", State: CheckPending, Required: true},
	}, []Review{review}, first)
	before := len(was.Events)

	got := was.Record(prAt("bbb"), []Check{
		{Name: "lint", State: CheckPassed, Required: true},
		{Name: "optional", State: CheckPending, Required: false},
		{Name: "test", State: CheckFailed, Required: true, URL: "https://ci/test"},
	}, []Review{review, {Author: "kim", State: ReviewChangesRequested, URL: "https://github.com/o/r/pull/42#r2"}}, second)

	at, firstAt := "2026-10-17T08:05:00Z", "2026-10-17T08:00:00Z"
	if got.HeadSHA != "bbb" || got.ReadAt != at {
		t.Errorf("Record head and read at = %s, %s; want bbb, %s", got.HeadSHA, got.ReadAt, at)
	}
	wantChecks := []CheckStatus{
		{Name: "lint", State: CheckPassed, Since: firstAt},
		{Name: "test", State: CheckFailed, URL: "https://ci/test", Since: at},
	}
	if !reflect.DeepEqual(got.Checks, wantChecks) {
		t.Errorf("Record checks =\n%+v\nwant\n%+v", got.Checks, wantChecks)
	}
	if len(got.Reviews) != 2 || got.Reviews[1].Author != "kim" {
		t.Errorf("Record reviews = %+v; want alex's then kim's", got.Reviews)
	}
	if !reflect.DeepEqual(got.Events[:before], was.Events) {
		t.Errorf("Record rewrote the events it kept:\n%+v\nwant\n%+v", got.Events[:before], was.Events)
	}
	if got.Events[0] != old {
		t.Errorf("Record lost the event from before the first read: %+v", got.Events[0])
	}
	wantNew := []StatusEvent{
		{At: at, Kind: EventHead, From: "aaa", To: "bbb", Commit: "bbb"},
		{At: at, Kind: EventCheck, Name: "test", From: "pending", To: "failed", Commit: "bbb"},
		{At: at, Kind: EventReview, Name: "kim", To: "changes-requested", Commit: "bbb"},
	}
	if !reflect.DeepEqual(got.Events[before:], wantNew) {
		t.Errorf("Record appended\n%+v\nwant\n%+v", got.Events[before:], wantNew)
	}
	if !reflect.DeepEqual(got.SentBack, []SendBack{sent}) {
		t.Errorf("Record send backs = %+v; want %+v carried over", got.SentBack, sent)
	}
}

func TestRecordAReviewThatChangesStateAppendsAnEvent(t *testing.T) {
	now := time.Date(2026, 10, 17, 8, 0, 0, 0, time.UTC)
	url := "https://github.com/o/r/pull/42#r1"
	was := Status{Story: "S-0394"}.Record(prAt("aaa"), nil, []Review{{Author: "kim", State: ReviewChangesRequested, URL: url}}, now)
	got := was.Record(prAt("aaa"), nil, []Review{{Author: "kim", State: ReviewDismissed, URL: url}}, now.Add(time.Minute))
	want := StatusEvent{At: "2026-10-17T08:01:00Z", Kind: EventReview, Name: "kim", From: "changes-requested", To: "dismissed", Commit: "aaa"}
	if len(got.Events) != len(was.Events)+1 || got.Events[len(got.Events)-1] != want {
		t.Errorf("Record events = %+v; want %+v appended", got.Events, want)
	}
}

func TestRecordDoesNotShareItsEventsWithTheStatusItWasMadeFrom(t *testing.T) {
	now := time.Date(2026, 10, 17, 8, 0, 0, 0, time.UTC)
	was := Status{Story: "S-0394", Events: make([]StatusEvent, 1, 8)}
	_ = was.Record(prAt("aaa"), nil, nil, now)
	b := was.Record(prAt("bbb"), nil, nil, now)
	if b.Events[1].To != "bbb" {
		t.Errorf("a second Record from the same status overwrote the first's events: %+v", b.Events)
	}
	if len(was.Events) != 1 || was.Events[:2][1] != (StatusEvent{}) {
		t.Errorf("Record wrote into the status it was made from: %+v", was.Events[:2])
	}
}

func TestKeysAndSentFindARecordedSendBack(t *testing.T) {
	if got := CheckKey("test", "aaa"); got != "check:test@aaa" {
		t.Errorf("CheckKey = %q", got)
	}
	if got := ReviewKey("https://github.com/o/r/pull/42#r1"); got != "review:https://github.com/o/r/pull/42#r1" {
		t.Errorf("ReviewKey = %q", got)
	}
	b := SendBack{Key: CheckKey("test", "aaa"), Kind: EventCheck, Name: "test", Commit: "aaa", Thread: "M-0001"}
	st := Status{SentBack: []SendBack{{Key: ReviewKey("x")}, b}}
	if got, ok := st.Sent(CheckKey("test", "aaa")); !ok || got != b {
		t.Errorf("Sent = %+v, %v; want %+v", got, ok, b)
	}
	if _, ok := st.Sent(CheckKey("test", "bbb")); ok {
		t.Error("Sent found a send back for a commit it was not recorded at")
	}
}
