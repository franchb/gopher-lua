package lua

import (
	"unsafe"
)

// iface is an internal representation of the go-interface.
type iface struct {
	itab unsafe.Pointer
	word unsafe.Pointer
}

// Preload cache for common integer values to avoid allocations
const preloadLimit LNumber = 256
const preloadNegativeLimit LNumber = 128

var preloads [int(preloadLimit)]LValue
var preloadsNegative [int(preloadNegativeLimit)]LValue

func init() {
	for i := range int(preloadLimit) {
		preloads[i] = LNumber(i)
	}
	for i := range int(preloadNegativeLimit) {
		preloadsNegative[i] = LNumber(-i - 1)
	}
}

// allocator is a fast bulk memory allocator for the LValue.
type allocator struct {
	fptrs []float64

	scratchValue  LValue
	scratchValueP *iface
}

func newAllocator(size int) *allocator {
	al := &allocator{
		fptrs: make([]float64, 0, size),
	}
	al.scratchValue = LNumber(0)
	al.scratchValueP = (*iface)(unsafe.Pointer(&al.scratchValue))

	return al
}

// LNumber2I takes a number value and returns an interface LValue representing the same number.
// Converting an LNumber to a LValue naively, by doing:
// `var val LValue = myLNumber`
// will result in an individual heap alloc of 8 bytes for the float value. LNumber2I amortizes the cost and memory
// overhead of these allocs by allocating blocks of floats instead.
// The downside of this is that all of the floats on a given block have to become eligible for gc before the block
// as a whole can be gc-ed.
func (al *allocator) LNumber2I(v LNumber) LValue {
	// first check for shared preloaded numbers (positive integers [0, 255])
	if v >= 0 && v < preloadLimit && float64(v) == float64(int64(v)) {
		return preloads[int(v)]
	}
	// check for shared preloaded negative numbers ([-128, -1])
	if v < 0 && v >= -preloadNegativeLimit && float64(v) == float64(int64(v)) {
		return preloadsNegative[int(-v)-1]
	}

	// check if we need a new alloc page
	if cap(al.fptrs) == len(al.fptrs) {
		al.fptrs = make([]float64, 0, cap(al.fptrs))
	}

	// alloc a new float, and store our value into it
	al.fptrs = append(al.fptrs, float64(v))
	fptr := &al.fptrs[len(al.fptrs)-1]

	// hack our scratch LValue to point to our allocated value
	// this scratch lvalue is copied when this function returns meaning the scratch value can be reused
	// on the next call
	al.scratchValueP.word = unsafe.Pointer(fptr)

	return al.scratchValue
}
