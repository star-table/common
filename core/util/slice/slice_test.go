package slice

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestSliceUniqueString(t *testing.T) {

	ss := []string{"a", "b", "a", "c", "aa"}
	ss = SliceUniqueString(ss)
	t.Log(ss)

	ii := []int64{1, 2, 4, 2, 3, 1}
	ii = SliceUniqueInt64(ii)
	assert.Equal(t, ii, []int64{1, 2, 4, 3})
	t.Log(ii)

}
