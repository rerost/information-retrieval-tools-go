package interleaving

import (
	"github.com/rerost/information-retrieval-tools-go/interleaving/internal"
)

type Item = internal.Item
type Ranking = []Item
type RawRanking = interface{}

type Interleaving interface {
	Perform() Ranking
}

type PrioritizeRanking interface {
	IsA() bool
}

func NewInterleaving(prioritizeRanking PrioritizeRanking, rawRankingA, rawRankingB RawRanking) (Interleaving, error) {
	rankingA, err := internal.ToRanking(rawRankingA)
	if err != nil {
		return nil, err
	}

	rankingB, err := internal.ToRanking(rawRankingB)
	if err != nil {
		return nil, err
	}

	return &interleavingImp{
		PrioritizeRanking:  prioritizeRanking,
		RankingA:           rankingA,
		RankingB:           rankingB,
		AllreadyUsedKeyMap: map[string]bool{}, // Do not have parallel safety
	}, nil
}

type interleavingImp struct {
	PrioritizeRanking  PrioritizeRanking
	RankingA           Ranking
	RankingB           Ranking
	AllreadyUsedKeyMap map[string]bool
}

func (i *interleavingImp) Perform() Ranking {
	rankingASize := int64(len(i.RankingA))
	rankingBSize := int64(len(i.RankingB))
	resultSize := rankingASize + rankingBSize

	result := make(Ranking, 0, resultSize)

	rankingAIter := int64(0)
	rankingBIter := int64(0)

	if rankingASize == 0 {
		return i.RankingB
	}
	if rankingBSize == 0 {
		return i.RankingA
	}

	for iter := int64(0); iter <= resultSize; iter += 2 {
		if i.PrioritizeRanking.IsA() {
			var ok bool
			rankingAIter, ok = i.pluck(i.RankingA, rankingAIter)
			if !ok {
				return result
			}
			result = append(
				result,
				i.RankingA[rankingAIter],
			)
			rankingBIter, ok = i.pluck(i.RankingB, rankingBIter)
			if !ok {
				return result
			}
			result = append(
				result,
				i.RankingB[rankingBIter],
			)
		} else {
			var ok bool
			rankingBIter, ok = i.pluck(i.RankingB, rankingBIter)
			if !ok {
				return result
			}
			result = append(
				result,
				i.RankingB[rankingBIter],
			)
			rankingAIter, ok = i.pluck(i.RankingA, rankingAIter)
			if !ok {
				return result
			}
			result = append(
				result,
				i.RankingA[rankingAIter],
			)
		}
		rankingAIter++
		rankingBIter++
		if rankingAIter >= rankingASize {
			break
		}
		if rankingBIter >= rankingBSize {
			break
		}
	}
	return result
}

func (i *interleavingImp) pluck(ranking Ranking, iter int64) (int64, bool) {
	for ; iter < int64(len(ranking)); iter++ {
		if _, ok := i.AllreadyUsedKeyMap[ranking[iter].Key()]; ok {
			continue
		}

		i.AllreadyUsedKeyMap[ranking[iter].Key()] = true
		return iter, true
	}

	return 0, false
}
