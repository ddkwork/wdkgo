package maps

import (
	"testing"

	"github.com/nalgeon/be"
)

func TestMap(t *testing.T) {
	t.Run("set and get", func(t *testing.T) {
		m := New[int, int](nil, 8)
		defer m.Free()
		m.Set(1, 10)
		m.Set(2, 20)
		m.Set(3, 30)
		be.Equal(t, m.Get(1), 10)
		be.Equal(t, m.Get(2), 20)
		be.Equal(t, m.Get(3), 30)
	})

	t.Run("has", func(t *testing.T) {
		m := New[int, int](nil, 8)
		defer m.Free()
		m.Set(1, 10)
		be.Equal(t, m.Has(1), true)
		be.Equal(t, m.Has(2), false)
	})

	t.Run("get missing", func(t *testing.T) {
		m := New[int, int](nil, 8)
		defer m.Free()
		be.Equal(t, m.Get(42), 0)
	})

	t.Run("overwrite", func(t *testing.T) {
		m := New[int, int](nil, 8)
		defer m.Free()
		m.Set(1, 10)
		m.Set(1, 20)
		be.Equal(t, m.Get(1), 20)
		be.Equal(t, m.Len(), 1)
	})

	t.Run("delete", func(t *testing.T) {
		m := New[int, int](nil, 8)
		defer m.Free()
		m.Set(1, 10)
		m.Set(2, 20)
		m.Delete(1)
		be.Equal(t, m.Has(1), false)
		be.Equal(t, m.Get(1), 0)
		be.Equal(t, m.Len(), 1)
		be.Equal(t, m.Get(2), 20)
	})

	t.Run("delete missing", func(t *testing.T) {
		m := New[int, int](nil, 8)
		defer m.Free()
		m.Set(1, 10)
		m.Delete(99) // no-op, no panic
		be.Equal(t, m.Len(), 1)
	})

	t.Run("len", func(t *testing.T) {
		m := New[int, int](nil, 8)
		defer m.Free()
		be.Equal(t, m.Len(), 0)
		m.Set(1, 10)
		be.Equal(t, m.Len(), 1)
		m.Set(2, 20)
		be.Equal(t, m.Len(), 2)
		m.Delete(1)
		be.Equal(t, m.Len(), 1)
	})

	t.Run("clear", func(t *testing.T) {
		m := New[int, int](nil, 8)
		defer m.Free()
		m.Set(1, 10)
		m.Set(2, 20)
		m.Clear()
		be.Equal(t, m.Len(), 0)
		be.Equal(t, m.Has(1), false)
		be.Equal(t, m.Has(2), false)
		// reusable after Clear
		m.Set(3, 30)
		be.Equal(t, m.Get(3), 30)
		be.Equal(t, m.Len(), 1)
	})

	t.Run("free", func(t *testing.T) {
		m := New[int, int](nil, 8)
		m.Free()
		m.Free() // double Free doesn't panic
	})

	t.Run("growth", func(t *testing.T) {
		m := New[int, int](nil, 8)
		defer m.Free()
		n := 1000
		for i := range n {
			m.Set(i, i*10)
		}
		be.Equal(t, m.Len(), n)
		for i := range n {
			be.Equal(t, m.Get(i), i*10)
		}
	})

	t.Run("string keys", func(t *testing.T) {
		m := New[string, int](nil, 8)
		defer m.Free()
		m.Set("hello", 1)
		m.Set("world", 2)
		be.Equal(t, m.Get("hello"), 1)
		be.Equal(t, m.Get("world"), 2)
		be.Equal(t, m.Has("missing"), false)
		m.Delete("hello")
		be.Equal(t, m.Has("hello"), false)
		be.Equal(t, m.Len(), 1)
	})

	t.Run("delete all", func(t *testing.T) {
		m := New[int, int](nil, 8)
		defer m.Free()
		n := 50
		for i := range n {
			m.Set(i, i)
		}
		for i := range n {
			m.Delete(i)
		}
		be.Equal(t, m.Len(), 0)
		// re-insert works
		m.Set(1, 100)
		be.Equal(t, m.Get(1), 100)
		be.Equal(t, m.Len(), 1)
	})
}
