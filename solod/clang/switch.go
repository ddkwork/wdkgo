package clang

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
)

func (g *Generator) emitSwitchStmt(stmt *ast.SwitchStmt) {
	w := g.state.writer
	if stmt.Init != nil {
		fmt.Fprintf(w, "%s{\n", g.indent())
		g.state.indent++
		ast.Walk(g, stmt.Init)
		g.emitSwitchBody(w, stmt)
		g.state.indent--
		fmt.Fprintf(w, "%s}\n", g.indent())
	} else {
		g.emitSwitchBody(w, stmt)
	}
}

func (g *Generator) emitSwitchBody(w io.Writer, stmt *ast.SwitchStmt) {
	var cases []*ast.CaseClause
	var def *ast.CaseClause
	for _, s := range stmt.Body.List {
		cc := s.(*ast.CaseClause)
		if cc.List == nil {
			def = cc
		} else {
			cases = append(cases, cc)
		}
	}

	if len(cases) == 0 && def == nil {
		fmt.Fprintf(w, "%sif (false) {\n%s}\n", g.indent(), g.indent())
		return
	}

	if len(cases) == 0 {
		g.walkStmts(def.Body)
		return
	}

	isString := stmt.Tag != nil && g.hasStringType(stmt.Tag)
	if isString || g.isKernelMode() {
		g.emitSwitchAsIfChain(w, stmt, cases, def)
		return
	}

	fmt.Fprintf(w, "%sswitch (", g.indent())
	if stmt.Tag == nil {
		fmt.Fprintf(w, "1")
	} else {
		g.emitExpr(stmt.Tag)
	}
	fmt.Fprintf(w, ") {\n")

	for _, cc := range cases {
		for _, expr := range cc.List {
			fmt.Fprintf(w, "%scase ", g.indent())
			g.emitExpr(expr)
			fmt.Fprintf(w, ":\n")
		}
		g.state.indent++
		g.walkStmts(cc.Body)
		if !g.endsWithBreakOrReturn(cc.Body) {
			fmt.Fprintf(w, "%sbreak;\n", g.indent())
		}
		g.state.indent--
	}

	if def != nil {
		fmt.Fprintf(w, "%sdefault:\n", g.indent())
		g.state.indent++
		g.walkStmts(def.Body)
		g.state.indent--
	}

	fmt.Fprintf(w, "%s}\n", g.indent())
}

func (g *Generator) emitSwitchAsIfChain(w io.Writer, stmt *ast.SwitchStmt, cases []*ast.CaseClause, def *ast.CaseClause) {
	isString := stmt.Tag != nil && g.hasStringType(stmt.Tag)
	useDoWhile := g.isKernelMode()
	if useDoWhile {
		fmt.Fprintf(w, "%sdo {\n", g.indent())
		g.state.indent++
	}
	for i, cc := range cases {
		if i == 0 {
			fmt.Fprintf(w, "%sif (", g.indent())
		} else {
			fmt.Fprintf(w, "%s} else if (", g.indent())
		}
		for j, expr := range cc.List {
			if j > 0 {
				fmt.Fprintf(w, " || ")
			}
			if isString {
				fmt.Fprintf(w, "so_string_eq(")
				g.emitExpr(stmt.Tag)
				fmt.Fprintf(w, ", ")
				g.emitExpr(expr)
				fmt.Fprintf(w, ")")
			} else {
				g.emitExpr(stmt.Tag)
				fmt.Fprintf(w, " == (")
				g.emitExpr(expr)
				fmt.Fprintf(w, ")")
			}
		}
		fmt.Fprintf(w, ") {\n")
		g.state.indent++
		g.walkStmts(cc.Body)
		g.state.indent--
	}

	if def != nil {
		fmt.Fprintf(w, "%s} else {\n", g.indent())
		g.state.indent++
		g.walkStmts(def.Body)
		g.state.indent--
	}
	fmt.Fprintf(w, "%s}\n", g.indent())
	if useDoWhile {
		g.state.indent--
		fmt.Fprintf(w, "%s} while (0);\n", g.indent())
	}
}

func (g *Generator) endsWithBreakOrReturn(stmts []ast.Stmt) bool {
	if len(stmts) == 0 {
		return false
	}
	last := stmts[len(stmts)-1]
	switch s := last.(type) {
	case *ast.BranchStmt:
		return s.Tok == token.BREAK || s.Tok == token.CONTINUE
	case *ast.ReturnStmt:
		return true
	case *ast.BlockStmt:
		return g.endsWithBreakOrReturn(s.List)
	}
	return false
}
