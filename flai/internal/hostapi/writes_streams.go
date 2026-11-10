package hostapi

import (
	"encoding/json"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/channel"
)

// streamSpecs are the writes to narratives, threads, documents, files, and ADRs.
func streamSpecs() map[string]spec {
	return map[string]spec{
		"stream.log": one(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID    string `json:"id"`
				Entry string `json:"entry"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			if strings.TrimSpace(in.Entry) == "" {
				return nil, "", bad("an entry is required")
			}
			return []string{"stream", "log", in.ID, "--", strings.TrimSpace(in.Entry)}, "", nil
		}),

		// stream.answer answers a hand-written open question (S-0090): one
		// with no thread of its own, so thread.reply does not already reach
		// it. It moves to Decisions (design/system/agent-narrative.md).
		"stream.answer": one(func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID       string `json:"id"`
				Question string `json:"question"`
				Answer   string `json:"answer"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			if strings.TrimSpace(in.Question) == "" {
				return nil, "", bad("question is required")
			}
			if strings.TrimSpace(in.Answer) == "" {
				return nil, "", bad("answer is required")
			}
			return []string{"stream", "answer", in.ID, "--by=" + owner(p), "--", strings.TrimSpace(in.Question), strings.TrimSpace(in.Answer)}, "", nil
		}),

		// stream.state writes a narrative's Current state and Next steps, or
		// one of them (S-0271): flai stream state, each text trimmed as flai
		// trims it, lines and all, and a blank one left out as flai leaves it.
		// A text of only "-" is refused, since flai would read it from standard
		// input. A story not in progress or in review, or with no narrative,
		// and a text the markdown lint rejects, are flai's refusals (exit 4).
		"stream.state": {exits: map[int]int{4: Refused}, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID      string `json:"id"`
				Current string `json:"current"`
				Next    string `json:"next"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needStory(in.ID); e != nil {
				return nil, "", e
			}
			var texts []string
			for _, f := range []struct{ flag, text string }{{"current", in.Current}, {"next", in.Next}} {
				t := strings.TrimSpace(f.text)
				switch t {
				case "":
					continue
				case "-":
					return nil, "", bad("%s is the text itself: - is not text to write", f.flag)
				}
				texts = append(texts, "--"+f.flag+"="+t)
			}
			if len(texts) == 0 {
				return nil, "", bad("current, next, or both are required: the text to write under ## Current state and ## Next steps")
			}
			return append([]string{"stream", "state", in.ID}, texts...), "", nil
		}},

		// stream.restore commits main's version of each path a rebased story
		// branch undoes onto the branch (ADR-0155): flai stream restore, as one
		// commit marked as the dashboard's. A story not in progress or in
		// review, a worktree not clean, and a path the story's record does not
		// hold are flai's refusals (exit 4), with what makes them so as data.
		"stream.restore": {exits: map[int]int{4: Refused}, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID    string   `json:"id"`
				Paths []string `json:"paths"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needStory(in.ID); e != nil {
				return nil, "", e
			}
			if len(in.Paths) == 0 {
				return nil, "", bad("paths are required: each path the branch undoes of main's change to restore, as flai accept --dry-run lists them under undoes_main")
			}
			args := []string{"stream", "restore", in.ID}
			for _, given := range in.Paths {
				path, e := undonePath("paths", given)
				if e != nil {
					return nil, "", e
				}
				args = append(args, "--path="+path)
			}
			return append(args, "--trailer="+Trailer), "", nil
		}},

		"thread.new": one(func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				On      string `json:"on"`
				Heading string `json:"heading"`
				Title   string `json:"title"`
				Text    string `json:"text"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.On == "" || text(in.Title) == "" || strings.TrimSpace(in.Text) == "" {
				return nil, "", bad("on, title, and text are required")
			}
			if !anyItemID.MatchString(in.On) {
				if _, e := docRel(p, in.On); e != nil {
					return nil, "", bad("a thread is opened on an item or on a document of the project")
				}
			}
			args := []string{"thread", "new", "--on=" + in.On, "--by=" + owner(p)}
			if in.Heading != "" {
				args = append(args, "--heading="+text(in.Heading))
			}
			return append(args, "--", text(in.Title), strings.TrimSpace(in.Text)), "", nil
		}),

		// thread.reply takes a recommendation's marks (ADR-0090): recommend,
		// and source, <path> or <path>#<heading>, which flai checks exist.
		"thread.reply": one(func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID        string `json:"id"`
				Text      string `json:"text"`
				Recommend bool   `json:"recommend"`
				Source    string `json:"source"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if !threadID.MatchString(in.ID) {
				return nil, "", bad("%q is not a thread", in.ID)
			}
			if strings.TrimSpace(in.Text) == "" {
				return nil, "", bad("text is required")
			}
			args := []string{"thread", "reply", in.ID, "--by=" + owner(p)}
			if in.Recommend {
				args = append(args, "--recommend")
			}
			if text(in.Source) != "" {
				args = append(args, "--source="+text(in.Source))
			}
			return append(args, "--", strings.TrimSpace(in.Text)), "", nil
		}),

		// thread.confirm makes a thread's pending recommendation the answer,
		// as the owner (ADR-0090).
		"thread.confirm": one(func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if !threadID.MatchString(in.ID) {
				return nil, "", bad("%q is not a thread", in.ID)
			}
			return []string{"thread", "confirm", in.ID, "--by=" + owner(p)}, "", nil
		}),

		"thread.resolve": one(func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID     string `json:"id"`
				Reason string `json:"reason"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if !threadID.MatchString(in.ID) {
				return nil, "", bad("%q is not a thread", in.ID)
			}
			args := []string{"thread", "resolve", in.ID, "--by=" + owner(p)}
			if text(in.Reason) != "" {
				args = append(args, "--reason="+text(in.Reason))
			}
			return args, "", nil
		}),

		// doc.show and doc.save read and save in a story's worktree when
		// story names one (S-0378): they run there, and flai doc refuses a
		// path under wip, which lives on the main checkout.
		"doc.show": {reads: true, dir: docDir, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Path  string `json:"path"`
				Story string `json:"story"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			rel, e := docRel(p, in.Path)
			if e != nil {
				return nil, "", e
			}
			args := []string{"doc", "show"}
			if in.Story != "" {
				if e := needStory(in.Story); e != nil {
					return nil, "", e
				}
				args = append(args, "--story="+in.Story)
			}
			return append(args, "--", rel), "", nil
		}},

		"doc.save": {exits: map[int]int{3: Conflict, 4: Refused}, dir: docDir, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Path    string  `json:"path"`
				Content *string `json:"content"`
				Hash    string  `json:"hash"`
				Message string  `json:"message"`
				Story   string  `json:"story"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.Path == "" || in.Content == nil || in.Hash == "" {
				return nil, "", bad("path, content, and hash are required")
			}
			rel, e := docRel(p, in.Path)
			if e != nil {
				return nil, "", e
			}
			if !docHash.MatchString(in.Hash) {
				return nil, "", bad("hash is the one the document was loaded with")
			}
			args := []string{"doc", "save", "--hash=" + in.Hash, "--trailer=" + Trailer}
			if text(in.Message) != "" {
				args = append(args, "--message="+text(in.Message))
			}
			if in.Story != "" {
				if e := needStory(in.Story); e != nil {
					return nil, "", e
				}
				args = append(args, "--story="+in.Story)
			}
			return append(args, "--", rel), *in.Content, nil
		}},

		// file.get and file.save read and save a text file that is not
		// markdown, such as the design token file, as doc.show and doc.save
		// do a document (S-0436, ADR-0147): in a story's worktree when story
		// names one, and a save of the token file writes and commits the
		// design stylesheet with it (ADR-0150).
		"file.get": {reads: true, dir: docDir, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Path  string `json:"path"`
				Story string `json:"story"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			rel, e := fileRel(p, in.Path)
			if e != nil {
				return nil, "", e
			}
			args := []string{"file", "show"}
			if in.Story != "" {
				if e := needStory(in.Story); e != nil {
					return nil, "", e
				}
				args = append(args, "--story="+in.Story)
			}
			return append(args, "--", rel), "", nil
		}},

		"file.save": {exits: map[int]int{3: Conflict, 4: Refused}, dir: docDir, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Path    string  `json:"path"`
				Content *string `json:"content"`
				Hash    string  `json:"hash"`
				Message string  `json:"message"`
				Story   string  `json:"story"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.Path == "" || in.Content == nil {
				return nil, "", bad("path and content are required")
			}
			rel, e := fileRel(p, in.Path)
			if e != nil {
				return nil, "", e
			}
			// An empty hash creates a file that is absent (ADR-0164).
			args := []string{"file", "save"}
			if in.Hash != "" {
				if !docHash.MatchString(in.Hash) {
					return nil, "", bad("hash is the one file.get gave for the content the file was loaded with, or empty to create a file that is absent")
				}
				args = append(args, "--hash="+in.Hash)
			}
			args = append(args, "--trailer="+Trailer)
			if text(in.Message) != "" {
				args = append(args, "--message="+text(in.Message))
			}
			if in.Story != "" {
				if e := needStory(in.Story); e != nil {
					return nil, "", e
				}
				args = append(args, "--story="+in.Story)
			}
			return append(args, "--", rel), *in.Content, nil
		}},

		"adr.new": {exits: map[int]int{4: Refused}, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Title      string   `json:"title"`
				Status     string   `json:"status"`
				Supersedes []string `json:"supersedes"`
				Refines    []string `json:"refines"`
				Body       string   `json:"body"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			title := text(in.Title)
			if title == "" {
				return nil, "", bad("a title is required: the decision, as a sentence")
			}
			if len([]rune(title)) > 200 {
				return nil, "", bad("the title is longer than 200 characters")
			}
			status := in.Status
			if status == "" {
				status = "proposed"
			}
			if status != "proposed" && status != "accepted" {
				return nil, "", bad("status must be proposed or accepted")
			}
			if strings.TrimSpace(in.Body) == "" {
				return nil, "", bad("write the decision: the body is empty")
			}
			args := []string{"adr", "new", "--status=" + status}
			for _, rel := range []struct {
				flag string
				list []string
			}{{"supersedes", in.Supersedes}, {"refines", in.Refines}} {
				for _, v := range rel.list {
					n, e := adrNo(v)
					if e != nil {
						return nil, "", e
					}
					args = append(args, "--"+rel.flag+"="+n)
				}
			}
			return append(args, "--body-stdin", "--autocommit", "--trailer="+Trailer, "--", title), in.Body, nil
		}},

		"adr.template": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			if _, e := decode[struct{}](raw); e != nil {
				return nil, "", e
			}
			return []string{"adr", "new", "--print-body"}, "", nil
		}),

		"adr.accept": {exits: map[int]int{4: Refused}, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			n, e := adrNo(in.ID)
			if e != nil {
				return nil, "", e
			}
			return []string{"adr", "accept", n, "--autocommit", "--trailer=" + Trailer}, "", nil
		}},
	}
}
