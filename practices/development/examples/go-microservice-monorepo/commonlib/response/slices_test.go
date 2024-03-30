package response

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNonNilSliceReturnsOriginalSliceIfFilled(t *testing.T) {
	originalSlice := []string{"1", "2", "3"}
	returnedSlice := NonNilSlice(originalSlice)

	assert.ElementsMatch(t, originalSlice, returnedSlice)
}

func TestNonNilSliceReturnsEmptySliceWhenNilIsPassed(t *testing.T) {
	var originalSlice []int
	returnedSlice := NonNilSlice(originalSlice)

	assert.NotNil(t, returnedSlice)
}
