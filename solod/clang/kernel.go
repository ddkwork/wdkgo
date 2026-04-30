package clang

import (
	"fmt"
	"go/ast"
	"io"
	"strings"
)

type KernelMode struct {
	Enabled    bool
	PackageDir string
}

func (g *Generator) isKernelMode() bool {
	return g.kernel.Enabled
}

func (g *Generator) hasWdkIncludes() bool {
	wdkHeaders := map[string]bool{
		"<ntddk.h>": true, "<ntifs.h>": true, "<ntdef.h>": true,
		"<wdm.h>": true, "<wdf.h>": true,
	}
	for _, inc := range g.includes {
		if wdkHeaders[inc] {
			return true
		}
	}
	return false
}

func (g *Generator) collectKernelDirectives() {
	for _, file := range g.pkg.Syntax {
		for _, cg := range file.Comments {
			for _, c := range cg.List {
				text := strings.TrimSpace(c.Text)
				if text == "//so:driver" {
					g.kernel.Enabled = true
				}
				if rest, ok := strings.CutPrefix(text, "//so:driver "); ok {
					g.kernel.Enabled = true
					g.kernel.PackageDir = strings.TrimSpace(rest)
				}
			}
		}
	}
}

func (g *Generator) emitKernelHeader(w io.Writer) {
	fmt.Fprintf(w, "#pragma once\n")

	for _, inc := range g.includes {
		fmt.Fprintf(w, "#include %s\n", inc)
	}

	if g.isKernelMode() {
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "#ifndef _%s_GEN_H_\n", strings.ToUpper(g.pkg.Name))
		fmt.Fprintf(w, "#define _%s_GEN_H_\n", strings.ToUpper(g.pkg.Name))
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "#ifdef ENV_WINDOWS\n")
		fmt.Fprintf(w, "#include <ntifs.h>\n")
		fmt.Fprintf(w, "#include <ntstrsafe.h>\n")
		fmt.Fprintf(w, "#endif\n")
	}

	fmt.Fprintf(w, "#include \"builtin_kernel.h\"\n")

	g.emitImports(w)
	g.emitHeaderDecls(w)

	if g.isKernelMode() {
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "#endif // _%s_GEN_H_\n", strings.ToUpper(g.pkg.Name))
	}
}

func (g *Generator) emitKernelImplHeader(w io.Writer) {
	if !g.isKernelMode() {
		return
	}

	fmt.Fprintf(w, "#include \"builtin_kernel.h\"\n")
	fmt.Fprintf(w, "#include \"%s.h\"\n", g.pkg.Name)
}

func (g *Generator) isDriverEntry(decl *ast.FuncDecl) bool {
	if !g.isKernelMode() {
		return false
	}
	return decl.Name.Name == "DriverEntry"
}
