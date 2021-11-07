package fetcher

type Fetcher interface {
	Do() error
}
