package interleaving_test

import (
	"fmt"
	"testing"

	"github.com/rerost/information-retrieval-tools-go/interleaving"
)

// func BenchmarkInterleaving(b *testing.B) {
// 	for _, inOutPair := range inOutPairs {
// 		inOutPair := inOutPair
// 		if inOutPair.Out.Error != nil {
// 			continue
// 		}
// 		b.Run(inOutPair.Name, func(b *testing.B) {
// 			b.ResetTimer()
// 			for i := 0; i < b.N; i++ {
// 				i, err := interleaving.NewInterleaving(inOutPair.In.PrioritizeRanking, inOutPair.In.RankingA, inOutPair.In.RankingB)
// 				if err != nil {
// 					b.Error(err)
// 					return
// 				}
// 				_ = i.Perform()
// 			}
// 		})
// 	}
// }
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

func CreateDummyArr(size int, weight int) []int {
	result := make([]int, size, size)
	for i := range result {
		result[i] = i * weight
	}
	return result
}

var largeInOutPairs = []struct {
	Name     string
	RankingA interleaving.Ranking
	RankingB interleaving.Ranking
}{
	{
		Name:     "1000",
		RankingA: ToRanking(CreateDummyArr(1000, 1)),
		RankingB: ToRanking(CreateDummyArr(1000, -1)),
	},
	{
		Name:     "10000",
		RankingA: ToRanking(CreateDummyArr(10000, 1)),
		RankingB: ToRanking(CreateDummyArr(10000, -1)),
	},
	{
		Name:     "100000",
		RankingA: ToRanking(CreateDummyArr(100000, 1)),
		RankingB: ToRanking(CreateDummyArr(100000, -1)),
	},
}

func BenchmarkNewInterleavingWithLarge(b *testing.B) {
	inOutPairs := largeInOutPairs

	for _, inOutPair := range inOutPairs {
		inOutPair := inOutPair
		b.Run(inOutPair.Name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := interleaving.NewInterleaving(DummyPrioritizeRanking{true}, inOutPair.RankingA, inOutPair.RankingB)
				if err != nil {
					b.Error(err)
					return
				}
			}
		})
	}
}

func BenchmarkPerform(b *testing.B) {
	inOutPairs := largeInOutPairs

	for _, inOutPair := range inOutPairs {
		inOutPair := inOutPair
		b.Run(inOutPair.Name, func(b *testing.B) {
			interleavingEngine, err := interleaving.NewInterleaving(DummyPrioritizeRanking{true}, inOutPair.RankingA, inOutPair.RankingB)
			if err != nil {
				b.Error(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				interleavingEngine.Perform()
			}
		})
	}
}
