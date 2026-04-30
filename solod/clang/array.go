package clang

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
)

// emitArrayLit emits a fixed-size array literal as a C initializer list.
// Example: [5]int{1, 2, 3, 4, 5} → {1, 2, 3, 4, 5}
// Wide string optimization: [...]uint16{'A', 'B', 0} → L"AB"
func (g *Generator) emitArrayLit(n *ast.CompositeLit) {
	w := g.state.writer

	if g.tryEmitWideStringLit(n) {
		return
	}

	fmt.Fprintf(w, "{")

	if hasKeyedElements(n) {
		g.emitSparseArrayValues(n)
	} else {
		for i, elt := range n.Elts {
			if i > 0 {
				fmt.Fprintf(w, ", ")
			}
			g.emitExpr(elt)
		}
	}

	fmt.Fprintf(w, "}")
}

// emitArrayCmpOperand emits an array comparison operand.
// Composite literals need a C compound literal prefix (e.g. (so_int[3]){...})
// wrapped in extra parentheses so commas inside braces don't split macro args.
func (g *Generator) emitArrayCmpOperand(expr ast.Expr, arr *types.Array) {
	w := g.state.writer
	if _, isLit := expr.(*ast.CompositeLit); isLit {
		elemType := g.mapType(expr, arr.Elem())
		fmt.Fprintf(w, "((%s%s)", elemType, arrayDims(arr))
		g.emitExpr(expr)
		fmt.Fprintf(w, ")")
		return
	}
	g.emitExpr(expr)
}

// emitSliceLit emits a slice literal as a so_Slice compound literal.
// Example: []int{1, 2, 3, 4} → {(so_int[4]){1, 2, 3, 4}, 4, 4}
func (g *Generator) emitSliceLit(n *ast.CompositeLit) {
	w := g.state.writer
	sl := g.types.TypeOf(n).Underlying().(*types.Slice)
	elemType := g.mapType(n, sl.Elem())
	size := len(n.Elts)
	if size == 0 {
		fmt.Fprintf(w, "(so_Slice){&so_Nil, 0, 0}")
		return
	}
	fmt.Fprintf(w, "(so_Slice){(%s[%d]){", elemType, size)
	for i, elt := range n.Elts {
		if i > 0 {
			fmt.Fprintf(w, ", ")
		}
		g.emitExpr(elt)
	}
	fmt.Fprintf(w, "}, %d, %d}", size, size)
}

// emitSparseArrayValues emits array values using C99 designated initializers
// for keyed elements. Example: [...]int{100, 3: 400, 500} → 100, [3] = 400, 500
func (g *Generator) emitSparseArrayValues(n *ast.CompositeLit) {
	w := g.state.writer
	for i, elt := range n.Elts {
		if i > 0 {
			fmt.Fprintf(w, ", ")
		}
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			fmt.Fprintf(w, "[")
			g.emitExpr(kv.Key)
			fmt.Fprintf(w, "] = ")
			g.emitExpr(kv.Value)
		} else {
			g.emitExpr(elt)
		}
	}
}

// emitSliceExpr emits a slice expression (e.g. nums[1:4]).
// For arrays: so_array_slice(T, arr, low, high, size).
// For slices: so_slice(T, s, low, high).
func (g *Generator) emitSliceExpr(n *ast.SliceExpr) {
	w := g.state.writer
	typ := g.types.TypeOf(n.X).Underlying()

	ptrDeref := false
	if ptr, ok := typ.(*types.Pointer); ok {
		if _, ok := ptr.Elem().Underlying().(*types.Array); ok {
			typ = ptr.Elem().Underlying()
			ptrDeref = true
		}
	}

	switch t := typ.(type) {
	case *types.Array:
		elemType := g.mapType(n, t.Elem())
		if n.Slice3 {
			fmt.Fprintf(w, "so_array_slice3(%s, ", elemType)
		} else {
			fmt.Fprintf(w, "so_array_slice(%s, ", elemType)
		}
		if ptrDeref {
			fmt.Fprintf(w, "(*")
			g.emitExpr(n.X)
			fmt.Fprintf(w, ")")
		} else {
			g.emitExpr(n.X)
		}
		fmt.Fprintf(w, ", ")
		if n.Low != nil {
			g.emitExpr(n.Low)
		} else {
			fmt.Fprintf(w, "0")
		}
		fmt.Fprintf(w, ", ")
		if n.High != nil {
			g.emitExpr(n.High)
		} else {
			fmt.Fprintf(w, "%d", t.Len())
		}
		if n.Slice3 {
			fmt.Fprintf(w, ", ")
			g.emitExpr(n.Max)
			fmt.Fprintf(w, ")")
		} else {
			fmt.Fprintf(w, ", %d)", t.Len())
		}

	case *types.Basic:
		if t.Kind() != types.String && t.Kind() != types.UntypedString {
			g.fail(n, "unsupported slice expression on basic type: %s", t)
			break
		}
		fmt.Fprintf(w, "so_string_slice(")
		g.emitExpr(n.X)
		fmt.Fprintf(w, ", ")
		if n.Low != nil {
			g.emitExpr(n.Low)
		} else {
			fmt.Fprintf(w, "0")
		}
		fmt.Fprintf(w, ", ")
		if n.High != nil {
			g.emitExpr(n.High)
		} else {
			g.emitExpr(n.X)
			fmt.Fprintf(w, ".len")
		}
		fmt.Fprintf(w, ")")

	case *types.Slice:
		elemType := g.mapType(n, t.Elem())
		if n.Slice3 {
			fmt.Fprintf(w, "so_slice3(%s, ", elemType)
		} else {
			fmt.Fprintf(w, "so_slice(%s, ", elemType)
		}
		g.emitExpr(n.X)
		fmt.Fprintf(w, ", ")
		if n.Low != nil {
			g.emitExpr(n.Low)
		} else {
			fmt.Fprintf(w, "0")
		}
		fmt.Fprintf(w, ", ")
		if n.High != nil {
			g.emitExpr(n.High)
		} else {
			g.emitExpr(n.X)
			fmt.Fprintf(w, ".len")
		}
		if n.Slice3 {
			fmt.Fprintf(w, ", ")
			g.emitExpr(n.Max)
		}
		fmt.Fprintf(w, ")")

	default:
		g.fail(n, "unsupported slice expression type: %T", t)
	}
}

