package fether

// Fetcher is used to fetch proxies from provider, one fetcher for one provider.
type Fetcher interface {
	Do() error
}
