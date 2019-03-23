package internal

import (
	"reflect"

	"github.com/rerost/information-retrieval-tools-go/interleaving/errors"
)

type Item interface {
	Key() string
}

func ToRanking(rawRanking interface{}) ([]Item, error) {
	if rawRanking == nil {
		return nil, errors.NotSliceError
	}
	if reflect.TypeOf(rawRanking).Kind() != reflect.Slice {
		return nil, errors.NotSliceError
	}
	rankingValue := reflect.ValueOf(rawRanking)
	rankingSize := rankingValue.Len()

	result := make([]Item, rankingSize, rankingSize)
	for i := 0; i < rankingSize; i++ {
		var ok bool
		result[i], ok = rankingValue.Index(i).Interface().(Item)
		if !ok {
			return nil, errors.NotImplementKeyError
		}
	}

	return result, nil
}
