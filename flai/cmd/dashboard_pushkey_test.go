package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

func TestSSHRemoteHost(t *testing.T) {
	cases := []struct {
		remote, host, port string
		ok                 bool
	}{
		{"git@github.com:bytepunx/system-flow.git", "github.com", "", true},
		{"github.com:owner/repo.git", "github.com", "", true},
		{"ssh://git@github.com/owner/repo.git", "github.com", "", true},
		{"ssh://git@git.example.org:22/owner/repo.git", "git.example.org", "", true},
		{"ssh://git@git.example.org:2222/owner/repo.git", "git.example.org", "2222", true},
		{"https://github.com/owner/repo.git", "", "", false},
		{"http://example.org/repo.git", "", "", false},
		{"/srv/git/repo.git", "", "", false},
		{"file:///srv/git/repo.git", "", "", false},
		{"", "", "", false},
	}
	for _, c := range cases {
		host, port, ok := sshRemoteHost(c.remote)
		if host != c.host || port != c.port || ok != c.ok {
			t.Errorf("%q: got %q %q %v, want %q %q %v", c.remote, host, port, ok, c.host, c.port, c.ok)
		}
	}
}

func pushKeyFile(t *testing.T, mode os.FileMode) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "id_test")
	if err := os.WriteFile(p, []byte("not really a key; the fake ssh-keygen decides\n"), mode); err != nil {
		t.Fatal(err)
	}
	return p
}

func sshRunner() *fakeRunner {
	return &fakeRunner{images: map[string]bool{}, running: map[string]bool{},
		remote: "git@github.com:owner/repo.git", knownHosts: "github.com ssh-ed25519 AAAAC3NzaHostKey"}
}

// ADR-0026: the container holds nothing unless the operator asks.
func TestDashboardGivesTheContainerNoKeyByDefault(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	f := sshRunner()
	out, errOut, code := runWith(t, root, f, "dashboard")
	if code != 0 {
		t.Fatal(errOut)
	}
	calls := strings.Join(f.calls, "\n")
	for _, never := range []string{"push_key", "GIT_SSH_COMMAND", "known_hosts", "/etc/passwd", "SSH_AUTH_SOCK", "ssh-keygen"} {
		if strings.Contains(calls, never) {
			t.Errorf("without the opt-in nothing about a key is run, mounted, or set; found %q in:\n%s", never, calls)
		}
	}
	if strings.Contains(out, "push key") {
		t.Errorf("nothing to say about a key:\n%s", out)
	}
	for _, file := range []string{knownHostsPath(root), passwdPath(root)} {
		if _, err := os.Stat(file); err == nil {
			t.Errorf("%s must not be written without the opt-in", file)
		}
	}
}

func TestDashboardMountsTheKeyTheOperatorNames(t *testing.T) {
	for _, how := range []string{"flag", "config"} {
		t.Run(how, func(t *testing.T) {
			t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
			root := tempProject(t)
			key := pushKeyFile(t, 0o600)
			f := sshRunner()
			args := []string{"dashboard"}
			if how == "flag" {
				args = append(args, "--push-key", key)
			} else if _, errOut, code := runWith(t, root, f, "config", "set", "dashboard.push_key", key); code != 0 {
				t.Fatal(errOut)
			}
			out, errOut, code := runWith(t, root, f, args...)
			if code != 0 {
				t.Fatal(errOut)
			}
			var run string
			for _, c := range f.calls {
				if strings.HasPrefix(c, "docker run ") {
					run = c
				}
			}
			for _, want := range []string{
				"--mount type=bind,source=" + key + ",target=/run/flaiover/push_key,readonly",
				"--mount type=bind,source=" + knownHostsPath(root) + ",target=/run/flaiover/known_hosts,readonly",
				"--env GIT_SSH_COMMAND=ssh -F /dev/null -i /run/flaiover/push_key -o IdentitiesOnly=yes -o IdentityAgent=none -o BatchMode=yes -o StrictHostKeyChecking=yes -o UserKnownHostsFile=/run/flaiover/known_hosts",
			} {
				if !strings.Contains(run, want) {
					t.Errorf("missing %q in:\n%s", want, run)
				}
			}
			if strings.Contains(run, "SSH_AUTH_SOCK") {
				t.Errorf("the agent socket is never forwarded:\n%s", run)
			}
			if os.Getuid() != 0 {
				if !strings.Contains(run, "--mount type=bind,source="+passwdPath(root)+",target=/etc/passwd,readonly") {
					t.Errorf("the host user needs a passwd entry for OpenSSH:\n%s", run)
				}
				data, _ := os.ReadFile(passwdPath(root))
				if !strings.Contains(string(data), "root:x:0:0:") || !strings.Contains(string(data), ":x:"+itoa(os.Getuid())+":"+itoa(os.Getgid())+":") {
					t.Errorf("passwd:\n%s", data)
				}
			}
			hosts, _ := os.ReadFile(knownHostsPath(root))
			if string(hosts) != "github.com ssh-ed25519 AAAAC3NzaHostKey\n" {
				t.Errorf("the host's own entry is pinned, without ssh-keygen's comment line: %q", hosts)
			}
			if st, _ := os.Stat(knownHostsPath(root)); st != nil && st.Mode().Perm() != 0o600 {
				t.Errorf("known_hosts mode %v", st.Mode().Perm())
			}
			for _, want := range []string{"push key: SHA256:fakefingerprint (alex at laptop)", "pushed to git@github.com:owner/repo.git", "Whoever holds the dashboard token can now publish a release"} {
				if !strings.Contains(out, want) {
					t.Errorf("the operator is told what the container holds; missing %q in:\n%s", want, out)
				}
			}
			if strings.Contains(out, "not really a key") {
				t.Error("the key itself is never printed")
			}
		})
	}
}

