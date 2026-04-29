package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"modernc.org/cc/v4"
)

func main() {
	wdkPath := flag.String("wdk", "E:\\Program Files\\Windows Kits\\10\\Include\\10.0.28000.0", "WDK include root path")
	output := flag.String("o", "", "Output directory")
	flag.Parse()

	if *output == "" {
		fmt.Fprintln(os.Stderr, "Usage: wdkgen -wdk <wdk-root> -o <output-dir>")
		os.Exit(1)
	}

	tmpDir := filepath.Join(os.TempDir(), "wdkgen_preprocess")
	os.MkdirAll(tmpDir, 0o755)

	kmInc := filepath.Join(*wdkPath, "km")
	sharedInc := filepath.Join(*wdkPath, "shared")

	if err := os.MkdirAll(*output, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output dir: %v\n", err)
		os.Exit(1)
	}

	allMacros := extractMacrosFromHeaders(*wdkPath)

	headerFiles := collectHeaderFiles(kmInc, sharedInc)
	sort.Strings(headerFiles)

	globalGen := NewGenerator(*output)
	globalGen.macros = allMacros

	globalGen.writeDriverGo()

	targetHeaders := map[string]bool{
		"wsk.h":      true,
		"ntifs.h":    true,
		"netioddk.h": true,
		"ntddk.h":    true,
	}

	for _, hFile := range headerFiles {
		relPath, _ := filepath.Rel(kmInc, hFile)
		if relPath == hFile {
			relPath, _ = filepath.Rel(sharedInc, hFile)
		}
		baseName := filepath.Base(hFile)
		if !targetHeaders[baseName] {
			continue
		}
		fmt.Fprintf(os.Stderr, "Processing: %s\n", relPath)

		cfg, err := cc.NewConfig("windows", "amd64")
		if err != nil {
			fmt.Fprintf(os.Stderr, "  SKIP: config error: %v\n", err)
			continue
		}
		ucrtInc := filepath.Join(*wdkPath, "ucrt")
		crtInc := filepath.Join(kmInc, "crt")
		cfg.IncludePaths = []string{kmInc, sharedInc, ucrtInc, crtInc}
		cfg.SysIncludePaths = []string{kmInc, sharedInc, ucrtInc, crtInc}
		cfg.Predefined += wdkPredefined()

		content, err := os.ReadFile(hFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  SKIP: read error: %v\n", err)
			continue
		}

		sources := []cc.Source{
			{Name: "<predefined>", Value: cfg.Predefined},
			{Name: "<builtin>", Value: cc.Builtin},
			{Name: filepath.Base(hFile), Value: string(content)},
		}
		ast, parseErr := cc.Parse(cfg, sources)
		if parseErr != nil || ast == nil {
			fmt.Fprintf(os.Stderr, "  SKIP: parse error: %v\n", parseErr)
			continue
		}

		fileGen := NewGenerator(*output)
		fileGen.macros = allMacros
		fileGen.walkAST(ast.TranslationUnit)

		if len(fileGen.types) > 0 || len(fileGen.funcs) > 0 || len(fileGen.macros) > 0 {
			fmt.Fprintf(os.Stderr, "  => types=%d funcs=%d macros=%d\n", len(fileGen.types), len(fileGen.funcs), len(fileGen.macros))
		}

		goFileName := sanitizeFileName(relPath) + ".go"
		fileGen.writeFileBindings(goFileName, globalGen)
	}

	totalTypes := len(globalGen.seenTypes)
	totalFuncs := len(globalGen.seenFuncs)
	totalMacros := len(globalGen.seenMacros)

	fmt.Printf("Generated WDK bindings in %s\n", *output)
	fmt.Printf("  Headers: %d, Types: %d, Functions: %d, Macros: %d\n", len(headerFiles), totalTypes, totalFuncs, totalMacros)
}

func collectHeaderFiles(dirs ...string) []string {
	var files []string
	visited := make(map[string]bool)
	var walkDir func(dir string)
	walkDir = func(dir string) {
		if visited[dir] {
			return
		}
		visited[dir] = true
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, entry := range entries {
			fullPath := filepath.Join(dir, entry.Name())
			if entry.IsDir() {
				walkDir(fullPath)
				continue
			}
			if strings.HasSuffix(strings.ToLower(entry.Name()), ".h") {
				files = append(files, fullPath)
			}
		}
	}
	for _, dir := range dirs {
		walkDir(dir)
	}
	return files
}

func sanitizeFileName(path string) string {
	path = strings.ReplaceAll(path, "\\", "_")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, ".", "_")
	path = regexp.MustCompile(`_+`).ReplaceAllString(path, "_")
	path = strings.Trim(path, "_")
	return path
}

