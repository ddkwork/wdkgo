package main

import (
	"strings"
	"testing"
)

func Benchmark_Clone(b *testing.B) {
	b.ReportAllocs()
	var str = strings.Repeat("a", 1024)
	for b.Loop() {
		sink = strings.Clone(str)
	}
}

func Benchmark_Compare(b *testing.B) {
	b.ReportAllocs()
	str1 := strings.Repeat("01234567890αβγδεζ", 64)
	str2 := strings.Repeat("01234567890αβγδεζ", 64)
	for b.Loop() {
		sinkInt = strings.Compare(str1, str2)
	}
}

func Benchmark_Fields(b *testing.B) {
	b.ReportAllocs()
	str := strings.Repeat("01234567890αβ γδεζ", 16)
	for b.Loop() {
		fields := strings.Fields(str)
		sink = fields[0]
	}
}

func Benchmark_Index(b *testing.B) {
	b.ReportAllocs()
	var sb strings.Builder
	for range 64 {
		sb.WriteString("01234567890αβγδεζ")
	}
	sb.WriteRune('ω')
	str := sb.String() // 1025 chars, search for last
	for b.Loop() {
		sinkInt = strings.Index(str, "ω")
	}
}

func Benchmark_IndexByte(b *testing.B) {
	b.ReportAllocs()
	var sb strings.Builder
	for range 64 {
		sb.WriteString("01234567890αβγδεζ")
	}
	sb.WriteByte('X')
	str := sb.String() // 1025 chars, search for last
	for b.Loop() {
		sinkInt = strings.Index(str, "X")
	}
}

func Benchmark_Repeat(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink = strings.Repeat("0123456789abcdef", 64)
	}
}

func Benchmark_ReplaceAll(b *testing.B) {
	b.ReportAllocs()
	str := strings.Repeat("0123456789abcdef", 16)
	for b.Loop() {
		sink = strings.ReplaceAll(str, "a", "A")
	}
}

func Benchmark_Split(b *testing.B) {
	b.ReportAllocs()
	str := strings.Repeat("01234567890αβγδεζ", 16)
	for b.Loop() {
		fields := strings.Split(str, "γ")
		sink = fields[0]
	}
}

func Benchmark_ToUpper(b *testing.B) {
	b.ReportAllocs()
	str := strings.Repeat("01234567890αβγδεζ", 16)
	for b.Loop() {
		sink = strings.ToUpper(str)
	}
}

func Benchmark_Trim(b *testing.B) {
	b.ReportAllocs()
	var sb strings.Builder
	sb.WriteString("ηθικλμνξοπρστυφχψω")
	for range 64 {
		sb.WriteString("01234567890αβγδεζ")
	}
	sb.WriteString("ηθικλμνξοπρστυφχψω")
	str := sb.String()
	for b.Loop() {
		sink = strings.Trim(str, "ωψχφυτσρποξνμλκιθη")
	}
}

func Benchmark_TrimSuffix(b *testing.B) {
	b.ReportAllocs()
	str := strings.Repeat("01234567890αβγδεζ", 16)
	suffix := "01234567890αβγδεζ"
	for b.Loop() {
		sink = strings.TrimSuffix(str, suffix)
	}
}
