package compiler

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"
)

//go:embed builtin/builtin.h builtin/builtin.c builtin/builtin_kernel.h builtin/builtin_kernel.c
var builtinFS embed.FS

//go:embed builtin/CMakeLists.txt builtin/build.bat builtin/wdk.h builtin/wdk.c
var cmakeFS embed.FS

func writeBuiltin(outDir string, kernelMode bool) error {
	if kernelMode {
		for _, name := range []string{"builtin_kernel.h", "builtin_kernel.c"} {
			data, err := builtinFS.ReadFile("builtin/" + name)
			if err != nil {
				return fmt.Errorf("read embedded builtin file %s: %w", name, err)
			}
			if err := os.WriteFile(filepath.Join(outDir, name), data, 0o644); err != nil {
				return fmt.Errorf("write builtin file %s: %w", name, err)
			}
		}
		return nil
	}
	for _, name := range []string{"builtin.h", "builtin.c"} {
		data, err := builtinFS.ReadFile("builtin/" + name)
		if err != nil {
			return fmt.Errorf("read embedded builtin file %s: %w", name, err)
		}
		if err := os.WriteFile(filepath.Join(outDir, name), data, 0o644); err != nil {
			return fmt.Errorf("write builtin file %s: %w", name, err)
		}
	}
	return nil
}

func writeCMake(outDir string) error {
	data, err := cmakeFS.ReadFile("builtin/CMakeLists.txt")
	if err != nil {
		return fmt.Errorf("read embedded CMakeLists.txt: %w", err)
	}
	content := string(data)
	content = strings.ReplaceAll(content, "${PROJECT_NAME}", "HyperDbgDriver")
	if err := os.WriteFile(filepath.Join(outDir, "CMakeLists.txt"), []byte(content), 0o644); err != nil {
		return fmt.Errorf("write CMakeLists.txt: %w", err)
	}

	if buildData, err := cmakeFS.ReadFile("builtin/build.bat"); err == nil {
		os.WriteFile(filepath.Join(outDir, "build.bat"), buildData, 0o644)
	}

	return nil
}

func copyAsmFiles(pkg *packages.Package, outDir string) error {
	if pkg == nil || len(pkg.GoFiles) == 0 {
		return nil
	}
	pkgDir := filepath.Dir(pkg.GoFiles[0])
	asmDir := filepath.Join(pkgDir, "asm")
	if _, err := os.Stat(asmDir); err == nil {
		if entries, err := os.ReadDir(asmDir); err == nil {
			for _, entry := range entries {
				name := entry.Name()
				if strings.HasSuffix(name, ".asm") {
					data, err := os.ReadFile(filepath.Join(asmDir, name))
					if err != nil {
						continue
					}
					content := string(data)
					os.WriteFile(filepath.Join(outDir, name), []byte(content), 0o644)
				}
			}
		}
	}
	return nil
}

func writeWdkBindings(outDir string) error {
	oldWdkDir := filepath.Join(outDir, "so", "wdk")
	os.RemoveAll(oldWdkDir)
	for _, name := range []string{"wdk.h", "wdk.c"} {
		data, err := cmakeFS.ReadFile("builtin/" + name)
		if err != nil {
			continue
		}
		os.WriteFile(filepath.Join(outDir, name), data, 0o644)
	}
	return nil
}

func copySoLibHeaders(outDir string) error {
	return nil
}

