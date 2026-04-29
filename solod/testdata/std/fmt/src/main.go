package main

import (
	"solod.dev/so/fmt"
	"solod.dev/so/strings"
)

func main() {
	{
		// Print.
		n, err := fmt.Print("hello", "world")
		if err != nil {
			panic("Print failed")
		}
		if n != 11 {
			panic("Print: wrong count")
		}
		fmt.Print("\n")
	}
	{
		// Println.
		n, err := fmt.Println("hello", "world")
		if err != nil {
			panic("Println failed")
		}
		if n != 12 {
			panic("Println: wrong count")
		}
	}
	{
		// Printf.
		s := "world"
		d := 42
		n, err := fmt.Printf("s = %s, d = %d\n", s, d)
		if err != nil {
			panic("Printf failed")
		}
		if n != 18 {
			panic("Printf: wrong count")
		}
	}
	{
		// Sprintf.
		buf := fmt.NewBuffer(32)
		s := "world"
		d := 42
		out := fmt.Sprintf(buf, "s = %s, d = %d", s, d)
		if out != "s = world, d = 42" {
			panic("Sprintf: wrong output")
		}
	}
	{
		// Fprintf.
		var sb strings.Builder
		var i int32 = 42
		s := "world"
		n, err := fmt.Fprintf(&sb, "hello %d %s", i, s)
		if err != nil {
			panic("Fprintf failed")
		}
		if n != 14 {
			panic("Fprintf: wrong count")
		}
		if sb.String() != "hello 42 world" {
			panic("Fprintf: wrong output")
		}
		sb.Free()
	}
	{
		// Sscanf.
		var n1, n2 int32
		buf := fmt.NewBuffer(32)
		n, err := fmt.Sscanf("5 1 gophers", "%d %d %s", &n1, &n2, buf.Ptr)
		if err != nil {
			panic("Sscanf failed")
		}
		s := buf.String()
		if n != 3 {
			panic("Sscanf: wrong count")
		}
		if n1 != 5 || n2 != 1 || s != "gophers" {
			panic("Sscanf: wrong values")
		}
	}
	{
		// Fscanf.
		var n1, n2 int32
		buf := fmt.NewBuffer(32)
		r := strings.NewReader("5 1 gophers")
		n, err := fmt.Fscanf(&r, "%d %d %s", &n1, &n2, buf.Ptr)
		if err != nil {
			panic("Fscanf failed")
		}
		s := buf.String()
		if n != 3 {
			panic("Fscanf: wrong count")
		}
		if n1 != 5 || n2 != 1 || s != "gophers" {
			panic("Fscanf: wrong values")
		}
	}
}
