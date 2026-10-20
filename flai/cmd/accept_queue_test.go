package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// enqueueWithBatch writes a queue state with S-0001 and S-0005 enqueued and
// in the batch under verification.
func enqueueWithBatch(t *testing.T, root string) {
	t.Helper()
	at := time.Date(2026, 9, 15, 20, 0, 0, 0, time.UTC)
	writeQueueState(t, root, mergequeue.State{
		Enqueued: []mergequeue.Enqueued{{Story: "S-0001", By: "orchestrator", At: at}, {Story: "S-0005", By: "orchestrator", At: at}},
		Batch:    &mergequeue.Batch{Stories: []string{"S-0001", "S-0005"}, Step: "verify"},
	})
}

// assertLeftTheQueue fails unless S-0001 is out of the queue, the batch that
// held it is dropped as accepted by hand by alex, and its outcome says so.
func assertLeftTheQueue(t *testing.T, root string) {
	t.Helper()
	st := readQueueState(t, root)
	if len(st.Enqueued) != 1 || st.Enqueued[0].Story != "S-0005" {
		t.Errorf("enqueued %+v, want S-0005 only", st.Enqueued)
	}
	if st.Batch == nil || st.Batch.Dropped == nil || st.Batch.Dropped.Reason != "S-0001 was accepted by hand by alex" || st.Batch.Dropped.By != "alex" {
		t.Errorf("batch %+v, want it dropped", st.Batch)
	}
	if len(st.Outcomes) != 1 || st.Outcomes[0].Story != "S-0001" || st.Outcomes[0].Outcome != mergequeue.OutcomeAcceptedByHand || st.Outcomes[0].Detail != "accepted by hand by alex" {
		t.Errorf("outcomes %+v", st.Outcomes)
	}
}

// S-0399 criterion 1, ADR-0173: flai accept by hand takes the story out of
// the queue, drops the batch that holds it, and records its outcome.
func TestAcceptByHandTakesTheStoryOutOfTheQueue(t *testing.T) {
	root, _ := orchestratedStoryInReview(t, false)
	enqueueWithBatch(t, root)
	if out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex"); code != 0 {
		t.Fatalf("flai accept: %d %s\n%s", code, errOut, out)
	}
	assertLeftTheQueue(t, root)
}

// S-0399 criterion 1: under pull-request acceptance too, the story merged
// through the remote leaves the queue.
func TestAcceptThroughTheRemoteTakesTheStoryOutOfTheQueue(t *testing.T) {
	root, _, _, _ := remoteProject(t, "")
	enqueueWithBatch(t, root)
	acceptJSON(t, root)
	assertLeftTheQueue(t, root)
}

