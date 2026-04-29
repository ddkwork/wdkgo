package mem

// A Tracker wraps an Allocator and tracks all
// allocations and deallocations made through it.
type Tracker struct {
	Allocator Allocator
	Stats     Stats
}

func (t *Tracker) Alloc(size int, align int) (any, error) {
	ptr, err := t.Allocator.Alloc(size, align)
	if err != nil {
		return nil, err
	}
	t.Stats.Alloc += uint64(size)
	t.Stats.TotalAlloc += uint64(size)
	t.Stats.Mallocs++
	return ptr, nil
}

func (t *Tracker) Realloc(ptr any, oldSize int, newSize int, align int) (any, error) {
	newPtr, err := t.Allocator.Realloc(ptr, oldSize, newSize, align)
	if err != nil {
		return nil, err
	}
	if newSize > oldSize {
		t.Stats.Alloc += uint64(newSize - oldSize)
		t.Stats.TotalAlloc += uint64(newSize - oldSize)
	} else {
		t.Stats.Alloc -= uint64(oldSize - newSize)
	}
	t.Stats.Mallocs++
	t.Stats.Frees++
	return newPtr, nil
}

func (t *Tracker) Free(ptr any, size int, align int) {
	t.Allocator.Free(ptr, size, align)
	t.Stats.Alloc -= uint64(size)
	t.Stats.Frees++
}
