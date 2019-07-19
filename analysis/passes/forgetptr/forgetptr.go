package forgetptr

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:             "forgetptr",
	Doc:              Doc,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	Run:              run,
	RunDespiteErrors: true,
}

const (
	Doc = "forgetptr validates the code that forget to set pointer receiver."
)

type state struct {
	pass       *analysis.Pass
	target     types.Object
	isPtr      bool
	isReturned bool
	reports    []report
}

type report struct {
	pos     token.Pos
	message string
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}
	inspector.Preorder(nodeFilter, func(n ast.Node) {
		s := &state{
			pass:    pass,
			reports: []report{},
		}
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv == nil {
				return
			}
			s.findRecv(x.Recv)
			if s.target == nil {
				return
			}
			s.checkStmt(x.Body)

			if !s.isReturned {
				for _, repo := range s.reports {
					s.pass.Reportf(repo.pos, repo.message)
				}
			}
		}
	})
	return nil, nil
}

func (s *state) checkStmt(stmt ast.Stmt) {
	switch x := stmt.(type) {
	case *ast.BlockStmt:
		for _, ss := range x.List {
			s.checkStmt(ss)
		}
	case *ast.AssignStmt:
		for _, expr := range x.Lhs {
			switch x := expr.(type) {
			case *ast.SelectorExpr:
				s.checkSelector(x)
			}
		}
	case *ast.IncDecStmt:
		expr, ok := x.X.(*ast.SelectorExpr)
		if !ok {
			return
		}
		s.checkSelector(expr)

	case *ast.ReturnStmt:
		for _, expr := range x.Results {
			id, ok := expr.(*ast.Ident)
			if !ok {
				continue
			}
			if s.pass.TypesInfo.ObjectOf(id).Type() == s.target.Type() {
				s.isReturned = true
			}
		}
	}
}

func (s *state) checkSelector(x *ast.SelectorExpr) {

	id, ok := x.X.(*ast.Ident)
	if !ok {
		return
	}
	v := s.pass.TypesInfo.ObjectOf(id)
	if !s.isPtr && v.Type() == s.target.Type() {
		s.reports = append(s.reports, report{x.Pos(), "this statement can not modify the value"})
	}
}

func (s *state) findRecv(recv *ast.FieldList) {
	f := recv.List[0]
	// maybe always length is #1
	if len(f.Names) == 0 {
		return
	}

	switch x := f.Type.(type) {
	case *ast.StarExpr:
		s.isPtr = true
		id, ok := x.X.(*ast.Ident)
		if !ok {
			return
		}
		s.target = s.pass.TypesInfo.ObjectOf(id)
	case *ast.Ident:
		s.target = s.pass.TypesInfo.ObjectOf(x)
	}
}
