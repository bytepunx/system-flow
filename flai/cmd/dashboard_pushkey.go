package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The push key (ADR-0026, S-0062). By default the dashboard container holds
// no credential, and an acceptance made from the board is committed and
// tagged in the clone and pushed by someone on the host. An operator who
// wants the board to publish names an SSH private key on the host
// (dashboard.push_key in the host config, or --push-key), and flai dashboard
// mounts it read-only with the remote's host keys pinned, so the bundled
// flai accept can push. It is off unless asked for, and everything that
// would make it fail at the first acceptance is refused at start instead.

const (
	pushKeyMountPath    = "/run/flaiover/push_key"
	knownHostsMountPath = "/run/flaiover/known_hosts"
	passwdMountPath     = "/etc/passwd"
	// The container's git uses this key and these host keys and nothing
	// else: no agent, no user configuration, no prompt, and never a host key
	// accepted on first use.
	pushSSHCommand = "ssh -F /dev/null -i " + pushKeyMountPath + " -o IdentitiesOnly=yes -o IdentityAgent=none -o BatchMode=yes -o StrictHostKeyChecking=yes -o UserKnownHostsFile=" + knownHostsMountPath
)

// pushKey is what flai dashboard worked out about the key it was given.
type pushKey struct {
	Path        string `json:"path"`
	Fingerprint string `json:"fingerprint"`
	Comment     string `json:"comment,omitempty"`
	Remote      string `json:"remote"`
	Host        string `json:"host"`
	knownHosts  string // file holding the pinned host keys
	passwd      string // file holding the passwd entries, empty when not needed
}

func knownHostsPath(root string) string {
	return filepath.Join(root, ".flai-cache", "dashboard.known_hosts")
}
func passwdPath(root string) string { return filepath.Join(root, ".flai-cache", "dashboard.passwd") }

// expandHome turns a leading ~ into the user's home directory.
func expandHome(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(p, "~"))
		}
	}
	return p
}

var scpLike = regexp.MustCompile(`^([^@/\s]+@)?([^:/\s]+):[^/].*$`)

// sshRemoteHost returns the host, and the port when it is not 22, of an SSH
// remote URL: ssh://[user@]host[:port]/path or the scp form [user@]host:path.
func sshRemoteHost(remote string) (host, port string, ok bool) {
	if strings.HasPrefix(remote, "ssh://") {
		u, err := url.Parse(remote)
		if err != nil || u.Hostname() == "" {
			return "", "", false
		}
		if p := u.Port(); p != "" && p != "22" {
			port = p
		}
		return u.Hostname(), port, true
	}
	if strings.Contains(remote, "://") {
		return "", "", false
	}
	if m := scpLike.FindStringSubmatch(remote); m != nil {
		return m[2], "", true
	}
	return "", "", false
}

