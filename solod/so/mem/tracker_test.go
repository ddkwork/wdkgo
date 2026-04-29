package mem

import (
	"testing"

	"github.com/nalgeon/be"
)

func TestTracker(t *testing.T) {
	t.Run("Alloc", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			tr := Tracker{Allocator: System}
			p, err := tr.Alloc(16, 8)
			be.Err(t, err, nil)
			if p == nil {
				t.Fatal("want non-nil pointer")
			}
			be.Equal(t, tr.Stats.Alloc, uint64(16))
			be.Equal(t, tr.Stats.TotalAlloc, uint64(16))
			be.Equal(t, tr.Stats.Mallocs, uint64(1))
			be.Equal(t, tr.Stats.Frees, uint64(0))
			tr.Free(p, 16, 8)
		})
		t.Run("multiple", func(t *testing.T) {
			tr := Tracker{Allocator: System}
			p1, _ := tr.Alloc(16, 8)
			p2, _ := tr.Alloc(32, 8)
			be.Equal(t, tr.Stats.Alloc, uint64(48))
			be.Equal(t, tr.Stats.TotalAlloc, uint64(48))
			be.Equal(t, tr.Stats.Mallocs, uint64(2))
			tr.Free(p1, 16, 8)
			tr.Free(p2, 32, 8)
		})
		t.Run("error", func(t *testing.T) {
			buf := make([]byte, 16)
			a := NewArena(buf)
			tr := Tracker{Allocator: &a}
			_, err := tr.Alloc(32, 8)
			be.Err(t, err, ErrOutOfMemory)
			be.Equal(t, tr.Stats.Alloc, uint64(0))
			be.Equal(t, tr.Stats.Mallocs, uint64(0))
		})
	})

	t.Run("Realloc", func(t *testing.T) {
		t.Run("grow", func(t *testing.T) {
			tr := Tracker{Allocator: System}
			p, _ := tr.Alloc(16, 8)

			p2, err := tr.Realloc(p, 16, 32, 8)
			be.Err(t, err, nil)
			if p2 == nil {
				t.Fatal("want non-nil pointer")
			}
			be.Equal(t, tr.Stats.Alloc, uint64(32))
			be.Equal(t, tr.Stats.TotalAlloc, uint64(32))
			be.Equal(t, tr.Stats.Mallocs, uint64(2))
			be.Equal(t, tr.Stats.Frees, uint64(1))
			tr.Free(p2, 32, 8)
		})
		t.Run("shrink", func(t *testing.T) {
			tr := Tracker{Allocator: System}
			p, _ := tr.Alloc(32, 8)

			p2, err := tr.Realloc(p, 32, 16, 8)
			be.Err(t, err, nil)
			be.Equal(t, tr.Stats.Alloc, uint64(16))
			be.Equal(t, tr.Stats.TotalAlloc, uint64(32)) // no increase
			be.Equal(t, tr.Stats.Mallocs, uint64(2))
			be.Equal(t, tr.Stats.Frees, uint64(1))
			tr.Free(p2, 16, 8)
		})
		t.Run("same size", func(t *testing.T) {
			tr := Tracker{Allocator: System}
			p, _ := tr.Alloc(16, 8)

			p2, err := tr.Realloc(p, 16, 16, 8)
			be.Err(t, err, nil)
			be.Equal(t, tr.Stats.Alloc, uint64(16))
			be.Equal(t, tr.Stats.TotalAlloc, uint64(16))
			tr.Free(p2, 16, 8)
		})
		t.Run("error", func(t *testing.T) {
			buf := make([]byte, 32)
			a := NewArena(buf)
			tr := Tracker{Allocator: &a}
			p, _ := tr.Alloc(16, 8)

			_, err := tr.Realloc(p, 16, 64, 8)
			be.Err(t, err, ErrOutOfMemory)
			// Stats unchanged after failed realloc.
			be.Equal(t, tr.Stats.Alloc, uint64(16))
			be.Equal(t, tr.Stats.TotalAlloc, uint64(16))
			be.Equal(t, tr.Stats.Mallocs, uint64(1))
			be.Equal(t, tr.Stats.Frees, uint64(0))
		})
	})

	t.Run("Free", func(t *testing.T) {
		tr := Tracker{Allocator: System}
		p, _ := tr.Alloc(16, 8)
		tr.Free(p, 16, 8)
		be.Equal(t, tr.Stats.Alloc, uint64(0))
		be.Equal(t, tr.Stats.TotalAlloc, uint64(16))
		be.Equal(t, tr.Stats.Mallocs, uint64(1))
		be.Equal(t, tr.Stats.Frees, uint64(1))
	})

	t.Run("lifecycle", func(t *testing.T) {
		tr := Tracker{Allocator: System}

		// Alloc 16 bytes.
		p1, _ := tr.Alloc(16, 8)
		// Alloc 32 bytes.
		p2, _ := tr.Alloc(32, 8)
		be.Equal(t, tr.Stats.Alloc, uint64(48))
		be.Equal(t, tr.Stats.TotalAlloc, uint64(48))
		be.Equal(t, tr.Stats.Mallocs, uint64(2))

		// Grow p1: 16 -> 64.
		p1, _ = tr.Realloc(p1, 16, 64, 8)
		be.Equal(t, tr.Stats.Alloc, uint64(96))
		be.Equal(t, tr.Stats.TotalAlloc, uint64(96))
		be.Equal(t, tr.Stats.Mallocs, uint64(3))
		be.Equal(t, tr.Stats.Frees, uint64(1))

		// Shrink p2: 32 -> 8.
		p2, _ = tr.Realloc(p2, 32, 8, 8)
		be.Equal(t, tr.Stats.Alloc, uint64(72))
		be.Equal(t, tr.Stats.TotalAlloc, uint64(96)) // unchanged
		be.Equal(t, tr.Stats.Mallocs, uint64(4))
		be.Equal(t, tr.Stats.Frees, uint64(2))

		// Free both.
		tr.Free(p1, 64, 8)
		tr.Free(p2, 8, 8)
		be.Equal(t, tr.Stats.Alloc, uint64(0))
		be.Equal(t, tr.Stats.TotalAlloc, uint64(96))
		be.Equal(t, tr.Stats.Mallocs, uint64(4))
		be.Equal(t, tr.Stats.Frees, uint64(4))
	})
}