func isArrayType(typ types.Type) bool {
	return arrayDims(typ) != ""
}

func hasKeyedElements(n *ast.CompositeLit) bool {
	for _, elt := range n.Elts {
		if _, ok := elt.(*ast.KeyValueExpr); ok {
			return true
		}
	}
	return false
}

func arrayDims(typ types.Type) string {
	typ = types.Unalias(typ)
	if _, ok := typ.(*types.Named); ok {
		return ""
	}
	var dims string
	for arr, ok := typ.(*types.Array); ok; arr, ok = arr.Elem().(*types.Array) {
		dims += fmt.Sprintf("[%d]", arr.Len())
	}
	return dims
}

func arraySize(typ types.Type) int64 {
	t := typ.Underlying()
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem().Underlying()
	}
	if arr, ok := t.(*types.Array); ok {
		return arr.Len()
	}
	return -1
}

func (g *Generator) tryEmitWideStringLit(n *ast.CompositeLit) bool {
	w := g.state.writer
	typ := g.types.TypeOf(n)
	arr, ok := typ.Underlying().(*types.Array)
	if !ok || len(n.Elts) == 0 {
		return false
	}
	elem := arr.Elem().Underlying()
	basic, ok := elem.(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Uint16, types.Uint, types.Int16, types.Rune:
	default:
		return false
	}
	var runes []rune
	for _, elt := range n.Elts {
		lit, ok := elt.(*ast.BasicLit)
		if !ok {
			return false
		}
		switch lit.Kind {
		case token.CHAR:
			r, err := parseRuneLit(lit.Value)
			if err != nil {
				return false
			}
			runes = append(runes, r)
		case token.INT:
			if lit.Value != "0" {
				return false
			}
			runes = append(runes, 0)
		default:
			return false
		}
	}
	fmt.Fprintf(w, "L\"")
	for _, r := range runes {
		if r == 0 {
			break
		}
		fmt.Fprintf(w, "%s", escapeWideChar(r))
	}
	fmt.Fprintf(w, "\"")
	return true
}

func parseRuneLit(s string) (rune, error) {
	if len(s) < 2 || s[0] != '\'' || s[len(s)-1] != '\'' {
		return 0, fmt.Errorf("invalid rune literal: %s", s)
	}
	content := s[1 : len(s)-1]
	if len(content) == 0 {
		return 0, fmt.Errorf("empty rune literal")
	}
	if content[0] != '\\' {
		return rune(content[0]), nil
	}
	switch content[1] {
	case 'n':
		return '\n', nil
	case 't':
		return '\t', nil
	case 'r':
		return '\r', nil
	case '\\':
		return '\\', nil
	case '\'':
		return '\'', nil
	case '"':
		return '"', nil
	case '0':
		return '\000', nil
	default:
		return rune(content[1]), nil
	}
}

func escapeWideChar(r rune) string {
	switch r {
	case '"':
		return `\"`
	case '\\':
		return `\\`
	case '\t':
		return `\t`
	case '\n':
		return `\n`
	case '\r':
		return `\r`
	default:
		if r < 32 || r > 127 {
			return fmt.Sprintf("\\u%04x", r)
		}
		return string(r)
	}
}
