package hostapi

import (
	"maps"
)

// The dashboard's writes (S-0075). Each method validates its arguments as
// data, builds the command line itself, and runs flai in the project's
// folder, so the rules stay where they are and are tested: in the commands.
// Nothing the dashboard sends is ever a flag or a command: values are given
// as --flag=value or after --, and IDs, states, and paths are checked first.

// Codes for what a caller can act on, beside the JSON-RPC ones.
const (
	Conflict = -32009 // the thing changed since it was read; Data has the current version
	Refused  = -32010 // flai check refused the content; Data has the findings
	Rule     = -32011 // a workflow rule said no; the message says which
	Disabled = -32012 // a host action the operator has not enabled; Data says what enables it
)

// Trailer marks a commit made for the dashboard, as it always has.
const Trailer = "Co-Authored-By: flaiover <flaiover@localhost>"

// itemSpecs are the dashboard's writes, but for the settings: the work
// items' in writes_items.go, the streams' and documents' in
// writes_streams.go, and the host's in writes_host.go.
func itemSpecs() map[string]spec {
	table := itemWriteSpecs()
	maps.Copy(table, streamSpecs())
	maps.Copy(table, hostSpecs())
	return table
}
