package mem

import (
	"testing"

	"github.com/nalgeon/be"
)

func TestNoAlloc(t *testing.T) {
	t.Run("Alloc", func(t *testing.T) {
		_, err := NoAlloc.Alloc(16, 8)
		be.Err(t, err, ErrNoAlloc)
	})
	t.Run("Realloc", func(t *testing.T) {
		_, err := NoAlloc.Realloc(nil, 0, 16, 8)
		be.Err(t, err, ErrNoAlloc)
	})
	t.Run("Free", func(t *testing.T) {
		NoAlloc.Free(nil, 0, 0) // should not panic
	})
}
