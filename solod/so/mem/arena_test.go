package mem

import (
	"testing"

	"github.com/nalgeon/be"
)

func TestArena(t *testing.T) {
	t.Run("Alloc", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			buf := make([]byte, 256)
			a := NewArena(buf)
			p, err := a.Alloc(16, 8)
			be.Err(t, err, nil)
			if p == nil {
				t.Fatal("want non-nil pointer")
			}
		})
		t.Run("invalid size", func(t *testing.T) {
			buf := make([]byte, 256)
			a := NewArena(buf)
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("want panic")
				}
			}()
			_, _ = a.Alloc(0, 8)
		})
		t.Run("invalid alignment", func(t *testing.T) {
			buf := make([]byte, 256)
			a := NewArena(buf)
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("want panic")
				}
			}()
			_, _ = a.Alloc(16, 3)
		})
		t.Run("out of memory", func(t *testing.T) {
			buf := make([]byte, 16)
			a := NewArena(buf)
			_, err := a.Alloc(32, 8)
			be.Err(t, err, ErrOutOfMemory)
		})
		t.Run("exact fit", func(t *testing.T) {
			buf := make([]byte, 256)
			a := NewArena(buf)
			p, err := a.Alloc(256, 1)
			be.Err(t, err, nil)
			if p == nil {
				t.Fatal("want non-nil pointer")
			}
			be.Equal(t, a.offset, 256)
		})
		t.Run("alignment", func(t *testing.T) {
			buf := make([]byte, 256)
			a := NewArena(buf)

			// Allocate 1 byte, then 8-aligned.
			// Offset should advance past alignment padding.
			_, err := a.Alloc(1, 1)
			be.Err(t, err, nil)
			be.Equal(t, a.offset, 1)

			_, err = a.Alloc(8, 8)
			be.Err(t, err, nil)
			be.Equal(t, a.offset, 16) // 1 -> aligned to 8, then +8
		})
	})

	t.Run("Realloc", func(t *testing.T) {
		t.Run("last grow", func(t *testing.T) {
			buf := make([]byte, 256)
			a := NewArena(buf)
			p, _ := a.Alloc(16, 8)
			buf[0] = 0xAA
			buf[15] = 0xBB

			p2, err := a.Realloc(p, 16, 32, 8)
			be.Err(t, err, nil)
			be.Equal(t, p2, p) // same pointer, extended in place
			be.Equal(t, a.offset, 32)
			be.Equal(t, buf[0], byte(0xAA))
			be.Equal(t, buf[15], byte(0xBB))
			be.Equal(t, buf[16], byte(0)) // new bytes zeroed
		})
		t.Run("last shrink", func(t *testing.T) {
			buf := make([]byte, 256)
			a := NewArena(buf)
			p, _ := a.Alloc(32, 8)

			p2, err := a.Realloc(p, 32, 16, 8)
			be.Err(t, err, nil)
			be.Equal(t, p2, p)
			be.Equal(t, a.offset, 16)
		})
		t.Run("last out of memory", func(t *testing.T) {
			buf := make([]byte, 32)
			a := NewArena(buf)
			p, _ := a.Alloc(16, 8)

			_, err := a.Realloc(p, 16, 64, 8)
			be.Err(t, err, ErrOutOfMemory)
		})
		t.Run("non-last grow", func(t *testing.T) {
			buf := make([]byte, 256)
			a := NewArena(buf)
			p1, _ := a.Alloc(16, 8)
			buf[0] = 0xAA
			buf[15] = 0xBB
			_, _ = a.Alloc(16, 8) // p2, now p1 is not last

			p3, err := a.Realloc(p1, 16, 32, 8)
			be.Err(t, err, nil)
			if p3 == p1 {
				t.Fatal("want different pointer")
			}
			// Verify data was copied.
			be.Equal(t, buf[a.lastStart], byte(0xAA))
			be.Equal(t, buf[a.lastStart+15], byte(0xBB))
		})
		t.Run("non-last shrink", func(t *testing.T) {
			buf := make([]byte, 256)
			a := NewArena(buf)
			p1, _ := a.Alloc(32, 8)
			_, _ = a.Alloc(16, 8) // now p1 is not last
			prevOffset := a.offset

			p2, err := a.Realloc(p1, 32, 16, 8)
			be.Err(t, err, nil)
			be.Equal(t, p2, p1) // same pointer, no new allocation
			be.Equal(t, a.offset, prevOffset)
		})
		t.Run("invalid size", func(t *testing.T) {
			buf := make([]byte, 256)
			a := NewArena(buf)
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("want panic")
				}
			}()
			_, _ = a.Realloc(nil, 0, 16, 8)
		})
		t.Run("invalid alignment", func(t *testing.T) {
			buf := make([]byte, 256)
			a := NewArena(buf)
			defer func() {
				if r := recover(); r == nil {
					t.Fatal("want panic")
				}
			}()
			_, _ = a.Realloc(nil, 16, 32, 3)
		})
	})

	t.Run("Free", func(t *testing.T) {
		t.Run("noop", func(t *testing.T) {
			buf := make([]byte, 256)
			a := NewArena(buf)
			p, _ := a.Alloc(16, 8)
			a.Free(p, 16, 8) // should not panic
		})
	})

	t.Run("Reset", func(t *testing.T) {
		buf := make([]byte, 256)
		a := NewArena(buf)

		_, err := a.Alloc(128, 8)
		be.Err(t, err, nil)
		be.Equal(t, a.offset, 128)

		a.Reset()
		be.Equal(t, a.offset, 0)

		// Can allocate again from start.
		_, err = a.Alloc(128, 8)
		be.Err(t, err, nil)
		be.Equal(t, a.offset, 128)
	})
}
