package metrics

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// coordNow is a Wednesday of ISO week 38; three weeks before it, the
// window's start, is a Wednesday of week 35.
var coordNow = time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC)

// coordStart is the start of a 21-day window ending at coordNow.
var coordStart = coordNow.Add(-21 * 24 * time.Hour)

// conv is a conversation between from and to with the status, shares, and
// entries given, each entry a heading's time, author, and story, "" for a
// closing entry, then its text.
func conv(id, from, to, status string, shares []workitem.Share, entries ...string) *messages.Conversation {
	var b strings.Builder
	for i := 0; i+3 < len(entries); i += 4 {
		head := entries[i] + " " + entries[i+1]
		if entries[i+2] != "" {
			head += " " + entries[i+2]
		}
		b.WriteString("\n### " + head + "\n\n" + entries[i+3] + "\n")
	}
	return &messages.Conversation{ID: id, From: from, To: to, Status: status, Shares: shares, Body: b.String()}
}

// shares is one share by holder with held, made at at.
func shares(holder, held, at string) []workitem.Share {
	return []workitem.Share{{Holder: holder, Held: held, Paths: []string{"flai"}, Split: "split", By: "agent", At: at}}
}

// coordItems are the stories the coordination fixtures name:
//   - S-0001 holds S-0002, which starts on its share on 31 August, until
//     S-0001 moves to review on 2 September;
//   - S-0003 holds S-0004, which starts on its share on 9 September, and is
//     still in progress, so the hold runs to now;
//   - S-0005 shared with S-0006 but was in review when S-0006 started, and
//     S-0008 was shared with only after it started: neither starts on a share.
func coordItems() []*workitem.Item {
	s, rd, ip, rv, done := workitem.Story, workitem.Ready, workitem.InProgress, workitem.Review, workitem.Done
	const day = "2026-08-20T00:00:00Z"
	return []*workitem.Item{
		claimItem("S-0001", s, day, nil, ip, "2026-08-27T00:00:00Z", rv, "2026-09-02T00:00:00Z", done, "2026-09-03T00:00:00Z"),
		claimItem("S-0002", s, day, nil, rd, "2026-08-27T00:00:00Z", ip, "2026-08-31T00:00:00Z"),
		claimItem("S-0003", s, day, nil, ip, "2026-09-07T00:00:00Z"),
		claimItem("S-0004", s, day, nil, rd, "2026-09-05T00:00:00Z", ip, "2026-09-09T00:00:00Z"),
		claimItem("S-0005", s, day, nil, ip, "2026-09-01T00:00:00Z", rv, "2026-09-05T00:00:00Z", done, "2026-09-13T00:00:00Z"),
		claimItem("S-0006", s, day, nil, rd, "2026-09-02T00:00:00Z", ip, "2026-09-10T00:00:00Z"),
		claimItem("S-0007", s, day, nil, rd, "2026-09-01T00:00:00Z"),
		claimItem("S-0008", s, day, nil, ip, "2026-09-12T00:00:00Z"),
	}
}

