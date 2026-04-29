package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: inspect <path-to-file>")
	}
	path := os.Args[len(os.Args)-1]

	// Read the source file.
	src, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, filepath.Base(path), src, parser.ParseComments)
	if err != nil {
		log.Fatalf("failed to parse file: %v", err)
	}

	// Print the AST.
	ast.Print(fset, f)

}
