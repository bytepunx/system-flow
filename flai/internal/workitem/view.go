package workitem

import "sort"

// Visit calls fn with every item, the archive's too when includeArchive, as
// the store keeps it, in no set order, and with whether it is archived. It
// copies nothing, so that a caller reading a few fields of every item does not
// pay for a copy of each (S-0606). The item is shared with every other caller
// and the store never changes it: neither fn nor anything it keeps the item
// for may change it, and what is handed on beyond the caller is a copy.
func (r *Repo) Visit(includeArchive bool, fn func(it *Item, archived bool)) error {
	scopes := []bool{false}
	if includeArchive {
		scopes = append(scopes, true)
	}
	for _, archived := range scopes {
		for _, typ := range Types {
			items, err := parsed.items(r.ItemDir(typ, archived))
			if err != nil {
				return err
			}
			for _, it := range items {
				fn(it, archived)
			}
		}
	}
	return nil
}

// ChildrenOf returns copies of the items whose parent is id, the archive's
// too when includeArchive, in ID order, with Archived set: what
// Children(List(includeArchive), id) returns, copying only the children.
func (r *Repo) ChildrenOf(id string, includeArchive bool) ([]*Item, error) {
	var out []*Item
	err := r.Visit(includeArchive, func(it *Item, archived bool) {
		if it.Parent != id {
			return
		}
		c := it.clone()
		c.Archived = archived
		out = append(out, c)
	})
	if err != nil {
		return nil, err
	}
	sort.SliceStable(out, func(i, j int) bool { return lessID(out[i].ID, out[j].ID) })
	return out, nil
}

// LessID orders two IDs as List does: epics, then stories, then tasks, each
// by number whatever its zero padding.
func LessID(a, b string) bool { return lessID(a, b) }