func wdkPredefined() string {
	return `
typedef struct _GUID { unsigned long Data1; unsigned short Data2; unsigned short Data3; unsigned char Data4[8]; } GUID;
typedef GUID *PGUID;
typedef GUID *LPGUID;
typedef const GUID *LPCGUID;
typedef struct _LUID { unsigned long LowPart; long HighPart; } LUID;
typedef LUID *PLUID;
#define __extension__
#define __declspec(x)
#define __inline inline
#define __forceinline inline
#define POINTER_64
#define INLINE __inline
#define __cdecl
#define __stdcall
#define __fastcall
#define THIS_CALL
#define __thiscall
#define __vectorcall
#define __unaligned
#define __w64
#define __ptr64
#define __ptr32
#define __uptr
#define __sptr
#define __int8 char
#define __int16 short
#define __int32 int
#define __int64 long long
#define __int128 __int128_t
#define __uint128 unsigned __int128_t
#define __unsigned_int128 unsigned __int128_t
#define _KERNEL_MODE 1
#define _WIN64 1
#define _AMD64_ 1
#define _M_X64 100
#define _M_AMD64 100
#define _X86_ 1
#define _M_IX86 600
#define _WIN32 1
#define NTDDI_VERSION 0x0A00000B
#define WINVER 0x0A00
#define _WIN32_WINNT 0x0A00
#define WINAPI_FAMILY_SYSTEM 4
#define WINAPI_PARTITION_SYSTEM 1
#define WINAPI_PARTITION_APP 1
#define WINAPI_PARTITION_DESKTOP 1
#define WINAPI_PARTITION_GAMES 1
#define WINAPI_PARTITION_SERVER 1
#define MSWSOCKDEF_INLINE inline
#define WS2TCPIP_INLINE inline
#define MSTCPIP_INLINE inline
#define __inline inline
#define __forceinline inline
#define INLINE __inline
#define FORCEINLINE __forceinline
#define __builtin_alignof(x) 4
#define __int64 long long
#define UNALIGNED
#define __unaligned
#define C_ASSERT(expr)
#define _PREFAST_ 1
#define POOL_NX_OPTIN 1
#define _MSC_VER 1944
#define _MSC_FULL_VER 194435207
#define _MSC_BUILD 0
#define __MSVC_RUNTIME_CHECKS 1
#define _Outptr_
#define _In_
#define IN
#define _In_range_(a, b)
#define _Out_writes_to_(a, b)
#define _Out_writes_(x)
#define _Ret_range_(a, b)
#define _At_(a, b)
#define _Inout_
#define _Out_
#define _In_reads_(x)
#define _Inexpressible_(x)
#define _Field_range_(a, b)
#define _Field_size_bytes_full_(...)
#define _Field_size_bytes_part_(...)
#define _Field_size_bytes_part_opt_(...)
#define _Field_size_bytes_(...)
#define _Field_size_bytes_opt_(...)
#define _Field_size_(...)
#define _When_(...)
#define _Success_(x)
#define _Analysis_assume_(...) ((void)0)
#define __assume(x)
#define __noop(...) ((void)0)
#define FAR
#define __declspec(x)
#define __pragma(x)
#define pragma
#define __analysis_noreturn
#define _Analysis_assume_(x) ((void)0)
#define __forceinline static inline
#define __annotation(...)
#define __break(x)
#define __emit(x)
#define NT_ASSERT(_exp) ((void)0)
#define NT_ASSERTMSG(_msg, _exp) ((void)0)
#define NT_ASSERTMSGW(_msg, _exp) ((void)0)
#define NT_FRE_ASSERT(_exp) ((void)0)
#define NT_FRE_ASSERTMSG(_msg, _exp) ((void)0)
#define NT_FRE_ASSERTMSGW(_msg, _exp) ((void)0)
#define NT_VERIFY(_exp) ((void)0)
#define NT_VERIFYMSG(_msg, _exp) ((void)0)
#define NT_VERIFYMSGW(_msg, _exp) ((void)0)
#define NTSYSAPI
#define NTAPI __stdcall
#define STDAPI __stdcall
#define STDMETHODIMP HRESULT __stdcall
#define STDMETHODIMP_(x) x __stdcall
#define CDECL __cdecl
#define CALLBACK __stdcall
#define WINAPI __stdcall
#define APIENTRY NTAPI
#define INLINE inline
#define FORCEINLINE static inline
#define MSTCPIP_INLINE inline
#define CONST const
#define VOLATILE volatile
#define AF_MAX 45
#define TRUE 1
#define FALSE 0
#define UCHAR unsigned char
#define USHORT unsigned short
#define ULONG unsigned long
#define LONG long
#define UINT unsigned int
#define INT int
#define VOID void
#define CHAR char
#define SHORT short
#ifndef NULL
#define NULL 0
#endif
#define MAXUINT ((unsigned int)-1)
#define MAXULONG ((unsigned long)-1)
#define MAXLONGLONG ((long long)-1)
#define MINLONGLONG ((long long)0x8000000000000000)
typedef unsigned long ULONG32;
typedef long NET_IF_OPER_STATUS;
typedef long NET_IF_ADMIN_STATUS;
typedef unsigned long NET_IF_OBJECT_ID;
typedef unsigned short UINT16;
typedef unsigned int UINT32;
typedef unsigned char UINT8;
typedef unsigned short USHORT;
typedef unsigned long ULONG;
typedef unsigned int UINT;
typedef unsigned int UINT32;
typedef UINT16 ICMP_HEADER;
typedef UINT16 ICMPV6_HEADER;
typedef UINT16 *PICMP_HEADER;
typedef struct _ICMP_MESSAGE {
    ICMP_HEADER Header;
    union {
        UINT32 Data32[1];
        UINT16 Data16[2];
        UINT8 Data8[4];
    } Data;
} ICMP_MESSAGE;
typedef ICMP_MESSAGE ICMPV6_MESSAGE, *PICMPV6_MESSAGE;
typedef NET_IF_OBJECT_ID *PNET_IF_OBJECT_ID;
#define IF_COUNTED_INTERFACE void
#define IF_LOGICALADDRESS void
#define DEFINE_GUID(name, l, w1, w2, b1, b2, b3, b4, b5, b6, b7, b8) \
    static const GUID name = { l, w1, w2, { b1, b2, b3, b4, b5, b6, b7, b8 } }
#define OPTIONAL
typedef GUID *LPGUID;
typedef const GUID *LPCGUID;
typedef struct _LUID { unsigned long LowPart; long HighPart; } LUID;
typedef LUID *PLUID;
typedef unsigned char UCHAR;
typedef unsigned short USHORT;
typedef unsigned long LONG;
typedef unsigned long long UINT64;
typedef unsigned short wchar_t;
typedef unsigned char *PUCHAR;
typedef unsigned long ULONG;
typedef unsigned long ULONG_PTR;
typedef unsigned long long ULONGLONG;
typedef unsigned long SIZE_T;
typedef unsigned long DWORD_PTR;
typedef unsigned long size_t;
typedef unsigned long ULONG64;
typedef unsigned short WCHAR;
typedef char CHAR;
typedef char PCHAR;
typedef short SHORT;
typedef short *PSHORT;
typedef long LONG;
typedef unsigned char BOOLEAN;
typedef unsigned char UINT8;
typedef unsigned short UINT16;
typedef unsigned long SCOPE_LEVEL;
typedef unsigned long SCOPE_ID;
typedef struct _SCOPE_ID { unsigned long Level; } _SCOPE_ID;
typedef WCHAR *PWSTR;
typedef void *HANDLE;
typedef struct sockaddr sockaddr;
typedef struct in_addr { unsigned long s_addr; } IN_ADDR;
typedef struct in6_addr { unsigned char u[16]; } IN6_ADDR;
typedef char *PSTR;
typedef unsigned short PORT;
typedef unsigned long *PULONG;
typedef char *PCSTR;
typedef unsigned short *PUSHORT;
typedef unsigned short *PCWSTR;
typedef unsigned short *LPCWSTR;
typedef USHORT ADDRESS_FAMILY;
typedef struct _PROCESSOR_NUMBER { USHORT Group; USHORT Number; USHORT Reserved; } PROCESSOR_NUMBER;
`
}

var lineMarkerRe = regexp.MustCompile(`^# \d+ ".*?"(?: \d+)*$`)
var lineMarkerMSVCRe = regexp.MustCompile(`^#line\s+\d+.*$`)
var pragmaRe = regexp.MustCompile(`^#\s*pragma\s`)

