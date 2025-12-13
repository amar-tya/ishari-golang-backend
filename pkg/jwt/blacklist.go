package jwt

import (
	"sync"
	"time"
)

// InMemoryBlacklist implements TokenBlacklist using in-memory storage
// For production, consider using Redis or similar distributed cache
type InMemoryBlacklist struct {
	mu     sync.RWMutex
	tokens map[string]time.Time
}

// NewInMemoryBlacklist creates a new in-memory blacklist
func NewInMemoryBlacklist() *InMemoryBlacklist {
	bl := &InMemoryBlacklist{
		tokens: make(map[string]time.Time),
	}

	// Start cleanup goroutine
	go bl.cleanup()

	return bl
}

// Add adds a token to the blacklist
func (bl *InMemoryBlacklist) Add(token string, expiry time.Time) error {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	bl.tokens[token] = expiry
	return nil
}

// IsBlacklisted checks if token is blacklisted
func (bl *InMemoryBlacklist) IsBlacklisted(token string) bool {
	bl.mu.RLock()
	defer bl.mu.RUnlock()

	_, exists := bl.tokens[token]
	return exists
}

// cleanup removes expired tokens from blacklist periodically
func (bl *InMemoryBlacklist) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		bl.mu.Lock()
		now := time.Now()
		for token, expiry := range bl.tokens {
			if now.After(expiry) {
				delete(bl.tokens, token)
			}
		}
		bl.mu.Unlock()
	}
}
