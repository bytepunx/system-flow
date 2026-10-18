package hostapi

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// minParamsSites is fewer params sites than the package has, by a margin:
// finding fewer means the test no longer sees how methods decode their
// params, and would pass by finding nothing.
const minParamsSites = 80

// paramsScan is what TestNoMethodTakesAParameterNamedProject reads of the
// package: its files, its named types, and the variables its functions
// declare, to find each type a channel method's params are decoded into.
type paramsScan struct {
	fset  *token.FileSet
	files []*ast.File
	// types are the package's named types, at any level, by name.
	types map[string]*ast.TypeSpec
	// vars are the variables declared in function bodies, by name, each
	// with the type it was declared with, in the order of their positions.
	vars map[string][]declared
}

// declared is one variable a function body declares, its type, and the
// innermost function that declares it, which is its scope.
type declared struct {
	pos   token.Pos
	typ   ast.Expr
	scope ast.Node
}

// I-0163: project in a call's params is the channel's project key alone,
// which the dashboard stamps on every call and the channel checks against
// the key it serves, so a method that reads project from its params gets
// the key, not its own argument. This reads the package's source for every
// type a method's params are decoded into, by decode[T] or params(raw, &v),
// and fails on any field whose JSON name is project, in any case, since
// encoding/json matches names without regard to case.
func TestNoMethodTakesAParameterNamedProject(t *testing.T) {
	s := scanPackage(t)
	sites := 0
	offending := map[string][]string{} // field's position: the sites that decode it
	var order []string
	for _, f := range s.files {
		for _, d := range f.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok && fd.Name.Name == "decode" {
				// decode's own params(raw, &v) is of its type parameter: its
				// callers' type arguments are checked instead
				continue
			}
			ast.Inspect(d, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				typ, found := s.paramsType(call)
				if !found {
					return true
				}
				sites++
				at := s.at(call.Pos())
				if typ == nil {
					t.Errorf("%s: cannot tell what type these params are decoded into; teach this test the form", at)
					return true
				}
				what := "the params decoded at " + at
				if id, ok := typ.(*ast.Ident); ok {
					what = id.Name + ", decoded at " + at
				}
				for _, field := range s.jsonFields(t, typ, at, map[string]bool{}) {
					if strings.EqualFold(field.name, "project") {
						where := s.at(field.pos)
						if _, ok := offending[where]; !ok {
							order = append(order, where)
						}
						offending[where] = append(offending[where], fmt.Sprintf("%s (json %q)", what, field.name))
					}
				}
				return true
			})
		}
	}
	for _, where := range order {
		t.Errorf("%s: a field named project in the params of %s: project is the channel's project key, which the dashboard stamps on every call; name the parameter otherwise, as the components methods take sub_project (I-0163)",
			where, strings.Join(offending[where], "; "))
	}
	if sites < minParamsSites {
		t.Errorf("found %d places that decode a method's params, want at least %d: the test no longer sees how they are decoded", sites, minParamsSites)
	}
}

// scanPackage parses the package's own Go files, its tests left out, and
// gathers its named types and the variables its functions declare.
func scanPackage(t *testing.T) *paramsScan {
	t.Helper()
	s := &paramsScan{fset: token.NewFileSet(), types: map[string]*ast.TypeSpec{}, vars: map[string][]declared{}}
	names, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(s.fset, name, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatal(err)
		}
		s.files = append(s.files, f)
	}
	if len(s.files) == 0 {
		t.Fatal("no Go files in the package's folder")
	}
	for _, f := range s.files {
		var stack []ast.Node // the nodes from the file down to the one inspected
		ast.Inspect(f, func(n ast.Node) bool {
			if n == nil {
				stack = stack[:len(stack)-1]
				return true
			}
			stack = append(stack, n)
			switch n := n.(type) {
			case *ast.TypeSpec:
				s.types[n.Name.Name] = n
			case *ast.ValueSpec:
				for i, name := range n.Names {
					typ := n.Type
					if typ == nil && i < len(n.Values) {
						typ = literalType(n.Values[i])
					}
					s.declare(stack, name, typ)
				}
			case *ast.AssignStmt:
				if n.Tok == token.DEFINE && len(n.Lhs) == len(n.Rhs) {
					for i, l := range n.Lhs {
						if id, ok := l.(*ast.Ident); ok {
							s.declare(stack, id, literalType(n.Rhs[i]))
						}
					}
				}
			}
			return true
		})
	}
	return s
}

