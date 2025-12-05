package service

import (
	model "UAS_BACKEND/domain/Model"
	"UAS_BACKEND/domain/repository"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrForbidden = errors.New("forbidden")

type RBACService struct {
	Repo  *repository.RBACRepository
	Cache *PermissionCache
}

func NewRBACService(repo *repository.RBACRepository) *RBACService {
	return &RBACService{
		Repo:  repo,
		Cache: NewPermissionCache(5 * time.Minute), // Cache 5 menit
	}
}

// GetRoleByID - Mendapatkan role berdasarkan ID
func (s *RBACService) GetRoleByID(roleID uuid.UUID) (*model.Roles, error) {
	return s.Repo.GetRoleByID(roleID)
}

// GetUserPermissions - Mendapatkan permissions user (dengan cache)
func (s *RBACService) GetUserPermissions(userID uuid.UUID, roleID uuid.UUID) ([]string, error) {
	// 3. Load user permissions dari cache/database
	
	// Cek cache dulu
	if cached, found := s.Cache.Get(userID); found {
		return cached, nil
	}

	// Jika tidak ada di cache, ambil dari database
	permissions, err := s.Repo.GetUserPermissions(roleID)
	if err != nil {
		return nil, err
	}

	// Simpan ke cache
	s.Cache.Set(userID, permissions)

	return permissions, nil
}

// UserHasPermission - Check apakah user memiliki permission tertentu
func (s *RBACService) UserHasPermission(userID uuid.UUID, roleID uuid.UUID, permName string) (bool, error) {
	permissions, err := s.GetUserPermissions(userID, roleID)
	if err != nil {
		return false, err
	}

	for _, perm := range permissions {
		if perm == permName {
			return true, nil
		}
	}

	return false, nil
}

// InvalidateCache - Hapus cache permissions untuk user tertentu
func (s *RBACService) InvalidateCache(userID uuid.UUID) {
	s.Cache.Delete(userID)
}

// InvalidateAllCache - Hapus semua cache
func (s *RBACService) InvalidateAllCache() {
	s.Cache.Clear()
}

// PermissionCache - Simple in-memory cache untuk permissions
type PermissionCache struct {
	data  map[uuid.UUID]cacheEntry
	mutex sync.RWMutex
	ttl   time.Duration
}

type cacheEntry struct {
	permissions []string
	expiry      time.Time
}

func NewPermissionCache(ttl time.Duration) *PermissionCache {
	cache := &PermissionCache{
		data: make(map[uuid.UUID]cacheEntry),
		ttl:  ttl,
	}

	// Cleanup goroutine untuk hapus expired entries
	go cache.cleanupLoop()

	return cache
}

func (c *PermissionCache) Get(userID uuid.UUID) ([]string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, found := c.data[userID]
	if !found {
		return nil, false
	}

	// Check apakah sudah expired
	if time.Now().After(entry.expiry) {
		return nil, false
	}

	return entry.permissions, true
}

func (c *PermissionCache) Set(userID uuid.UUID, permissions []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[userID] = cacheEntry{
		permissions: permissions,
		expiry:      time.Now().Add(c.ttl),
	}
}

func (c *PermissionCache) Delete(userID uuid.UUID) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, userID)
}

func (c *PermissionCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[uuid.UUID]cacheEntry)
}

func (c *PermissionCache) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

func (c *PermissionCache) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for userID, entry := range c.data {
		if now.After(entry.expiry) {
			delete(c.data, userID)
		}
	}
}
