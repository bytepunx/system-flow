package verify

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path"
	"slices"
	"strconv"
	"strings"
)

// importers are the packages pkgs names and every package of the Go module
// in dir that imports one of them, directly or through other packages, in its
// own files or its test files, each as ./x/y or ., relative to dir, sorted;
// pkgs as they are when dir holds no go.mod naming its module (I-0128).
func importers(fsys fs.FS, dir string, pkgs []string) []string {
	mod := modulePath(env{dir: dir, fsys: fsys})
	if mod == "" {
		return pkgs
	}
	reverse := reverseImports(fsys, dir, mod)
	seen := map[string]bool{}
	var next []string
	for _, p := range pkgs {
		next = append(next, strings.TrimPrefix(p, "./"))
	}
	var out []string
	for ; len(next) > 0; next = next[1:] {
		d := next[0]
		if seen[d] {
			continue
		}
		seen[d] = true
		out = append(out, packageArg(d))
		next = append(next, reverse[d]...)
	}
	slices.Sort(out)
	return out
}

// reverseImports maps each package of the module mod in dir, relative to
// dir, to the packages of the module that import it in their own files or
// their test files. It skips the directories the go tool ignores and those
// that hold a module of their own, and ignores build constraints. A file
// that does not parse, or a directory that cannot be read, counts as
// importing nothing rather than failing, so that selection never fails: the
// tier's own command reports the broken file.
func reverseImports(fsys fs.FS, dir, mod string) map[string][]string {
	root := fsysDir(dir)
	reverse := map[string][]string{}
	_ = fs.WalkDir(fsys, root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			// only a directory that cannot be read comes with an error:
			// leave it out
			return fs.SkipDir
		}
		if d.IsDir() {
			if p != root && (ignoredDir(d.Name()) || isFile(fsys, path.Join(p, "go.mod"))) {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(p, ".go") {
			return nil
		}
		importer := relDir(root, path.Dir(p))
		for _, imported := range moduleImports(fsys, p, mod) {
			if imported != importer {
				reverse[imported] = append(reverse[imported], importer)
			}
		}
		return nil
	})
	return reverse
}

// moduleImports are the packages of the module mod that the .go file at p
// imports, each relative to the module's root; none when it does not parse.
func moduleImports(fsys fs.FS, p, mod string) []string {
	src, err := fs.ReadFile(fsys, p)
	if err != nil {
		return nil
	}
	f, err := parser.ParseFile(token.NewFileSet(), p, src, parser.ImportsOnly)
	if err != nil {
		return nil
	}
	var out []string
	for _, spec := range f.Imports {
		imp, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		if imp == mod {
			out = append(out, ".")
		} else if rest, ok := strings.CutPrefix(imp, mod+"/"); ok {
			out = append(out, rest)
		}
	}
	return out
}

// ignoredDir reports whether the go tool ignores a directory of this name.
func ignoredDir(name string) bool {
	return name == "testdata" || strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}

// isFile reports whether p names a file of fsys.
func isFile(fsys fs.FS, p string) bool {
	info, err := fs.Stat(fsys, p)
	return err == nil && !info.IsDir()
}

// fsysDir is a tier's cleaned dir as fsys names it: "." for the root.
func fsysDir(dir string) string {
	if dir == "" {
		return "."
	}
	return dir
}

// relDir is the directory d of fsys relative to root, "." for root itself.
func relDir(root, d string) string {
	if root == "." {
		return d
	}
	if d == root {
		return "."
	}
	return strings.TrimPrefix(d, root+"/")
}

// packageArg is a package directory relative to a tier's dir as the go tool
// takes it: ./x/y, or . for the dir itself.
func packageArg(d string) string {
	if d == "." {
		return "."
	}
	return "./" + d
}
