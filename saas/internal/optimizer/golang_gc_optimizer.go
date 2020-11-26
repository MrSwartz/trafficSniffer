package optimizer

import "runtime"

/*
	This function requests a block of memory from the OS,
	which is immediately freed and becomes available to the
	scheduler in order to reduce the number of small allocations

	This function should be used rationally and
	only after testing for memory consumption
*/

func Allocator(size int) {
	a := make([]byte, 1024 * 1024 * 1024 * size)
	for i := range a {
		a[i] = 'c'
	}
	runtime.GC()
}
