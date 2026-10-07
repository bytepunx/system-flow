package host

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

// The control API. Every request but /_health needs the host's token as a
// bearer; any request carrying an Origin is refused before that, so that no
// web page a browser on this machine loads can reach it.

// ProcessRequest asks the host to start, stop, or restart serve, the MCP
// servers, or all of them.
type ProcessRequest struct {
	Process string `json:"process"`
	Action  string `json:"action"`
}

// MCPRequest names every project whose MCP server the host is to keep.
type MCPRequest struct {
	Roots []string `json:"roots"`
}

// UpgradeRequest names the published flai release to install; no body, or
// no version, is the newest (S-0298).
type UpgradeRequest struct {
	Version string `json:"version,omitempty"`
}

// apiError is how the API says no.
type apiError struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *host) api() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /_health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, Health{PID: os.Getpid(), Version: h.o.Version, Config: h.o.Config, Started: h.started})
	})
	mux.HandleFunc("GET /status", h.authed(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, h.status())
	}))
	mux.HandleFunc("POST /process", h.authed(func(w http.ResponseWriter, r *http.Request) {
		var in ProcessRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSON(w, http.StatusBadRequest, apiError{"the body is a JSON object with process and action"})
			return
		}
		if err := h.act(in.Process, in.Action); err != nil {
			writeJSON(w, http.StatusBadRequest, apiError{err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, h.status())
	}))
	mux.HandleFunc("PUT /mcp", h.authed(func(w http.ResponseWriter, r *http.Request) {
		var in MCPRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSON(w, http.StatusBadRequest, apiError{"the body is a JSON object with roots"})
			return
		}
		h.wantMCP(in.Roots)
		writeJSON(w, http.StatusOK, h.status())
	}))
	mux.HandleFunc("GET /check", h.authed(func(w http.ResponseWriter, r *http.Request) {
		if h.o.Check == nil {
			writeJSON(w, http.StatusNotImplemented, apiError{"this host cannot check for a newer flai"})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
		defer cancel()
		res, err := h.o.Check(ctx)
		if err != nil {
			writeJSON(w, http.StatusBadGateway, apiError{err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, res)
	}))
	mux.HandleFunc("GET /versions", h.authed(func(w http.ResponseWriter, r *http.Request) {
		if h.o.Versions == nil {
			writeJSON(w, http.StatusNotImplemented, apiError{"this host cannot list the published flai releases"})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
		defer cancel()
		res, err := h.o.Versions(ctx)
		if err != nil {
			writeJSON(w, http.StatusBadGateway, apiError{err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, res)
	}))
	mux.HandleFunc("POST /upgrade", h.authed(func(w http.ResponseWriter, r *http.Request) {
		in, err := upgradeRequest(r.Body)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, apiError{err.Error()})
			return
		}
		// the upgrade goes on if whoever asked goes away: it replaces a binary
		ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), 5*time.Minute)
		defer cancel()
		res, err := h.upgrade(ctx, in.Version)
		if err != nil {
			writeJSON(w, http.StatusBadGateway, apiError{err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, res)
	}))
	return refuseOrigin(mux)
}

// upgradeRequest reads POST /upgrade's body: none, or a JSON object whose
// version, when given, is a bare X.Y.Z, refused before anything runs.
func upgradeRequest(body io.Reader) (UpgradeRequest, error) {
	var in UpgradeRequest
	if err := json.NewDecoder(body).Decode(&in); err != nil && !errors.Is(err, io.EOF) {
		return in, errors.New(`the body is empty for the newest flai, or a JSON object {"version": "X.Y.Z"} naming a published release`)
	}
	if in.Version != "" && !bareVersion(in.Version) {
		return in, fmt.Errorf("%q is not a release version: give X.Y.Z, such as 1.2.3, or no version for the newest flai", in.Version)
	}
	return in, nil
}

// bareVersion says whether v is X.Y.Z, with no v, suffix, or leading zeros.
func bareVersion(v string) bool {
	n, ok := buildinfo.Semver(v)
	return ok && fmt.Sprintf("%d.%d.%d", n[0], n[1], n[2]) == v
}

func refuseOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") != "" {
			writeJSON(w, http.StatusForbidden, apiError{"flai host answers no web page"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *host) authed(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		got, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || subtle.ConstantTimeCompare([]byte(got), []byte(h.token)) != 1 {
			writeJSON(w, http.StatusUnauthorized, apiError{"the host's token is required"})
			return
		}
		next(w, r)
	}
}
