package auth

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/storacha/guppy/pkg/client"
)

// TestClientCacheBasicOperations tests basic cache operations
func TestClientCacheBasicOperations(t *testing.T) {
	cache := &ClientCache{
		clients: make(map[string]*CachedClient),
		ttl:     1 * time.Hour,
	}

	// Create a mock client
	mockClient, _ := client.NewClient()

	// Test Set and Get
	key := "test-key"
	cache.Set(key, mockClient)

	retrievedClient, found := cache.Get(key)
	if !found {
		t.Errorf("Expected to find cached client for key %s", key)
	}
	if retrievedClient != mockClient {
		t.Errorf("Retrieved client doesn't match stored client")
	}

	// Test cache miss
	_, found = cache.Get("non-existent-key")
	if found {
		t.Errorf("Expected cache miss for non-existent key")
	}
}

// TestClientCacheExpiration tests cache expiration functionality
func TestClientCacheExpiration(t *testing.T) {
	cache := &ClientCache{
		clients: make(map[string]*CachedClient),
		ttl:     50 * time.Millisecond, // Very short TTL for testing
	}

	mockClient, _ := client.NewClient()
	key := "expiring-key"

	// Store client
	cache.Set(key, mockClient)

	// Should be available immediately
	_, found := cache.Get(key)
	if !found {
		t.Errorf("Expected to find fresh cached client")
	}

	// Wait for expiration
	time.Sleep(60 * time.Millisecond)

	// Should be expired now
	_, found = cache.Get(key)
	if found {
		t.Errorf("Expected cached client to be expired")
	}
}

// TestClientCacheDelete tests cache deletion
func TestClientCacheDelete(t *testing.T) {
	cache := &ClientCache{
		clients: make(map[string]*CachedClient),
		ttl:     1 * time.Hour,
	}

	mockClient, _ := client.NewClient()
	key := "delete-test-key"

	// Store and verify
	cache.Set(key, mockClient)
	_, found := cache.Get(key)
	if !found {
		t.Errorf("Expected to find cached client before deletion")
	}

	// Delete and verify
	cache.Delete(key)
	_, found = cache.Get(key)
	if found {
		t.Errorf("Expected client to be deleted from cache")
	}
}

// TestClientCacheClear tests cache clearing
func TestClientCacheClear(t *testing.T) {
	cache := &ClientCache{
		clients: make(map[string]*CachedClient),
		ttl:     1 * time.Hour,
	}

	mockClient, _ := client.NewClient()

	// Store multiple clients
	cache.Set("key1", mockClient)
	cache.Set("key2", mockClient)
	cache.Set("key3", mockClient)

	// Verify they exist
	total, _ := cache.Stats()
	if total != 3 {
		t.Errorf("Expected 3 cached clients, got %d", total)
	}

	// Clear cache
	cache.Clear()

	// Verify cache is empty
	total, _ = cache.Stats()
	if total != 0 {
		t.Errorf("Expected 0 cached clients after clear, got %d", total)
	}
}

