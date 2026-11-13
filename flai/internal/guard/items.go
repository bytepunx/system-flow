package guard

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// itemFile is the ID a work item file's path gives, and whether the path is
// one: <ID>-<slug>.md in a folder epics, stories, or tasks under a folder
// kanban, the archive's included, with an ID whose prefix is the folder's
// type, in any checkout.
func itemFile(path string) (string, bool) {
	parts := strings.Split(filepath.ToSlash(filepath.Clean(path)), "/")
	if len(parts) < 3 || parts[len(parts)-3] != "kanban" {
		return "", false
	}
	name, ok := strings.CutSuffix(parts[len(parts)-1], ".md")
	if !ok {
		return "", false
	}
	prefix, rest, ok := strings.Cut(name, "-")
	if !ok {
		return "", false
	}
	digits, _, _ := strings.Cut(rest, "-")
	id := prefix + "-" + digits
	typ := workitem.TypeOfID(id)
	if typ == "" || digits == "" || strings.Trim(digits, "0123456789") != "" || parts[len(parts)-2] != workitem.Folder(typ) {
		return "", false
	}
	return id, true
}

// itemWrite refuses a call that replaces a work item file whole with content
// that is not that item: no front matter, front matter that does not parse,
// or another id or type than the file's name and folder give (I-0352,
// ADR-0298). Every other call passes, an edit of part of the file among them.
func itemWrite(c Call) Verdict {
	if c.Kind != FileWrite || !c.Whole {
		return Verdict{}
	}
	id, ok := itemFile(c.Path)
	if !ok {
		return Verdict{}
	}
	why := notItem(c.Content, id)
	if why == "" {
		return Verdict{}
	}
	return deny(fmt.Sprintf("%s cannot replace %s whole: the new content is not work item %s, because %s, and every flai command that lists items would fail on the file until it is put back (I-0352, ADR-0298). Change its body with Edit, and its fields with flai edit or the MCP tool item_edit.", c.Name, c.Path, id, why))
}

// notItem says why content is not the work item id, or "" when it is.
func notItem(content, id string) string {
	it, err := workitem.ParseItem(content)
	switch {
	case err != nil:
		return fmt.Sprintf("its front matter cannot be read (%v)", err)
	case workitem.CanonicalID(it.ID) != workitem.CanonicalID(id):
		return fmt.Sprintf("its front matter names id %q", it.ID)
	case it.Type != workitem.TypeOfID(id):
		return fmt.Sprintf("its front matter names type %q, not %s", it.Type, workitem.TypeOfID(id))
	}
	return ""
}
