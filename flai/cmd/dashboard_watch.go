package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/host"
)

// dashboardState probes the dashboard container: gone when docker does not
// run it, running when its published port answers /_health, not answering
// otherwise (S-0184). health is docker's own verdict from the image's
// HEALTHCHECK, "" when the image has none.
func (a *app) dashboardState(name string) (state, health string, err error) {
	running, err := a.containerRunning(name)
	if err != nil {
		return "", "", err
	}
	if !running {
		return host.DashboardGone, "", nil
	}
	health = a.dockerHealth(name)
	addr, err := a.probeAddress(name)
	if err == nil && a.healthProbe("http://"+probeHost(addr)+"/_health") {
		return host.DashboardRunning, health, nil
	}
	return host.DashboardNotAnswering, health, nil
}

// dockerHealth is the container's health as docker's HEALTHCHECK last judged
// it (starting, healthy, unhealthy), or "" when its image has none.
func (a *app) dockerHealth(name string) string {
	out, err := a.runner.Run("", "docker", "inspect", "--format", "{{if .State.Health}}{{.State.Health.Status}}{{end}}", name)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

// probeHost is addr with a wildcard host (0.0.0.0, ::) replaced by 127.0.0.1,
// so that a port published on every interface is probed from this host.
func probeHost(addr string) string {
	h, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	if h == "" || h == "0.0.0.0" || h == "::" {
		h = "127.0.0.1"
	}
	return net.JoinHostPort(h, port)
}

// dashboardRecord is the shared container as flai dashboard last started it,
// kept beside its token so that flai host can start it again as it was: the
// image reference it ran, never a fresh pull (S-0184).
type dashboardRecord struct {
	Name    string `json:"name"`
	Ref     string `json:"ref"`
	Publish string `json:"publish"`
}

const dashboardRecordFile = "dashboard-container.json"

func readDashboardRecord(dir string) (dashboardRecord, bool) {
	var r dashboardRecord
	data, err := os.ReadFile(filepath.Join(dir, dashboardRecordFile))
	if err != nil || json.Unmarshal(data, &r) != nil || r.Name == "" || r.Ref == "" || r.Publish == "" {
		return dashboardRecord{}, false
	}
	return r, true
}

func writeDashboardRecord(dir string, r dashboardRecord) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	tmp := filepath.Join(dir, dashboardRecordFile+".tmp")
	if err := os.WriteFile(tmp, append(data, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, dashboardRecordFile))
}

// forgetDashboard drops the record, so that flai host no longer watches a
// dashboard the operator stopped.
func forgetDashboard(dir string) {
	_ = os.Remove(filepath.Join(dir, dashboardRecordFile))
}

// startDashboard starts the shared container from ref, published as s says,
// and records it for flai host's watch.
func (a *app) startDashboard(dir string, s dashboardSettings, ref string) (string, error) {
	publish := fmt.Sprintf("%s:%d:%d", s.Bind, s.Port, containerPort)
	id, err := a.startContainer(dir, s.Name, publish, ref)
	if err != nil {
		return id, err
	}
	a.recordDashboard(dir, dashboardRecord{Name: s.Name, Ref: ref, Publish: publish})
	return id, nil
}

func (a *app) recordDashboard(dir string, r dashboardRecord) {
	if err := writeDashboardRecord(dir, r); err != nil {
		a.logger().Warn("dashboard not recorded: flai host will not restart it", "component", "dashboard", "err", err.Error())
	}
}

// dashboardWatch is flai host's watch of the dashboard (S-0184): it looks at
// the container flai dashboard recorded, and starts it again as recorded.
func (a *app) dashboardWatch() *host.Watch {
	return &host.Watch{Look: a.lookDashboard, Restart: a.restartDashboard}
}

// lookDashboard is off when dashboard.no_restart is set or no dashboard is
// recorded, and otherwise how it stands. When docker cannot say, it is gone:
// the restart that follows fails and says why.
func (a *app) lookDashboard(context.Context) string {
	cfg, _, err := a.loadConfig()
	if err != nil || cfg.Dashboard.NoRestart {
		return host.DashboardOff
	}
	r, ok := readDashboardRecord(string(a.serveDir()))
	if !ok {
		return host.DashboardOff
	}
	state, _, err := a.dashboardState(r.Name)
	if err != nil {
		return host.DashboardGone
	}
	return state
}

// restartDashboard removes the recorded container, if anything is left of it,
// and starts it again from the image and publish it was recorded with.
func (a *app) restartDashboard(context.Context) error {
	dir := string(a.serveDir())
	r, ok := readDashboardRecord(dir)
	if !ok {
		return errors.New("no dashboard is recorded; flai dashboard starts one")
	}
	_, _ = a.runner.Run("", "docker", "rm", "-f", r.Name) // a container that does not answer is still there
	if _, err := a.startContainer(dir, r.Name, r.Publish, r.Ref); err != nil {
		return fmt.Errorf("start %s from %s: %w", r.Name, r.Ref, err)
	}
	return nil
}
