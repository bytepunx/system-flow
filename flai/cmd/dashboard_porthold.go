package cmd

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// procRoot is where the process table is read on Linux; findPortHolder takes
// it as an argument so that tests give it a fake tree (I-0243).
const procRoot = "/proc"

// portProbeTimeout bounds the dial that asks whether a host port is held.
const portProbeTimeout = 2 * time.Second

// portHolder is a docker-proxy process that holds a host port: what flai names
// when the dashboard's recorded port answers although docker publishes nothing
// for its container (I-0243).
type portHolder struct {
	PID           int    `json:"pid"`
	HostIP        string `json:"host_ip,omitempty"`
	HostPort      int    `json:"host_port"`
	ContainerIP   string `json:"container_ip,omitempty"`
	ContainerPort string `json:"container_port,omitempty"`
}

// String names the holder for a message, such as
// "docker-proxy (pid 188220) for 172.17.0.2:3000".
func (h portHolder) String() string {
	s := fmt.Sprintf("docker-proxy (pid %d)", h.PID)
	if h.ContainerIP != "" && h.ContainerPort != "" {
		s += " for " + net.JoinHostPort(h.ContainerIP, h.ContainerPort)
	}
	return s
}

// findPortHolder reads each <root>/<pid>/cmdline and returns the first
// docker-proxy whose arguments carry -host-port port. It finds none, and
// fails nothing, when root does not exist, when a process vanishes or cannot
// be read mid-scan, or when no process matches. The command line is read
// because any user may read it, where the sockets of a root process need root.
func findPortHolder(root string, port int) (portHolder, bool) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return portHolder{}, false
	}
	want := strconv.Itoa(port)
	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil || pid <= 0 {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, e.Name(), "cmdline"))
		if err != nil || len(data) == 0 {
			continue
		}
		args := strings.Split(strings.TrimRight(string(data), "\x00"), "\x00")
		if filepath.Base(args[0]) != "docker-proxy" {
			continue
		}
		flags := proxyFlags(args[1:])
		if flags["host-port"] != want {
			continue
		}
		return portHolder{
			PID:           pid,
			HostIP:        flags["host-ip"],
			HostPort:      port,
			ContainerIP:   flags["container-ip"],
			ContainerPort: flags["container-port"],
		}, true
	}
	return portHolder{}, false
}

// proxyFlags reads docker-proxy's flags, given as -name value or -name=value,
// into a map from name to value.
func proxyFlags(args []string) map[string]string {
	flags := map[string]string{}
	for i := 0; i < len(args); i++ {
		a := strings.TrimLeft(args[i], "-")
		if a == args[i] || a == "" {
			continue
		}
		if name, value, ok := strings.Cut(a, "="); ok {
			flags[name] = value
			continue
		}
		if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
			flags[a] = args[i+1]
			i++
		}
	}
	return flags
}

// publishAddress is the host address to dial for a docker publish of
// <bind>:<port>:<container port>, such as 0.0.0.0:4242:3000 or
// [::]:4242:3000, with a wildcard bind dialled at 127.0.0.1, and its port.
func publishAddress(publish string) (addr string, port int, err error) {
	i := strings.LastIndex(publish, ":")
	if i < 0 {
		return "", 0, fmt.Errorf("publish %q is not <bind>:<port>:<container port>", publish)
	}
	hostPart := publish[:i]
	j := strings.LastIndex(hostPart, ":")
	if j < 0 {
		return "", 0, fmt.Errorf("publish %q is not <bind>:<port>:<container port>", publish)
	}
	bind, p := hostPart[:j], hostPart[j+1:]
	port, err = strconv.Atoi(p)
	if err != nil || port <= 0 || port > 65535 {
		return "", 0, fmt.Errorf("publish %q has no host port", publish)
	}
	bind = strings.TrimSuffix(strings.TrimPrefix(bind, "["), "]")
	return probeHost(net.JoinHostPort(bind, p)), port, nil
}

// publishAnswers says whether the host port of a docker publish accepts a TCP
// connection within portProbeTimeout, whatever holds it.
func publishAnswers(publish string) bool {
	addr, _, err := publishAddress(publish)
	return err == nil && tcpAnswers(addr, portProbeTimeout)
}

// publishHolder is the docker-proxy in the host's process table, /proc, that
// holds the host port of a docker publish, when one is found.
func publishHolder(publish string) (portHolder, bool) {
	_, port, err := publishAddress(publish)
	if err != nil {
		return portHolder{}, false
	}
	return findPortHolder(procRoot, port)
}

// tcpAnswers says whether addr accepts a TCP connection within timeout.
func tcpAnswers(addr string, timeout time.Duration) bool {
	c, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	_ = c.Close()
	return true
}
