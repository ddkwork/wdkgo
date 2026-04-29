// Package mem provides memory allocation facilities.
package mem

import (
	"unsafe"

	"solod.dev/so/c"
	"solod.dev/so/errors"
	"solod.dev/so/math"
)

// ErrOutOfMemory is returned when a memory allocation
// fails due to insufficient memory.
var ErrOutOfMemory = errors.New("out of memory")

//so:extern so_Nil
var Nil byte

// Alloc allocates a single value of type T using allocator a.
// Returns a pointer to the allocated memory or panics on failure.
// Whether new memory is zeroed depends on the allocator.
// If the allocator is nil, uses the system allocator.
//
//so:inline
func Alloc[T any](a Allocator) *T {
	_ptr, _err := TryAlloc[T](a)
	if _err != nil {
		panic(_err)
	}
	return _ptr
}

// TryAlloc is like [Alloc] but returns an error
// instead of panicking on failure.
//
//so:inline
func TryAlloc[T any](a Allocator) (*T, error) {
	_a := a
	if _a == nil {
		_a = System
	}
	_ptr, _err := _a.Alloc(c.Sizeof[T](), c.Alignof[T]())
	return c.PtrAs[T](_ptr), _err
}

// Free frees a value previously allocated with [Alloc] or [TryAlloc].
// If the allocator is nil, uses the system allocator.
//
//so:inline
func Free[T any](a Allocator, ptr *T) {
	_a := a
	if _a == nil {
		_a = System
	}
	_a.Free(ptr, c.Sizeof[T](), c.Alignof[T]())
}

// AllocSlice allocates a slice of type T with given length
// and capacity using allocator a.
// Returns a slice of the allocated memory or panics on failure.
// Whether new memory is zeroed depends on the allocator.
// If the allocator is nil, uses the system allocator.
//
//so:inline
func AllocSlice[T any](a Allocator, len int, cap int) []T {
	_s, _err := TryAllocSlice[T](a, len, cap)
	if _err != nil {
		panic(_err)
	}
	return _s
}

// TryAllocSlice is like [AllocSlice] but returns an error
// instead of panicking on allocation failure.
//
//so:inline
func TryAllocSlice[T any](a Allocator, len int, cap int) ([]T, error) {
	_a := a
	if _a == nil {
		_a = System
	}

	_len, _cap := len, cap
	_esize, _align := c.Sizeof[T](), c.Alignof[T]()

	c.Assert(_len >= 0, "mem: negative length")
	c.Assert(_cap >= 0, "mem: negative capacity")
	c.Assert(_len <= _cap, "mem: length exceeds capacity")
	c.Assert(_cap < math.MaxInt/_esize, "mem: capacity overflow")

	var _ptr any
	var _err error
	if _cap > 0 {
		_ptr, _err = _a.Alloc(_esize*_cap, _align)
	} else {
		_ptr = &Nil
	}

	var _ts []T
	if _err == nil {
		_ts = c.Slice(c.PtrAs[T](_ptr), _len, _cap)
	}
	return _ts, _err
}

// ReallocSlice reallocates a slice of type T with new length and capacity
// using allocator a. Preserves contents up to the old capacity.
// Returns the reallocated slice or panics on failure.
// Whether new memory is zeroed depends on the allocator.
// If the allocator is nil, uses the system allocator.
//
//so:inline
func ReallocSlice[T any](a Allocator, slice []T, newLen int, newCap int) []T {
	_s, _err := TryReallocSlice(a, slice, newLen, newCap)
	if _err != nil {
		panic(_err)
	}
	return _s
}

