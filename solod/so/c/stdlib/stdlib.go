// Package stdlib wraps the C <stdlib.h> header.
// It offers process control, memory management,
// string-to-number conversion, and environment access.
package stdlib

import _ "embed"

//so:embed stdlib.h
var stdlib_h string

// ExitSuccess indicates successful program termination.
//
//so:extern
const ExitSuccess int = 0

// ExitFailure indicates unsuccessful program termination.
//
//so:extern
const ExitFailure int = 1

// Exit terminates the program with the given status code.
// All C streams are flushed and closed; temporary files are removed.
//
//so:extern
func Exit(status int) { _ = status }

// Malloc allocates size bytes of uninitialized memory.
// Returns a pointer to the allocated memory, or nil on failure.
//
// If size is zero, the behavior is implementation-defined,
// either returning a null or non-null pointer.
//
// If Malloc returns a non-null pointer, the caller is responsible for
// eventually calling [Free] to deallocate the memory.
//
//so:extern
func Malloc(size uintptr) any { _ = size; return nil }

// Calloc allocates memory for count elements of size bytes each,
// initialized to zero. Returns a pointer to the allocated memory,
// or nil on failure.
//
// If count or size is zero, the behavior is implementation-defined,
// either returning a null or non-null pointer.
//
// If Calloc returns a non-null pointer, the caller is responsible for
// eventually calling [Free] to deallocate the memory.
//
//so:extern
func Calloc(count uintptr, size uintptr) any { _, _ = count, size; return nil }

// Realloc changes the size of the memory block pointed to by ptr
// to size bytes. Returns a pointer to the reallocated memory,
// or nil on failure.

// If ptr is not nil, it must have been returned by an earlier call to
// [Malloc], [Calloc], or [Realloc], and not yet deallocated by a call to [Free].
//
// If ptr is nil, Realloc behaves like Malloc(size).
// If size is zero, the behavior is implementation-defined (<C23) or undefined (C23+).
//
//so:extern
func Realloc(ptr any, size uintptr) any { _, _ = ptr, size; return nil }

// Free deallocates the memory previously allocated by
// [Malloc], [Calloc], or [Realloc].
//
// If ptr is a null pointer, does nothing. If ptr has already been
// deallocated, or was not allocated by a call to an allocation function,
// the behavior is undefined.
//
//so:extern
func Free(ptr any) { _ = ptr }

// Atoi converts the string str to an int value. The implied radix is always 10.
// Returns the converted value, or zero if conversion failed.
//
// Discards any whitespace characters until the first non-whitespace character
// is found, then takes as many characters as possible to form a valid integer
// number representation ([+-]?[0-9]+) and converts them to an integer value.
//
//so:extern
func Atoi(str string) int { _ = str; return 0 }

// Atof converts the string str to a float64 value.
// Returns the converted value, or zero if conversion failed.
//
// Discards any whitespace characters until the first non-whitespace character
// is found, then takes as many characters as possible to form a valid floating-point
// number representation and converts them to a float64 value.
//
// Valid representations include:
//   - [+-]?[0-9]+(\.[0-9]*)?([eE][+-]?[0-9]+)?
//   - [+-]?0[xX][0-9a-fA-F]+(\.[0-9a-fA-F]*)?[pP][+-]?[0-9]+
//   - [+-]?INF(INITY)?
//   - [+-]?NAN
//
//so:extern
func Atof(str string) float64 { _ = str; return 0 }

// Getenv returns a pointer to the value of the environment variable
// named by name, or nil if the variable is not set.
//
//so:extern
func Getenv(name string) *byte { _ = name; return nil }
