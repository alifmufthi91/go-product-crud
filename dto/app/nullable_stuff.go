package app

type NullableStuff[T any] struct {
	Stuff T
	Valid bool
}

func NewNullableStuff[T any](stuff T) NullableStuff[T] {
	return NullableStuff[T]{
		Stuff: stuff,
		Valid: true,
	}
}
