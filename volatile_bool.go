package concurrent

import "sync/atomic"

const (
	volatileBoolTrue = iota
	volatileBoolFalse
)

// VolatileBool creates a volatile bool.
//
// TODO(pedge): Is v even needed? Need to understand go memory model better.
type VolatileBool interface {
	Value() bool
	// return old value == new value
	CompareAndSwap(oldBool bool, newBool bool) bool
}

// NewVolatileBool creates a new VolatileBool.
func NewVolatileBool(initialBool bool) VolatileBool {
	return newVolatileBool(initialBool)
}

type volatileBool struct {
	int32Value int32
}

func newVolatileBool(initialBool bool) *volatileBool {
	return &volatileBool{boolToVolatileBoolValue(initialBool)}
}

func (v *volatileBool) Value() bool {
	return volatileBoolValueToBool(atomic.LoadInt32(&v.int32Value))
}

// return old value == new value
func (v *volatileBool) CompareAndSwap(oldBool bool, newBool bool) bool {
	return atomic.CompareAndSwapInt32(
		&v.int32Value,
		boolToVolatileBoolValue(oldBool),
		boolToVolatileBoolValue(newBool),
	)
}

func boolToVolatileBoolValue(b bool) int32 {
	if b {
		return volatileBoolTrue
	}
	return volatileBoolFalse
}

func volatileBoolValueToBool(volatileBoolValue int32) bool {
	switch int(volatileBoolValue) {
	case volatileBoolTrue:
		return true
	case volatileBoolFalse:
		return false
	default:
		panic("concurrent: unknown volatileBoolValue")
	}
}
