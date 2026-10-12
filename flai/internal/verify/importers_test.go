package verify

import (
	"reflect"
	"testing"
	"testing/fstest"
)

// moduleFS is a module example.com/m under mod, shaped as S-0338 found
// flai: b builds on a's type, c imports only b, and d imports a only in an
// external test package, so a change to a can break c and d though neither
// imports a in its own files. e imports other modules, f does not parse, and
// the rest import a from where the go tool does not look or from a nested
// module.
func moduleFS() fstest.MapFS {
	fsys := fstest.MapFS{}
	for _, f := range [][2]string{
		{"mod/go.mod", "module example.com/m\n\ngo 1.26\n"},
		{"mod/main.go", "package main\n\nimport \"fmt\"\n\nfunc main() { fmt.Println() }\n"},
		{"mod/a/a.go", "package a\n\n// T is a value.\ntype T struct{ N int }\n"},
		{"mod/b/b.go", "package b\n\nimport \"example.com/m/a\"\n\n// New is a T.\nfunc New() a.T { return a.T{N: 1} }\n"},
		{"mod/c/c.go", "package c\n\nimport \"example.com/m/b\"\n\n// N is b's number.\nfunc N() int { return b.New().N }\n"},
		{"mod/d/d.go", "package d\n\nimport \"strings\"\n\n// Up is s in capitals.\nfunc Up(s string) string { return strings.ToUpper(s) }\n"},
		{"mod/d/d_test.go", "package d_test\n\nimport (\n\t\"testing\"\n\n\t\"example.com/m/a\"\n)\n\nfunc TestT(t *testing.T) { _ = a.T{} }\n"},
		{"mod/e/e.go", "package e\n\nimport (\n\t\"example.com/mother/a\"\n\t\"github.com/x/a\"\n)\n"},
		{"mod/f/f.go", "package f\n\nimport (\n\t\"example.com/m/a\"\n\tbroken\n)\n"},
		{"mod/a/testdata/t.go", "package t\n\nimport \"example.com/m/a\"\n"},
		{"mod/_x/x.go", "package x\n\nimport \"example.com/m/a\"\n"},
		{"mod/.y/y.go", "package y\n\nimport \"example.com/m/a\"\n"},
		{"mod/nested/go.mod", "module example.com/m/nested\n"},
		{"mod/nested/n.go", "package nested\n\nimport \"example.com/m/a\"\n"},
		{"mod/nested/deep/k.go", "package deep\n\nimport \"example.com/m/a\"\n"},
		{"mod/root/r.go", "package root\n\nimport \"example.com/m\"\n"},
	} {
		fsys[f[0]] = &fstest.MapFile{Data: []byte(f[1])}
	}
	return fsys
}

// I-0128: a change to a package reaches every package of its module that
// imports it, directly, through others, or in a test file only; nothing the
// go tool ignores, no nested module, and no import of another module.
func TestImportersCloseOverTheModulesReverseImports(t *testing.T) {
	fsys := moduleFS()
	for _, tc := range []struct {
		name string
		dir  string
		pkgs []string
		want []string
	}{
		{name: "a's importers, through b and in d's test file", dir: "mod", pkgs: []string{"./a"}, want: []string{"./a", "./b", "./c", "./d"}},
		{name: "b's importers", dir: "mod", pkgs: []string{"./b"}, want: []string{"./b", "./c"}},
		{name: "a package no other imports", dir: "mod", pkgs: []string{"./c"}, want: []string{"./c"}},
		{name: "the module's root package", dir: "mod", pkgs: []string{"."}, want: []string{".", "./root"}},
		{name: "two packages give each importer once", dir: "mod", pkgs: []string{"./a", "./b"}, want: []string{"./a", "./b", "./c", "./d"}},
		{name: "no packages give none", dir: "mod", pkgs: nil, want: nil},
		{name: "no go.mod gives the packages", dir: "mod/a", pkgs: []string{"."}, want: []string{"."}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := importers(fsys, tc.dir, tc.pkgs); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("importers(%q, %q) = %q, want %q", tc.dir, tc.pkgs, got, tc.want)
			}
		})
	}
}

// A module at the checkout's root names its packages relative to the root.
func TestImportersOfAModuleAtTheRoot(t *testing.T) {
	fsys := fstest.MapFS{}
	for name, f := range moduleFS() {
		fsys[name[len("mod/"):]] = f
	}
	want := []string{"./a", "./b", "./c", "./d"}
	if got := importers(fsys, "", []string{"./a"}); !reflect.DeepEqual(got, want) {
		t.Errorf("importers = %q, want %q", got, want)
	}
}
