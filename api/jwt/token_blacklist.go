package jwt

import (
	"sync"
	"time"
)

type blacklistToken struct {
	token        string
	creationTime int64
}

type BlacklistTokenMap struct {
	blacklist map[string]*blacklistToken
	lock      sync.Mutex
	ttl       int64
}

// New Creates a new blacklist for invalid tokens. The timeout for the element is given as integer
func New(maxTimeToLive int, interval time.Duration) (m *BlacklistTokenMap) {
	m = &BlacklistTokenMap{
		blacklist: make(map[string]*blacklistToken),
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
func (m *BlacklistTokenMap) UpdateList() {
	m.lock.Lock()
	defer m.lock.Unlock()
	for token, mapToken := range m.blacklist {
		if m.IsExpired(mapToken) {
			delete(m.blacklist, token)
		}
	}
}

// Puts an element into the blacklist
func (m *BlacklistTokenMap) Put(token string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	it := &blacklistToken{
		token:        token,
		creationTime: time.Now().Unix(),
	}
	m.blacklist[token] = it
}

// IsBlocked checks whether a certain token is blacklisted. When the requested token is expired it does not count as blocked.
func (m *BlacklistTokenMap) IsBlocked(token string) bool {
	mapToken, ok := m.blacklist[token]
	if !ok {
		return false
	}
	return !m.IsExpired(mapToken)
}

// IsExpired returns true when the token is expired
func (m *BlacklistTokenMap) IsExpired(mapToken *blacklistToken) bool {
	return time.Now().Unix()-mapToken.creationTime > m.ttl
}

// Remove removes a specific token from the map
func (m *BlacklistTokenMap) Remove(token string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.blacklist[token]
	if ok {
		delete(m.blacklist, token)
	}
}
