package text

type Filter interface {
	Filter(tokens []string) ([]string, error)
	Cleanup()
}
