package text

import (
	"strings"
	"unicode"
)

type Tokenizer interface {
	Tokenize(text string) ([]string, error)
	Cleanup()
}

type tokenizer struct {
	filters []Filter
}

var _ = (Tokenizer)((*tokenizer)(nil))

func NewTokenizer(language Language) Tokenizer {
	return &tokenizer{
		filters: []Filter{NewStemFilter(language), StopWordsFilter{}},
	}
}

func (t *tokenizer) Cleanup() {
	for _, v := range t.filters {
		v.Cleanup()
	}
}

func (t *tokenizer) Tokenize(text string) ([]string, error) {
	tokens := tokenize(text)
	return t.applyFilters(tokens)
}

func (t *tokenizer) applyFilters(tokens []string) ([]string, error) {
	var err error
	for _, f := range t.filters {
		tokens, err = f.Filter(tokens)
		if err != nil {
			return tokens, err
		}
	}
	return tokens, nil
}

func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Split on any character that is not a letter or a number.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}