func TestDashboardRefusesAKeyThatCannotWork(t *testing.T) {
	cases := []struct {
		name string
		prep func(t *testing.T, f *fakeRunner) string // returns the key path
		want []string
	}{
		{"a missing file", func(t *testing.T, f *fakeRunner) string { return filepath.Join(t.TempDir(), "gone") }, []string{"cannot be read", "--push-key"}},
		{"a directory", func(t *testing.T, f *fakeRunner) string { return t.TempDir() }, []string{"not a regular file"}},
		{"a key others can read", func(t *testing.T, f *fakeRunner) string { return pushKeyFile(t, 0o644) }, []string{"readable by others", "chmod 600"}},
		{"a key with a passphrase", func(t *testing.T, f *fakeRunner) string { f.keyKind = "passphrase"; return pushKeyFile(t, 0o600) },
			[]string{"protected by a passphrase", "never use it unattended", "every key in the agent", "deploy key"}},
		{"a file that is not a private key", func(t *testing.T, f *fakeRunner) string { f.keyKind = "notkey"; return pushKeyFile(t, 0o600) }, []string{"not an SSH private key", ".pub"}},
		{"an HTTPS remote", func(t *testing.T, f *fakeRunner) string {
			f.remote = "https://github.com/owner/repo.git"
			return pushKeyFile(t, 0o600)
		}, []string{"not an SSH URL", "git remote set-url origin"}},
		{"no remote", func(t *testing.T, f *fakeRunner) string { f.remote = ""; return pushKeyFile(t, 0o600) }, []string{"no git remote named origin"}},
		{"a host this machine has never verified", func(t *testing.T, f *fakeRunner) string { f.knownHosts = ""; return pushKeyFile(t, 0o600) },
			[]string{"no known_hosts entry for github.com", "must not accept a host key on first use", "ssh -T git@github.com"}},
		{"a path docker cannot mount", func(t *testing.T, f *fakeRunner) string {
			dir := filepath.Join(t.TempDir(), "a,b")
			_ = os.MkdirAll(dir, 0o755)
			p := filepath.Join(dir, "key")
			_ = os.WriteFile(p, []byte("k"), 0o600)
			return p
		}, []string{"contains a comma"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
			root := tempProject(t)
			f := sshRunner()
			key := c.prep(t, f)
			_, errOut, code := runWith(t, root, f, "dashboard", "--push-key", key)
			if code == 0 {
				t.Fatal("must refuse")
			}
			for _, w := range c.want {
				if !strings.Contains(errOut, w) {
					t.Errorf("the refusal says why and what to do; missing %q in:\n%s", w, errOut)
				}
			}
			for _, call := range f.calls {
				if strings.HasPrefix(call, "docker run ") {
					t.Errorf("refused before anything is started, but ran: %s", call)
				}
			}
		})
	}
}

