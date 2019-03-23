package internal_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rerost/information-retrieval-tools-go/interleaving/internal"
)

var (
	ErrComparer = cmp.Comparer(func(x, y error) bool {
		if x == nil || y == nil {
			return x == nil && y == nil
		}
		return x.Error() == y.Error()
	})
)

func TestToRanking(t *testing.T) {
	inOutPairs := []struct {
		Name string
		In   interface{}
		Out  []interface{}
		Err  error
	}{
		{
			Name: "Normal case",
			In:   []int{1, 2, 3},
			Out:  []interface{}{1, 2, 3},
			Err:  nil,
		},
		{
			Name: "Passed not slice",
			In:   1,
			Out:  nil,
			Err:  internal.NotSliceError,
		},
		{
			Name: "Passed interface slice",
			In:   []interface{}{1, "test"},
			Out:  []interface{}{1, "test"},
			Err:  nil,
		},
		{
			Name: "Nil case",
			In:   nil,
			Out:  nil,
			Err:  internal.NotSliceError,
		},
	}
	for _, inOutPair := range inOutPairs {
		inOutPair := inOutPair
		t.Run(inOutPair.Name, func(t *testing.T) {
			result, err := internal.ToRanking(inOutPair.In)
			if diff := cmp.Diff(err, inOutPair.Err, ErrComparer); diff != "" {
				t.Error(diff)
			}
			if inOutPair.Err != nil {
				if diff := cmp.Diff(result, inOutPair.Out); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}
