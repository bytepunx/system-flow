package guard

import (
	"os"
	"path/filepath"
	"testing"
)

// I-0052, S-0368, ADR-0141: written finds the file a simple command writes
// outside the temporary folder, through a redirection or a program that
// writes its operands, past wrappers and in a shell's script; a file under
// /tmp or /dev, a descriptor duplication, and a read write none.
func TestWrittenFindsAShellCommandsWrite(t *testing.T) {
	for _, c := range []struct{ line, want string }{
		// redirections, attached or apart
		{"echo x > f", "f"},
		{"echo x>f", "f"},
		{"echo x >> f", "f"},
		{"echo x >|f", "f"},
		{"cmd &> f", "f"},
		{"cmd &>>f", "f"},
		{"cmd 2>err.txt", "err.txt"},
		{"cmd >& out", "out"},
		{"cmd > /tmp/../home/x", "/tmp/../home/x"},
		{"cmd > $HOME/x", "HOME/x"},
		{"cmd 2>&1", ""},
		{"cmd >&2", ""},
		{"cmd 2>&-", ""},
		{"cmd > /tmp/out 2>/dev/null", ""},
		{"cmd < in.txt", ""},
		{"cat <<< text", ""},
		{"echo '>' x", ""},
		{"echo $((1>2))", ""},
		// programs that write their operands
		{"tee -a log.txt", "log.txt"},
		{"tee /tmp/a", ""},
		{"sed -i s/a/b/ x.go", "x.go"},
		{"sed -i.bak -e s/a/b/ x.go", "x.go"},
		{"sed -ni s/a/b/p x.go", "x.go"},
		{"sed --in-place=.bak s/a/b/ x.go", "x.go"},
		{"sed -n s/a/b/p x.go", ""},
		{"sed --posix -n s/a/b/p x.go", ""},
		{"perl -i -pe s/a/b/ x.go", "x.go"},
		{"perl -pi.bak -e s/a/b/ x.go", "x.go"},
		{"perl -ne print x.go", ""},
		{"perl -Ilib -i script.pl x.go", "x.go"},
		{"cp a.go b.go", "b.go"},
		{"cp a.go /tmp/b.go", ""},
		{"cp -t dir a b", "dir"},
		{"cp --target-directory=/tmp a b", ""},
		{"cp a /tmp", ""},
		{"cp a /tmpx", "/tmpx"},
		{"install -m 644 a /usr/local/bin/a", "/usr/local/bin/a"},
		{"ln -s /tmp/x", "."},
		{"mv /tmp/a /tmp/b", ""},
		{"mv x.go /tmp/x.go", "x.go"},
		{"rm -rf /tmp/scratch", ""},
		{"rm -r build", "build"},
		{"touch -r ref.go new.go", "new.go"},
		{"mkdir -p -m 755 out", "out"},
		{"rmdir old", "old"},
		{"truncate -s 0 log", "log"},
		{"patch -p1 < fix.diff", patchNames},
		{"patch -p1 --dry-run < fix.diff", ""},
		{"patch x.go fix.diff", "x.go"},
		{"patch -o /tmp/x.go x.go fix.diff", ""},
		{"dd if=/dev/zero of=disk.img", "disk.img"},
		{"dd if=x of=/dev/null", ""},
		{"gofmt -l -w .", "."},
		{"gofmt -l .", ""},
		{"goimports -w x.go", "x.go"},
		{"go fmt ./...", "./..."},
		{"go fmt", "."},
		{"go vet ./...", ""},
		{"prettier --write src", "src"},
		{"prettier --check src", ""},
		{"eslint --fix", "."},
		{"golangci-lint run --fix ./...", "./..."},
		{"golangci-lint run ./...", ""},
		// the program found past assignments and wrappers, and in scripts
		{"FOO=1 env -u X sudo -u me timeout 30 rm x", "x"},
		{"xargs -I {} cp {} out/", "out/"},
		{"find . -name '*.go' -exec gofmt -w {} +", "{}"},
		{"find . -name '*.go' -exec grep -l x {} +", ""},
		{"bash -c 'cp a b'", "b"},
		{"sh -ec 'go test > out'", "out"},
		{"eval rm x", "x"},
		{"grep -r rm .", ""},
		{"echo cp a b", ""},
	} {
		got := ""
		for _, s := range parse(c.line) {
			if got = written(s); got != "" {
				break
			}
		}
		if got != c.want {
			t.Errorf("%q writes %q, want %q", c.line, got, c.want)
		}
	}
}

// S-0368, ADR-0141: a sub-agent is read-only when its definition's tools
// line names neither Edit nor Write, or when it is Explore or Plan; one with
// Edit, *, or no tools line, one with no definition, and a type that is not
// a file name are not.
func TestReadOnlyAgentsReadsTheDefinitions(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".claude", "agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, tools := range map[string]string{
		"verifier": "tools: Read, Grep, Glob, Bash, mcp__flai__prime\n",
		"listed":   "tools: [Read, Bash]\n",
		"editor":   "tools: Read, Edit, Bash\n",
		"multi":    "tools: Read, MultiEdit\n",
		"writer":   `tools: ["Read", "Write"]` + "\n",
		"all":      "tools: *\n",
		"open":     "",
		"block":    "tools:\n  - Read\n",
	} {
		doc := "---\nname: " + name + "\ndescription: x\n" + tools + "model: haiku\n---\n\nBody: tools: Read\n"
		if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte(doc), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, outside := range []string{filepath.Join(root, ".claude", "x.md"), filepath.Join(root, "x.md")} {
		if err := os.WriteFile(outside, []byte("---\ntools: Read\n---\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	readOnly := ReadOnlyAgents(root)
	for agent, want := range map[string]bool{
		"verifier": true, "listed": true, "Explore": true, "Plan": true,
		"editor": false, "multi": false, "writer": false, "all": false, "open": false, "block": false,
		"missing": false, "general-purpose": false, "": false, "../x": false, "../../x": false, "a/b": false,
	} {
		if got := readOnly(agent); got != want {
			t.Errorf("%q read-only %v, want %v", agent, got, want)
		}
	}
	none := ReadOnlyAgents("")
	if !none("Explore") || none("verifier") {
		t.Errorf("with no project only Explore and Plan are read-only")
	}
}