// preparePushKey checks the key and writes the files mounted with it. Every
// refusal says what is wrong and what to do; none of them leaves anything
// half set up, because nothing has been started yet.
func (a *app) preparePushKey(repo *workitem.Repo, path, hostsFile, remoteName string) (*pushKey, error) {
	if runtime.GOOS == "windows" {
		return nil, fmt.Errorf("a push key is not supported on a Windows host yet: the container cannot run as your user there, and the key's permissions and the passwd entry OpenSSH needs have not been worked out (ADR-0026); push from a shell instead")
	}
	path = expandHome(path)
	if !filepath.IsAbs(path) {
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}
	}
	if strings.Contains(path, ",") {
		return nil, fmt.Errorf("the push key path %s contains a comma, which docker cannot mount; move or link the key", path)
	}
	st, err := os.Stat(path)
	switch {
	case err != nil:
		return nil, fmt.Errorf("the push key %s cannot be read: %w; name an SSH private key with --push-key or flai config set dashboard.push_key, or unset it to run without one", path, err)
	case !st.Mode().IsRegular():
		return nil, fmt.Errorf("the push key %s is not a regular file", path)
	case st.Mode().Perm()&0o077 != 0:
		return nil, fmt.Errorf("the push key %s is readable by others (mode %04o); ssh refuses such a key and so does flai: chmod 600 %s", path, st.Mode().Perm(), path)
	}
	if _, err := a.runner.LookPath("ssh-keygen"); err != nil {
		return nil, fmt.Errorf("ssh-keygen is needed to check the push key and is not installed on this host")
	}
	// An empty passphrase must open it: nobody is there to type one when an
	// acceptance is pushed from the board.
	if out, err := a.runner.Run("", "ssh-keygen", "-y", "-P", "", "-f", path); err != nil {
		if strings.Contains(strings.ToLower(out+err.Error()), "passphrase") {
			return nil, fmt.Errorf("the push key %s is protected by a passphrase, so the container could never use it unattended. Forwarding your SSH agent instead is not offered: it would lend the container every key in the agent. Make a key for this repository without a passphrase and add it as a deploy key with write access (ssh-keygen -t ed25519 -N \"\" -C \"flaiover push\" -f ~/.ssh/flaiover-push; docs/operators), or push from a shell", path)
		}
		return nil, fmt.Errorf("%s is not an SSH private key ssh-keygen can read; name the private key, not the .pub file", path)
	}
	pk := &pushKey{Path: path}
	if out, err := a.runner.Run("", "ssh-keygen", "-l", "-f", path); err == nil {
		// "256 SHA256:abc… comment words (ED25519)"
		f := strings.Fields(strings.TrimSpace(out))
		if len(f) >= 2 {
			pk.Fingerprint = f[1]
		}
		if len(f) > 3 {
			pk.Comment = strings.Join(f[2:len(f)-1], " ")
		}
	}
	remote, err := a.runner.Run(repo.MainRoot, "git", "remote", "get-url", "--push", remoteName)
	remote = strings.TrimSpace(remote)
	if err != nil || remote == "" {
		return nil, fmt.Errorf("there is no git remote named %s to push to; add one, or unset dashboard.push_key", remoteName)
	}
	host, port, ok := sshRemoteHost(remote)
	if !ok {
		return nil, fmt.Errorf("the push key is an SSH key and the remote %s is %s, which is not an SSH URL; set the remote to its SSH form (git remote set-url %s git@host:owner/repo.git), or unset dashboard.push_key", remoteName, remote, remoteName)
	}
	pk.Remote, pk.Host = remote, host
	lookup := host
	if port != "" {
		lookup = "[" + host + "]:" + port
	}
	// The host's own verified entries, hashed or not; never a key fetched
	// now. The operator's file when they name one (a curated list of the
	// provider's published keys, say), else the user's known_hosts and then
	// the system's.
	searches := [][]string{{"-F", lookup}}
	if _, err := os.Stat("/etc/ssh/ssh_known_hosts"); err == nil {
		searches = append(searches, []string{"-F", lookup, "-f", "/etc/ssh/ssh_known_hosts"})
	}
	if hostsFile != "" {
		hostsFile = expandHome(hostsFile)
		if _, err := os.Stat(hostsFile); err != nil {
			return nil, fmt.Errorf("the known hosts file %s cannot be read: %w", hostsFile, err)
		}
		searches = [][]string{{"-F", lookup, "-f", hostsFile}}
	}
	var lines []string
	for _, args := range searches {
		found, _ := a.runner.Run("", "ssh-keygen", args...)
		for _, l := range strings.Split(found, "\n") {
			l = strings.TrimSpace(l)
			if l != "" && !strings.HasPrefix(l, "#") {
				lines = append(lines, l)
			}
		}
		if len(lines) > 0 {
			break
		}
	}
	if len(lines) == 0 {
		where := "this host has no known_hosts entry"
		if hostsFile != "" {
			where = hostsFile + " has no entry"
		}
		return nil, fmt.Errorf("%s for %s, and the container must not accept a host key on first use. Connect once from a shell (ssh -T git@%s), check the fingerprint it shows against the ones your provider publishes, accept it, and start the dashboard again", where, lookup, host)
	}
	if err := os.MkdirAll(filepath.Join(repo.MainRoot, ".flai-cache"), 0o755); err != nil {
		return nil, err
	}
	pk.knownHosts = knownHostsPath(repo.MainRoot)
	if err := os.WriteFile(pk.knownHosts, []byte(strings.Join(lines, "\n")+"\n"), 0o600); err != nil {
		return nil, err
	}
	// OpenSSH refuses to run for a user ID with no passwd entry ("No user
	// exists for uid"), and the container runs as the host's user, whom the
	// image need not know. Root needs none.
	if uid := os.Getuid(); uid != 0 {
		pk.passwd = passwdPath(repo.MainRoot)
		entries := fmt.Sprintf("root:x:0:0:root:/root:/sbin/nologin\nflaiover:x:%d:%d:flaiover dashboard:/tmp:/sbin/nologin\n", uid, os.Getgid())
		if err := os.WriteFile(pk.passwd, []byte(entries), 0o644); err != nil {
			return nil, err
		}
	}
	return pk, nil
}

