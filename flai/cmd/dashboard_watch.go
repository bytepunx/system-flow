package cmd

import (
	"net"
	"strings"
)

// How the dashboard container stands, as flai dashboard status and flai
// host's watch report it (S-0184).
const (
	dashboardRunning      = "running"       // docker runs it and /_health answers
	dashboardNotAnswering = "not-answering" // docker runs it and /_health does not answer
	dashboardGone         = "gone"          // docker does not run it
)

// dashboardState probes the dashboard container: gone when docker does not
// run it, running when its published port answers /_health, not answering
// otherwise. health is docker's own verdict from the image's HEALTHCHECK,
// "" when the image has none.
func (a *app) dashboardState(name string) (state, health string, err error) {
	running, err := a.containerRunning(name)
	if err != nil {
		return "", "", err
	}
	if !running {
		return dashboardGone, "", nil
	}
	health = a.dockerHealth(name)
	addr, err := a.probeAddress(name)
	if err == nil && a.healthProbe("http://"+probeHost(addr)+"/_health") {
		return dashboardRunning, health, nil
	}
	return dashboardNotAnswering, health, nil
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
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}
	return net.JoinHostPort(host, port)
}
