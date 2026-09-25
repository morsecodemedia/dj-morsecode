package metadata

import "sync"

type EnrichmentStore struct {
	mu sync.RWMutex

	matches map[string]EnrichmentMatch
}

func NewEnrichmentStore() *EnrichmentStore {

	return &EnrichmentStore{
		matches: make(
			map[string]EnrichmentMatch,
		),
	}

}

func (s *EnrichmentStore) Get(
	artist string,
	title string,
) (EnrichmentMatch, bool) {

	key := EnrichmentCacheKey(
		artist,
		title,
	)

	s.mu.RLock()
	defer s.mu.RUnlock()

	match, ok := s.matches[key]

	return match, ok

}

func (s *EnrichmentStore) Put(
	artist string,
	title string,
	match EnrichmentMatch,
) {

	key := EnrichmentCacheKey(
		artist,
		title,
	)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.matches[key] = match

}