func cleanPreprocessed(s string) string {
	lines := strings.Split(s, "\n")
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || lineMarkerRe.MatchString(line) || pragmaRe.MatchString(line) {
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

func cleanPreprocessedMSVC(s string) string {
	lines := strings.Split(s, "\n")
	var result []string
	inInlineFunc := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || lineMarkerMSVCRe.MatchString(line) || pragmaRe.MatchString(line) {
			continue
		}
		if strings.Contains(line, "__asm") || strings.Contains(line, "__try") ||
			strings.Contains(line, "__except") || strings.Contains(line, "__finally") ||
			strings.Contains(line, "__leave") || strings.Contains(line, "__declspec(allocate(") ||
			strings.Contains(line, "__declspec(code_seg(") || strings.Contains(line, "__declspec(deprecated(") ||
			strings.Contains(line, "__pragma") || strings.Contains(line, "__C_ASSERT__") ||
			strings.Contains(line, "#include") {
			continue
		}
		if strings.HasPrefix(line, "__readgsqword") || strings.HasPrefix(line, "__readfsdword") ||
			strings.HasPrefix(line, "__writegsqword") || strings.HasPrefix(line, "__writefsdword") ||
			strings.HasPrefix(line, "__readcr0") || strings.HasPrefix(line, "__writecr0") ||
			strings.HasPrefix(line, "__readcr2") || strings.HasPrefix(line, "__writecr2") ||
			strings.HasPrefix(line, "__readcr3") || strings.HasPrefix(line, "__writecr3") ||
			strings.HasPrefix(line, "__readcr4") || strings.HasPrefix(line, "__writecr4") ||
			strings.HasPrefix(line, "__readcr8") || strings.HasPrefix(line, "__writecr8") ||
			strings.HasPrefix(line, "__readmsr") || strings.HasPrefix(line, "__writemsr") ||
			strings.HasPrefix(line, "__readpmc") || strings.HasPrefix(line, "__readtsc") {
			continue
		}
		if strings.HasPrefix(line, "__inline") || strings.HasPrefix(line, "__forceinline") ||
			strings.HasPrefix(line, "MSTCPIP_INLINE") || strings.HasPrefix(line, "MSWSOCKDEF_INLINE") ||
			strings.HasPrefix(line, "WS2TCPIP_INLINE") || line == "inline" || strings.HasPrefix(line, "inline ") {
			inInlineFunc = true
			continue
		}
		if inInlineFunc {
			continue
		}
		result = append(result, line)
	}
	cleaned := strings.Join(result, "\n")
	replacements := map[string]string{
		"__int64":       "long long",
		"__int128":      "__int128_t",
		"__int32":       "int",
		"__int16":       "short",
		"__int8":        "char",
		"__int3264":     "intptr_t",
		"__wchar_t":     "wchar_t",
		"__cdecl":       "",
		"__stdcall":     "",
		"__fastcall":    "",
		"__thiscall":    "",
		"__vectorcall":  "",
		"__forceinline": "inline",
		"__inline":      "inline",
		"__declspec":    "",
		"__unaligned":   "",
		"__ptr64":       "",
		"__ptr32":       "",
		"__w64":         "",
		"__uptr":        "",
		"__sptr":        "",
		"CONST":         "const",
	}
	for old, new := range replacements {
		cleaned = strings.ReplaceAll(cleaned, old, new)
	}
	re := regexp.MustCompile(`__declspec\s*\(.*?\)`)
	cleaned = re.ReplaceAllString(cleaned, "")
	return cleaned
}

func extractMacrosFromDM(s string) map[string]string {
	macros := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(s))
	defineRe := regexp.MustCompile(`^#define\s+(\w+)\s*(.*)$`)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		m := defineRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		name := m[1]
		value := strings.TrimSpace(m[2])
		if strings.HasPrefix(name, "__") || (strings.HasPrefix(name, "_") && len(name) <= 3) {
			continue
		}
		if isReserved(name) {
			continue
		}
		if strings.Contains(value, "#") || strings.Contains(value, "##") {
			continue
		}
		if len(value) > 500 {
			continue
		}
		if strings.Contains(value, "__builtin") || strings.Contains(value, "__attribute__") {
			continue
		}
		macros[name] = value
	}
	return macros
}

func extractMacrosFromHeaders(wdkPath string) map[string]string {
	macros := make(map[string]string)
	defineRe := regexp.MustCompile(`^\s*#define\s+(\w+)(?:\(([^)]*)\))?\s*(.*)?$`)
	dirs := []string{
		filepath.Join(wdkPath, "km"),
		filepath.Join(wdkPath, "shared"),
		filepath.Join(wdkPath, "um"),
	}
	visited := make(map[string]bool)
	var walkDir func(dir string)
	walkDir = func(dir string) {
		if visited[dir] {
			return
		}
		visited[dir] = true
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, entry := range entries {
			fullPath := filepath.Join(dir, entry.Name())
			if entry.IsDir() {
				walkDir(fullPath)
				continue
			}
			if !strings.HasSuffix(strings.ToLower(entry.Name()), ".h") {
				continue
			}
			data, err := os.ReadFile(fullPath)
			if err != nil {
				continue
			}
			scanner := bufio.NewScanner(strings.NewReader(string(data)))
			for scanner.Scan() {
				line := scanner.Text()
				m := defineRe.FindStringSubmatch(line)
				if m == nil {
					continue
				}
				name := m[1]
				params := m[2]
				value := strings.TrimSpace(m[3])
				if params != "" {
					continue
				}
				if !isWDKConstant(name) && !strings.HasPrefix(strings.ToUpper(name), "NTDDI") &&
					!strings.HasPrefix(strings.ToUpper(name), "WINVER") &&
					!strings.HasPrefix(strings.ToUpper(name), "NTSTATUS") {
					continue
				}
				if len(name) > 64 {
					continue
				}
				if _, exists := macros[name]; exists {
					continue
				}
				if value == "" {
					value = "1"
				}
				if strings.Contains(value, "#") || strings.Contains(value, "##") {
					continue
				}
				if len(value) > 200 {
					continue
				}
				if strings.TrimSpace(value) == "\\" || strings.HasSuffix(strings.TrimSpace(value), "\\") {
					continue
				}
				trimmedVal := strings.TrimSpace(value)
				isStringLiteral := strings.HasPrefix(trimmedVal, "\"") || strings.HasPrefix(trimmedVal, "L\"") || len(trimmedVal) > 0 && trimmedVal[0] == '\''
				if isStringLiteral && !isComplexExpression(cleanMacroValue(value)) {
					continue
				}
				if strings.Contains(trimmedVal, "(") && !strings.HasSuffix(trimmedVal, ")") {
					continue
				}
				if strings.Contains(trimmedVal, "/*") || strings.Contains(trimmedVal, "*/") {
					continue
				}
				cleanedValue := cleanMacroValue(value)
				if !isValidGoConst(cleanedValue) {
					continue
				}
				macros[name] = cleanedValue
			}
		}
	}
	for _, dir := range dirs {
		walkDir(dir)
	}
	return macros
}

func isReserved(s string) bool {
	switch s {
	case "true", "false", "NULL", "null", "errno", "stdin", "stdout", "stderr":
		return true
	}
	return false
}

type FuncInfo struct {
	Name       string
	ReturnType string
	Params     []ParamInfo
	Comment    string
}

type ParamInfo struct {
	Name string
	Type string
}

type TypeInfo struct {
	Name     string
	Kind     string
	CName    string
	GoType   string
	Members  []MemberInfo
	Size     int
	IsOpaque bool
	Comment  string
}

type MemberInfo struct {
	Name     string
	Type     string
	Offset   int
	BitField bool
	BitWidth int
}

type Generator struct {
	outDir     string
	macros     map[string]string
	types      map[string]*TypeInfo
	funcs      map[string]*FuncInfo
	seenTypes  map[string]bool
	seenFuncs  map[string]bool
	seenMacros map[string]bool
}

func NewGenerator(outDir string) *Generator {
	return &Generator{
		outDir:     outDir,
		macros:     make(map[string]string),
		types:      make(map[string]*TypeInfo),
		funcs:      make(map[string]*FuncInfo),
		seenTypes:  make(map[string]bool),
		seenFuncs:  make(map[string]bool),
		seenMacros: make(map[string]bool),
	}
}

func (g *Generator) walkAST(tu *cc.TranslationUnit) {
	if tu == nil {
		return
	}
	g.walkTranslationUnit(tu)
	g.walkTranslationUnit(tu.TranslationUnit)
}

func (g *Generator) walkTranslationUnit(tu *cc.TranslationUnit) {
	if tu == nil || tu.ExternalDeclaration == nil {
		return
	}
	ed := tu.ExternalDeclaration
	switch ed.Case {
	case cc.ExternalDeclarationDecl:
		if ed.Declaration != nil {
			g.processDeclaration(ed.Declaration)
		}
	case cc.ExternalDeclarationFuncDef:
		if ed.FunctionDefinition != nil {
			g.processFunctionDef(ed.FunctionDefinition)
		}
	}
}

