package serve

import (
	"slices"
	"strings"
	"testing"
)

const (
	bindInfo    = "info --format {{.OperatingSystem}}"
	bindInspect = "network inspect bridge --format " + gatewayFormat
)

// bindDocker is a docker whose engine is engine, answering network inspect
// bridge with the gateway 172.17.0.1; an empty engine is docker info failing.
func bindDocker(engine string) *fakeDocker {
	d := &fakeDocker{out: map[string]string{bindInspect: "172.17.0.1 \n"}}
	if engine != "" {
		d.out[bindInfo] = engine + "\n"
	} else {
		d.out[bindInfo] = "Cannot connect to the Docker daemon"
		d.fail = map[string]bool{bindInfo: true}
	}
	return d
}

func TestResolveMCPBinds(t *testing.T) {
	cases := []struct {
		name    string
		engine  string
		network string
		binds   []string
		want    []string
		wantErr string
		calls   []string
	}{
		{"docker desktop listens on loopback alone", DockerDesktop, "", []string{"auto"}, []string{"127.0.0.1"}, "", []string{"docker " + bindInfo}},
		{"a native engine adds the bridge gateway", "Ubuntu 24.04 LTS", "", []string{"auto"}, []string{"127.0.0.1", "172.17.0.1"}, "", []string{"docker " + bindInfo, "docker " + bindInspect}},
		{"bridge named outright is the same", "Ubuntu 24.04 LTS", "bridge", []string{"auto"}, []string{"127.0.0.1", "172.17.0.1"}, "", []string{"docker " + bindInfo, "docker " + bindInspect}},
		{"the host network needs loopback alone", "Ubuntu 24.04 LTS", "host", []string{"auto"}, []string{"127.0.0.1"}, "", nil},
		{"an unreadable engine is loopback with an error", "", "", []string{"auto"}, []string{"127.0.0.1"}, "Cannot connect to the Docker daemon", []string{"docker " + bindInfo}},
		{"a list without auto asks docker nothing", "Ubuntu 24.04 LTS", "", []string{"127.0.0.1", "10.0.0.5"}, []string{"127.0.0.1", "10.0.0.5"}, "", nil},
		{"auto and its gateway written out are deduplicated", "Ubuntu 24.04 LTS", "", []string{"auto", "172.17.0.1"}, []string{"127.0.0.1", "172.17.0.1"}, "", []string{"docker " + bindInfo, "docker " + bindInspect}},
		{"auto twice asks docker once", "Ubuntu 24.04 LTS", "", []string{"auto", "auto"}, []string{"127.0.0.1", "172.17.0.1"}, "", []string{"docker " + bindInfo, "docker " + bindInspect}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := bindDocker(c.engine)
			got, err := ResolveMCPBinds(d, "/project", "docker", c.network, c.binds)
			if !slices.Equal(got, c.want) {
				t.Errorf("binds = %q, want %q", got, c.want)
			}
			switch {
			case c.wantErr == "" && err != nil:
				t.Errorf("err = %v, want none", err)
			case c.wantErr != "" && (err == nil || !strings.Contains(err.Error(), c.wantErr)):
				t.Errorf("err = %v, want one saying %q", err, c.wantErr)
			}
			if !slices.Equal(d.calls, c.calls) {
				t.Errorf("docker ran %q, want %q", d.calls, c.calls)
			}
		})
	}
}

func TestResolveMCPBindsGatewayUnreadable(t *testing.T) {
	d := bindDocker("Ubuntu 24.04 LTS")
	d.out[bindInspect] = "Error: No such network: bridge"
	d.fail = map[string]bool{bindInspect: true}
	got, err := ResolveMCPBinds(d, "/project", "docker", "", []string{"auto"})
	if !slices.Equal(got, []string{"127.0.0.1"}) {
		t.Errorf("binds = %q, want loopback alone", got)
	}
	if err == nil || !strings.Contains(err.Error(), "No such network") || !strings.Contains(err.Error(), "Ubuntu 24.04 LTS") {
		t.Errorf("err = %v, want one naming the engine and docker's answer", err)
	}
}

func TestDockerEngine(t *testing.T) {
	if got, err := DockerEngine(bindDocker("Ubuntu 24.04 LTS"), "/project", "docker"); err != nil || got != "Ubuntu 24.04 LTS" {
		t.Errorf("engine = %q, %v, want the trimmed answer", got, err)
	}
	if _, err := DockerEngine(bindDocker(""), "/project", "docker"); err == nil || !strings.Contains(err.Error(), "Cannot connect to the Docker daemon") {
		t.Errorf("err = %v, want docker's answer", err)
	}
	blank := &fakeDocker{out: map[string]string{bindInfo: "\n"}}
	if _, err := DockerEngine(blank, "/project", "docker"); err == nil {
		t.Error("an empty answer is no engine, want an error")
	}
}

func TestNetworkGateway(t *testing.T) {
	d := &fakeDocker{out: map[string]string{bindInspect: "172.17.0.1 fd00::1 \n"}}
	d.out["network inspect none --format "+gatewayFormat] = "\n"
	if got, err := NetworkGateway(d, "/project", "docker", ""); err != nil || got != "172.17.0.1" {
		t.Errorf("gateway = %q, %v, want bridge's first", got, err)
	}
	if _, err := NetworkGateway(d, "/project", "docker", "none"); err == nil || !strings.Contains(err.Error(), "names no gateway") {
		t.Errorf("err = %v, want one saying none has no gateway", err)
	}
	if _, err := NetworkGateway(d, "/project", "docker", "missing"); err == nil || !strings.Contains(err.Error(), "network inspect missing") {
		t.Errorf("err = %v, want one naming the network", err)
	}
}

func TestContainerGateways(t *testing.T) {
	d := &fakeDocker{out: map[string]string{"network ls --format {{.Name}}": "bridge\nhost\nnone\nflai\n"}}
	d.out["network inspect bridge host none flai --format "+gatewayFormat] = "172.17.0.1 \n\n\n172.18.0.1 172.17.0.1 \n"
	got, err := ContainerGateways(d, "/project", "docker")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"172.17.0.1", "172.18.0.1"}; !slices.Equal(got, want) {
		t.Errorf("gateways = %q, want %q", got, want)
	}
	if len(d.calls) != 2 {
		t.Errorf("docker ran %q, want ls and one inspect", d.calls)
	}
	if _, err := ContainerGateways(&fakeDocker{}, "/project", "docker"); err == nil || !strings.Contains(err.Error(), "network ls") {
		t.Errorf("err = %v, want one naming network ls", err)
	}
}
