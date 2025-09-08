package infrastructure

import (
	"sync"
	"sync/atomic"
	"time"


	"github.com/go-clean/platform/logger"
)

// CacheEntry represents a local cache entry for rate limiting
type CacheEntry struct {
	Count     int64 // Use int64 for atomic operations
	Limit     int
	ResetTime int64 // Use int64 for atomic time operations (Unix nano)
}

// HybridRateLimitRepository implements rate limiting with local cache and Redis fallback
type HybridRateLimitRepository struct {
	logger          logger.Logger
	redisRepository *RedisRateLimitRepository
	localCache      sync.Map // Use sync.Map for lock-free reads
	windowSize      time.Duration
}

// NewHybridRateLimitRepository creates a new hybrid rate limit repository
func NewHybridRateLimitRepository(
	logger logger.Logger,
	redisRepository *RedisRateLimitRepository,
) *HybridRateLimitRepository {
	return &HybridRateLimitRepository{
		logger:          logger,
		redisRepository: redisRepository,
		localCache:      sync.Map{},
		windowSize:      time.Minute, // Default 1-minute window
	}
}

// RateLimit checks rate limit using local cache first, then Redis for atomic updates
func (h *HybridRateLimitRepository) RateLimit(userId string, limit int) bool {
	h.logger.Debug().Str("user_id", userId).Int("limit", limit).Msg("Checking hybrid rate limit")

	// First check local cache
	if !h.checkLocalCache(userId, limit) {
		h.logger.Debug().Str("user_id", userId).Msg("Rate limit exceeded in local cache")
		return false
	}

	// Local cache allows, now call Redis for atomic update with detail
	currentCount, ttl, err := h.redisRepository.RateLimitWithDetail(userId, limit)
	if err != nil {
		h.logger.Error().Str("user_id", userId).Err(err).Msg("Redis rate limit with detail failed, falling back to local cache")
		// Fallback: increment local cache and check
		h.incrementLocalCache(userId, limit)
		// Check local cache again after increment
		value, exists := h.localCache.Load(userId)
		if !exists {
			return true // Should not happen after increment, but safe fallback
		}
		entry := value.(*CacheEntry)
		count := atomic.LoadInt64(&entry.Count)
		return count < int64(limit)
	}

	// Update local cache with Redis values
	h.updateLocalCacheWithRedisValues(userId, limit, currentCount, ttl)

	return currentCount <= limit
}

// checkLocalCache checks if the request is allowed based on local cache
func (h *HybridRateLimitRepository) checkLocalCache(userId string, limit int) bool {
	value, exists := h.localCache.Load(userId)
	if !exists {
		// No local cache entry, allow and let Redis handle the actual check
		return true
	}

	entry := value.(*CacheEntry)
	now := time.Now().UnixNano()
	resetTime := atomic.LoadInt64(&entry.ResetTime)

	// Check if the window has expired
	if now > resetTime {
		// Window expired, allow and let Redis handle reset
		return true
	}

	// Check if limit is exceeded in local cache
	count := atomic.LoadInt64(&entry.Count)
	return count < int64(limit)
}

// incrementLocalCache increments the local cache counter for each request
// updateLocalCacheWithRedisValues updates the local cache with values from Redis
func (h *HybridRateLimitRepository) updateLocalCacheWithRedisValues(userId string, limit int, currentCount int, ttl time.Duration) {
	now := time.Now().UnixNano()
	resetTime := now + ttl.Nanoseconds()

	// Load or create cache entry
	value, exists := h.localCache.Load(userId)
	if !exists {
		// Create new entry
		entry := &CacheEntry{
			Count:     int64(currentCount),
			Limit:     limit,
			ResetTime: resetTime,
		}
		h.localCache.Store(userId, entry)
		h.logger.Debug().Str("user_id", userId).Int("count", currentCount).Int("limit", limit).Int64("reset_time", resetTime).Msg("Created new local cache entry with Redis values")
		return
	}

	// Update existing entry
	entry := value.(*CacheEntry)
	atomic.StoreInt64(&entry.Count, int64(currentCount))
	entry.Limit = limit
	atomic.StoreInt64(&entry.ResetTime, resetTime)

	h.logger.Debug().Str("user_id", userId).Int("count", currentCount).Int("limit", limit).Int64("reset_time", resetTime).Msg("Updated local cache entry with Redis values")
}

func (h *HybridRateLimitRepository) incrementLocalCache(userId string, limit int) {
	value, exists := h.localCache.Load(userId)
	if !exists {
		// Create new entry with count 1 (this request)
		entry := &CacheEntry{
			Count:     1,
			Limit:     limit,
			ResetTime: time.Now().Add(h.windowSize).UnixNano(),
		}
		h.localCache.Store(userId, entry)
		return
	}

	entry := value.(*CacheEntry)
	now := time.Now()
	nowNano := now.UnixNano()
	resetTime := entry.ResetTime

	// Check if window expired and reset if needed
	if nowNano > resetTime {
		// Create new entry with count 1 (this request)

		atomic.StoreInt64(&entry.Count, 1)
		entry.Limit = limit
		atomic.StoreInt64(&entry.ResetTime, time.Now().Add(h.windowSize).UnixNano())
		return
	}

	// Always increment count since each call represents a request
	atomic.AddInt64(&entry.Count, 1)
}

// RateLimitWithDetail checks rate limit using local cache first, then Redis for atomic updates with detailed info
func (h *HybridRateLimitRepository) RateLimitWithDetail(userId string, limit int) (int, time.Duration, error) {
	h.logger.Debug().Str("user_id", userId).Int("limit", limit).Msg("Checking hybrid rate limit with detail")

	// First check local cache
	if !h.checkLocalCache(userId, limit) {
		h.logger.Debug().Str("user_id", userId).Msg("Rate limit exceeded in local cache")
		// Return 0 remaining and get TTL from local cache if possible
		value, exists := h.localCache.Load(userId)
		if exists {
			entry := value.(*CacheEntry)
			resetTime := atomic.LoadInt64(&entry.ResetTime)
			now := time.Now().UnixNano()
			if resetTime > now {
				ttl := time.Duration(resetTime - now)
				return 0, ttl, nil
			}
		}
		return 0, 0, nil
	}

	// Local cache allows, now call Redis for atomic update with detail
	remaining, ttl, err := h.redisRepository.RateLimitWithDetail(userId, limit)
	if err != nil {
		h.logger.Error().Str("user_id", userId).Err(err).Msg("Redis rate limit with detail failed")
		return 0, 0, err
	}

	// Always increment local cache counter since each call represents a request
	h.incrementLocalCache(userId, limit)

	return remaining, ttl, nil
}

// CleanupExpiredEntries removes expired entries from local cache
func (h *HybridRateLimitRepository) CleanupExpiredEntries() {
	now := time.Now().UnixNano()
	count := 0

	h.localCache.Range(func(key, value interface{}) bool {
		userId := key.(string)
		entry := value.(*CacheEntry)
		resetTime := atomic.LoadInt64(&entry.ResetTime)

		if now > resetTime {
			h.localCache.Delete(userId)
		} else {
			count++
		}
		return true // Continue iteration
	})

	h.logger.Debug().Int("remaining_entries", count).Msg("Cleaned up expired cache entries")
}