func (g *Generator) processDeclaration(d *cc.Declaration) {
	if d == nil {
		return
	}
	switch d.Case {
	case cc.DeclarationDecl:
		g.processInitDeclaratorList(d.DeclarationSpecifiers, d.InitDeclaratorList)
	case cc.DeclarationAssert:
		return
	case cc.DeclarationAuto:
		if d.DeclarationSpecifiers != nil {
			g.processTypeOnlyDecl(d.DeclarationSpecifiers)
		}
	}
}

func (g *Generator) processInitDeclaratorList(ds *cc.DeclarationSpecifiers, idl *cc.InitDeclaratorList) {
	if ds == nil || idl == nil {
		return
	}
	for idl != nil {
		if idl.InitDeclarator != nil && idl.InitDeclarator.Declarator != nil {
			decl := idl.InitDeclarator.Declarator
			name := g.getDeclaratorName(decl)
			if name == "" || isReserved(name) {
				idl = idl.InitDeclaratorList
				continue
			}
			isTypedef := g.isTypedef(ds)
			typeStr := g.typeSpecToString(ds)
			if isTypedef {
				g.types[name] = &TypeInfo{
					Name:     name,
					Kind:     "typedef",
					CName:    typeStr,
					GoType:   cTypeToGo(typeStr),
					IsOpaque: false,
				}
			} else if g.isFunctionPointer(decl) || strings.Contains(typeStr, "(*)") {
				g.funcs[name] = &FuncInfo{
					Name:       name,
					ReturnType: cTypeToGo(typeStr),
					Comment:    fmt.Sprintf("// function pointer or complex type: %s", typeStr),
				}
			} else if decl.DirectDeclarator != nil &&
				decl.DirectDeclarator.Case == cc.DirectDeclaratorFuncParam ||
				decl.DirectDeclarator != nil &&
					(decl.DirectDeclarator.Case == cc.DirectDeclaratorFuncIdent) {
				if g.hasFunctionParameters(decl) {
					params := g.extractParams(decl)
					g.funcs[name] = &FuncInfo{
						Name:       name,
						ReturnType: cTypeToGo(typeStr),
						Params:     params,
					}
				}
			}
		}
		idl = idl.InitDeclaratorList
	}
}

func (g *Generator) processTypeOnlyDecl(ds *cc.DeclarationSpecifiers) {
	if ds == nil {
		return
	}
	ts := g.findTypeSpecifier(ds)
	if ts == nil {
		return
	}
	switch ts.Case {
	case cc.TypeSpecifierStructOrUnion:
		if ts.StructOrUnionSpecifier != nil {
			sus := ts.StructOrUnionSpecifier
			name := g.getStructOrUnionName(sus)
			if name != "" {
				kind := "struct"
				if sus.StructOrUnion != nil && string(sus.StructOrUnion.Token.Src()) == "union" {
					kind = "union"
				}
				members := g.extractMembers(sus)
				g.types[name] = &TypeInfo{
					Name:     name,
					Kind:     kind,
					GoType:   name,
					Members:  members,
					IsOpaque: len(members) == 0,
				}
			}
		}
	case cc.TypeSpecifierEnum:
		if ts.EnumSpecifier != nil {
			es := ts.EnumSpecifier
			name := g.getEnumName(es)
			if name != "" {
				g.types[name] = &TypeInfo{
					Name:   name,
					Kind:   "enum",
					GoType: "int32",
				}
			}
			if es.EnumeratorList != nil {
				g.processEnumerators(es.EnumeratorList)
			}
		}
	}
}

func (g *Generator) processFunctionDef(fd *cc.FunctionDefinition) {
	if fd == nil || fd.Declarator == nil {
		return
	}
	name := g.getDeclaratorName(fd.Declarator)
	if name == "" || isReserved(name) {
		return
	}
	typeStr := ""
	if fd.DeclarationSpecifiers != nil {
		typeStr = g.typeSpecToString(fd.DeclarationSpecifiers)
	}
	params := g.extractParams(fd.Declarator)
	g.funcs[name] = &FuncInfo{
		Name:       name,
		ReturnType: cTypeToGo(typeStr),
		Params:     params,
	}
}

func (g *Generator) processEnumerators(el *cc.EnumeratorList) {
	for el != nil {
		if el.Enumerator != nil {
			name := ""
			if el.Enumerator.Case == cc.EnumeratorIdent && string(el.Enumerator.Token.Src()) != "" {
				name = string(el.Enumerator.Token.Src())
			} else if el.Enumerator.Case == cc.EnumeratorExpr && string(el.Enumerator.Token.Src()) != "" {
				name = string(el.Enumerator.Token.Src())
			}
			if name != "" && !isReserved(name) {
				val := fmt.Sprintf("%d", len(g.macros))
				if _, ok := g.macros[name]; !ok {
					g.macros[name] = val
				}
			}
		}
		el = el.EnumeratorList
	}
}

func (g *Generator) getDeclaratorName(d *cc.Declarator) string {
	if d == nil || d.DirectDeclarator == nil {
		return ""
	}
	dd := d.DirectDeclarator
	for dd != nil {
		switch dd.Case {
		case cc.DirectDeclaratorIdent:
			if string(dd.Token.Src()) != "" {
				return string(dd.Token.Src())
			}
		case cc.DirectDeclaratorDecl:
			dd = dd.DirectDeclarator
			continue
		case cc.DirectDeclaratorFuncParam:
			if dd.DirectDeclarator != nil {
				dd = dd.DirectDeclarator
				continue
			}
		case cc.DirectDeclaratorFuncIdent:
			if dd.DirectDeclarator != nil {
				dd = dd.DirectDeclarator
				continue
			}
		}
		break
	}
	return ""
}

func (g *Generator) getStructOrUnionName(sus *cc.StructOrUnionSpecifier) string {
	if sus == nil {
		return ""
	}
	if string(sus.Token.Src()) != "" {
		return string(sus.Token.Src())
	}
	return ""
}

func (g *Generator) getEnumName(es *cc.EnumSpecifier) string {
	if es == nil {
		return ""
	}
	if es.Case == cc.EnumSpecifierTag && string(es.Token.Src()) != "" {
		return string(es.Token.Src())
	}
	return ""
}

func (g *Generator) findTypeSpecifier(ds *cc.DeclarationSpecifiers) *cc.TypeSpecifier {
	for ds != nil {
		if ds.TypeSpecifier != nil {
			return ds.TypeSpecifier
		}
		ds = ds.DeclarationSpecifiers
	}
	return nil
}

func (g *Generator) isTypedef(ds *cc.DeclarationSpecifiers) bool {
	for ds != nil {
		if ds.StorageClassSpecifier != nil {
			s := string(ds.StorageClassSpecifier.Token.Src())
			if s == "typedef" {
				return true
			}
		}
		ds = ds.DeclarationSpecifiers
	}
	return false
}

func (g *Generator) hasFunctionParameters(d *cc.Declarator) bool {
	if d == nil || d.DirectDeclarator == nil {
		return false
	}
	dd := d.DirectDeclarator
	for dd != nil {
		if dd.Case == cc.DirectDeclaratorFuncParam {
			return true
		}
		if dd.DirectDeclarator != nil {
			dd = dd.DirectDeclarator
		} else {
			break
		}
	}
	return false
}

