package main

import (
	"strings"
	"testing"
)

func Benchmark_WriteB_AutoGrow(b *testing.B) {
	someBytes := []byte(someStr)
	b.ReportAllocs()
	for b.Loop() {
		var buf strings.Builder
		for range numWrite {
			buf.Write(someBytes)
		}
		sink = buf.String()
	}
}

func Benchmark_WriteB_PreGrow(b *testing.B) {
	someBytes := []byte(someStr)
	b.ReportAllocs()
	for b.Loop() {
		var buf strings.Builder
		buf.Grow(len(someBytes) * numWrite)
		for range numWrite {
			buf.Write(someBytes)
		}
		sink = buf.String()
	}
}

func Benchmark_WriteS_AutoGrow(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		var buf strings.Builder
		for range numWrite {
			buf.WriteString(someStr)
		}
		sink = buf.String()
	}
}

func Benchmark_WriteS_PreGrow(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		var buf strings.Builder
		buf.Grow(len(someStr) * numWrite)
		for range numWrite {
			buf.WriteString(someStr)
		}
		sink = buf.String()
	}
}
