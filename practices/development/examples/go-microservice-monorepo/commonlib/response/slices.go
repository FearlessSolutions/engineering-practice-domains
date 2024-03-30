package response

// NonNilSlice takes a slice value and returns a zero-length slice if the incoming
// slice is nil, returning the passed slice otherwise
func NonNilSlice[T any](slice []T) []T {
	if slice == nil {
		slice = make([]T, 0)
	}

	return slice
}
