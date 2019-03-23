package interleaving_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rerost/information-retrieval-tools-go/interleaving"
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

func TestInterleaving(t *testing.T) {
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

	inOutPairs := []InOutPair{
		{
			Name: "Normal case",
			In: In{
				PrioritizeRanking: NewDummyPrioritizeRanking(true),
				RankingA:          []int{1, 3, 5},
				RankingB:          []int{2, 4, 6},
			},
			Out: Out{
				Error:   nil,
				Ranking: []interface{}{1, 2, 3, 4, 5, 6},
			},
		},
		{
			Name: "Duplicate case",
			In: In{
				PrioritizeRanking: NewDummyPrioritizeRanking(true),
				RankingA:          []int{1, 2, 3},
				RankingB:          []int{1, 2, 3},
			},
			Out: Out{
				Error:   nil,
				Ranking: []interface{}{1, 2, 3},
			},
		},
		{
			Name: "Size is different",
			In: In{
				PrioritizeRanking: NewDummyPrioritizeRanking(true),
				RankingA:          []int{1, 2, 3},
				RankingB:          []int{},
			},
			Out: Out{
				Error:   nil,
				Ranking: []interface{}{1, 2, 3},
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
				Error:   nil,
				Ranking: nil,
			},
		},
	}
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