func fixAsmFunctions(outDir string) error {
	mainC := filepath.Join(outDir, "main.c")
	data, err := os.ReadFile(mainC)
	if err != nil {
		return nil
	}
	content := string(data)

	mainH := filepath.Join(outDir, "main.h")
	dataH, err := os.ReadFile(mainH)
	if err == nil {
		contentH := string(dataH)
		for _, fn := range []string{"GetGdtBase", "GetIdtBase", "GetGdtLimit", "GetIdtLimit", "GetRflags", "GetEs", "GetCs", "GetSs", "GetDs", "GetFs", "GetGs", "GetTr", "GetLdtr", "GetFsBase", "GetGsBase", "VmexitHandlerAddr", "VmxOn", "VmxOff", "VmxVmlaunch", "VmxVmread", "VmxVmwrite", "Invept", "Invvpid", "InveptAllContexts", "InvvpidAllContexts", "Vmcall", "SetSs", "SetDs", "SetEs", "SetFs", "SetGs", "ReloadGdtr", "ReloadIdtr", "StiInstruction", "CliInstruction", "ReadCr4", "ReadCr3", "ReadCr0", "VmxVmcall"} {
			contentH = strings.ReplaceAll(contentH, "main_Asm"+fn, "Asm"+fn)
		}
		contentH = strings.ReplaceAll(contentH, "main_", "")
		contentH = strings.ReplaceAll(contentH, "typedef uintptr_t PVOID;\n", "")
		os.WriteFile(mainH, []byte(contentH), 0o644)
	}

	asmPublicRe := regexp.MustCompile(`(?m)^PUBLIC\s+(Asm\w+)`)
	asmProcRe := regexp.MustCompile(`(?m)^(\w+)\s+PROC`)
	asmFuncs := make(map[string]bool)
	asmFiles, _ := filepath.Glob(filepath.Join(outDir, "*.asm"))
	for _, af := range asmFiles {
		data, err := os.ReadFile(af)
		if err != nil {
			continue
		}
		src := string(data)
		matches := asmPublicRe.FindAllStringSubmatch(src, -1)
		for _, m := range matches {
			asmFuncs[m[1]] = true
		}
		procMatches := asmProcRe.FindAllStringSubmatch(src, -1)
		for _, m := range procMatches {
			asmFuncs[m[1]] = true
		}
	}

	if len(asmFuncs) > 0 {
		var patterns []*regexp.Regexp
		for fn := range asmFuncs {
			patterns = append(patterns, regexp.MustCompile(regexp.QuoteMeta(fn)+`\([^)]*\)\s*\{(?:[^{}]|\{[^{}]*\})*\}`))
		}
		combined := regexp.MustCompile("(?m)(" + strings.Join(func() []string {
			var ps []string
			for _, p := range patterns {
				ps = append(ps, p.String())
			}
			return ps
		}(), "|") + ")")
		content = combined.ReplaceAllStringFunc(content, func(match string) string {
			sigEnd := strings.Index(match, "{")
			if sigEnd < 0 {
				return match
			}
			return match[:sigEnd] + ";"
		})
	}

	content = strings.ReplaceAll(content, "main_", "")

	if strings.Contains(content, "_solod_kmod_init") {
		content = strings.Replace(content, "DriverEntry(DRIVER_OBJECT* driverObj, UNICODE_STRING* registryPath) {", "DriverEntry(DRIVER_OBJECT* driverObj, UNICODE_STRING* registryPath) {\n    _solod_kmod_init();", 1)
		content = strings.Replace(content, "#include \"builtin_kernel.h\"\n", "#include \"builtin_kernel.h\"\nvoid _solod_kmod_init(void);\n", 1)
	}

	os.WriteFile(mainC, []byte(content), 0o644)
	return nil
}

type typeDef struct {
	full string
	name string
	body string
}

func fixTypeOrder(outDir string) error {
	mainH := filepath.Join(outDir, "main.h")
	data, err := os.ReadFile(mainH)
	if err != nil {
		return nil
	}
	content := string(data)

	typesMarker := "\n// -- Types --\n"
	typesStartIdx := strings.Index(content, typesMarker)
	if typesStartIdx == -1 {
		return nil
	}
	typesBodyStart := typesStartIdx + len(typesMarker)

	nextMarkerIdx := len(content)
	for _, marker := range []string{"\n// -- Variables", "\n// -- Functions", "\n// -- Extern"} {
		if idx := strings.Index(content[typesBodyStart:], marker); idx != -1 {
			if typesBodyStart+idx < nextMarkerIdx {
				nextMarkerIdx = typesBodyStart + idx
			}
		}
	}

	typesBody := content[typesBodyStart:nextMarkerIdx]
	typeDefs, typeNames := parseTypeDefs(typesBody)
	typeDefs = topoSortTypeDefs(typeDefs, typeNames)

	var forwardDecls strings.Builder
	forwardDecls.WriteString("\n// -- Forward declarations --\n")
	for _, td := range typeDefs {
		if strings.Contains(td.full, "typedef struct ") && strings.Contains(td.full, "{") {
			forwardDecls.WriteString("typedef struct ")
			forwardDecls.WriteString(td.name)
			forwardDecls.WriteString(" ")
			forwardDecls.WriteString(td.name)
			forwardDecls.WriteString(";\n")
		}
	}

	var buf strings.Builder
	buf.WriteString("// -- Types --\n")
	for _, td := range typeDefs {
		buf.WriteString(td.full)
	}

	newContent := content[:typesStartIdx] + "\n" + forwardDecls.String() + buf.String() + content[nextMarkerIdx:]
	os.WriteFile(mainH, []byte(newContent), 0o644)
	return nil
}

