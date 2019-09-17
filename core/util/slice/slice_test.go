package slice

import "testing"

func TestSliceUniqueString(t *testing.T) {

	ss := []string{"a", "b", "a"}
	ss = SliceUniqueString(ss)
	t.Log(ss)

	ii := []int64{1, 2, 3, 1}
	ii = SliceUniqueInt64(ii)
	t.Log(ii)

}
