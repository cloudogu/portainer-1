package jwt

import (
	"sync"
	"time"
)

type blocklistToken struct {
	token        string
	creationTime int64
}

type BlocklistTokenMap struct {
	blocklist map[string]*blocklistToken
	lock      sync.Mutex
	ttl       int64
}

// New Creates a new blocklist for invalid tokens. The timeout for the element is given as time.Duration
func NewBlocklistTokenMap(maxTimeToLive int, interval time.Duration) (m *BlocklistTokenMap) {
	m = &BlocklistTokenMap{
		blocklist: make(map[string]*blocklistToken),
		ttl:       int64(maxTimeToLive),
	}
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				m.UpdateList()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return m
}

// UpdateList updates the list and removes old elements
func (m *BlocklistTokenMap) UpdateList() {
	m.lock.Lock()
	defer m.lock.Unlock()
	for token, mapToken := range m.blocklist {
		if m.IsExpired(mapToken) {
			delete(m.blocklist, token)
		}
	}
}

// Puts an element into the blocklist
func (m *BlocklistTokenMap) Put(token string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	it := &blocklistToken{
		token:        token,
		creationTime: time.Now().Unix(),
	}
	m.blocklist[token] = it
}

// IsBlocked checks whether a certain token is blocked. When the requested token is expired it does not count as blocked.
func (m *BlocklistTokenMap) IsBlocked(token string) bool {
	mapToken, ok := m.blocklist[token]
	if !ok {
		return false
	}
	return !m.IsExpired(mapToken)
}

// IsExpired returns true when the token is expired
func (m *BlocklistTokenMap) IsExpired(token *blocklistToken) bool {
	return time.Now().Unix()-token.creationTime > m.ttl
}

// Remove removes a specific token from the map
func (m *BlocklistTokenMap) Remove(token string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.blocklist[token]
	if ok {
		delete(m.blocklist, token)
	}
}