// args are the docker run arguments that hand the key to the container.
func (pk *pushKey) args() []string {
	if pk == nil {
		return nil
	}
	out := []string{
		"--mount", "type=bind,source=" + pk.Path + ",target=" + pushKeyMountPath + ",readonly",
		"--mount", "type=bind,source=" + pk.knownHosts + ",target=" + knownHostsMountPath + ",readonly",
		"--env", "GIT_SSH_COMMAND=" + pushSSHCommand,
	}
	if pk.passwd != "" {
		out = append(out, "--mount", "type=bind,source="+pk.passwd+",target="+passwdMountPath+",readonly")
	}
	return out
}

// describe is what the operator is told the container holds, and what that
// means. The key itself is never printed.
func (pk *pushKey) describe() string {
	who := pk.Fingerprint
	if pk.Comment != "" {
		who += " (" + pk.Comment + ")"
	}
	return fmt.Sprintf("  push key: %s, mounted read-only from %s\n    acceptances made from the board are pushed to %s with it, release tags included.\n    Whoever holds the dashboard token can now publish a release by accepting a story, and a\n    compromise of the container yields this key and everything it opens. A key made for this\n    repository alone opens the least (docs/operators). Unset dashboard.push_key to stop.\n", who, pk.Path, pk.Remote)
}

// pushKeyPath is the key the operator asked for: the flag, else the host
// config. Empty means none, which is the default.
func (a *app) pushKeyPath(flag string) string {
	if flag != "" {
		return flag
	}
	cfg, _, err := a.loadConfig()
	if err != nil {
		return ""
	}
	return cfg.Dashboard.PushKey
}

// pushKnownHosts is the known hosts file the operator named, if any.
func (a *app) pushKnownHosts(flag string) string {
	if flag != "" {
		return flag
	}
	cfg, _, err := a.loadConfig()
	if err != nil {
		return ""
	}
	return cfg.Dashboard.PushKnownHosts
}

// runningPushKey says which key a running dashboard container holds, from
// its mounts: status must describe what is running, not what is configured.
func (a *app) runningPushKey(container string) (path, fingerprint string) {
	out, err := a.runner.Run("", "docker", "inspect", "--format", "{{range .Mounts}}{{if eq .Destination \""+pushKeyMountPath+"\"}}{{.Source}}{{end}}{{end}}", container)
	path = strings.TrimSpace(out)
	if err != nil || path == "" || !strings.HasPrefix(path, "/") {
		return "", ""
	}
	if l, err := a.runner.Run("", "ssh-keygen", "-l", "-f", path); err == nil {
		if f := strings.Fields(l); len(f) >= 2 {
			fingerprint = f[1]
		}
	}
	return path, fingerprint
}