// A queue state that cannot be written after the merge is warned of, and the
// acceptance is not failed for it.
func TestAcceptWarnsOfAQueueStateItCannotWrite(t *testing.T) {
	root, _ := orchestratedStoryInReview(t, false)
	path := mergequeue.StatePath(filepath.Join(root, ".flai-cache"))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"enqueued": [`), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex")
	if code != 0 || !strings.Contains(out, "accepted S-0001: done") {
		t.Fatalf("flai accept: %d %s\n%s", code, errOut, out)
	}
	var warned bool
	for _, line := range strings.Split(errOut, "\n") {
		if strings.Contains(line, path) && strings.Contains(line, `"level":"WARN"`) && strings.Contains(line, `"component":"accept"`) && strings.Contains(line, `"item":"S-0001"`) {
			warned = true
		}
	}
	if !warned {
		t.Errorf("no warning naming the queue's state:\n%s", errOut)
	}
}

// landingBatch is the ID of the merge queue batch landingState writes.
const landingBatch = "20260915-200000"

// landedAt is when the landings and the acceptances they are compared with
// run, two days after the story moved to review.
var landedAt = time.Date(2026, 9, 17, 9, 0, 0, 0, time.UTC)

// landingState writes a queue state with S-0001, enqueued by by, in merge
// queue batch landingBatch, which passed on main's head and is landing it
// next: its heads are head, the branch's, and a copy of it rebased onto
// main, which it returns. edit, when not nil, changes the batch first.
func landingState(t *testing.T, root, head, by string, edit func(*mergequeue.Batch)) string {
	t.Helper()
	main := gitIn(t, root, "rev-parse", "main")
	rebased := gitIn(t, root, "commit-tree", "story/S-0001^{tree}", "-p", main, "-m", "feat: [S-0001] ship it")
	at := time.Date(2026, 9, 15, 20, 0, 0, 0, time.UTC)
	b := &mergequeue.Batch{
		ID: landingBatch, Branch: "merge-queue/" + landingBatch, Base: main, Stories: []string{"S-0001"}, StartedAt: at,
		Step:    mergequeue.StepLanding + " S-0001",
		Result:  &mergequeue.Result{Passed: true, FinishedAt: at},
		Heads:   []mergequeue.Head{{Story: "S-0001", Verified: head, Rebased: rebased}},
		Landing: &mergequeue.Landing{Stories: []string{"S-0001"}, Landed: []string{}, Head: main},
	}
	if edit != nil {
		edit(b)
	}
	writeQueueState(t, root, mergequeue.State{Enqueued: []mergequeue.Enqueued{{Story: "S-0001", By: by, At: at}}, Batch: b})
	return rebased
}

// acceptedS0001 is S-0001 as the project at root holds it, archived or not.
func acceptedS0001(t *testing.T, root string) *workitem.Item {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	it, err := repo.Get("S-0001")
	if err != nil {
		t.Fatal(err)
	}
	return it
}

// assertLanded fails unless S-0001 is out of the queue, landed by
// landingBatch as by, and the batch is kept, not dropped, with S-0001
// landed in its landing.
func assertLanded(t *testing.T, root, by string) {
	t.Helper()
	st := readQueueState(t, root)
	if len(st.Enqueued) != 0 {
		t.Errorf("enqueued %+v, want none: a landed story leaves the queue", st.Enqueued)
	}
	if st.Batch == nil || st.Batch.ID != landingBatch || st.Batch.Dropped != nil || st.Batch.Landing == nil || len(st.Batch.Landing.Landed) != 1 || st.Batch.Landing.Landed[0] != "S-0001" {
		t.Errorf("batch %+v, want %s kept with S-0001 landed", st.Batch, landingBatch)
	}
	detail := "landed by merge queue batch " + landingBatch + " as " + by
	if o := st.Outcomes; len(o) != 1 || o[0].Story != "S-0001" || o[0].Outcome != mergequeue.OutcomeLanded || !o[0].At.Equal(landedAt) || o[0].Detail != detail {
		t.Errorf("outcomes %+v, want S-0001 %s at %s, %q", o, mergequeue.OutcomeLanded, landedAt, detail)
	}
}

// S-0400 criterion 1, ADR-0176: a story landed with flai accept --batch reads
// as one accepted by hand: the same done transition, by whoever enqueued it,
// the same acceptance commit changing the same files, the same archive, and
// the same flai stats figures. The queue records it landed and keeps the
// batch, where an acceptance by hand drops it.
func TestALandingIsAnAcceptanceByHand(t *testing.T) {
	byHand, _ := orchestratedStoryInReview(t, false)
	if out, errOut, code := runInAt(t, byHand, landedAt, "accept", "S-0001", "--by", "alex"); code != 0 {
		t.Fatalf("flai accept by hand: %d %s\n%s", code, errOut, out)
	}
	landed, head := orchestratedStoryInReview(t, false)
	rebased := landingState(t, landed, head, "alex", nil)
	out, errOut, code := runInAt(t, landed, landedAt, "accept", "S-0001", "--by", "alex", "--batch", landingBatch, "--json")
	if code != 0 {
		t.Fatalf("flai accept --batch: %d %s\n%s", code, errOut, out)
	}
	var res accepted
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if res.Batch != landingBatch || res.Verified != rebased || res.By != "alex" || !res.Merged || res.Status != workitem.Done {
		t.Errorf("--json must carry the batch and its rebased head as the commit verified: %+v", res.Acceptance)
	}

	hand, land := acceptedS0001(t, byHand), acceptedS0001(t, landed)
	if !hand.Archived || !land.Archived || mustRel(t, byHand, hand.Path) != mustRel(t, landed, land.Path) || !strings.HasPrefix(filepath.ToSlash(mustRel(t, landed, land.Path)), "wip/archive/") {
		t.Errorf("both are archived to the same place: %s and %s", hand.Path, land.Path)
	}
	last := land.Transitions[len(land.Transitions)-1]
	if last.To != workitem.Done || last.By != "alex" || last.At != landedAt.Format(time.RFC3339) {
		t.Errorf("the done transition is alex's, at the landing: %+v", last)
	}
	if h, l := hand.Transitions, land.Transitions; !slices.Equal(h, l) {
		t.Errorf("the transitions differ:\nby hand %+v\nlanded  %+v", h, l)
	}
	for _, args := range [][]string{{"log", "-1", "--format=%s"}, {"show", "--name-status", "--format=", "HEAD"}, {"show", "main:cli/feature.go"}} {
		if h, l := gitIn(t, byHand, args...), gitIn(t, landed, args...); h != l {
			t.Errorf("git %v differs:\nby hand %s\nlanded  %s", args, h, l)
		}
	}
	if subject := gitIn(t, landed, "log", "-1", "--format=%s"); subject != "chore: [S-0001] accept and archive, with E-0001" {
		t.Errorf("the acceptance commit's subject is %q", subject)
	}
	handStats, _, hc := runInAt(t, byHand, landedAt, "stats", "--json")
	landStats, _, lc := runInAt(t, landed, landedAt, "stats", "--json")
	var figures struct {
		Summary struct {
			Completed int `json:"completed"`
		} `json:"summary"`
	}
	if hc != 0 || lc != 0 || json.Unmarshal([]byte(landStats), &figures) != nil || figures.Summary.Completed != 1 {
		t.Fatalf("flai stats --json: %d %d\n%s", hc, lc, landStats)
	}
	// The landing writes a landed line to the queue's event log (ADR-0177),
	// so only the landed project reports a queue; the rest is the same.
	var handFig, landFig map[string]any
	if json.Unmarshal([]byte(handStats), &handFig) != nil || json.Unmarshal([]byte(landStats), &landFig) != nil {
		t.Fatalf("flai stats --json:\nby hand %s\nlanded  %s", handStats, landStats)
	}
	var queue struct {
		Weeks []struct {
			Landed int `json:"landed"`
		} `json:"weeks"`
	}
	raw, _ := json.Marshal(landFig["queue"])
	if _, ok := handFig["queue"]; ok || json.Unmarshal(raw, &queue) != nil || len(queue.Weeks) == 0 || queue.Weeks[len(queue.Weeks)-1].Landed != 1 {
		t.Errorf("only the landing reports a queue, with the story landed this week:\nby hand %s\nlanded  %s", handStats, landStats)
	}
	delete(landFig, "queue")
	if h, l := string(mustJSON(t, handFig)), string(mustJSON(t, landFig)); h != l {
		t.Errorf("flai stats --json differs:\nby hand %s\nlanded  %s", h, l)
	}
	assertLanded(t, landed, "alex")
}

// S-0400 criteria 1 and 2, ADR-0176: for a story the orchestrator enqueued,
// the landing is its acceptance, taking the batch's verification in place of
// --verified and needing no evidence; its Notes name the story's rebased head
// in the batch and the batch, with the evidence when one is given.
func TestTheOrchestratorsLandingTakesTheBatchsVerification(t *testing.T) {
	evidence := "Verdict: pass, tests and lint clean\n- 1: `cli/feature.go`"
	for _, ev := range []string{"", evidence} {
		root, head := orchestratedStoryInReview(t, true)
		rebased := landingState(t, root, head, "orchestrator", nil)
		args := []string{"accept", "S-0001", "--by", "orchestrator", "--batch", landingBatch}
		if ev != "" {
			args = append(args, "--evidence", evidenceFile(t, ev+"\n"))
		}
		out, errOut, code := runInAt(t, root, landedAt, args...)
		if code != 0 {
			t.Fatalf("the orchestrator's landing, evidence %q: %d %s\n%s", ev, code, errOut, out)
		}
		if want := "landed from merge queue batch " + landingBatch + ", verified at " + rebased[:12] + "\n"; !strings.Contains(out, want) {
			t.Errorf("the output must name the batch and its verified commit, %q:\n%s", want, out)
		}
		story := acceptedS0001(t, root)
		section := "## Notes\n\n### Accepted by the orchestrator\n\n- Verified: " + rebased + "\n- Batch: " + landingBatch + "\n- At: 2026-09-17T09:00:00Z\n"
		if ev != "" {
			section += "\n" + ev + "\n"
		}
		if !story.Archived || !strings.HasSuffix(story.Body, section) {
			t.Errorf("the Notes must end with the rebased head, the batch, and the evidence given:\n%s\nwant\n%s", story.Body, section)
		}
		if last := story.Transitions[len(story.Transitions)-1]; last.To != workitem.Done || last.By != "orchestrator" {
			t.Errorf("the done transition is the orchestrator's: %+v", last)
		}
		assertLanded(t, root, "orchestrator")
	}
}

// S-0400 criterion 2, ADR-0176: flai accept --batch is refused, changing
// nothing, unless the queue's batch is the one named, is not dropped, and
// lands the story next, the story's branch is at the head the batch rebased,
// and main is where the landing left it; and --verified is not given with it.
func TestALandingIsRefusedUnlessTheBatchVerifiedWhatLands(t *testing.T) {
	for name, c := range map[string]struct {
		edit  func(*mergequeue.Batch)
		move  func(t *testing.T, root string)
		extra []string
		want  string
	}{
		"another batch": {extra: []string{"--batch", "20260915-190000"}, want: "the merge queue's batch is " + landingBatch + ", not 20260915-190000, so S-0001 cannot land as part of 20260915-190000"},
		"dropped": {edit: func(b *mergequeue.Batch) {
			b.Dropped = &mergequeue.Dropped{Reason: "S-0005 was removed from the queue by alex", By: "alex", At: landedAt}
		}, want: "merge queue batch " + landingBatch + " was dropped by alex: S-0005 was removed from the queue by alex"},
		"not next": {edit: func(b *mergequeue.Batch) {
			b.Stories, b.Landing.Stories = []string{"S-0005", "S-0001"}, []string{"S-0005", "S-0001"}
		}, want: "S-0001 is not the next story merge queue batch " + landingBatch + " lands: S-0005 is"},
		"branch moved": {move: func(t *testing.T, root string) {
			wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
			if err := os.WriteFile(filepath.Join(wt, "cli", "more.go"), []byte("package cli\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			gitIn(t, wt, "add", "-A")
			gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] more")
		}, want: "the branch moved since, so S-0001 lands with a batch built from its new head"},
		"main moved": {move: func(t *testing.T, root string) {
			if err := os.WriteFile(filepath.Join(root, "design", "note.md"), []byte("# Note\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			gitIn(t, root, "add", "-A")
			gitIn(t, root, "commit", "-q", "-m", "docs: a note")
		}, want: "main moved other than by the batch's landing, so flai serve drops the batch and builds it again from the new main"},
		"--verified beside it": {extra: []string{"--by", "orchestrator", "--verified", "HEAD"}, want: "rule: --batch " + landingBatch + " takes the merge queue batch's verification of S-0001's rebased head in place of --verified's: give one or the other (ADR-0176)"},
	} {
		t.Run(name, func(t *testing.T) {
			root, head := orchestratedStoryInReview(t, true)
			landingState(t, root, head, "alex", c.edit)
			if c.move != nil {
				c.move(t, root)
			}
			before := readQueueState(t, root)
			mainBefore := gitIn(t, root, "rev-parse", "main")
			args := append([]string{"accept", "S-0001", "--by", "alex", "--batch", landingBatch}, c.extra...)
			_, errOut, code := runInAt(t, root, landedAt, args...)
			if code == 0 || !strings.Contains(errOut, c.want) {
				t.Errorf("flai %v must be refused with %q: %d\n%s", args, c.want, code, errOut)
			}
			assertNothingAccepted(t, root, mainBefore)
			if after := readQueueState(t, root); after.Batch == nil || len(after.Batch.Landing.Landed) != 0 || len(after.Enqueued) != len(before.Enqueued) || len(after.Outcomes) != 0 {
				t.Errorf("a refused landing must leave the queue as it was: %+v", after)
			}
		})
	}
}

// S-0400 criterion 2, ADR-0173: under pull-request acceptance there is no
// merge queue, so --batch is refused before anything changes.
func TestALandingIsRefusedUnderPullRequestAcceptance(t *testing.T) {
	root, _ := forgeGatedStoryInReview(t, "pull-request")
	mainBefore := gitIn(t, root, "rev-parse", "main")
	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "alex", "--batch", landingBatch)
	if want := "rule: --batch " + landingBatch + " is the merge queue's landing, and " + mergequeue.OffUnderPullRequest; code == 0 || !strings.Contains(errOut, want) {
		t.Errorf("--batch under pull-request acceptance must be refused with %q: %d\n%s", want, code, errOut)
	}
	assertNothingAccepted(t, root, mainBefore)
}
