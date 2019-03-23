package internal

import (
	"errors"
	"reflect"
)

var (
	NotSliceError = errors.New("Passed ranking is not slice")
)

func ToRanking(rawRanking interface{}) ([]interface{}, error) {
	if rawRanking == nil {
		return nil, NotSliceError
	}
	if reflect.TypeOf(rawRanking).Kind() != reflect.Slice {
		return nil, NotSliceError
	}
	rankingValue := reflect.ValueOf(rawRanking)
	rankingSize := rankingValue.Len()

	result := make([]interface{}, rankingSize, rankingSize)
	for i := 0; i < rankingSize; i++ {
		result[i] = rankingValue.Index(i).Interface()
	}

	return result, nil
}