// coordConversations are the conversations of the fixtures, the week of
// each event noted above it.
func coordConversations() []*messages.Conversation {
	open, closed := messages.StatusOpen, messages.StatusClosed
	return []*messages.Conversation{
		// opened and agreed by a share in week 35, closed with S-0001's
		// acceptance in week 36
		conv("MS-0001", "S-0001", "S-0002", closed, shares("S-0001", "S-0002", "2026-08-30T00:00:00Z"),
			"2026-08-29T10:00:00Z", "flai", "S-0002", "S-0001 holds S-0002 on overlap.",
			"2026-09-03T00:00:00Z", "alex", "", "Closed: S-0001 was accepted"),
		// opened and escalated in week 36, escalated again and agreed by a
		// share in week 37: each counted once
		conv("MS-0002", "S-0003", "S-0004", open, shares("S-0003", "S-0004", "2026-09-08T00:00:00Z"),
			"2026-09-01T09:00:00Z", "flai", "S-0003", "story/S-0003 and story/S-0004 conflict when merged.",
			"2026-09-02T09:00:00Z", "agent", "S-0004", "Asked the operator on TH-0001, `wip/threads/TH-0001.md`, since we do not agree.",
			"2026-09-09T09:00:00Z", "agent", "S-0003", "Asked the operator on TH-0002, `wip/threads/TH-0002.md`, since we do not agree."),
		// opened in week 37, agreed by a clean merge in week 38; S-0005 is
		// done, but a closing entry for another reason is not closed with it
		conv("MS-0003", "S-0005", "S-0006", closed, nil,
			"2026-09-10T08:00:00Z", "flai", "S-0005", "story/S-0005 and story/S-0006 conflict when merged.",
			"2026-09-14T08:00:00Z", "flai", "", "Closed: story/S-0005 and story/S-0006 merge cleanly at the sync of S-0006"),
		// opened in week 36 and stored open with no closing entry, but S-0001
		// is done: closed with it at its done transition, in week 36
		conv("MS-0004", "S-0007", "S-0001", open, nil,
			"2026-09-01T00:00:00Z", "agent", "S-0007", "May I?"),
		// opened before the window, closed with two stories archived in week 38
		conv("MS-0005", "S-0001", "S-0005", closed, nil,
			"2026-08-20T00:00:00Z", "agent", "S-0001", "Before the window.",
			"2026-09-15T00:00:00Z", "alex", "", "Closed: S-0001 and S-0005 were archived"),
		// opened and closed in week 37, its other story no longer open
		conv("MS-0006", "S-0003", "S-0007", closed, nil,
			"2026-09-11T00:00:00Z", "flai", "S-0003", "story/S-0003 and story/S-0007 conflict when merged.",
			"2026-09-12T00:00:00Z", "flai", "", "Closed: S-0007 is ready, no longer open, at the sync of S-0003"),
		// opened in week 37 and closed for another reason: neither agreed
		// nor closed with its story
		conv("MS-0007", "S-0003", "S-0007", closed, nil,
			"2026-09-11T01:00:00Z", "agent", "S-0003", "Shall we?",
			"2026-09-12T00:00:00Z", "agent", "", "Closed: we agreed in a reply"),
		// opened and agreed by a share in week 37, made after S-0008 started
		conv("MS-0008", "S-0003", "S-0008", open, shares("S-0003", "S-0008", "2026-09-13T00:00:00Z"),
			"2026-09-12T06:00:00Z", "flai", "S-0008", "S-0003 holds S-0008 on overlap."),
		// opened and agreed by a share in week 36; S-0005 is done in week 37,
		// so it reads as closed with it then
		conv("MS-0009", "S-0005", "S-0006", open, shares("S-0005", "S-0006", "2026-09-04T00:00:00Z"),
			"2026-09-03T12:00:00Z", "flai", "S-0006", "S-0005 holds S-0006 on overlap."),
	}
}

// coordConflicts are the conflict log's lines of the fixtures: one before the
// window, one after now, and one of no known kind, none counted.
func coordConflicts() []storygit.Conflict {
	at := parseTime
	return []storygit.Conflict{
		{At: at("2026-08-20T00:00:00Z"), Story: "S-0001", Other: "S-0002", Where: storygit.WhereTrial, Paths: []string{"a"}},
		{At: at("2026-08-28T00:00:00Z"), Story: "S-0001", Other: "S-0002", Where: storygit.WhereTrial, Paths: []string{"a", "b"}},
		{At: at("2026-09-01T00:00:00Z"), Story: "S-0003", Other: "S-0004", Where: storygit.WhereTrial, Paths: []string{"a"}},
		{At: at("2026-09-01T00:00:00Z"), Story: "S-0003", Where: "elsewhere", Paths: []string{"a"}},
		{At: at("2026-09-03T00:00:00Z"), Story: "S-0002", Where: storygit.WhereSync, Paths: []string{"a", "b", "c"}},
		{At: at("2026-09-15T00:00:00Z"), Story: "S-0004", Where: storygit.WhereAccept, Paths: []string{"a"}},
		{At: at("2026-09-17T00:00:00Z"), Story: "S-0004", Where: storygit.WhereAccept, Paths: []string{"a"}},
	}
}

