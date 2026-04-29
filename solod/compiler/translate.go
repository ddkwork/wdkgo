package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"golang.org/x/tools/go/packages"

	"solod.dev/clang"
)

// Translate loads all Go packages from srcDir (including So stdlib dependencies),
// translates them to C, and writes the output to outDir.
func Translate(srcDir string, outDir string) error {
	// Clean output directory before generating
	if err := os.RemoveAll(outDir); err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	pkgs, err := loadPackages(srcDir)
	if err != nil {
		return err
	}
	if len(pkgs) == 0 {
		return fmt.Errorf("no packages found")
	}

	entry := pkgs[0]

	var entryModulePath string
	if entry.Module != nil {
		entryModulePath = entry.Module.Path
	}

	var soModulePath string
	if info, ok := debug.ReadBuildInfo(); ok {
		soModulePath = info.Main.Path
	}

	// Walk import graph and collect transpilable packages in topological order
	ordered := topoSort(entry, entryModulePath, soModulePath)

	// Translate each package
	for _, pkg := range ordered {
		pkgOutDir := packageOutDir(pkg, entry, outDir)
		if err := clang.Emit(clang.EmitOptions{
			Pkg:    pkg,
			OutDir: pkgOutDir,
		}); err != nil {
			return err
		}
	}

	// Write embedded builtin files into the output directory (flat, no subdirs)
	if err := writeBuiltin(outDir, hasKernelDirective(entry)); err != nil {
		return err
	}

	// Copy .s and .asm files from the source package to output
	if err := copyAsmFiles(entry, outDir); err != nil {
		return err
	}

	// Write WDK bindings (so/wdk) manually to flatten to gen/
	if err := writeWdkBindings(outDir); err != nil {
		return err
	}

	// Copy so lib headers (slices.h, atomic.h) that Solod references but doesn't emit
	if err := copySoLibHeaders(outDir); err != nil {
		return err
	}

	// Reorder type definitions in main.h to resolve forward references
	if err := fixTypeOrder(outDir); err != nil {
		return err
	}

	// Fix asm function declarations in main.c to be extern
	if err := fixAsmFunctions(outDir); err != nil {
		return err
	}

	// Write CMakeLists.txt template
	if err := writeCMake(outDir); err != nil {
		return err
	}

	return nil
}

func hasKernelDirective(pkg *packages.Package) bool {
	for _, file := range pkg.Syntax {
		for _, cg := range file.Comments {
			for _, c := range cg.List {
				text := strings.TrimSpace(c.Text)
				if text == "//so:driver" || strings.HasPrefix(text, "//so:driver ") {
					return true
				}
			}
		}
	}
	return false
}

// loadPackages uses go/packages to load the entry package and all dependencies.
func loadPackages(dir string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedImports | packages.NeedDeps |
			packages.NeedModule | packages.NeedTypesInfo,
		Dir: dir,
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return nil, fmt.Errorf("load packages: %w", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("packages contain errors")
	}
	return pkgs, nil
}

// topoSort walks the import graph from entry and returns transpilable packages
// (module-internal + So stdlib) in topological order (dependencies before dependents).
func topoSort(entry *packages.Package, entryModulePath, soModulePath string) []*packages.Package {
	var ordered []*packages.Package
	visited := make(map[string]bool)

	var walk func(pkg *packages.Package)
	walk = func(pkg *packages.Package) {
		if visited[pkg.PkgPath] {
			return
		}
		visited[pkg.PkgPath] = true

		// Visit dependencies first (post-order)
		for _, dep := range pkg.Imports {
			if shouldTranspile(dep, entryModulePath, soModulePath) {
				walk(dep)
			}
		}
		ordered = append(ordered, pkg)
	}
	walk(entry)
	return ordered
}

// packageOutDir returns the output directory for a package.
// Entry package goes to outDir directly.
// Other packages strip their module prefix (e.g. solod.dev/math -> math).
func packageOutDir(pkg, entry *packages.Package, outDir string) string {
	if pkg.PkgPath == entry.PkgPath {
		return outDir
	}
	relPath := strings.TrimPrefix(pkg.PkgPath, pkg.Module.Path+"/")
	// Flatten so/ packages to output root (so/wdk -> outDir)
	if strings.HasPrefix(relPath, "so/") {
		return outDir
	}
	return filepath.Join(outDir, relPath)
}

// shouldTranspile returns true if a package should be transpiled to C.
// This only includes packages from the entry module.
func shouldTranspile(pkg *packages.Package, entryModulePath, soModulePath string) bool {
	if pkg.Module == nil {
		return false
	}
	if pkg.Module.Path == entryModulePath {
		return true
	}
	if pkg.Module.Path == soModulePath && !strings.HasPrefix(pkg.PkgPath, soModulePath+"/so/wdk") {
		return true
	}
	return false
}
