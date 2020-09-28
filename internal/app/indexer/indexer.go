package indexer

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/vasileknik76/dummysearch/internal/app/text"
)

var DefaultUpdatePeriod time.Duration = 180 * time.Second

type Index struct {
	cfg       *IndexConfig
	T         *Thesaurus
	Tokenizer text.Tokenizer
	ctx       context.Context

	docsMu sync.RWMutex
	Docs   map[int]*Document
	docIds map[string]int
	// key: DocId; value: Map: key: WordId; value: TFIDF
	tFIDF   map[int]map[int]float64
	tfidfMu sync.RWMutex
	// key: WordId; value: IDF
	IDF   map[int]float64
	idfMu sync.RWMutex
	// key: WordId; value: Document count
	Freq   map[int]int
	freqMu sync.RWMutex

	lastDocIdMu sync.Mutex
	LastDocId   int
	hasChanges  bool
}

type Document struct {
	TotalWords int
	// key: WordId; value: TF
	TF        map[int]float64
	MaxWordID int
	Meta      interface{}
}

func NewIndex(ctx context.Context, config *IndexConfig) *Index {
	i := &Index{
		cfg:       config,
		T:         NewThesaurus(),
		Tokenizer: text.NewTokenizer(config.Language),
		Docs:      make(map[int]*Document),
		tFIDF:     make(map[int]map[int]float64),
		IDF:       make(map[int]float64),
		Freq:      make(map[int]int),
		docIds:    make(map[string]int),
		ctx:       ctx,
		LastDocId: 1,
	}
	if config.AutoUpdate {
		go i.updateWorker()
	}
	return i
}

func (i *Index) updateWorker() {
	period := DefaultUpdatePeriod
	if i.cfg.UpdatePeriod.Seconds() != 0 {
		period = i.cfg.UpdatePeriod
	}
	t := time.NewTicker(period)
	for {
		select {
		case <-i.ctx.Done():
			return
		case <-t.C:
			if i.hasChanges {
				i.UpdateTFIDF()
				i.hasChanges = false
			}
		}
	}
}

func (i *Index) newDocumentId(id string) int {
	var docID int
	var ok bool
	i.docsMu.RLock()
	docID, ok = i.docIds[id]
	i.docsMu.RUnlock()
	if ok {
		return docID
	}

	i.lastDocIdMu.Lock()
	docID = i.LastDocId
	i.LastDocId++
	i.lastDocIdMu.Unlock()

	i.docsMu.Lock()
	i.docIds[id] = docID
	i.docsMu.Unlock()

	return docID
}

func (i *Index) getInternalID(id string) int {
	i.docsMu.RLock()
	defer i.docsMu.RUnlock()
	return i.docIds[id]
}

func (i *Index) addDocument(id int, doc *Document) {
	i.docsMu.Lock()
	i.Docs[id] = doc
	i.hasChanges = true
	i.docsMu.Unlock()
}

func (i *Index) GetDocument(id string) *Document {
	i.docsMu.RLock()
	defer i.docsMu.RUnlock()
	return i.Docs[i.getInternalID(id)]
}

func (i *Index) DeleteDocument(id string) {
	i.docsMu.Lock()
	delete(i.Docs, i.getInternalID(id))
	i.hasChanges = true
	i.docsMu.Unlock()
}

func (i *Index) TFIDFVal(wordIds []int) map[int]float64 {
	tf := float64(1) / float64(len(wordIds))
	tfidf := make(map[int]float64)
	i.idfMu.RLock()
	for _, wordId := range wordIds {
		tfidf[wordId] = i.IDF[wordId] * tf
	}
	i.idfMu.RUnlock()
	return tfidf
}

func (i *Index) AddDocument(id string, text string, meta interface{}) string {
	tokens, _ := i.Tokenizer.Tokenize(text)
	docID := i.newDocumentId(id)
	doc := &Document{len(tokens), make(map[int]float64), 0, meta}
	wordsComplete := make(map[int]bool)
	freq := make(map[int]int)
	for _, token := range tokens {
		wordId := i.T.Add(token)
		if wordId > doc.MaxWordID {
			doc.MaxWordID = wordId
		}
		_, ok := freq[wordId]
		if ok {
			freq[wordId]++
		} else {
			freq[wordId] = 1
		}

		if _, ok := wordsComplete[wordId]; !ok {
			i.freqMu.Lock()
			i.Freq[wordId]++
			i.freqMu.Unlock()
			wordsComplete[wordId] = true
		}
		doc.TF[wordId] = float64(freq[wordId]) / float64(doc.TotalWords)
	}
	i.addDocument(docID, doc)
	// i.UpdateTFIDF()
	return id
}

func (i *Index) updateIDF() {
	i.docsMu.RLock()
	d := float64(len(i.Docs))
	i.docsMu.RUnlock()
	i.freqMu.RLock()
	i.idfMu.Lock()
	for wordID := 0; wordID < i.T.NextID; wordID++ {
		i.IDF[wordID] = math.Log(d / float64(i.Freq[wordID]))
	}

	i.idfMu.Unlock()
	i.freqMu.RUnlock()
}

func (i *Index) UpdateTFIDF() {
	i.Tokenizer.Cleanup()
	i.updateIDF()
	var wg sync.WaitGroup

	i.docsMu.RLock()
	for docID, doc := range i.Docs {
		doc := doc
		docID := docID

		l := i.T.NextID
		if doc.MaxWordID < l {
			l = doc.MaxWordID
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			a := make(map[int]float64)
			i.idfMu.RLock()
			for wordID, tf := range doc.TF {
				a[wordID] = tf * i.IDF[wordID]
			}
			i.idfMu.RUnlock()
			i.tfidfMu.Lock()
			i.tFIDF[docID] = a
			i.tfidfMu.Unlock()
		}()
	}
	i.docsMu.RUnlock()
	wg.Wait()
}

func (i *Index) TFIDFGet(docID string) map[int]float64 {
	i.tfidfMu.RLock()
	defer i.tfidfMu.RUnlock()
	cp := make(map[int]float64)

	for k, v := range i.tFIDF[i.getInternalID(docID)] {
		cp[k] = v
	}
	return cp
}

func (i *Index) HasDoc(id string) bool {
	i.docsMu.RLock()
	defer i.docsMu.RUnlock()
	_, ok := i.Docs[i.getInternalID(id)]
	return ok
}

func (i *Index) DocsIter(f func(docID string, d *Document)) {
	i.docsMu.RLock()
	defer i.docsMu.RUnlock()

	for docID, _ := range i.docIds {
		iId := i.getInternalID(docID)
		f(docID, i.Docs[iId])
	}
}

func (i *Index) WordsLen() int {
	return i.T.NextID
}
