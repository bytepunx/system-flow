package serve

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
)

// The records of detached host writes (S-0109) lie beside flai serve's
// state, so that a write that ends the flai serve taking it (a restart of
// serve, an upgrade) is still known to the next one when the dashboard asks
// again.

func (d Dir) requests() string { return filepath.Join(string(d), "requests.json") }

var requestsMu sync.Mutex

// Requests changes the records of detached writes, read from the file and
// written back whole. One flai serve runs per directory, so a lock within
// the process is enough.
func (d Dir) Requests(change func(map[string]hostapi.Request)) error {
	requestsMu.Lock()
	defer requestsMu.Unlock()
	all := map[string]hostapi.Request{}
	data, err := os.ReadFile(d.requests())
	switch {
	case errors.Is(err, os.ErrNotExist):
	case err != nil:
		return err
	default:
		if err := json.Unmarshal(data, &all); err != nil {
			return fmt.Errorf("%s is not the records flai serve wrote (%w); remove it, and repeats of requests made before are done again", d.requests(), err)
		}
		if all == nil {
			all = map[string]hostapi.Request{}
		}
	}
	change(all)
	return d.write(d.requests(), all)
}