func (g *Generator) isFunctionPointer(d *cc.Declarator) bool {
	if d == nil || d.DirectDeclarator == nil {
		return false
	}
	dd := d.DirectDeclarator
	for dd != nil {
		if dd.Case == cc.DirectDeclaratorFuncParam && dd.ParameterTypeList != nil {
			ptl := dd.ParameterTypeList
			if ptl != nil && ptl.ParameterList != nil {
				pl := ptl.ParameterList
				for pl != nil {
					if pl.ParameterDeclaration != nil {
						innerDecl := pl.ParameterDeclaration.Declarator
						if innerDecl != nil && innerDecl.Pointer != nil && innerDecl.DirectDeclarator != nil {
							innerDD := innerDecl.DirectDeclarator
							for innerDD != nil {
								if innerDD.Case == cc.DirectDeclaratorFuncParam {
									return true
								}
								if innerDD.DirectDeclarator != nil {
									innerDD = innerDD.DirectDeclarator
								} else {
									break
								}
							}
						}
					}
					pl = pl.ParameterList
				}
			}
		}
		if dd.DirectDeclarator != nil {
			dd = dd.DirectDeclarator
		} else {
			break
		}
	}
	return false
}

func (g *Generator) typeSpecToString(ds *cc.DeclarationSpecifiers) string {
	if ds == nil {
		return "int"
	}
	var parts []string
	isUnsigned := false
	isLong := false
	isLongLong := false
	baseType := "int"
	for ds != nil {
		if ds.TypeSpecifier != nil {
			ts := ds.TypeSpecifier
			switch ts.Case {
			case cc.TypeSpecifierVoid:
				baseType = "void"
			case cc.TypeSpecifierChar:
				baseType = "char"
			case cc.TypeSpecifierShort:
				parts = append(parts, "short")
			case cc.TypeSpecifierInt:
				baseType = "int"
			case cc.TypeSpecifierLong:
				if isLong {
					isLongLong = true
				} else {
					isLong = true
				}
			case cc.TypeSpecifierFloat:
				baseType = "float"
			case cc.TypeSpecifierDouble:
				baseType = "double"
			case cc.TypeSpecifierSigned:
			case cc.TypeSpecifierUnsigned:
				isUnsigned = true
			case cc.TypeSpecifierBool:
				baseType = "_Bool"
			case cc.TypeSpecifierComplex:
				baseType = "complex"
			case cc.TypeSpecifierStructOrUnion:
				if ts.StructOrUnionSpecifier != nil {
					name := g.getStructOrUnionName(ts.StructOrUnionSpecifier)
					if name != "" {
						kind := "struct"
						if ts.StructOrUnionSpecifier.StructOrUnion != nil &&
							string(ts.StructOrUnionSpecifier.StructOrUnion.Token.Src()) == "union" {
							kind = "union"
						}
						baseType = kind + " " + name
					} else {
						baseType = "struct"
					}
				}
			case cc.TypeSpecifierEnum:
				if ts.EnumSpecifier != nil {
					name := g.getEnumName(ts.EnumSpecifier)
					if name != "" {
						baseType = "enum " + name
					} else {
						baseType = "int"
					}
				}
			case cc.TypeSpecifierTypeName:
				if len(ts.Token.Src()) > 0 {
					baseType = string(ts.Token.Src())
				}
			default:
				if len(ts.Token.Src()) > 0 {
					baseType = string(ts.Token.Src())
				}
			}
		}
		ds = ds.DeclarationSpecifiers
	}
	if isUnsigned {
		parts = append([]string{"unsigned"}, parts...)
	}
	if isLongLong {
		parts = append(parts, "long", "long")
	} else if isLong {
		parts = append(parts, "long")
	}
	if baseType != "int" || len(parts) == 0 {
		parts = append(parts, baseType)
	}
	result := strings.Join(parts, " ")
	result = strings.ReplaceAll(result, "long long", "long long")
	result = strings.ReplaceAll(result, "unsigned int", "unsigned")
	result = strings.ReplaceAll(result, "unsigned long", "unsigned long")
	result = strings.ReplaceAll(result, "unsigned long long", "unsigned long long")
	result = strings.ReplaceAll(result, "signed int", "signed")
	result = strings.ReplaceAll(result, "signed long", "signed long")
	result = strings.ReplaceAll(result, "signed long long", "signed long long")
	return result
}

func (g *Generator) extractParams(d *cc.Declarator) []ParamInfo {
	if d == nil || d.DirectDeclarator == nil {
		return nil
	}
	dd := d.DirectDeclarator
	for dd != nil {
		if dd.Case == cc.DirectDeclaratorFuncParam {
			if dd.ParameterTypeList != nil && dd.ParameterTypeList.ParameterList != nil {
				return g.parseParameterList(dd.ParameterTypeList.ParameterList)
			}
			return nil
		}
		if dd.DirectDeclarator != nil {
			dd = dd.DirectDeclarator
		} else {
			break
		}
	}
	return nil
}

func (g *Generator) parseParameterList(pl *cc.ParameterList) []ParamInfo {
	var result []ParamInfo
	for pl != nil {
		if pl.ParameterDeclaration != nil {
			pd := pl.ParameterDeclaration
			if pd.DeclarationSpecifiers != nil && pd.Declarator != nil {
				typeStr := g.typeSpecToString(pd.DeclarationSpecifiers)
				name := g.getDeclaratorName(pd.Declarator)
				goType := cTypeToGo(typeStr)
				if goType == "" {
					goType = "uintptr"
				}
				result = append(result, ParamInfo{Name: name, Type: goType})
			} else if pd.DeclarationSpecifiers != nil {
				typeStr := g.typeSpecToString(pd.DeclarationSpecifiers)
				result = append(result, ParamInfo{Name: "", Type: cTypeToGo(typeStr)})
			}
		}
		pl = pl.ParameterList
	}
	return result
}

func (g *Generator) extractMembers(sus *cc.StructOrUnionSpecifier) []MemberInfo {
	var members []MemberInfo
	if sus == nil || sus.StructDeclarationList == nil {
		return members
	}
	sdl := sus.StructDeclarationList
	for sdl != nil {
		if sdl.StructDeclaration != nil {
			sd := sdl.StructDeclaration
			if sd.SpecifierQualifierList != nil && sd.StructDeclaratorList != nil {
				typeStr := g.specQualListToString(sd.SpecifierQualifierList)
				sdl2 := sd.StructDeclaratorList
				for sdl2 != nil {
					if sdl2.StructDeclarator != nil {
						sd2 := sdl2.StructDeclarator
						memName := ""
						if sd2.Declarator != nil {
							memName = g.getDeclaratorName(sd2.Declarator)
						}
						if memName != "" && !isReserved(memName) {
							members = append(members, MemberInfo{
								Name: memName,
								Type: typeStr,
							})
						}
					}
					sdl2 = sdl2.StructDeclaratorList
				}
			}
		}
		sdl = sdl.StructDeclarationList
	}
	return members
}

func (g *Generator) typeNameToString(tn *cc.TypeName) string {
	if tn == nil || tn.SpecifierQualifierList == nil {
		return ""
	}
	return g.specQualListToString(tn.SpecifierQualifierList)
}

