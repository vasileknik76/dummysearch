package text

import (
	"sync"

	"github.com/kljensen/snowball"
)

type stemFilter struct {
	mu        sync.Mutex
	stemCache map[string]string
	language  string
}

func NewStemFilter(language Language) Filter {
	if language == "" {
		language = LanguageEnglish
	}

	return &stemFilter{
		language:  string(language),
		stemCache: make(map[string]string),
	}
}

func (f *stemFilter) Cleanup() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.stemCache = make(map[string]string)
}

func (f *stemFilter) Filter(tokens []string) ([]string, error) {
	var res []string
	for _, token := range tokens {
		var stemmed string
		var ok bool
		var err error
		f.mu.Lock()
		stemmed, ok = f.stemCache[token]
		f.mu.Unlock()
		if !ok {
			stemmed, err = snowball.Stem(token, f.language, false)
			if err != nil {
				return res, err
			}
			f.mu.Lock()
			f.stemCache[token] = stemmed
			f.mu.Unlock()
		}
		res = append(res, stemmed)
	}
	return res, nil
}