// S-0337: each value of coordination, per ISO week of the window, from the
// conversations, the conflict log, and the transitions.
func TestCoordinationCountsEachValuePerWeek(t *testing.T) {
	got := coordination(coordItems(), coordConversations(), coordConflicts(), true, coordStart, coordNow)
	want := Coordination{Logged: true, Weeks: []CoordinationWeek{
		{
			Week: "2026-W35", Start: "2026-08-24",
			Conversations: ConversationCounts{Opened: 1, Agreed: 1},
			Conflicts:     &ConflictCounts{Found: 1, FoundPaths: 2},
		},
		{
			Week: "2026-W36", Start: "2026-08-31",
			Conversations: ConversationCounts{Opened: 3, Agreed: 1, Escalated: 1, ClosedWithStory: 2},
			Conflicts:     &ConflictCounts{Found: 1, FoundPaths: 1, MetSync: 1, MetSyncPaths: 3},
			// S-0002 held until S-0001 moved to review, two days on
			Shares: ShareCounts{Started: 1, HeldSeconds: 48 * 3600},
		},
		{
			Week: "2026-W37", Start: "2026-09-07",
			Conversations: ConversationCounts{Opened: 4, Agreed: 2, ClosedWithStory: 2},
			Conflicts:     &ConflictCounts{},
			// S-0004 held to now, its holder still in progress
			Shares: ShareCounts{Started: 1, HeldSeconds: (7*24 + 12) * 3600},
		},
		{
			Week: "2026-W38", Start: "2026-09-14",
			Conversations: ConversationCounts{Agreed: 1, ClosedWithStory: 1},
			Conflicts:     &ConflictCounts{MetAccept: 1, MetAcceptPaths: 1},
		},
	}}
	if !reflect.DeepEqual(got, want) {
		gotJSON, _ := json.Marshal(got)
		wantJSON, _ := json.Marshal(want)
		t.Errorf("coordination =\n%s\nwant\n%s", gotJSON, wantJSON)
	}
}

// S-0337: the JSON names are those of metrics.md § Coordination, which the
// dashboard charts.
func TestCoordinationJSONNames(t *testing.T) {
	c := coordination(nil, nil, nil, true, coordNow.Add(-time.Hour), coordNow)
	got, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"logged":true,"weeks":[{"week":"2026-W38","start":"2026-09-14",` +
		`"conversations":{"opened":0,"agreed":0,"escalated":0,"closed_with_story":0},` +
		`"conflicts":{"found":0,"found_paths":0,"met_sync":0,"met_sync_paths":0,"met_accept":0,"met_accept_paths":0},` +
		`"shares":{"started":0,"held_seconds":0}}]}`
	if string(got) != want {
		t.Errorf("coordination JSON =\n%s\nwant\n%s", got, want)
	}
}

// S-0337: a window with nothing in it, and a missing conflict log, which
// means no conflicts: every week is present with every count 0.
func TestCoordinationOfAnEmptyWindowAndAMissingLog(t *testing.T) {
	got := coordination(coordItems(), nil, nil, true, coordStart, coordNow)
	if !got.Logged || len(got.Weeks) != 4 || got.Weeks[0].Week != "2026-W35" || got.Weeks[3].Week != "2026-W38" {
		t.Fatalf("weeks of an empty window: %+v", got)
	}
	for _, w := range got.Weeks {
		if w.Conversations != (ConversationCounts{}) || w.Shares != (ShareCounts{}) || w.Conflicts == nil || *w.Conflicts != (ConflictCounts{}) {
			t.Errorf("week %s of an empty window: %+v", w.Week, w)
		}
	}
}

// S-0337: a conflict log that could not be read leaves logged false and the
// weeks without conflicts, and the rest is still counted.
func TestCoordinationWithAnUnreadableLog(t *testing.T) {
	got := coordination(coordItems(), coordConversations(), coordConflicts(), false, coordStart, coordNow)
	if got.Logged {
		t.Error("logged is true for a log that could not be read")
	}
	for _, w := range got.Weeks {
		if w.Conflicts != nil {
			t.Errorf("week %s carries conflicts: %+v", w.Week, *w.Conflicts)
		}
	}
	if got.Weeks[1].Conversations.Opened != 3 || got.Weeks[1].Shares.Started != 1 {
		t.Errorf("week 36 without the log: %+v", got.Weeks[1])
	}
	raw, err := json.Marshal(got.Weeks[0])
	if err != nil || strings.Contains(string(raw), `"conflicts"`) {
		t.Errorf("a week without the log: %v %s", err, raw)
	}
}

// S-0337: Compute reports coordination from the options' conversations and
// conflict log, over every story whatever the report's type.
func TestComputeReportsCoordination(t *testing.T) {
	rep := Compute(coordItems(), Options{
		Now: coordNow, Since: 21 * 24 * time.Hour, Type: workitem.Task,
		Conversations: coordConversations(), Conflicts: coordConflicts(),
	})
	c := rep.Coordination
	if !c.Logged || len(c.Weeks) != 4 || c.Weeks[1].Shares.Started != 1 || c.Weeks[1].Conflicts == nil || c.Weeks[1].Conflicts.MetSync != 1 {
		t.Errorf("rep.Coordination = %+v", c)
	}
}