func (g *Generator) specQualListToString(sql *cc.SpecifierQualifierList) string {
	if sql == nil {
		return ""
	}
	var parts []string
	hasType := false
	for sql != nil {
		if sql.TypeSpecifier != nil {
			ts := sql.TypeSpecifier
			switch ts.Case {
			case cc.TypeSpecifierVoid:
				parts = append(parts, "void")
				hasType = true
			case cc.TypeSpecifierChar:
				parts = append(parts, "char")
				hasType = true
			case cc.TypeSpecifierShort:
				parts = append(parts, "short")
			case cc.TypeSpecifierInt:
				parts = append(parts, "int")
				hasType = true
			case cc.TypeSpecifierLong:
				parts = append(parts, "long")
			case cc.TypeSpecifierFloat:
				parts = append(parts, "float")
				hasType = true
			case cc.TypeSpecifierDouble:
				parts = append(parts, "double")
				hasType = true
			case cc.TypeSpecifierSigned:
				parts = append(parts, "signed")
			case cc.TypeSpecifierUnsigned:
				parts = append(parts, "unsigned")
			case cc.TypeSpecifierBool:
				parts = append(parts, "_Bool")
				hasType = true
			case cc.TypeSpecifierStructOrUnion:
				if ts.StructOrUnionSpecifier != nil {
					name := g.getStructOrUnionName(ts.StructOrUnionSpecifier)
					if name != "" {
						kind := "struct"
						if ts.StructOrUnionSpecifier.StructOrUnion != nil &&
							string(ts.StructOrUnionSpecifier.StructOrUnion.Token.Src()) == "union" {
							kind = "union"
						}
						parts = append(parts, kind+" "+name)
						hasType = true
					}
				}
			case cc.TypeSpecifierEnum:
				if ts.EnumSpecifier != nil {
					name := g.getEnumName(ts.EnumSpecifier)
					if name != "" {
						parts = append(parts, "enum "+name)
						hasType = true
					}
				}
			case cc.TypeSpecifierTypeName:
				if len(ts.Token.Src()) > 0 {
					parts = append(parts, string(ts.Token.Src()))
					hasType = true
				}
			default:
				if len(ts.Token.Src()) > 0 {
					parts = append(parts, string(ts.Token.Src()))
				}
			}
		}
		sql = sql.SpecifierQualifierList
	}
	if !hasType {
		parts = append(parts, "int")
	}
	return strings.Join(parts, " ")
}

func cTypeToGo(cType string) string {
	cType = strings.TrimSpace(cType)
	cType = strings.ReplaceAll(cType, "\t", " ")
	for strings.Contains(cType, "  ") {
		cType = strings.ReplaceAll(cType, "  ", " ")
	}
	cType = strings.TrimPrefix(cType, "const ")
	cType = strings.TrimPrefix(cType, "volatile ")
	cType = strings.TrimPrefix(cType, "struct ")
	cType = strings.TrimPrefix(cType, "union ")
	cType = strings.TrimPrefix(cType, "enum ")
	cType = strings.TrimPrefix(cType, "unsigned ")
	cType = strings.TrimPrefix(cType, "signed ")
	cType = strings.TrimSpace(cType)
	typeMap := map[string]string{
		"void":            "",
		"char":            "int8",
		"int":             "int32",
		"short":           "int16",
		"long":            "intptr_t",
		"long long":       "int64",
		"float":           "float32",
		"double":          "float64",
		"_Bool":           "bool",
		"size_t":          "uint64",
		"ssize_t":         "int64",
		"ptrdiff_t":       "int64",
		"wchar_t":         "uint16",
		"uintptr_t":       "uintptr",
		"intptr_t":        "uintptr",
		"int8_t":          "int8",
		"uint8_t":         "uint8",
		"int16_t":         "int16",
		"uint16_t":        "uint16",
		"int32_t":         "int32",
		"uint32_t":        "uint32",
		"int64_t":         "int64",
		"uint64_t":        "uint64",
		"__int8":          "int8",
		"__int16":         "int16",
		"__int32":         "int32",
		"__int64":         "int64",
		"HANDLE":          "uintptr",
		"PVOID":           "uintptr",
		"VOID":            "",
		"BYTE":            "uint8",
		"WORD":            "uint16",
		"DWORD":           "uint32",
		"DWORD32":         "uint32",
		"DWORD64":         "uint64",
		"BOOL":            "int32",
		"BOOLEAN":         "byte",
		"CHAR":            "int8",
		"UCHAR":           "uint8",
		"SHORT":           "int16",
		"USHORT":          "uint16",
		"INT":             "int32",
		"UINT":            "uint32",
		"LONG":            "int32",
		"ULONG":           "uint32",
		"LONGLONG":        "int64",
		"ULONGLONG":       "uint64",
		"LONG_PTR":        "intptr",
		"ULONG_PTR":       "uintptr",
		"LONG64":          "int64",
		"ULONG64":         "uint64",
		"SIZE_T":          "uint64",
		"SSIZE_T":         "int64",
		"NTSTATUS":        "int32",
		"HRESULT":         "int32",
		"WCHAR":           "uint16",
		"TCHAR":           "uint16",
		"LANGID":          "uint16",
		"LCID":            "uint32",
		"COLORREF":        "uint32",
		"LRESULT":         "intptr",
		"WPARAM":          "uintptr",
		"LPARAM":          "intptr",
		"ATOM":            "uint16",
		"HINSTANCE":       "uintptr",
		"HWND":            "uintptr",
		"HDC":             "uintptr",
		"HMENU":           "uintptr",
		"HICON":           "uintptr",
		"HBRUSH":          "uintptr",
		"HCURSOR":         "uintptr",
		"HFONT":           "uintptr",
		"HBITMAP":         "uintptr",
		"HRGN":            "uintptr",
		"HKL":             "uintptr",
		"PROC":            "uintptr",
		"FARPROC":         "uintptr",
		"NEARPROC":        "uintptr",
		"WNDPROC":         "uintptr",
		"KAFFINITY":       "uintptr",
		"KAFFINITY_EX":    "uint64",
		"ACCESS_MASK":     "uint32",
		"ULONG32":         "uint32",
		"LONG32":          "int32",
		"INT_PTR":         "intptr",
		"UINT_PTR":        "uintptr",
		"POINTER_64":      "uint64",
		"POINTER_32":      "uint32",
		"VOID*":           "uintptr",
		"PVOID*":          "uintptr",
		"void*":           "uintptr",
		"char*":           "uintptr",
		"int*":            "uintptr",
		"UCHAR*":          "uintptr",
		"USHORT*":         "uintptr",
		"ULONG*":          "uintptr",
		"DWORD*":          "uintptr",
		"BYTE*":           "uintptr",
		"WORD*":           "uintptr",
		"LONG*":           "uintptr",
		"BOOL*":           "uintptr",
		"HANDLE*":         "uintptr",
		"LPCSTR":          "uintptr",
		"LPSTR":           "uintptr",
		"LPCWSTR":         "uintptr",
		"LPWSTR":          "uintptr",
		"LPCVOID":         "uintptr",
		"LPVOID":          "uintptr",
		"LPCTSTR":         "uintptr",
		"LPTSTR":          "uintptr",
		"PSTR":            "uintptr",
		"PCSTR":           "uintptr",
		"PWSTR":           "uintptr",
		"PCWSTR":          "uintptr",
		"PTSTR":           "uintptr",
		"PCTSTR":          "uintptr",
		"PBOOL":           "uintptr",
		"PBOOLEAN":        "uintptr",
		"PBYTE":           "uintptr",
		"PCHAR":           "uintptr",
		"PCHAR8":          "uintptr",
		"PUCHAR":          "uintptr",
		"PSHORT":          "uintptr",
		"PSHORT16":        "uintptr",
		"PUSHORT":         "uintptr",
		"PLONG":           "uintptr",
		"PLONG32":         "uintptr",
		"PLONG64":         "uintptr",
		"PULONG":          "uintptr",
		"PULONG32":        "uintptr",
		"PULONG64":        "uintptr",
		"PHANDLE":         "uintptr",
		"PHALF_PTR":       "uintptr",
		"PINT":            "uintptr",
		"PINT8":           "uintptr",
		"PINT16":          "uintptr",
		"PINT32":          "uintptr",
		"PINT64":          "uintptr",
		"PUINT":           "uintptr",
		"PUINT8":          "uintptr",
		"PUINT16":         "uintptr",
		"PUINT32":         "uintptr",
		"PUINT64":         "uintptr",
		"PLONGLONG":       "uintptr",
		"PULONGLONG":      "uintptr",
		"PDWORD":          "uintptr",
		"PDWORD32":        "uintptr",
		"PDWORD64":        "uintptr",
		"PWORD":           "uintptr",
		"PDWORDLONG":      "uintptr",
		"PLARGE_INTEGER":  "uintptr",
		"PEARLY_STRUCT":   "uintptr",
		"PULARGE_INTEGER": "uintptr",
		"PCOLORREF":       "uintptr",
		"PHKEY":           "uintptr",
		"KUINT_PTR":       "uintptr",
		"KPOINTER":        "uintptr",
		"PACCESS_TOKEN":   "uintptr",
		"PEPROCESS":       "uintptr",
		"PETHREAD":        "uintptr",
	}
	if goType, ok := typeMap[cType]; ok {
		return goType
	}
	if strings.HasSuffix(cType, "*") {
		base := cType[:len(cType)-1]
		base = strings.TrimSpace(base)
		if baseType, ok := typeMap[base]; ok {
			if baseType == "" {
				return "unsafe.Pointer"
			}
			return "*" + baseType
		}
		return "*"
	}
	if strings.ContainsAny(cType, "[]") {
		return ""
	}
	if strings.Contains(cType, "(") {
		return ""
	}
	if cType == "" {
		return ""
	}
	if unicode.IsUpper(rune(cType[0])) || (len(cType) > 1 && cType[0] == '_' && unicode.IsLetter(rune(cType[1]))) {
		return cType
	}
	return ""
}

