package pragma

// DoNotImplement can be embedded in an interface to prevent
// trivial implementations of the interface.
//
// This is useful to prevent unauthorized implementations of an
// interface so that it can be extended in the future for any changes or
// ensure certain guarantees for type conversion.
//
//nolint:inamedparam // not needed here
type DoNotImplement interface {
	TestoInternal(DoNotImplement)
}
