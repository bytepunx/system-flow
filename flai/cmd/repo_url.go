package cmd

import (
	"fmt"
	"net/url"
	"strings"
)

// The template's repo_url is the repository's web address (S-0117): an
// import offers the origin remote made one, and a value that is not one is
// refused, rather than written to system-flow.yaml as typed.

// repoURLVar is the template variable that holds the repository's address.
const repoURLVar = "repo_url"

// webURL is the https address of a git remote: git@host:owner/repo.git,
// ssh://git@host[:port]/owner/repo.git, git://host/owner/repo, and
// http(s)://[user@]host/owner/repo.git all become https://host/owner/repo,
// without credentials or port. A local path or file:// remote has none: "".
func webURL(remote string) string {
	remote = strings.TrimSpace(remote)
	var host, path string
	if !strings.Contains(remote, "://") {
		// scp-like, [user@]host:path, as git reads it: no slash before the colon
		at, rest, ok := strings.Cut(remote, ":")
		if !ok || strings.Contains(at, "/") {
			return ""
		}
		if _, h, ok := strings.Cut(at, "@"); ok {
			at = h
		}
		host, path = at, rest
	} else if u, err := url.Parse(remote); err == nil {
		switch u.Scheme {
		case "https", "http", "ssh", "git", "git+ssh", "ssh+git":
			host, path = u.Hostname(), u.Path
		}
	}
	path = strings.TrimSuffix(strings.Trim(path, "/"), ".git")
	if host == "" || path == "" {
		return ""
	}
	return "https://" + host + "/" + path
}

// checkRepoURL refuses a repo_url that is not an http or https address with a
// host and a path; empty is allowed, as the template does not require one.
func checkRepoURL(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	const want = "give the repository's web address, such as https://github.com/owner/repo"
	u, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("%s %q is not a URL (%s); %s", repoURLVar, s, strings.TrimPrefix(err.Error(), "parse "+`"`+s+`": `), want)
	}
	switch {
	case u.Scheme != "http" && u.Scheme != "https":
		return fmt.Errorf("%s %q is not an http or https URL; %s", repoURLVar, s, want)
	case u.Host == "" || u.Hostname() == "":
		return fmt.Errorf("%s %q has no host; %s", repoURLVar, s, want)
	case strings.Trim(u.Path, "/") == "":
		return fmt.Errorf("%s %q names no repository on %s; %s", repoURLVar, s, u.Host, want)
	}
	return nil
}

// originURL is the web address of the origin remote of the git repository at
// root, or "" when it has none that makes one.
func (a *app) originURL(root string) string {
	out, err := a.runner.Run(root, "git", "remote", "get-url", "origin")
	if err != nil {
		return ""
	}
	return webURL(out)
}