func (g *Generator) writeDriverGo() {
	var sb strings.Builder
	sb.WriteString("//so:include <ntddk.h>\n//so:include <ntifs.h>\n\npackage wdk\n\nimport \"unsafe\"\n\n")

	sb.WriteString("func CTL_CODE(DeviceType uint32, Function uint32, Method uint32, Access uint32) uint32 {\n\treturn ((DeviceType) << 16) | ((Access) << 14) | ((Function) << 2) | (Method)\n}\n\n")

	os.WriteFile(filepath.Join(g.outDir, "driver.go"), []byte(sb.String()), 0o644)
}

func (g *Generator) writeFileBindings(goFileName string, globalGen *Generator) {
	var sb strings.Builder
	sb.WriteString("package wdk\n\nimport \"unsafe\"\n\n")

	newTypes := 0
	newFuncs := 0
	newMacros := 0

	names := make([]string, 0, len(g.types))
	for name := range g.types {
		if !globalGen.seenTypes[name] {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		ti := g.types[name]
		if !isValidGoType(name, ti.GoType) {
			continue
		}
		switch ti.Kind {
		case "typedef":
			sb.WriteString(fmt.Sprintf("type %s = %s\n", name, ti.GoType))
		case "struct":
			if ti.IsOpaque || len(ti.Members) == 0 {
				sb.WriteString(fmt.Sprintf("type %s struct { _ [0]byte }\n", name))
			} else {
				sb.WriteString(fmt.Sprintf("type %s struct {\n", name))
				for _, m := range ti.Members {
					goType := cTypeToGo(m.Type)
					if goType == "" {
						goType = "uintptr"
					}
					fieldName := m.Name
					if fieldName == "" {
						fieldName = "_"
					}
					sb.WriteString(fmt.Sprintf("\t%s %s\n", fieldName, goType))
				}
				sb.WriteString("}\n")
			}
		case "union":
			sb.WriteString(fmt.Sprintf("type %s struct { _ [0]byte }\n", name))
		case "enum":
			sb.WriteString(fmt.Sprintf("type %s = %s\n", name, ti.GoType))
		default:
			sb.WriteString(fmt.Sprintf("type %s = %s\n", name, ti.GoType))
		}
		globalGen.seenTypes[name] = true
		newTypes++
	}

	funcNames := make([]string, 0, len(g.funcs))
	for name := range g.funcs {
		if !globalGen.seenFuncs[name] {
			funcNames = append(funcNames, name)
		}
	}
	sort.Strings(funcNames)
	if len(funcNames) > 0 {
		sb.WriteString("\n")
	}
	for _, name := range funcNames {
		fi := g.funcs[name]
		retType := fi.ReturnType
		if retType == "" {
			retType = "int32"
		}
		var paramStr string
		if len(fi.Params) > 0 {
			var parts []string
			for _, p := range fi.Params {
				goType := p.Type
				if goType == "" {
					goType = "uintptr"
				}
				parts = append(parts, fmt.Sprintf("%s %s", p.Name, goType))
			}
			paramStr = strings.Join(parts, ", ")
		}
		sb.WriteString(fmt.Sprintf("//so:extern func %s(%s) %s\n", name, paramStr, retType))
		globalGen.seenFuncs[name] = true
		newFuncs++
	}

	macroNames := make([]string, 0, len(g.macros))
	for name := range g.macros {
		if !globalGen.seenMacros[name] && isWDKConstant(name) {
			macroNames = append(macroNames, name)
		}
	}
	sort.Strings(macroNames)
	if len(macroNames) > 0 {
		sb.WriteString("\nconst (\n")
		for _, name := range macroNames {
			val := g.macros[name]
			val = cleanMacroValue(val)
			sb.WriteString(fmt.Sprintf("\t%s = %s\n", name, val))
			globalGen.seenMacros[name] = true
			newMacros++
		}
		sb.WriteString(")\n")
	}

	if newTypes == 0 && newFuncs == 0 && newMacros == 0 {
		return
	}

	outPath := filepath.Join(g.outDir, goFileName)
	os.WriteFile(outPath, []byte(sb.String()), 0o644)
	fmt.Fprintf(os.Stderr, "  -> %s (+%d types, +%d funcs, +%d macros)\n", goFileName, newTypes, newFuncs, newMacros)
}

func isWDKConstant(name string) bool {
	if name == "TRUE" || name == "FALSE" || name == "NULL" || name == "NO_ERROR" || name == "ERROR_SUCCESS" {
		return true
	}
	prefixes := []string{
		"STATUS_", "NT_", "FILE_", "IOCTL_",
		"WMIGUID_", "GUID_", "PCI_", "USB_",
		"NDIS_", "SCSI_", "IRP_", "SL_",
		"DEBUG_", "PF_", "Se", "Ps",
		"Ke", "Mm", "Io", "Ex", "Ob", "Zw",
		"Rtl", "Hal", "Po", "Cm", "FsRtl", "Cc",
		"Pp", "Pi", "Iop", "Inbv", "Wd", "Dbg", "Kd", "Vf",
		"VER_", "WINVER", "NTDDI", "WINNT", "NTINST", "NT_WIN",
		"DBG", "POOL_", "MM_", "FLG_", "KPCR_", "PRCB_",
		"IRQL_", "PASSIVE_", "APC_", "DISPATCH_", "HIGH_",
		"IMAGE_", "TLS_", "MSVCRT_", "CRT_", "UCRTBASE_",
		"__WIDL_", "__MINGW_", "_MSC_VER", "_ATL_",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	if strings.HasPrefix(name, "_") && len(name) > 4 {
		return true
	}
	return false
}

func cleanMacroValue(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "0"
	}
	if idx := strings.Index(v, "//"); idx >= 0 {
		v = strings.TrimSpace(v[:idx])
	}
	castRe := regexp.MustCompile(`^\(\([A-Za-z_][A-Za-z0-9_]*\s*\)\s*(.+?)\)`)
	for {
		m := castRe.FindStringSubmatch(v)
		if m == nil {
			break
		}
		v = strings.TrimSpace(m[1])
	}
	singleCastRe := regexp.MustCompile(`^\([A-Za-z_][A-Za-z0-9_]*\s*\)\s*(.+)`)
	for {
		m := singleCastRe.FindStringSubmatch(v)
		if m == nil {
			break
		}
		v = strings.TrimSpace(m[2])
	}
	v = strings.ReplaceAll(v, "UL", "")
	v = strings.ReplaceAll(v, "ULL", "")
	v = strings.TrimSuffix(v, "U")
	v = strings.TrimSuffix(v, "LL")
	v = strings.TrimSuffix(v, "L")
	v = strings.TrimSuffix(v, "l")
	v = strings.TrimSuffix(v, "u64")
	v = strings.TrimSuffix(v, "u32")
	v = strings.TrimSuffix(v, "u16")
	v = strings.TrimSuffix(v, "u8")
	v = strings.ReplaceAll(v, "(void*)0", "0")
	v = strings.ReplaceAll(v, "((ULONG)-1)", "0xFFFFFFFF")
	v = strings.ReplaceAll(v, "((LONG)-1)", "0xFFFFFFFF")
	v = strings.ReplaceAll(v, "((unsigned int)-1)", "0xFFFFFFFF")
	v = strings.ReplaceAll(v, "((unsigned long)-1)", "0xFFFFFFFF")
	v = strings.ReplaceAll(v, "((UINT)-1)", "0xFFFFFFFF")
	v = strings.ReplaceAll(v, "(ULONG)(LONG_PTR)-1", "0xFFFFFFFFFFFFFFFF")
	v = strings.ReplaceAll(v, "(1i64<<63)", "0x8000000000000000")
	v = strings.ReplaceAll(v, "(1i64", "(1")
	v = strings.ReplaceAll(v, "i64", "")
	v = strings.ReplaceAll(v, "u64", "")
	v = strings.ReplaceAll(v, "u32", "")
	v = strings.ReplaceAll(v, "u16", "")
	v = strings.ReplaceAll(v, "u8", "")
	v = strings.TrimSpace(v)
	if isComplexExpression(v) {
		return "0"
	}
	return v
}

func stripCast(s string) string {
	s = strings.TrimSpace(s)
	for {
		if !strings.HasPrefix(s, "(") {
			return s
		}
		re := regexp.MustCompile(`^\(([A-Za-z_]\w*)\)\s*(.*)$`)
		m := re.FindStringSubmatch(s)
		if m != nil && isCTypename(m[1]) {
			s = strings.TrimSpace(m[2])
			continue
		}
		re2 := regexp.MustCompile(`^\(\(([A-Za-z_]\w*)\)\s*(.*)\)\s*$`)
		m2 := re2.FindStringSubmatch(s)
		if m2 != nil && isCTypename(m2[1]) {
			s = strings.TrimSpace(m2[2])
			continue
		}
		re3 := regexp.MustCompile(`^\(\(\w[\w\s]*\)\s*[^)]+\)$`)
		if re3.MatchString(s) {
			inner := s[1 : len(s)-1]
			s = inner
			continue
		}
		return s
	}
}

func isCTypename(s string) bool {
	s = strings.TrimSpace(s)
	typeKeywords := map[string]bool{
		"int": true, "char": true, "short": true, "long": true,
		"float": true, "double": true, "void": true, "bool": true,
		"signed": true, "unsigned": true, "const": true, "volatile": true,
		"ULONG": true, "LONG": true, "USHORT": true, "SHORT": true,
		"UCHAR": true, "CHAR": true, "BOOL": true, "BOOLEAN": true,
		"DWORD": true, "WORD": true, "BYTE": true, "HANDLE": true,
		"PVOID": true, "INT": true, "UINT": true, "SIZE_T": true,
		"NTSTATUS": true, "HRESULT": true,
		"WCHAR": true, "INT8": true, "INT16": true, "INT32": true, "INT64": true,
		"UINT8": true, "UINT16": true, "UINT32": true, "UINT64": true,
		"intptr_t": true, "uintptr_t": true, "size_t": true, "ssize_t": true,
		"ptrdiff_t": true, "wchar_t": true, "int8_t": true, "uint8_t": true,
		"int16_t": true, "uint16_t": true, "int32_t": true, "uint32_t": true,
		"int64_t": true, "uint64_t": true, "LONGLONG": true, "ULONGLONG": true,
		"ATOM": true, "LANGID": true, "LCID": true, "COLORREF": true,
		"LRESULT": true, "WPARAM": true, "LPARAM": true, "HINSTANCE": true,
		"HWND": true, "HDC": true, "HMENU": true, "HICON": true, "HBRUSH": true,
		"HCURSOR": true, "HFONT": true, "HBITMAP": true, "HRGN": true, "HKL": true,
		"PROC": true, "FARPROC": true, "NEARPROC": true, "WNDPROC": true,
		"ULONG_PTR": true, "LONG_PTR": true, "KUINT_PTR": true, "KPOINTER": true,
		"PACCESS_TOKEN": true, "PEPROCESS": true, "PETHREAD": true,
		"PHYSICAL_ADDRESS": true, "LARGE_INTEGER": true, "ULARGE_INTEGER": true,
		"KAFFINITY": true, "PKERNEL_ROUTINE": true, "MODE": true,
		"ONG": true, "ONGL": true, "ONG64": true,
	}
	if typeKeywords[s] {
		return true
	}
	return false
}

func isValidGoConst(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" || v == "0" {
		return true
	}
	if regexp.MustCompile(`^-?0x[0-9a-fA-F]+$`).MatchString(v) {
		return true
	}
	if regexp.MustCompile(`^-?[0-9]+$`).MatchString(v) {
		return true
	}
	return false
}

