package internal_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rerost/information-retrieval-tools-go/interleaving/errors"
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

type DummyItem int

func (e DummyItem) Key() string {
	return fmt.Sprintf("%v", e)
}

type DummyStringItem string

func (e DummyStringItem) Key() string {
	return fmt.Sprintf("%v", e)
}

func TestToRanking(t *testing.T) {
	inOutPairs := []struct {
		Name string
		In   interface{}
		Out  []internal.Item
		Err  error
	}{
		{
			Name: "Normal case",
			In:   []DummyItem{1, 2, 3},
			Out:  []internal.Item{DummyItem(1), DummyItem(2), DummyItem(3)},
			Err:  nil,
		},
		{
			Name: "Passed not slice",
			In:   1,
			Out:  nil,
			Err:  errors.NotSliceError,
		},
		{
			Name: "Passed interface slice",
			In:   []internal.Item{DummyItem(1), DummyStringItem("test")},
			Out:  []internal.Item{DummyItem(1), DummyStringItem("test")},
			Err:  nil,
		},
		{
			Name: "Nil case",
			In:   nil,
			Out:  nil,
			Err:  errors.NotSliceError,
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