// declare records name, declared with typ inside the innermost function in
// stack; one declared outside any function, or with a type this cannot
// tell, is left out.
func (s *paramsScan) declare(stack []ast.Node, name *ast.Ident, typ ast.Expr) {
	if typ == nil {
		return
	}
	for i := len(stack) - 1; i >= 0; i-- {
		switch stack[i].(type) {
		case *ast.FuncDecl, *ast.FuncLit:
			s.vars[name.Name] = append(s.vars[name.Name], declared{name.Pos(), typ, stack[i]})
			return
		}
	}
}

// literalType is the type of x when it is a composite literal, T{...} or
// &T{...}, and nil otherwise.
func literalType(x ast.Expr) ast.Expr {
	if u, ok := x.(*ast.UnaryExpr); ok && u.Op == token.AND {
		x = u.X
	}
	if lit, ok := x.(*ast.CompositeLit); ok {
		return lit.Type
	}
	return nil
}

// paramsType is the type call decodes a method's params into, and whether
// call decodes params at all: decode[T](raw) decodes into T, and
// params(raw, &v) into v's type, or the composite literal's. The type is
// nil when the call decodes params but this cannot tell into what.
func (s *paramsScan) paramsType(call *ast.CallExpr) (ast.Expr, bool) {
	switch fun := call.Fun.(type) {
	case *ast.IndexExpr:
		if id, ok := fun.X.(*ast.Ident); ok && id.Name == "decode" {
			return fun.Index, true
		}
	case *ast.Ident:
		if fun.Name != "params" || len(call.Args) != 2 {
			return nil, false
		}
		ref, ok := call.Args[1].(*ast.UnaryExpr)
		if !ok || ref.Op != token.AND {
			return nil, true
		}
		switch x := ref.X.(type) {
		case *ast.CompositeLit:
			return x.Type, true
		case *ast.Ident:
			// the nearest declaration of that name before the call, in a
			// function that holds the call
			var typ ast.Expr
			for _, d := range s.vars[x.Name] {
				if d.pos < call.Pos() && d.scope.Pos() <= call.Pos() && call.End() <= d.scope.End() {
					typ = d.typ
				}
			}
			return typ, true
		}
		return nil, true
	}
	return nil, false
}

// jsonField is one field a params type decodes, by its JSON name.
type jsonField struct {
	name string
	pos  token.Pos
}

// jsonFields are the fields encoding/json decodes into typ by name: each
// exported field by its tag's name, or its Go name without one, and the
// fields of an untagged embedded struct as its own. A type this cannot
// read fails the test, naming at, the site that decodes it.
func (s *paramsScan) jsonFields(t *testing.T, typ ast.Expr, at string, seen map[string]bool) []jsonField {
	t.Helper()
	switch x := typ.(type) {
	case *ast.StarExpr:
		return s.jsonFields(t, x.X, at, seen)
	case *ast.Ident:
		ts, ok := s.types[x.Name]
		if !ok {
			t.Errorf("%s: params of type %s, which is not declared in the package; teach this test to read it", at, x.Name)
			return nil
		}
		if seen[x.Name] {
			return nil
		}
		seen[x.Name] = true
		return s.jsonFields(t, ts.Type, at, seen)
	case *ast.StructType:
		var out []jsonField
		for _, f := range x.Fields.List {
			name, skip := "", false
			if f.Tag != nil {
				tag, err := strconv.Unquote(f.Tag.Value)
				if err != nil {
					t.Errorf("%s: %v", s.at(f.Tag.Pos()), err)
					continue
				}
				tagged := reflect.StructTag(tag).Get("json")
				if tagged == "-" {
					skip = true
				}
				name, _, _ = strings.Cut(tagged, ",")
			}
			if skip {
				continue
			}
			if len(f.Names) == 0 {
				if name == "" {
					out = append(out, s.jsonFields(t, f.Type, at, seen)...)
				} else {
					out = append(out, jsonField{name, f.Pos()})
				}
				continue
			}
			for _, n := range f.Names {
				if !n.IsExported() {
					continue
				}
				if name != "" {
					out = append(out, jsonField{name, n.Pos()})
				} else {
					out = append(out, jsonField{n.Name, n.Pos()})
				}
			}
		}
		return out
	}
	t.Errorf("%s: params of a type this test cannot read, %T; teach it to", at, typ)
	return nil
}

// at is pos as file:line.
func (s *paramsScan) at(pos token.Pos) string {
	p := s.fset.Position(pos)
	return fmt.Sprintf("%s:%d", p.Filename, p.Line)
}
