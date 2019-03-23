package interleaving_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rerost/information-retrieval-tools-go/interleaving"
	"github.com/rerost/information-retrieval-tools-go/interleaving/errors"
)

var (
	ErrComparer = cmp.Comparer(func(x, y error) bool {
		if x == nil || y == nil {
			return x == nil && y == nil
		}
		return x.Error() == y.Error()
	})
)

type DummyPrioritizeRanking struct {
	isA bool
}

func NewDummyPrioritizeRanking(isA bool) interleaving.PrioritizeRanking {
	return DummyPrioritizeRanking{isA}
}
func (d DummyPrioritizeRanking) IsA() bool {
	return d.isA
}

type DummyItem int

func (e DummyItem) Key() string {
	return fmt.Sprintf("%v", e)
}

func ToRanking(arr []int) interleaving.Ranking {
	result := make(interleaving.Ranking, len(arr), len(arr))

	for i, elem := range arr {
		result[i] = DummyItem(elem)
	}

	return result
}

type In struct {
	PrioritizeRanking interleaving.PrioritizeRanking
	RankingA          interface{}
	RankingB          interface{}
}
type Out struct {
	Error   error
	Ranking interleaving.Ranking
}
type InOutPair struct {
	Name string
	In   In
	Out  Out
}

var inOutPairs = []InOutPair{
	{
		Name: "Normal case",
		In: In{
			PrioritizeRanking: NewDummyPrioritizeRanking(true),
			RankingA:          []DummyItem{1, 3, 5},
			RankingB:          []DummyItem{2, 4, 6},
		},
		Out: Out{
			Error:   nil,
			Ranking: ToRanking([]int{1, 2, 3, 4, 5, 6}),
		},
	},
	{
		Name: "Use b case",
		In: In{
			PrioritizeRanking: NewDummyPrioritizeRanking(false),
			RankingA:          []DummyItem{2, 4, 6},
			RankingB:          []DummyItem{1, 3, 5},
		},
		Out: Out{
			Error:   nil,
			Ranking: ToRanking([]int{1, 2, 3, 4, 5, 6}),
		},
	},
	{
		Name: "Duplicate case",
		In: In{
			PrioritizeRanking: NewDummyPrioritizeRanking(true),
			RankingA:          []DummyItem{1, 2, 3},
			RankingB:          []DummyItem{1, 2, 3},
		},
		Out: Out{
			Error:   nil,
			Ranking: ToRanking([]int{1, 2, 3}),
		},
	},
	{
		Name: "Size is different",
		In: In{
			PrioritizeRanking: NewDummyPrioritizeRanking(true),
			RankingA:          ToRanking([]int{1, 2, 3}),
			RankingB:          ToRanking([]int{}),
		},
		Out: Out{
			Error:   nil,
			Ranking: ToRanking([]int{1, 2, 3}),
		},
	},
	{
		Name: "Not slice passed",
		In: In{
			PrioritizeRanking: NewDummyPrioritizeRanking(true),
			RankingA:          nil,
			RankingB:          nil,
		},
		Out: Out{
			Error:   errors.NotSliceError,
			Ranking: nil,
		},
	},
}

func TestInterleaving(t *testing.T) {
	for _, inOutPair := range inOutPairs {
		inOutPair := inOutPair
		t.Run(inOutPair.Name, func(t *testing.T) {
			i, err := interleaving.NewInterleaving(inOutPair.In.PrioritizeRanking, inOutPair.In.RankingA, inOutPair.In.RankingB)
			if diff := cmp.Diff(err, inOutPair.Out.Error, ErrComparer); diff != "" {
				t.Error(diff)
			}
			if err != nil {
				return
			}
			result := i.Perform()
			if diff := cmp.Diff(result, inOutPair.Out.Ranking); diff != "" {
				t.Error(diff)
			}
		})
	}
}
