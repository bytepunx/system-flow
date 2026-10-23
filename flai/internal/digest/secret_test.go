package digest

import (
	"strings"
	"testing"
)

// ADR-0191: each shape of a key or token is found, once, on its line, and
// described without its value.
func TestSecretsFindsEachShape(t *testing.T) {
	gh := "ghp_" + strings.Repeat("a1B2", 9)
	for _, c := range []struct{ text, kind, value string }{
		{"-----BEGIN RSA PRIVATE KEY-----", "private-key", ""},
		{"-----BEGIN PRIVATE KEY-----", "private-key", ""},
		{"-----BEGIN OPENSSH PRIVATE KEY-----", "private-key", ""},
		{"use sk-ant-api03-AbCdEfGh1234567890IjKlMnOp_qr here", "api-key", "sk-ant-api03-AbCdEfGh1234567890IjKlMnOp_qr"},
		{"OPENAI_API_KEY=sk-proj-Ab12Cd34Ef56Gh78Ij90Kl", "api-key", "sk-proj-Ab12Cd34Ef56Gh78Ij90Kl"},
		{"token " + gh, "github-token", gh},
		{"github_pat_11ABCDEFG0123456789_abcdefghijklmnop", "github-token", "github_pat_11ABCDEFG0123456789_abcdefghijklmnop"},
		{"xoxb-1234567890-abcdefghij", "slack-token", "xoxb-1234567890-abcdefghij"},
		{"AKIAIOSFODNN7EXAMPLE", "aws-access-key", "AKIAIOSFODNN7EXAMPLE"},
		{"ASIAABCDEFGHIJ012345", "aws-access-key", "ASIAABCDEFGHIJ012345"},
		{"Bearer eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U", "jwt", "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"},
		{"password: Tr0ub4dor3xQ9zLmK7pWvB2", "assigned-secret", "Tr0ub4dor3xQ9zLmK7pWvB2"},
		{"api_token = 3f9a6c1e8b2d4f7a0c5e9b1d3f6a8c2e", "assigned-secret", "3f9a6c1e8b2d4f7a0c5e9b1d3f6a8c2e"},
		{`{"client_secret": "c2VjcmV0LXZhbHVlLWhlcmUxMjM0NQ=="}`, "assigned-secret", "c2VjcmV0LXZhbHVlLWhlcmUxMjM0NQ=="},
		{"DB_PASSWORD=abcd1234efgh5678ijkl", "assigned-secret", "abcd1234efgh5678ijkl"},
	} {
		got := Secrets("first line\n" + c.text + "\n")
		if len(got) != 1 {
			t.Errorf("%q: want one match, got %+v", c.text, got)
			continue
		}
		m := got[0]
		if m.Line != 2 || m.Kind != c.kind {
			t.Errorf("%q: got %+v, want line 2 kind %s", c.text, m, c.kind)
		}
		if c.value != "" && (strings.Contains(m.Description, c.value) || !strings.Contains(m.Description, c.value[:4]+"…")) {
			t.Errorf("%q: description %q shows the value, or not its first characters", c.text, m.Description)
		}
	}
}

// ADR-0191: what only looks near a secret passes: the word in prose, an
// environment variable's name, a path, a short value, a slug, a public key,
// and a commit hash given to no key or token.
func TestSecretsPassesNearMisses(t *testing.T) {
	for _, text := range []string{
		"the token is in ~/.flai/config.json",
		"FLAI_TOKEN=$(cat ~/.flai/token)",
		"set FLAI_TOKEN in the environment, as FLAI_AGENT is",
		"token: FLAI_TOKEN_FROM_THE_ENVIRONMENT",
		"key: design/system/agent-memory.md",
		"api_key_path: flai/internal/digest/secret.go",
		"password: hunter2",
		"secret_key: s-0428-the-digest-entry-name-for-priming",
		"commit 3f9a6c1e8b2d4f7a0c5e9b1d3f6a8c2e4b7d9f01 fixed it",
		"sk-this-is-not-a-key-at-all-really",
		"the task-ab12cd34ef56gh78ij90kl branch",
		"ghp_short and github_pat_tooshort",
		"xoxb- is how a bot token starts; xoxb-abcdefghijkl has no digit",
		"AKIAEXAMPLE is too short",
		"eyJhbGciOiJIUzI1NiJ9.short",
		"-----BEGIN PUBLIC KEY-----",
		"monkey: the keyboard shortcut",
	} {
		if got := Secrets(text); len(got) != 0 {
			t.Errorf("%q: want no match, got %+v", text, got)
		}
	}
}
