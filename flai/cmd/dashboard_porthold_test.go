package cmd

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// writeProc writes a fake process <root>/<pid> whose cmdline is args, NUL
// separated as the kernel gives it; nil args writes the folder alone.
func writeProc(t *testing.T, root, pid string, args []string) {
	t.Helper()
	dir := filepath.Join(root, pid)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if args == nil {
		return
	}
	data := strings.Join(args, "\x00") + "\x00"
	if err := os.WriteFile(filepath.Join(dir, "cmdline"), []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fakeProcWithProxies(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeProc(t, root, "1", []string{"/sbin/init"})
	writeProc(t, root, "4100", []string{"/usr/bin/docker-proxy", "-proto", "tcp", "-host-ip", "0.0.0.0", "-host-port", "8080", "-container-ip", "172.17.0.3", "-container-port", "80"})
	writeProc(t, root, "188220", []string{"/usr/bin/docker-proxy", "-proto", "tcp", "-host-ip", "0.0.0.0", "-host-port", "4242", "-container-ip", "172.17.0.2", "-container-port", "3000"})
	writeProc(t, root, "188237", []string{"node", "server.js", "-host-port", "4242"})
	writeProc(t, root, "200000", nil) // a process that vanished before its cmdline was read
	if err := os.MkdirAll(filepath.Join(root, "self"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "uptime"), []byte("1 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestFindPortHolderNamesTheDockerProxyOfThePort(t *testing.T) {
	t.Parallel()
	root := fakeProcWithProxies(t)
	h, ok := findPortHolder(root, 4242)
	if !ok {
		t.Fatal("no holder found for 4242")
	}
	want := portHolder{PID: 188220, HostIP: "0.0.0.0", HostPort: 4242, ContainerIP: "172.17.0.2", ContainerPort: "3000"}
	if h != want {
		t.Fatalf("holder = %+v, want %+v", h, want)
	}
	if got := h.String(); got != "docker-proxy (pid 188220) for 172.17.0.2:3000" {
		t.Fatalf("String() = %q", got)
	}
}

func TestFindPortHolderTakesFlagsWithEquals(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeProc(t, root, "77", []string{"docker-proxy", "-proto=tcp", "-host-ip=127.0.0.1", "-host-port=4242", "-container-ip=172.17.0.9", "-container-port=3000"})
	h, ok := findPortHolder(root, 4242)
	if !ok || h.PID != 77 || h.HostIP != "127.0.0.1" || h.ContainerIP != "172.17.0.9" || h.ContainerPort != "3000" {
		t.Fatalf("holder = %+v, %v", h, ok)
	}
}

func TestFindPortHolderFindsNoneForAPortNoProcessNames(t *testing.T) {
	t.Parallel()
	if h, ok := findPortHolder(fakeProcWithProxies(t), 9999); ok {
		t.Fatalf("holder = %+v, want none", h)
	}
}

func TestFindPortHolderFindsNoneWithoutAProcessTable(t *testing.T) {
	t.Parallel()
	if h, ok := findPortHolder(filepath.Join(t.TempDir(), "no-proc"), 4242); ok {
		t.Fatalf("holder = %+v, want none", h)
	}
}

func TestPortHolderStringWithoutContainerAddress(t *testing.T) {
	t.Parallel()
	if got := (portHolder{PID: 5, HostPort: 4242}).String(); got != "docker-proxy (pid 5)" {
		t.Fatalf("String() = %q", got)
	}
}

func TestPublishAddressDialsWildcardsAtLoopback(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		publish, addr string
		port          int
	}{
		{"0.0.0.0:4242:3000", "127.0.0.1:4242", 4242},
		{"127.0.0.1:4242:3000", "127.0.0.1:4242", 4242},
		{"[::]:4242:3000", "127.0.0.1:4242", 4242},
		{":4242:3000", "127.0.0.1:4242", 4242},
		{"192.168.1.5:8080:3000", "192.168.1.5:8080", 8080},
		{"[::1]:4242:3000", "[::1]:4242", 4242},
	} {
		addr, port, err := publishAddress(tc.publish)
		if err != nil || addr != tc.addr || port != tc.port {
			t.Errorf("publishAddress(%q) = %q, %d, %v; want %q, %d", tc.publish, addr, port, err, tc.addr, tc.port)
		}
	}
	for _, bad := range []string{"", "4242", "0.0.0.0:x:3000", "0.0.0.0:0:3000"} {
		if _, _, err := publishAddress(bad); err == nil {
			t.Errorf("publishAddress(%q) gave no error", bad)
		}
	}
}

func TestTCPAnswersAListenerAndNotAClosedOne(t *testing.T) {
	t.Parallel()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = l.Close() }()
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			_ = c.Close()
		}
	}()
	if !tcpAnswers(l.Addr().String(), time.Second) {
		t.Fatal("an open listener does not answer")
	}

	closed, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := closed.Addr().String()
	_ = closed.Close()
	if tcpAnswers(addr, time.Second) {
		t.Fatal("a closed listener answers")
	}
	_, port, _ := net.SplitHostPort(addr)
	if publishAnswers("127.0.0.1:" + port + ":3000") {
		t.Fatal("a closed listener's publish answers")
	}
	_, open, _ := net.SplitHostPort(l.Addr().String())
	if !publishAnswers("0.0.0.0:" + open + ":3000") {
		t.Fatal("an open listener's wildcard publish does not answer")
	}
	if publishAnswers("not-a-publish") {
		t.Fatal("a malformed publish answers")
	}
}

func TestPublishHolderFindsNoneForAMalformedPublish(t *testing.T) {
	t.Parallel()
	if h, ok := publishHolder("not-a-publish"); ok {
		t.Fatalf("holder = %+v, want none", h)
	}
}
