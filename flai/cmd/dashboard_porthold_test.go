package cmd

import (
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
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

// TestTCPAnswersAListener dials for real only a listener the test holds open.
func TestTCPAnswersAListener(t *testing.T) {
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
	_, open, _ := net.SplitHostPort(l.Addr().String())
	if !publishAnswers("0.0.0.0:" + open + ":3000") {
		t.Fatal("an open listener's wildcard publish does not answer")
	}
}

// fakeDialer records the addresses dialled and answers each as answer says:
// a connection it closes itself, or a refusal.
type fakeDialer struct {
	answer bool
	dialed []string
	conns  []net.Conn
}

func (d *fakeDialer) dial(network, addr string, _ time.Duration) (net.Conn, error) {
	d.dialed = append(d.dialed, network+" "+addr)
	if !d.answer {
		return nil, &net.OpError{Op: "dial", Net: network, Err: syscall.ECONNREFUSED}
	}
	ours, theirs := net.Pipe()
	d.conns = append(d.conns, theirs)
	return ours, nil
}

// TestTCPAnswersNotARefusedPort asserts the closed-port case through a dialer
// that refuses, not a real closed port: on a WSL2 host with mirrored
// networking a port free in Linux can be a Windows listener's, so a real dial
// to it may connect (I-0321).
func TestTCPAnswersNotARefusedPort(t *testing.T) {
	t.Parallel()
	refuse := &fakeDialer{}
	if tcpAnswersWith(refuse.dial, "127.0.0.1:4242", time.Second) {
		t.Fatal("a refused port answers")
	}
	if publishAnswersWith(refuse.dial, "0.0.0.0:4242:3000") {
		t.Fatal("a refused port's publish answers")
	}
	if publishAnswersWith(refuse.dial, "not-a-publish") {
		t.Fatal("a malformed publish answers")
	}
	want := []string{"tcp 127.0.0.1:4242", "tcp 127.0.0.1:4242"}
	if strings.Join(refuse.dialed, ",") != strings.Join(want, ",") {
		t.Fatalf("dialled %q, want %q: a malformed publish dials nothing", refuse.dialed, want)
	}
}

// TestTCPAnswersWhateverAcceptsThePort shows what a real dial meets on such a
// host: a port something else accepts answers, and the probe closes the
// connection it opened.
func TestTCPAnswersWhateverAcceptsThePort(t *testing.T) {
	t.Parallel()
	relay := &fakeDialer{answer: true}
	if !publishAnswersWith(relay.dial, "127.0.0.1:4242:3000") {
		t.Fatal("a port something accepts does not answer")
	}
	if len(relay.conns) != 1 {
		t.Fatalf("dialled %d times, want 1", len(relay.conns))
	}
	if _, err := relay.conns[0].Read(make([]byte, 1)); !errors.Is(err, io.EOF) {
		t.Fatalf("read on the far end = %v, want io.EOF: the probe left its connection open", err)
	}
}

func TestPublishHolderFindsNoneForAMalformedPublish(t *testing.T) {
	t.Parallel()
	if h, ok := publishHolder("not-a-publish"); ok {
		t.Fatalf("holder = %+v, want none", h)
	}
}