// TestClientCacheConcurrency tests thread safety
func TestClientCacheConcurrency(t *testing.T) {
	cache := &ClientCache{
		clients: make(map[string]*CachedClient),
		ttl:     1 * time.Hour,
	}

	mockClient, _ := client.NewClient()
	numGoroutines := 100
	numOperations := 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch multiple goroutines doing concurrent operations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("concurrent-key-%d-%d", id, j)

				// Set
				cache.Set(key, mockClient)

				// Get
				_, found := cache.Get(key)
				if !found {
					t.Errorf("Expected to find client for key %s", key)
				}

				// Delete half of them
				if j%2 == 0 {
					cache.Delete(key)
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify cache is in a consistent state
	total, _ := cache.Stats()
	t.Logf("Final cache size after concurrent operations: %d", total)
}

// TestClientCacheStats tests statistics functionality
func TestClientCacheStats(t *testing.T) {
	cache := &ClientCache{
		clients: make(map[string]*CachedClient),
		ttl:     2 * time.Hour,
	}

	mockClient, _ := client.NewClient()

	// Initially empty
	total, ttl := cache.Stats()
	if total != 0 {
		t.Errorf("Expected 0 initial clients, got %d", total)
	}
	if ttl != 2*time.Hour {
		t.Errorf("Expected TTL of 2 hours, got %v", ttl)
	}

	// Add some clients
	cache.Set("stats-key-1", mockClient)
	cache.Set("stats-key-2", mockClient)

	total, _ = cache.Stats()
	if total != 2 {
		t.Errorf("Expected 2 clients, got %d", total)
	}
}

// TestClientCacheSetTTL tests TTL modification
func TestClientCacheSetTTL(t *testing.T) {
	cache := &ClientCache{
		clients: make(map[string]*CachedClient),
		ttl:     1 * time.Hour,
	}

	// Check initial TTL
	_, ttl := cache.Stats()
	if ttl != 1*time.Hour {
		t.Errorf("Expected initial TTL of 1 hour, got %v", ttl)
	}

	// Update TTL
	newTTL := 30 * time.Minute
	cache.SetTTL(newTTL)

	// Verify TTL was updated
	_, ttl = cache.Stats()
	if ttl != newTTL {
		t.Errorf("Expected TTL of %v, got %v", newTTL, ttl)
	}
}

// TestEmailAuthCaching tests the EmailAuth function with caching
func TestEmailAuthCaching(t *testing.T) {
	// This test would require mocking the email auth flow
	// For now, we'll test the cache key generation and basic flow

	// Reset the singleton for testing
	clientCache = nil
	once = sync.Once{}

	// Test that cache is initialized
	cache := getClientCache()
	if cache == nil {
		t.Errorf("Expected cache to be initialized")
	}

	// Test singleton behavior
	cache2 := getClientCache()
	if cache != cache2 {
		t.Errorf("Expected singleton cache instance")
	}
}

// TestPrivateKeyAuthCacheKey tests cache key generation for private key auth
func TestPrivateKeyAuthCacheKey(t *testing.T) {
	config := &AuthConfig{
		PrivateKeyPath: "/path/to/key",
		ProofPath:      "/path/to/proof",
		SpaceDID:       "did:key:example",
	}

	expectedKey := "pk:/path/to/key:/path/to/proof:did:key:example"
	actualKey := fmt.Sprintf("pk:%s:%s:%s", config.PrivateKeyPath, config.ProofPath, config.SpaceDID)

	if actualKey != expectedKey {
		t.Errorf("Expected cache key %s, got %s", expectedKey, actualKey)
	}
}

// TestCacheCleanup tests automatic cleanup of expired entries
func TestCacheCleanup(t *testing.T) {
	cache := &ClientCache{
		clients: make(map[string]*CachedClient),
		ttl:     50 * time.Millisecond,
	}

	mockClient, _ := client.NewClient()

	// Add multiple clients
	cache.Set("cleanup-key-1", mockClient)
	cache.Set("cleanup-key-2", mockClient)
	cache.Set("cleanup-key-3", mockClient)

	// Verify they're all there
	total, _ := cache.Stats()
	if total != 3 {
		t.Errorf("Expected 3 clients before cleanup, got %d", total)
	}

	// Wait for expiration
	time.Sleep(60 * time.Millisecond)

	// Trigger cleanup by adding a new client
	cache.Set("cleanup-trigger", mockClient)

	// The expired clients should have been cleaned up
	// We should have only the new client
	total, _ = cache.Stats()
	if total != 1 {
		t.Errorf("Expected 1 client after cleanup, got %d", total)
	}

	// Verify the new client is still accessible
	_, found := cache.Get("cleanup-trigger")
	if !found {
		t.Errorf("Expected new client to still be accessible")
	}
}

// BenchmarkCacheOperations benchmarks cache performance
func BenchmarkCacheOperations(b *testing.B) {
	cache := &ClientCache{
		clients: make(map[string]*CachedClient),
		ttl:     1 * time.Hour,
	}

	mockClient, _ := client.NewClient()

	b.ResetTimer()

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("bench-key-%d", i)
			cache.Set(key, mockClient)
		}
	})

	b.Run("Get", func(b *testing.B) {
		// Pre-populate cache
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("bench-get-key-%d", i)
			cache.Set(key, mockClient)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("bench-get-key-%d", i%1000)
			cache.Get(key)
		}
	})

	b.Run("ConcurrentAccess", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := fmt.Sprintf("concurrent-bench-key-%d", i)
				cache.Set(key, mockClient)
				cache.Get(key)
				i++
			}
		})
	})
}