// TryReallocSlice is like [ReallocSlice] but returns an error
// instead of panicking on allocation failure.
//
//so:inline
func TryReallocSlice[T any](a Allocator, slice []T, newLen int, newCap int) ([]T, error) {
	_a := a
	if _a == nil {
		_a = System
	}

	_oldCap := cap(slice)
	_newLen, _newCap := newLen, newCap
	_esize, _align := c.Sizeof[T](), c.Alignof[T]()

	c.Assert(_newLen >= 0, "mem: negative length")
	c.Assert(_newCap >= 0, "mem: negative capacity")
	c.Assert(_newLen <= _newCap, "mem: length exceeds capacity")
	c.Assert(_newCap < math.MaxInt/_esize, "mem: capacity overflow")

	var _newPtr any
	var _err error
	if _newCap == 0 {
		if _oldCap > 0 {
			_a.Free(unsafe.SliceData(slice), _esize*_oldCap, _align)
		}
		_newPtr = &Nil
	} else if _oldCap == 0 {
		_newPtr, _err = _a.Alloc(_esize*_newCap, _align)
	} else {
		ptr := unsafe.SliceData(slice)
		_newPtr, _err = _a.Realloc(ptr, _esize*_oldCap, _esize*_newCap, _align)
	}

	var _s []T
	if _err == nil {
		_s = c.Slice(c.PtrAs[T](_newPtr), _newLen, _newCap)
	}
	return _s, _err
}

// FreeSlice frees a slice previously allocated with [AllocSlice] or [TryAllocSlice].
// If the allocator is nil, uses the system allocator.
// Calling FreeSlice on an empty or nil slice is a no-op.
//
//so:inline
func FreeSlice[T any](a Allocator, slice []T) {
	_a := a
	if _a == nil {
		_a = System
	}

	_s := slice
	_cap := cap(_s)
	if _cap > 0 {
		_a.Free(unsafe.SliceData(_s), c.Sizeof[T]()*_cap, c.Alignof[T]())
	}
}

// FreeString frees a heap-allocated string.
// If the allocator is nil, uses the system allocator.
func FreeString(a Allocator, s string) {
	if len(s) == 0 {
		return
	}
	Free(a, unsafe.StringData(s))
}

// Clear zeroes size bytes starting at ptr.
//
//so:inline
func Clear(ptr any, size int) {
	c.Assert(ptr != nil, "mem: nil pointer")
	c.Assert(size >= 0, "mem: negative size")
	memset(ptr, 0, uintptr(size))
}

// Copy copies n bytes from src to dst. Returns dst.
// The memory areas must not overlap.
// Panics if either dst or src is nil.
//
//so:inline
func Copy(dst any, src any, n int) any {
	c.Assert(dst != nil, "mem: nil pointer")
	c.Assert(src != nil, "mem: nil pointer")
	c.Assert(n >= 0, "mem: negative size")
	return memcpy(dst, src, uintptr(n))
}

// Move copies n bytes from src to dst. Returns dst.
// The memory areas may overlap.
// Panics if either dst or src is nil.
//
//so:inline
func Move(dst any, src any, n int) any {
	c.Assert(dst != nil, "mem: nil pointer")
	c.Assert(src != nil, "mem: nil pointer")
	c.Assert(n >= 0, "mem: negative size")
	return memmove(dst, src, uintptr(n))
}

// void* memset(void *dest, int ch, size_t count);
//
//so:extern
func memset(ptr any, ch int, count uintptr) any {
	s := unsafe.Slice((*byte)(ptrVal(ptr)), int(count))
	for i := range s {
		s[i] = byte(ch)
	}
	return ptr
}

// void* memcpy(void* dest, const void* src, size_t count);
//
//so:extern
func memcpy(dst any, src any, count uintptr) any {
	return memmove(dst, src, count)
}

// void* memmove(void* dest, const void* src, size_t count);
//
//so:extern
func memmove(dst any, src any, count uintptr) any {
	dstSlice := unsafe.Slice((*byte)(ptrVal(dst)), int(count))
	srcSlice := unsafe.Slice((*byte)(ptrVal(src)), int(count))
	copy(dstSlice, srcSlice)
	return dst
}
