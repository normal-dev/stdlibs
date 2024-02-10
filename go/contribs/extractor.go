package main

import (
	"contribs-go/model"

	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
)

type extractor struct {
	Error error

	fset  *token.FileSet
	nodes ast.Node

	info *types.Info
}

func NewExtractor(src []byte) *extractor {
	var ex = &extractor{}

	srcFile, fset, err := parse(src)
	if err != nil {
		ex.Error = err
		return ex
	}

	ex.fset = fset
	ex.nodes = srcFile

	conf := types.Config{Importer: importer.ForCompiler(fset, "source", nil)}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}
	defer func() {
		if r := recover(); r != nil {
			switch typ := r.(type) {
			case error:
				ex.Error = typ

			case string:
				ex.Error = errors.New(typ)
			}
		}
	}()
	_, err = conf.Check("main", fset, []*ast.File{srcFile}, info)
	if err != nil {
		ex.Error = err
		return ex
	}
	ex.info = info

	return ex
}

func (ex *extractor) Extract() map[model.API]struct{} {
	var apis = make(map[model.API]struct{})

	if ex.Error != nil || ex.info == nil {
		return apis
	}

	var (
		newAPI = func(pos token.Pos, imporSpec *ast.ImportSpec, idents ...string) model.API {
			tokPos := ex.fset.Position(pos)
			pkg := strings.Trim(imporSpec.Path.Value, "\"")
			return model.API{
				Ident: fmt.Sprintf("%s.%s", pkg, strings.Join(idents, ".")),
				Line:  tokPos.Line,
			}
		}

		resolveSelExpr = func(expr *ast.SelectorExpr) (*ast.Ident, *ast.Ident) {
			switch typ := expr.X.(type) {
			case *ast.Ident:
				return typ, expr.Sel
			}

			return nil, expr.Sel
		}
	)

	for expr := range ex.info.Types {
		switch typ := expr.(type) {
		case *ast.SelectorExpr:
			x, sel := resolveSelExpr(typ)
			switch x {
			case nil:
				continue
			}

			switch x.Obj {
			// Exclude references
			case nil:
				imporSpec := ex.findImport(x)
				switch imporSpec {
				case nil:
					// Reference must be in another file
					continue
				}

				if !isStdImport(imporSpec) {
					continue
				}

				api := newAPI(sel.Pos(), imporSpec, sel.Name)
				apis[api] = struct{}{}
			}
		}
	}

	return apis
}

func (ex *extractor) findImport(x *ast.Ident) *ast.ImportSpec {
	var impSpec *ast.ImportSpec

	var (
		eqPath = func(lit *ast.BasicLit, name string) bool {
			if lit == nil {
				return false
			}
			val := strings.Trim(lit.Value, "\"")
			spl := strings.Split(val, "/")
			if len(spl) > 1 {
				val = spl[len(spl)-1]
			}
			return val == name
		}
		eqName = func(ident *ast.Ident, name string) bool {
			if ident == nil {
				return false
			}
			return ident.Name == name
		}
	)

	ast.Inspect(ex.nodes, func(n ast.Node) bool {
		switch typ := n.(type) {
		case *ast.ImportSpec:
			// Alias
			if eqPath(typ.Path, x.Name) {
				impSpec = typ
				return false
			}

			// Normal
			if eqName(typ.Name, x.Name) {
				impSpec = typ
				return false
			}
		}
		return true
	})

	return impSpec
}

func isStdImport(importSpec *ast.ImportSpec) bool {
	_, ok := gopkgs[strings.Trim(importSpec.Path.Value, "\"")]
	return ok
}

func parse(src []byte) (*ast.File, *token.FileSet, error) {
	fset := token.NewFileSet()
	srcFile, err := parser.ParseFile(fset, "main.go", src, 0)
	if err != nil {
		return nil, nil, err
	}
	return srcFile, fset, nil
}