func parseTypeDefs(body string) ([]typeDef, map[string]bool) {
	var defs []typeDef
	typeNames := make(map[string]bool)

	lines := strings.Split(body, "\n")
	var current strings.Builder
	braceDepth := 0
	inTypedef := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" && braceDepth == 0 && current.Len() == 0 {
			continue
		}

		current.WriteString(line)
		current.WriteString("\n")

		for _, ch := range trimmed {
			if ch == '{' {
				braceDepth++
			} else if ch == '}' {
				braceDepth--
			}
		}

		if strings.HasPrefix(trimmed, "typedef ") {
			inTypedef = true
		}

		if braceDepth == 0 && inTypedef && strings.HasSuffix(trimmed, ";") {
			text := current.String()
			name := extractTypeName(text)
			if name != "" {
				typeNames[name] = true
				defs = append(defs, typeDef{
					full: text,
					name: name,
					body: text,
				})
			}
			current.Reset()
			inTypedef = false
		}
	}

	return defs, typeNames
}

func extractTypeName(text string) string {
	if strings.Contains(text, "}") {
		if idx := strings.LastIndex(text, "} "); idx != -1 {
			name := strings.TrimSpace(text[idx+2:])
			name = strings.TrimSuffix(name, ";")
			name = strings.TrimSpace(name)
			return name
		}
	}
	lines := strings.Split(text, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || trimmed == ";" {
			continue
		}
		trimmed = strings.TrimSuffix(trimmed, ";")
		parts := strings.Fields(trimmed)
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}
	return ""
}

func topoSortTypeDefs(defs []typeDef, typeNames map[string]bool) []typeDef {
	nameToIdx := make(map[string]int)
	for i, td := range defs {
		nameToIdx[td.name] = i
	}

	deps := make([]map[int]bool, len(defs))
	for i := range deps {
		deps[i] = make(map[int]bool)
	}

	for i, td := range defs {
		for name := range typeNames {
			if name == td.name {
				continue
			}
			if isWordInBody(td.body, name) {
				if idx, ok := nameToIdx[name]; ok {
					deps[i][idx] = true
				}
			}
		}
	}

	inDegree := make([]int, len(defs))
	for i := range defs {
		inDegree[i] = len(deps[i])
	}

	var queue []int
	for i := range defs {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	var sorted []typeDef
	placed := make(map[int]bool)
	for len(queue) > 0 {
		idx := queue[0]
		queue = queue[1:]
		sorted = append(sorted, defs[idx])
		placed[idx] = true

		for i := range defs {
			if placed[i] {
				continue
			}
			if deps[i][idx] {
				delete(deps[i], idx)
				inDegree[i] = len(deps[i])
				if inDegree[i] == 0 {
					queue = append(queue, i)
				}
			}
		}
	}

	for i := range defs {
		if !placed[i] {
			sorted = append(sorted, defs[i])
		}
	}

	return sorted
}

func isWordInBody(body, word string) bool {
	for i := 0; i <= len(body)-len(word); i++ {
		if body[i:i+len(word)] == word {
			before := i == 0 || !isAlphaNum(body[i-1])
			after := i+len(word) == len(body) || !isAlphaNum(body[i+len(word)])
			if before && after {
				return true
			}
		}
	}
	return false
}

func isAlphaNum(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_'
}
