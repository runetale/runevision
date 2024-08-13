package atomics

import "sync/atomic"

// todo (snt) FIX
// AtomicValue is the generic version of [atomic.Value].
type AtomicValue[T any] struct {
	v atomic.Value
}

type wrappedValue[T any] struct{ v T }

func (v *AtomicValue[T]) Load() T {
	x, _ := v.LoadOk()
	return x
}

func (v *AtomicValue[T]) LoadOk() (_ T, ok bool) {
	x := v.v.Load()
	if x != nil {
		return x.(wrappedValue[T]).v, true
	}
	var zero T
	return zero, false
}

func (v *AtomicValue[T]) Store(x T) {
	v.v.Store(wrappedValue[T]{x})
}

func (v *AtomicValue[T]) Swap(x T) (old T) {
	oldV := v.v.Swap(wrappedValue[T]{x})
	if oldV != nil {
		return oldV.(wrappedValue[T]).v
	}
	return old
}

func (v *AtomicValue[T]) CompareAndSwap(oldV, newV T) (swapped bool) {
	return v.v.CompareAndSwap(wrappedValue[T]{oldV}, wrappedValue[T]{newV})
}
