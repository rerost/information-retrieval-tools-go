package interleaving

import (
	"github.com/rerost/information-retrieval-tools-go/interleaving/internal"
)

type Ranking = []interface{}
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
		PrioritizeRanking: prioritizeRanking,
		RankingA:          rankingA,
		RankingB:          rankingB,
	}, nil
}

type interleavingImp struct {
	PrioritizeRanking PrioritizeRanking
	RankingA          Ranking
	RankingB          Ranking
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
			result = append(
				result,
				i.RankingA[rankingAIter],
				i.RankingB[rankingBIter],
			)
		} else {
			result = append(
				result,
				i.RankingB[rankingBIter],
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
