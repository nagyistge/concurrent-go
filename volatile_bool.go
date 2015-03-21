package concurrent

import "sync/atomic"

const (
	volatileBoolTrue = iota
	volatileBoolFalse
)

// TODO(pedge): is this even needed? need to understand go memory model better
type VolatileBool interface {
	Value() bool
	// return old value == new value
	CompareAndSwap(oldBool bool, newBool bool) bool
}

func NewVolatileBool(initialBool bool) VolatileBool {
	return newVolatileBool(initialBool)
}

type volatileBool struct {
	value_ int32
}

func newVolatileBool(initialBool bool) *volatileBool {
	return &volatileBool{boolToVolatileBoolValue(initialBool)}
}

func (this *volatileBool) Value() bool {
	return volatileBoolValueToBool(atomic.LoadInt32(&this.value_))
}

// return old value == new value
func (this *volatileBool) CompareAndSwap(oldBool bool, newBool bool) bool {
	return atomic.CompareAndSwapInt32(
		&this.value_,
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
