package main

import "solod.dev/so/mem"

func arenaTest() {
	{
		// Arena allocator.
		buf := make([]byte, 1024)
		arena := mem.NewArena(buf)
		var a mem.Allocator = &arena

		// Allocate a Point.
		p, err := mem.TryAlloc[Point](a)
		if err != nil {
			panic("initial allocation failed")
		}
		p.x = 11
		p.y = 22
		if p.x != 11 || p.y != 22 {
			panic("unexpected p.x or p.y")
		}

		// Free is a no-op.
		mem.Free(a, p)

		// Reset and reallocate.
		arena.Reset()
		p2, err := mem.TryAlloc[Point](a)
		if err != nil {
			panic("allocation after reset failed")
		}
		// Memory should be zeroed.
		if p2.x != 0 || p2.y != 0 {
			panic("memory not zeroed after reset")
		}
		p2.x = 33
		p2.y = 44
	}
}