func TestDashboardTakesHostKeysFromTheFileTheOperatorNames(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	hosts := filepath.Join(t.TempDir(), "curated_hosts")
	_ = os.WriteFile(hosts, []byte("github.com ssh-ed25519 AAAAcurated\n"), 0o644)
	f := sshRunner()
	if _, errOut, code := runWith(t, root, f, "dashboard", "--push-key", pushKeyFile(t, 0o600), "--push-known-hosts", hosts); code != 0 {
		t.Fatal(errOut)
	}
	if calls := strings.Join(f.calls, "\n"); !strings.Contains(calls, "ssh-keygen -F github.com -f "+hosts) || strings.Contains(calls, "ssh-keygen -F github.com\n") {
		t.Errorf("only the named file is searched:\n%s", calls)
	}
	_, errOut, code := runWith(t, tempProject(t), sshRunner(), "dashboard", "--push-key", pushKeyFile(t, 0o600), "--push-known-hosts", filepath.Join(t.TempDir(), "gone"))
	if code == 0 || !strings.Contains(errOut, "cannot be read") {
		t.Errorf("a missing known hosts file is refused: %d %s", code, errOut)
	}
}

func TestDashboardStatusSaysWhatTheContainerHolds(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	f := sshRunner()
	f.running["flaiover-t"] = true
	out, errOut, code := runWith(t, root, f, "dashboard", "status")
	if code != 0 {
		t.Fatal(errOut)
	}
	if !strings.Contains(out, "holds no credential") {
		t.Errorf("status without a key:\n%s", out)
	}
	f.pushMount = "/home/alex/.ssh/flaiover-push"
	out, _, _ = runWith(t, root, f, "dashboard", "status")
	if !strings.Contains(out, "holds a push key: SHA256:fakefingerprint from /home/alex/.ssh/flaiover-push") {
		t.Errorf("status with a key:\n%s", out)
	}
	out, _, _ = runWith(t, root, f, "dashboard", "status", "--json")
	var st struct {
		PushKey map[string]string `json:"push_key"`
	}
	if err := json.Unmarshal([]byte(out), &st); err != nil || st.PushKey["fingerprint"] != "SHA256:fakefingerprint" {
		t.Errorf("status --json: %v\n%s", err, out)
	}
}

// The fake runner decides what ssh-keygen says everywhere else; this one
// asks the real tool, so the refusals do not rest on remembered messages.
func TestPushKeyChecksWithTheRealSSHKeygen(t *testing.T) {
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		t.Skip("ssh-keygen not installed")
	}
	dir := t.TempDir()
	gen := func(name, passphrase string) string {
		p := filepath.Join(dir, name)
		if out, err := exec.Command("ssh-keygen", "-q", "-t", "ed25519", "-N", passphrase, "-C", "flai test "+name, "-f", p).CombinedOutput(); err != nil {
			t.Fatalf("ssh-keygen: %v %s", err, out)
		}
		return p
	}
	open, locked := gen("open", ""), gen("locked", "correct horse")
	r := execx.System{}
	if _, err := r.Run("", "ssh-keygen", "-y", "-P", "", "-f", open); err != nil {
		t.Errorf("a key without a passphrase opens: %v", err)
	}
	out, err := r.Run("", "ssh-keygen", "-y", "-P", "", "-f", locked)
	if err == nil || !strings.Contains(strings.ToLower(out+err.Error()), "passphrase") {
		t.Errorf("a locked key is recognised by the word passphrase: %q %v", out, err)
	}
	out, err = r.Run("", "ssh-keygen", "-y", "-P", "", "-f", open+".pub")
	if err == nil {
		t.Errorf("a public key is not a private key: %q", out)
	} else if strings.Contains(strings.ToLower(out+err.Error()), "passphrase") {
		t.Errorf("a public key must not be mistaken for a locked key: %q %v", out, err)
	}
	line, err := r.Run("", "ssh-keygen", "-l", "-f", open)
	if f := strings.Fields(line); err != nil || len(f) < 4 || !strings.HasPrefix(f[1], "SHA256:") || strings.Join(f[2:len(f)-1], " ") != "flai test open" {
		t.Errorf("fingerprint line: %q %v", line, err)
	}
}

func itoa(n int) string { return strconv.Itoa(n) }