func isValidGoType(name, goType string) bool {
	name = strings.TrimSpace(name)
	goType = strings.TrimSpace(goType)
	if name == "" || goType == "" {
		return false
	}
	if strings.Contains(name, ",") || strings.Contains(name, " ") || strings.Contains(name, "(") || strings.Contains(name, ")") {
		return false
	}
	if strings.Contains(goType, ",") {
		return false
	}
	if strings.Contains(goType, "(") && !strings.HasSuffix(goType, ")") {
		return false
	}
	validIdent := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	if !validIdent.MatchString(name) {
		return false
	}
	return true
}

func isComplexExpression(v string) bool {
	if strings.Contains(v, "|") && !strings.HasPrefix(v, "0x") {
		return true
	}
	if strings.Contains(v, "<<") || strings.Contains(v, ">>") {
		return true
	}
	if strings.Contains(v, "{") || strings.Contains(v, "}") {
		return true
	}
	if len(v) > 60 {
		return true
	}
	if strings.Contains(v, "sizeof") {
		return true
	}
	if strings.Contains(v, ";") {
		return true
	}
	if strings.Contains(v, "_asm") || strings.Contains(v, "__asm") {
		return true
	}
	if strings.Contains(v, "return") {
		return true
	}
	if regexp.MustCompile(`^\w+\s+\w+$`).MatchString(v) && !regexp.MustCompile(`^0x`).MatchString(v) {
		return true
	}
	return false
}
