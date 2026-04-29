// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings_test

import (
	"testing"

	. "solod.dev/so/strings"
)

var longString = "a" + string(make([]byte, 1<<16)) + "z"

var longSpaces = func() string {
	b := make([]byte, 200)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}()

var RepeatTests = []struct {
	in, out string
	count   int
}{
	{"", "", 0},
	{"", "", 1},
	{"", "", 2},
	{"-", "", 0},
	{"-", "-", 1},
	{"-", "----------", 10},
	{"abc ", "abc abc abc ", 3},
	{" ", " ", 1},
	{"--", "----", 2},
	{"===", "======", 2},
	{"000", "000000000", 3},
	{"\t\t\t\t", "\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t", 4},
	{" ", longSpaces, len(longSpaces)},
	// Tests for results over the chunkLimit
	{string(rune(0)), string(make([]byte, 1<<16)), 1 << 16},
	{longString, longString + longString, 2},
}

func TestRepeat(t *testing.T) {
	for _, tt := range RepeatTests {
		a := Repeat(nil, tt.in, tt.count)
		if !equal("Repeat(s)", a, tt.out, t) {
			t.Errorf("Repeat(%v, %d) = %v; want %v", tt.in, tt.count, a, tt.out)
			continue
		}
	}
}

// See Issue golang.org/issue/16237
func TestRepeatCatchesOverflow(t *testing.T) {
	type testCase struct {
		s      string
		count  int
		errStr string
	}

	runTestCases := func(prefix string, tests []testCase) {
		for i, tt := range tests {
			errStr := repeat(tt.s, tt.count)
			if tt.errStr == "" {
				if errStr != "" {
					t.Errorf("#%d panicked %v", i, errStr)
				}
				continue
			}

			if errStr == "" || !Contains(errStr, tt.errStr) {
				t.Errorf("%s#%d got %q want %q", prefix, i, errStr, tt.errStr)
			}
		}
	}

	const maxInt = int(^uint(0) >> 1)

	runTestCases("", []testCase{
		0: {"--", -2147483647, "negative"},
		1: {"", maxInt, ""},
		2: {"-", 10, ""},
		3: {"gopher", 0, ""},
		4: {"-", -1, "negative"},
		5: {"--", -102, "negative"},
		6: {string(make([]byte, 255)), int((^uint(0))/255 + 1), "overflow"},
	})

	const is64Bit = 1<<(^uintptr(0)>>63)/2 != 0
	if !is64Bit {
		return
	}

	runTestCases("64-bit", []testCase{
		0: {"-", maxInt, "capacity overflow"},
	})
}

func repeat(s string, count int) (err string) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v.Error()
			default:
				err = v.(string)
			}
		}
	}()

	Repeat(nil, s, count)

	return
}
