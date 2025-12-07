package utils

import (
	"sync"
	"time"
)

// TokenBlacklist menyimpan token yang sudah di-logout
type TokenBlacklist struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

var blacklist = &TokenBlacklist{
	tokens: make(map[string]time.Time),
}

// AddToBlacklist menambahkan token ke blacklist
func AddToBlacklist(token string, expiryTime time.Time) {
	blacklist.mu.Lock()
	defer blacklist.mu.Unlock()
	blacklist.tokens[token] = expiryTime
}

// IsBlacklisted mengecek apakah token sudah di-logout
func IsBlacklisted(token string) bool {
	blacklist.mu.RLock()
	defer blacklist.mu.RUnlock()
	
	expiryTime, exists := blacklist.tokens[token]
	if !exists {
		return false
	}
	
	// Jika token belum expired, token masih di-blacklist
	if time.Now().Before(expiryTime) {
		return true
	}
	
	// Jika sudah expired, hapus dari blacklist
	delete(blacklist.tokens, expiryTime)
	return false
}

// CleanupExpiredTokens membersihkan token yang sudah expired dari blacklist
func CleanupExpiredTokens() {
	blacklist.mu.Lock()
	defer blacklist.mu.Unlock()
	
	now := time.Now()
	for token, expiryTime := range blacklist.tokens {
		if now.After(expiryTime) {
			delete(blacklist.tokens, token)
		}
	}
}
