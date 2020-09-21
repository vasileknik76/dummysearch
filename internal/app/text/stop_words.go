package text

type StopWordsFilter struct {
}

var _ = (Filter)((*StopWordsFilter)(nil))

var stopWords = map[string]bool{
	"а":  true,
	"и":  true,
	"с":  true,
	"в":  true,
	"на": true,
	"к":  true,
	"от": true,
	"у":  true,
	"во": true,
	"из": true,
	"о":  true,
	"об": true,
}

func (StopWordsFilter) Cleanup() {

}

func (StopWordsFilter) Filter(tokens []string) ([]string, error) {
	var res []string
	for _, token := range tokens {
		if _, ok := stopWords[token]; ok {
			continue
		}
		res = append(res, token)
	}
	return res, nil
}
