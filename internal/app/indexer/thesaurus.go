package indexer

import "sync"

type Thesaurus struct {
	mu      sync.Mutex
	Reverse map[string]int
	Words   map[int]string
	NextID  int
}

func NewThesaurus() *Thesaurus {
	return &Thesaurus{
		Reverse: make(map[string]int),
		Words:   make(map[int]string),
		NextID:  0,
	}
}

func (t *Thesaurus) Add(word string) int {
	wordId := t.GetWordId(word)
	if wordId != -1 {
		return wordId
	}
	t.mu.Lock()
	t.Reverse[word] = t.NextID
	t.Words[t.NextID] = word
	t.mu.Unlock()
	t.NextID++
	return t.NextID - 1
}

func (t *Thesaurus) GetWordId(word string) int {
	t.mu.Lock()
	val, ok := t.Reverse[word]
	t.mu.Unlock()
	if ok {
		return val
	}
	return -1
}

func (t *Thesaurus) GetWord(wordId int) string {
	t.mu.Lock()
	val, ok := t.Words[wordId]
	t.mu.Unlock()
	if ok {
		return val
	}
	return ""
}

func (t *Thesaurus) Size() int {
	return t.NextID
}
