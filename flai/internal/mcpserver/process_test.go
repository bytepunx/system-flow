package mcpserver

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// processExempt names the test functions that may set the process, and why.
// Every other test of the package runs in parallel and gives its server what
// it needs through Options (I-0234). Each entry is exempt from the rules its
// flags name.
var processExempt = map[string]struct{ env, assign, parallel bool }{
	// TestMain runs before every test and once: it clears the session and
	// isolates git for the whole package.
	"TestMain": {env: true, assign: true, parallel: true},
	// isolateGit is called by TestMain alone, before m.Run.
	"isolateGit": {env: true, assign: true, parallel: true},
	// SessionFromEnv reads the process environment, so its one test must set
	// it, and cannot run beside the others.
	"TestSessionFromEnvReadsEachVariable": {env: true, parallel: true},
}

// TestNoTestSetsTheProcess keeps flai/internal/mcpserver's tests off the
// process, which made 95 of them run serially or in groups before the rest
// and the package take 136 to 162 s under -race (I-0234). It fails, naming
// file and line, on a test that sets the environment (Setenv, os.Unsetenv,
// gittest.Identity), assigns a package-level variable of this or an imported
// package, or is a top-level test that does not call t.Parallel.
func TestNoTestSetsTheProcess(t *testing.T) {
	t.Parallel()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	var tests []*ast.File
	vars := map[string]bool{}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") {
			continue
		}
		f, err := parser.ParseFile(fset, name, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range f.Decls {
			if g, ok := d.(*ast.GenDecl); ok && g.Tok == token.VAR {
				for _, s := range g.Specs {
					for _, n := range s.(*ast.ValueSpec).Names {
						if n.Name != "_" {
							vars[n.Name] = true
						}
					}
				}
			}
		}
		if strings.HasSuffix(name, "_test.go") {
			tests = append(tests, f)
		}
	}
	var found []string
	for _, f := range tests {
		imports := importNames(f)
		for _, d := range f.Decls {
			fn, ok := d.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			ex := processExempt[fn.Name.Name]
			if fn.Recv != nil {
				ex = struct{ env, assign, parallel bool }{}
			}
			report := func(n ast.Node, what string) {
				found = append(found, fset.Position(n.Pos()).String()+": "+fn.Name.Name+" "+what)
			}
			locals := declaredIn(fn)
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				switch n := n.(type) {
				case *ast.CallExpr:
					if what := setsEnv(n); what != "" && !ex.env {
						report(n, "calls "+what)
					}
				case *ast.AssignStmt:
					if n.Tok != token.DEFINE && !ex.assign {
						for _, l := range n.Lhs {
							if what := packageTarget(l, vars, imports, locals); what != "" {
								report(l, "assigns "+what)
							}
						}
					}
				case *ast.IncDecStmt:
					if what := packageTarget(n.X, vars, imports, locals); what != "" && !ex.assign {
						report(n, "assigns "+what)
					}
				}
				return true
			})
			if isTopLevelTest(fn) && !ex.parallel && !callsParallel(fn) {
				report(fn, "does not call t.Parallel()")
			}
		}
	}
	if len(found) > 0 {
		sort.Strings(found)
		t.Errorf("%d findings: a test that sets the process or does not run in parallel slows the package (I-0234); give the server what it needs through Options: %s", len(found), strings.Join(found, "\n"))
	}
}

// setsEnv names the call when it sets the process environment.
func setsEnv(c *ast.CallExpr) string {
	sel, ok := c.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	x, _ := sel.X.(*ast.Ident)
	switch {
	case sel.Sel.Name == "Setenv" && x != nil:
		return x.Name + ".Setenv"
	case sel.Sel.Name == "Setenv":
		return "Setenv"
	case sel.Sel.Name == "Unsetenv" && x != nil && x.Name == "os":
		return "os.Unsetenv"
	case sel.Sel.Name == "Identity" && x != nil && x.Name == "gittest":
		return "gittest.Identity"
	}
	return ""
}

// packageTarget names what an assignment to e changes when it is a
// package-level variable of this package, or a variable of an imported one,
// and not a local of the function: "" otherwise.
func packageTarget(e ast.Expr, vars, imports, locals map[string]bool) string {
	for {
		switch x := e.(type) {
		case *ast.ParenExpr:
			e = x.X
		case *ast.StarExpr:
			e = x.X
		case *ast.IndexExpr:
			e = x.X
		case *ast.SelectorExpr:
			if id, ok := x.X.(*ast.Ident); ok && imports[id.Name] && !locals[id.Name] {
				return id.Name + "." + x.Sel.Name
			}
			e = x.X
		case *ast.Ident:
			if vars[x.Name] && !locals[x.Name] {
				return x.Name
			}
			return ""
		default:
			return ""
		}
	}
}

// importNames are the names f refers to its imports by.
func importNames(f *ast.File) map[string]bool {
	names := map[string]bool{}
	for _, im := range f.Imports {
		p, err := strconv.Unquote(im.Path.Value)
		if err != nil {
			continue
		}
		name := path.Base(p)
		if im.Name != nil {
			name = im.Name.Name
		}
		names[name] = true
	}
	return names
}

// declaredIn is every name fn declares, anywhere in it: its receiver,
// parameters, and results, and those of its function literals, and each
// variable it defines. A name declared anywhere in fn counts as a local
// throughout it, which can only miss an assignment, never invent one.
func declaredIn(fn *ast.FuncDecl) map[string]bool {
	names := map[string]bool{}
	fields := func(l *ast.FieldList) {
		if l == nil {
			return
		}
		for _, f := range l.List {
			for _, n := range f.Names {
				names[n.Name] = true
			}
		}
	}
	fields(fn.Recv)
	fields(fn.Type.Params)
	fields(fn.Type.Results)
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			fields(n.Type.Params)
			fields(n.Type.Results)
		case *ast.AssignStmt:
			if n.Tok == token.DEFINE {
				for _, l := range n.Lhs {
					if id, ok := l.(*ast.Ident); ok {
						names[id.Name] = true
					}
				}
			}
		case *ast.RangeStmt:
			if n.Tok == token.DEFINE {
				for _, l := range []ast.Expr{n.Key, n.Value} {
					if id, ok := l.(*ast.Ident); ok {
						names[id.Name] = true
					}
				}
			}
		case *ast.ValueSpec:
			for _, id := range n.Names {
				names[id.Name] = true
			}
		}
		return true
	})
	return names
}

// isTopLevelTest says whether fn is a test go test runs: TestXxx(t *testing.T).
func isTopLevelTest(fn *ast.FuncDecl) bool {
	if fn.Recv != nil || !strings.HasPrefix(fn.Name.Name, "Test") || fn.Name.Name == "TestMain" {
		return false
	}
	ps := fn.Type.Params.List
	if len(ps) != 1 || len(ps[0].Names) != 1 {
		return false
	}
	star, ok := ps[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "T"
}

// callsParallel says whether fn calls t.Parallel() as a statement of its own
// body, on its parameter.
func callsParallel(fn *ast.FuncDecl) bool {
	t := fn.Type.Params.List[0].Names[0].Name
	for _, s := range fn.Body.List {
		es, ok := s.(*ast.ExprStmt)
		if !ok {
			continue
		}
		c, ok := es.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		if sel, ok := c.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Parallel" {
			if id, ok := sel.X.(*ast.Ident); ok && id.Name == t {
				return true
			}
		}
	}
	return false
}
