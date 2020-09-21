package search

import (
	"errors"
	"math"

	"github.com/vasileknik76/dummysearch/internal/app/indexer"
)

type Searcher struct {
	i *indexer.Index
}

func NewSearcher(i *indexer.Index) Searcher {
	return Searcher{
		i: i,
	}
}

func vectorLen(v map[int]float64) float64 {
	var d float64 = 0

	for _, val := range v {
		d += val * val
	}
	return math.Sqrt(d)
}

func CosineDistance(a map[int]float64, b map[int]float64) float64 {
	r := float64(0)
	less := a
	if len(a) > len(b) {
		less = b
	}
	for k, _ := range less {
		r += a[k] * b[k]
	}
	return r / (vectorLen(a) * vectorLen(b))
}

func (s Searcher) Score(d1 int, d2 int) float64 {

	return CosineDistance(s.i.TFIDFGet(d1), s.i.TFIDFGet(d2))
}

func (s Searcher) Search(query string) (map[int]float64, error) {

	t := s.i.Tokenizer
	tokens, err := t.Tokenize(query)
	if err != nil {
		return nil, errors.New("Can't tokenize query")
	}

	var wordIds []int
	for _, token := range tokens {
		wordId := s.i.T.GetWordId(token)
		if wordId == -1 {
			continue
		}
		wordIds = append(wordIds, wordId)
	}

	tfidf := s.i.TFIDFVal(wordIds)

	res := make(map[int]float64)

	if len(tfidf) == 0 {
		return res, nil
	}

	s.i.DocsIter(func(docID int, doc *indexer.Document) {
		score := CosineDistance(tfidf, s.i.TFIDFGet(docID))
		if score == 0 || math.IsNaN(score) {
			return
		}
		res[docID] = score
	})
	return res, nil
}
